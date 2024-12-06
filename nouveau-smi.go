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
	"log"
	"path/filepath"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

func getCodename() string {
	cmd := exec.Command("grep", "-i", "chipset", "/var/log/Xorg.0.log")
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

		// Search for the "Chipset" line and extract the chipset model
		if strings.Contains(line, "Chipset:") {
			// Extract the part after "Chipset: "
			chipset := strings.TrimSpace(strings.Split(line, "Chipset:")[1])
			// Remove "NVIDIA " if it exists and also remove surrounding quotes
			codename = strings.Replace(chipset, "NVIDIA ", "", 1)
			codename = strings.Replace(codename, "\"", "", -1)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}

	if codename == "" {
		fmt.Println("Error: Chipset information not found.")
		return ""
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

		// Extract memory details if available
		if strings.Contains(line, "Dedicated video memory:") {
			dedicatedMemory = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Currently available dedicated video memory:") {
			availableMemory = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}

	// If either memory value is empty, print a warning
	if dedicatedMemory == "" || availableMemory == "" {
		fmt.Println("Warning: Missing memory information.")
		return "Unknown memory"
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

func parseLspciOutput() (string, string, error) {
	cmd := exec.Command("lspci", "-k", "-d", "::03xx")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("Error executing lspci: %v", err)
	}

	// Convert output to string
	outputStr := out.String()

	// Define regex pattern to extract NVIDIA chipset model and GPU name
	re2 := regexp.MustCompile(`(?m)NVIDIA\s+Corporation\s+([A-Za-z0-9]+)\s+\[(.*?)\]`)
	matches := re2.FindAllStringSubmatch(outputStr, -1)

	if len(matches) == 0 {
		return "", "", fmt.Errorf("No NVIDIA GPU found")
	}

	// Extract chipsetModel and gpuName
	chipsetModel := matches[0][1] // The first capture group contains the chipset model
	gpuName := matches[0][2]     // The second capture group contains the GPU name

	return chipsetModel, gpuName, nil
}

func getChipset() string {
	chipsetModel, _, err := parseLspciOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return chipsetModel
}

func getGpuName() string {
	_, gpuName, err := parseLspciOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return gpuName
}

func getFamilyName(codename string) string {
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

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.VisitAll(func(f *flag.Flag) {
			// Show the short and long flags together
			var flagStr string
			if len(f.Name) > 1 {
				flagStr = fmt.Sprintf("--%s", f.Name) // Long flag
			} else {
				flagStr = fmt.Sprintf("-%s", f.Name) // Short flag
			}

			// Add short flag to the usage message if present
			if f.Name == "auto" {
				fmt.Fprintf(os.Stderr, "  -a, %s   %s\n", flagStr, f.Usage)
			} else if f.Name == "fan" {
				fmt.Fprintf(os.Stderr, "  -f, %s   %s\n", flagStr, f.Usage)
			}
		})
	}

	// Define the fan speed flag (long version)
	speedFlag := flag.Int("fan", 0, "Set the fan speed (range: 40 to 80).")
	autoFlag := flag.Bool("auto", false, "Set fan control to AUTO mode.")

	// Manually add short flags
	flag.IntVar(speedFlag, "f", *speedFlag, "Set the fan speed (range: 40 to 80).")    // -f short version
	flag.BoolVar(autoFlag, "a", *autoFlag, "Set fan control to AUTO mode.")           // -a short version

	// Parse the flags
	flag.Parse()

	// If the user provided a fan speed, change it
	if *autoFlag {
		setAutoMode()
	} else if *speedFlag > 0 {
		changeFanSpeed(*speedFlag)
	}

	// Print out system information
	fmt.Println(printDate())
	
	// Create a new table writer
	tCombined := table.NewWriter()
	tCombined.SetOutputMirror(os.Stdout)
	tCombined.SetStyle(table.StyleDefault)
	tCombined.AppendHeader(table.Row{"GPU NAME", "FAMILY CODE NAME", "CODE NAME", "GPU CHIPSET"})
	tCombined.AppendRow(table.Row{getGpuName(), getFamilyName(getCodename()), getCodename(), getChipset()})
	tCombined.AppendSeparator()

	// Calculate column widths dynamically
	columnWidths := []int{
		len("GPU NAME"),
		len("FAMILY CODE NAME"),
		len("CODE NAME"),
		len("GPU CHIPSET"),
	}

	fanMode, speed := getFanspeed()

	// Add the maximum widths from the data rows
	dataRows := []table.Row{
		{getGpuName(), getFamilyName(getCodename()), getCodename(), getChipset()},
		{getTemp(), getDram(), fanMode, speed},
	}

	for _, row := range dataRows {
		for i, col := range row {
			// Convert col to string and then calculate length
			colStr := fmt.Sprintf("%v", col)
			if len(colStr) > columnWidths[i] {
				columnWidths[i] = len(colStr)
			}
		}
	}

	// Create a separator row with `=` that matches the column widths
	separator := []string{}
	for _, width := range columnWidths {
		separator = append(separator, strings.Repeat("·", width))
	}

	// Add the separator row
	tCombined.AppendRow(table.Row{separator[0], separator[1], separator[2], separator[3]})
	tCombined.AppendSeparator()

	// Add second section
	tCombined.AppendRow(table.Row{"TEMPERATURE", "DRAM", "FAN STATUS", "FAN SPEED"})
	tCombined.AppendSeparator()
	tCombined.AppendRow(table.Row{getTemp(), getDram(), fanMode, speed})

	// Define column configurations to center the text in each column
	tCombined.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:  "GPU NAME",
			Align: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:  "FAMILY CODE NAME",
			Align: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:  "CODE NAME",
			Align: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:  "GPU CHIPSET",
			Align: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
	})

	// Render combined table
	tCombined.Render()
}