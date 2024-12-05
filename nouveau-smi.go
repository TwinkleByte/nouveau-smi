package main

import (
	"flag"
	"bufio"
	"bytes"
	"strconv"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
	"strings"
	"io/ioutil"
	"log"
	"github.com/olekukonko/tablewriter"
)

// Map of NVIDIA families with code names as keys and family names as values
var nvidiaFamilies = map[string]string{
	"NV04": "NV04 family (Fahrenheit)",
	"NV05": "NV04 family (Fahrenheit)",
	"NV0A": "NV04 family (Fahrenheit)",
	"NV10": "NV10 family (Celsius)",
	"NV11": "NV10 family (Celsius)",
	"NV15": "NV10 family (Celsius)",
	"NV17": "NV10 family (Celsius)",
	"NV18": "NV10 family (Celsius)",
	"NV1A": "NV10 family (Celsius)",
	"NV1F": "NV10 family (Celsius)",
	"NV19": "NV10 family (Celsius)",
	"NV20": "NV20 family (Kelvin)",
	"NV25": "NV20 family (Kelvin)",
	"NV28": "NV20 family (Kelvin)",
	"NV2A": "NV20 family (Kelvin)",
	"NV30": "NV30 family (Rankine)",
	"NV31": "NV30 family (Rankine)",
	"NV34": "NV30 family (Rankine)",
	"NV35": "NV30 family (Rankine)",
	"NV36": "NV30 family (Rankine)",
	"NV37": "NV30 family (Rankine)",
	"NV39": "NV30 family (Rankine)",
	"NV38": "NV30 family (Rankine)",
	"NV40": "NV40 family (Curie)",
	"NV41": "NV40 family (Curie)",
	"NV42": "NV40 family (Curie)",
	"NV43": "NV40 family (Curie)",
	"NV44": "NV40 family (Curie)",
	"NV46": "NV40 family (Curie)",
	"NV47": "NV40 family (Curie)",
	"NV49": "NV40 family (Curie)",
	"NV4A": "NV40 family (Curie)",
	"NV4B": "NV40 family (Curie)",
	"NV4C": "NV40 family (Curie)",
	"NV4E": "NV40 family (Curie)",
	"NV63": "NV40 family (Curie)",
	"NV67": "NV40 family (Curie)",
	"NV68": "NV40 family (Curie)",
	"NV50": "NV50 family (Tesla)",
	"NV84": "NV50 family (Tesla)",
	"NV86": "NV50 family (Tesla)",
	"NV92": "NV50 family (Tesla)",
	"NV94": "NV50 family (Tesla)",
	"NV96": "NV50 family (Tesla)",
	"NV98": "NV50 family (Tesla)",
	"NVA0": "NV50 family (Tesla)",
	"NVA3": "NV50 family (Tesla)",
	"NVA5": "NV50 family (Tesla)",
	"NVA8": "NV50 family (Tesla)",
	"NVAA": "NV50 family (Tesla)",
	"NVAC": "NV50 family (Tesla)",
	"NVAF": "NV50 family (Tesla)",
	"NVC0": "NVC0 family (Fermi)",
	"NVC1": "NVC0 family (Fermi)",
	"NVC3": "NVC0 family (Fermi)",
	"NVC4": "NVC0 family (Fermi)",
	"NVC8": "NVC0 family (Fermi)",
	"NVCE": "NVC0 family (Fermi)",
	"NVCF": "NVC0 family (Fermi)",
	"NVD7": "NVC0 family (Fermi)",
	"NVD9": "NVC0 family (Fermi)",
	"NVE0": "NVE0 family (familyName)",
	"NVE4": "NVE0 family (Kepler)",
	"NVE7": "NVE0 family (Kepler)",
	"NVE6": "NVE0 family (Kepler)",
	"NVF0": "NVE0 family (Kepler)",
	"NVF1": "NVE0 family (Kepler)",
	"NV106": "NVE0 family (Kepler)",
	"NV108": "NVE0 family (Kepler)",
	"NVEA": "NVE0 family (Kepler)",
	"NV110": "NV110 family (Maxwell)",
	"NV117": "NV110 family (Maxwell)",
	"NV118": "NV110 family (Maxwell)",
	"NV120": "NV110 family (Maxwell)",
	"NV124": "NV110 family (Maxwell)",
	"NV126": "NV110 family (Maxwell)",
	"NV12B": "NV110 family (Maxwell)",
	"NV130": "NV130 family (Pascal)",
	"NV132": "NV130 family (Pascal)",
	"NV134": "NV130 family (Pascal)",
	"NV136": "NV130 family (Pascal)",
	"NV137": "NV130 family (Pascal)",
	"NV138": "NV130 family (Pascal)",
	"NV140": "NV140 family (Volta)",
	"NV160": "NV160 family (Turing)",
	"NV162": "NV160 family (Turing)",
	"NV164": "NV160 family (Turing)",
	"NV166": "NV160 family (Turing)",
	"NV168": "NV160 family (Turing)",
	"NV167": "NV160 family (Turing)",
	"NV170": "NV170 family (Ampere)",
	"NV172": "NV170 family (Ampere)",
	"NV174": "NV170 family (Ampere)",
	"NV176": "NV170 family (Ampere)",
	"NV177": "NV170 family (Ampere)",
	"NV190": "NV190 family (Ada Lovelace)",
	"NV192": "NV190 family (Ada Lovelace)",
	"NV193": "NV190 family (Ada Lovelace)",
	"NV194": "NV190 family (Ada Lovelace)",
	"NV196": "NV190 family (Ada Lovelace)",
	"NV197": "NV190 family (Ada Lovelace)",
}

