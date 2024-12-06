package main

import (
	"encoding/json"
	"flag"
	"bufio"
	"bytes"
	"strconv"
	"fmt"
	"os"
	"io"
	"os/exec"
	"regexp"
	"time"
	"strings"
	"log"
	"path/filepath"
	"github.com/olekukonko/tablewriter"
)

func printDate() string {
	return time.Now().Format("Mon Jan 2 15:04:05 2006")
}

func getCodename() string {
	cmd := exec.Command("glxinfo")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting the command:", err)
		return ""
	}
	defer stdout.Close()

	scanner := bufio.NewScanner(stdout)
	var codename string

	for scanner.Scan() {
		line := scanner.Text()

		// Extract the codename
		if strings.Contains(line, "OpenGL renderer string:") {
			codename = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}

	return codename
}

func getDram() string {
	cmd := exec.Command("glxinfo")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting the command:", err)
		return ""
	}
	defer stdout.Close()

	scanner := bufio.NewScanner(stdout)
	var dedicatedMemory, availableMemory string

	for scanner.Scan() {
		line := scanner.Text()

		// Extract memory details
		if strings.Contains(line, "Dedicated video memory:") {
			dedicatedMemory = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Currently available dedicated video memory:") {
			availableMemory = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}

	// Function to convert MB to MiB
	convertToMiB := func(memory string) float64 {
		mb, err := strconv.ParseFloat(strings.TrimSuffix(memory, " MB"), 64)
		if err != nil {
			fmt.Println("Error parsing memory:", err)
			return 0
		}
		return mb * 0.953674
	}

	// Convert dedicated and available memory to MiB
	dedicatedMemoryMiB := convertToMiB(dedicatedMemory)
	availableMemoryMiB := convertToMiB(availableMemory)
	// Calculate used memory (dedicated - available)
	usedMemoryMiB := dedicatedMemoryMiB - availableMemoryMiB

	// Combine used and dedicated memory into a single string (dram)
	dram := fmt.Sprintf("%.0fMiB / %.0fMiB", usedMemoryMiB, dedicatedMemoryMiB)
	return dram
}

func getChipset() string {
	cmd := exec.Command("lspci", "-k", "-d", "::03xx")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing lspci:", err)
		return ""
	}

	// Convert output to string
	outputStr := out.String()

	// Define regex pattern to extract NVIDIA chipset model
	re2 := regexp.MustCompile(`(?m)NVIDIA\s+Corporation\s+([A-Za-z0-9]+)\s+\[(.*?)\]`)
	matches := re2.FindAllStringSubmatch(outputStr, -1)

	if len(matches) == 0 {
		fmt.Println("No NVIDIA GPU found.")
		return ""
	}

	// Extract chipsetModel
	chipsetModel := matches[0][1] // The first capture group contains the chipset model
	return chipsetModel
}

func getGpuName() string {
	cmd := exec.Command("lspci", "-k", "-d", "::03xx")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing lspci:", err)
		return ""
	}

	// Convert output to string
	outputStr := out.String()

	// Define regex pattern to extract NVIDIA GPU name
	re2 := regexp.MustCompile(`(?m)NVIDIA\s+Corporation\s+([A-Za-z0-9]+)\s+\[(.*?)\]`)
	matches := re2.FindAllStringSubmatch(outputStr, -1)

	if len(matches) == 0 {
		fmt.Println("No NVIDIA GPU found.")
		return ""
	}

	// Extract gpuName
	gpuName := matches[0][2] // The second capture group contains the GPU name
	return gpuName
}

func loadNvidiaFamilies(filename string) (map[string]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Read the JSON file
    fileData, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var nvidiaFamilies map[string]string
    err = json.Unmarshal(fileData, &nvidiaFamilies)
    if err != nil {
        return nil, err
    }

    return nvidiaFamilies, nil
}

func getFamilyName(codename string, nvidiaFamilies map[string]string) string {
    familyName, found := nvidiaFamilies[codename]
    if !found {
        familyName = "Unknown Family"
    }
    return familyName
}

func findHwmonPath(driverName string) (string, error) {
	basePath := "/sys/class/hwmon"
	hwmonDirs, err := filepath.Glob(basePath + "/hwmon*")
	if err != nil {
		return "", fmt.Errorf("error finding hwmon directories: %v", err)
	}

	for _, dir := range hwmonDirs {
		nameFile := filepath.Join(dir, "name")
		nameData, err := os.ReadFile(nameFile)
		if err != nil {
			continue // Skip directories we can't read
		}

		if strings.TrimSpace(string(nameData)) == driverName {
			return dir, nil
		}
	}

	return "", fmt.Errorf("could not find hwmon directory for %s driver", driverName)
}

func getTemp() string {
	hwmonPath, err := findHwmonPath("nouveau")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	tempFilePath := filepath.Join(hwmonPath, "temp1_input")
	tempData, err := os.ReadFile(tempFilePath)
	if err != nil {
		fmt.Println("Error reading temp1_input file:", err)
		return ""
	}

	tempMilli, err := strconv.Atoi(strings.TrimSpace(string(tempData)))
	if err != nil {
		fmt.Println("Error converting temperature value:", err)
		return ""
	}

	return fmt.Sprintf("%.1f°C", float64(tempMilli)/1000.0)
}