func printDate() string {
	return time.Now().Format("Mon Jan 2 15:04:05 2006")
}

func printTable() {
	// Run the sensors command and capture the output
	cmd := exec.Command("sensors")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running sensors:", err)
		return
	}	

	// Convert the output to a string
	sensorsOutput := string(output)

	// Define the regex pattern for extracting GPU temperature
	re := regexp.MustCompile(`nouveau-pci-[^\n]+.*temp1:\s*\+([0-9.]+)°C`)

	// Find the first match in the output
	match := re.FindStringSubmatch(sensorsOutput)
	if len(match) < 2 {
		// Adjust for possible line break after temp1: 
		re = regexp.MustCompile(`nouveau-pci-[^\n]+[\s\S]*temp1:\s*\+([0-9.]+)°C`)
		match = re.FindStringSubmatch(sensorsOutput)
	}

	if len(match) < 2 {
		fmt.Println("GPU temperature not found")
		return
	}

	// Assign value directly to gpuTemp
	gpuTemp := match[1] + "°C"

	cmd = exec.Command("glxinfo")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting the command:", err)
		return
	}
	defer stdout.Close()

	scanner := bufio.NewScanner(stdout)
	var codename, dedicatedMemory, availableMemory string

	for scanner.Scan() {
		line := scanner.Text()

		// Extract the codename
		if strings.Contains(line, "OpenGL renderer string:") {
			codename = strings.TrimSpace(strings.Split(line, ":")[1])

		// Extract memory details
		} else if strings.Contains(line, "Dedicated video memory:") {
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
		var mb float64
		fmt.Sscanf(memory, "%f MB", &mb)
		return mb * 0.953674
	}

	// Convert dedicated and available memory to MiB
	dedicatedMemoryMiB := convertToMiB(dedicatedMemory)
	availableMemoryMiB := convertToMiB(availableMemory)
	// Calculate used memory (dedicated - available)
	usedMemoryMiB := dedicatedMemoryMiB - availableMemoryMiB

	// Combine used and dedicated memory into a single string (dram)
	dram := fmt.Sprintf("%.0fMiB / %.0fMiB", usedMemoryMiB, dedicatedMemoryMiB)
	
	// Get family name using codename from the nvidiaFamilies map
	familyName, found := nvidiaFamilies[codename]
	if !found {
		familyName = "Unknown Family"
	}

	// Run the lspci command with the specified filter
	cmd = exec.Command("lspci", "-k", "-d", "::03xx")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing lspci:", err)
		return
	}

	// Convert output to string
	outputStr := out.String()

	// Define regex pattern to extract NVIDIA chipset model and GPU name
	re2 := regexp.MustCompile(`(?m)NVIDIA\s+Corporation\s+([A-Za-z0-9]+)\s+\[(.*?)\]`)
	matches := re2.FindAllStringSubmatch(outputStr, -1)

	if len(matches) == 0 {
		fmt.Println("No NVIDIA GPU found.")
		return
	}

	// Extract chipsetModel and gpuName
	chipsetModel := matches[0][1] // The first capture group contains the chipset model
	gpuName := matches[0][2]      // The second capture group contains the GPU name

	// Path to the pwm1_enable file
	filePath := "/sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0/hwmon/hwmon1/pwm1_enable"

	// Read the pwm1_enable file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading pwm1_enable file: %v", err)
	}

	// Trim whitespace from the result
	status := strings.TrimSpace(string(data))

	// Determine the fan control mode based on the value of pwm1_enable
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

	// Path to the pwm1 value (you might need to adjust this based on your system)
	pwm1Path := "/sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0/hwmon/hwmon1/pwm1"

	// Read the file contents
	data, err = ioutil.ReadFile(pwm1Path)
	if err != nil {
		log.Fatalf("Error reading pwm1 file: %v", err)
	}

	// Remove any extra spaces or newline characters
	speed := strings.TrimSpace(string(data))

	// First table
	table1 := tablewriter.NewWriter(os.Stdout)
	table1.SetHeader([]string{"GPU NAME       ", "FAMILY CODE NAME", "CODE NAME", "GPU CHIPSET"})
	table1.SetBorder(true)   // Add borders around the table
	table1.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table1.SetRowLine(true)
	table1.Append([]string{gpuName, familyName, codename, chipsetModel})
	table1.Render()

	// Second table
	table2 := tablewriter.NewWriter(os.Stdout)
	table2.SetHeader([]string{"TEMPERATURE", "DRAM                ", "  FAN STATUS ", " FAN SPEED "})
	table2.SetBorder(true) // Add borders around the table
	table2.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table2.SetRowLine(true)

	// Add rows of data to the second table
	table2.Append([]string{gpuTemp, dram, fanMode, speed})

	// Print the second table
	table2.Render()
}

func changeFanSpeed(speed int) {
	// Path to the pwm1 file for manual control
	pwm1Path := "/sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0/hwmon/hwmon1/pwm1"

	// Convert the speed to a string
	speedStr := strconv.Itoa(speed)

	// Write the new fan speed to the pwm1 file
	err := ioutil.WriteFile(pwm1Path, []byte(speedStr), 0644)
	if err != nil {
		log.Fatalf("Error setting fan speed: %v", err)
	}
	fmt.Printf("Fan speed set to %d%%\n", speed)
}


func setAutoMode() {
	// Path to the pwm1_enable file to enable AUTO mode
	pwm1EnablePath := "/sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0/hwmon/hwmon1/pwm1_enable"

	// Set the pwm1_enable value to AUTO (2)
	err := ioutil.WriteFile(pwm1EnablePath, []byte("2"), 0644)
	if err != nil {
		log.Fatalf("Error setting fan control to AUTO: %v", err)
	}
	fmt.Println("Fan control set to AUTO")
}

func main() {
	// Define the fan speed flag
	speedFlag := flag.Int("f", 0, "Set the fan speed (0 to 100)")
	autoFlag := flag.Bool("auto", false, "Set fan control to AUTO mode")
	flag.Parse()

	// If the user provided a fan speed, change it
	if *autoFlag {
		setAutoMode()
	} else if *speedFlag > 0 {
		changeFanSpeed(*speedFlag)
	}
	fmt.Println(printDate())
	printTable()
}