func getFanspeed() (string, string) {
	hwmonPath, err := findHwmonPath("nouveau")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	pwm1EnablePath := filepath.Join(hwmonPath, "pwm1_enable")
	pwm1Path := filepath.Join(hwmonPath, "pwm1")

	pwm1EnableData, err := os.ReadFile(pwm1EnablePath)
	if err != nil {
		log.Fatalf("Error reading pwm1_enable file: %v", err)
	}

	status := strings.TrimSpace(string(pwm1EnableData))
	var fanMode string
	switch status {
	case "0":
		fanMode = "NONE"
	case "1":
		fanMode = "MANUAL"
	case "2":
		fanMode = "AUTO"
	default:
		fanMode = "UNKNOWN"
	}

	pwm1Data, err := os.ReadFile(pwm1Path)
	if err != nil {
		log.Fatalf("Error reading pwm1 file: %v", err)
	}

	speed := strings.TrimSpace(string(pwm1Data))
	return fanMode, speed
}

func changeFanSpeed(speed int) {
	hwmonPath, err := findHwmonPath("nouveau")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	pwm1EnablePath := filepath.Join(hwmonPath, "pwm1_enable")
	pwm1Path := filepath.Join(hwmonPath, "pwm1")

	enableCmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo 1 > %s", pwm1EnablePath))
	err = enableCmd.Run()
	if err != nil {
		log.Fatalf("Error enabling manual fan control: %v", err)
	}

	speedStr := strconv.Itoa(speed)
	err = os.WriteFile(pwm1Path, []byte(speedStr), 0644)
	if err != nil {
		log.Fatalf("Error setting fan speed: %v", err)
	}
	fmt.Printf("Fan speed set to %d%%\n", speed)
}

func setAutoMode() {
	hwmonPath, err := findHwmonPath("nouveau")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	pwm1EnablePath := filepath.Join(hwmonPath, "pwm1_enable")

	err = os.WriteFile(pwm1EnablePath, []byte("2"), 0644)
	if err != nil {
		log.Fatalf("Error setting fan control to AUTO: %v", err)
	}
	fmt.Println("Fan control set to AUTO")
}

func printTable(gpuTemp string, codename string, dram string, chipsetModel string, gpuName string, familyName string, fanMode string, speed string) {
	// First table
	table1 := tablewriter.NewWriter(os.Stdout)
	table1.SetHeader([]string{"GPU NAME", "FAMILY CODE NAME", "CODE NAME", "GPU CHIPSET"})
	table1.SetBorder(true)   // Add borders around the table
	table1.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table1.SetRowLine(true)
	table1.Append([]string{gpuName, familyName, codename, chipsetModel})
	table1.Render()

	// Second table
	table2 := tablewriter.NewWriter(os.Stdout)
	table2.SetHeader([]string{"TEMPERATURE", "DRAM", "FAN STATUS", "FAN SPEED"})
	table2.SetBorder(true) // Add borders around the table
	table2.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table2.SetRowLine(true)

	// Add rows of data to the second table
	table2.Append([]string{gpuTemp, dram, fanMode, speed})

	// Print the second table
	table2.Render()
}

func main() {
	// Load the nvidiaFamilies from the JSON file
	nvidiaFamilies, err := loadNvidiaFamilies("nvidia_families.json")
	if err != nil {
		log.Fatalf("Error loading NVIDIA families: %v", err)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.VisitAll(func(f *flag.Flag) {
			name := f.Name
			if len(name) > 1 {
				name = "--" + name // Make sure to print long flags with --
			} else {
				name = "-" + name // For short flags, leave them as-is
			}
			fmt.Fprintf(os.Stderr, "  %-10s %s\n", name, f.Usage)
		})
	}	
	// Define the fan speed flag (long version)
	speedFlag := flag.Int("fan", 0, "Set the fan speed (40 to 80)")
	autoFlag := flag.Bool("auto", false, "Set fan control to AUTO mode")

	// Manually add short flags
	flag.IntVar(speedFlag, "f", *speedFlag, "Set the fan speed (40 to 80)")    // -f short version
	flag.BoolVar(autoFlag, "a", *autoFlag, "Set fan control to AUTO mode")      // -a short version

	// Parse the flags
	flag.Parse()

	// If the user provided a fan speed, change it
	if *autoFlag {
		setAutoMode()
	} else if *speedFlag > 0 {
		changeFanSpeed(*speedFlag)
	}
	fmt.Println(printDate())
	gpuTemp := getTemp()
	dram := getDram()
	codename := getCodename()
	gpuName := getGpuName()
	chipsetModel := getChipset()
	familyName := getFamilyName(codename, nvidiaFamilies)
	fanMode, speed := getFanspeed()
	printTable(gpuTemp, codename, dram, chipsetModel, gpuName, familyName, fanMode, speed)
}
