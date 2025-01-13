package main

import (
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
	"github.com/spf13/cobra"
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

// findDrmDevicePath searches for the DRM device directory related to the Nouveau driver.
func findDrmDevicePath(driverName string) (string, error) {
	basePath := "/sys/class/drm"
	cardDirs, err := filepath.Glob(basePath + "/card*/device")
	if err != nil {
		return "", fmt.Errorf("error finding DRM card directories: %v", err)
	}

	for _, dir := range cardDirs {
		driverPath := filepath.Join(dir, "driver")
		driverLink, err := os.Readlink(driverPath)
		if err != nil {
			continue // Skip directories we can't read
		}

		if strings.Contains(driverLink, driverName) {
			return dir, nil
		}
	}

	return "", fmt.Errorf("could not find DRM device directory for %s driver", driverName)
}

func getBusID() string {
	// Find the DRM device path for the nouveau driver
	drmDevicePath, err := findDrmDevicePath("nouveau")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Read the uevent file to get the PCI_SLOT_NAME
	ueventFilePath := filepath.Join(drmDevicePath, "uevent")
	ueventData, err := os.ReadFile(ueventFilePath)
	if err != nil {
		fmt.Println("Error reading uevent file:", err)
		return ""
	}

	// Find the PCI_SLOT_NAME in the uevent data
	lines := strings.Split(string(ueventData), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PCI_SLOT_NAME=") {
			return strings.TrimPrefix(line, "PCI_SLOT_NAME=")
		}
	}

	return "PCI_SLOT_NAME not found"
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

func setMaxFanSpeed(speed int) {
	if speed < 10 || speed > 100 {
		log.Fatalf("Error: Invalid max fan speed. It must be between 10 and 100.\n")
	}

	hwmonPath, err := findHwmonPath("nouveau")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	pwm1MaxPath := filepath.Join(hwmonPath, "pwm1_max")

	speedStr := strconv.Itoa(speed)
	err = os.WriteFile(pwm1MaxPath, []byte(speedStr), 0644)
	if err != nil {
		log.Fatalf("Error setting max fan speed: %v", err)
	}
	fmt.Printf("Max fan speed set to %d%%\n", speed)
}

func setMinFanSpeed(speed int) {
	if speed < 10 || speed > 100 {
		log.Fatalf("Error: Invalid min fan speed. It must be between 10 and 100.\n")
	}

	hwmonPath, err := findHwmonPath("nouveau")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	pwm1MinPath := filepath.Join(hwmonPath, "pwm1_min")

	// Check if min speed is greater than max speed
	maxSpeedPath := filepath.Join(hwmonPath, "pwm1_max")
	maxSpeedData, err := os.ReadFile(maxSpeedPath)
	if err != nil {
		log.Fatalf("Error reading max fan speed: %v", err)
	}

	maxSpeed, err := strconv.Atoi(strings.TrimSpace(string(maxSpeedData)))
	if err != nil {
		log.Fatalf("Error parsing max fan speed: %v", err)
	}

	if speed > maxSpeed {
		log.Fatalf("Error: Min fan speed cannot be greater than max fan speed. Either lower your value or change max fan speed.\n")
	}

	speedStr := strconv.Itoa(speed)
	err = os.WriteFile(pwm1MinPath, []byte(speedStr), 0644)
	if err != nil {
		log.Fatalf("Error setting min fan speed: %v", err)
	}
	fmt.Printf("Min fan speed set to %d%%\n", speed)
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

var (
	maxSpeedFlag int
	minSpeedFlag int
	speedFlag int
	autoFlag  bool
)

func init() {
	// Set up the root command
	var rootCmd = &cobra.Command{
		Use:   "nouveau-smi",
		Short: "CLI Tool for Monitoring Nvidia GPU Using Nouveau Driver",
		Long: ` Simple Fast CLI Tool for Monitoring Nvidia GPU Using Nouveau Driver Written in Go `,
		Run: func(cmd *cobra.Command, args []string) {
			// If the user provided the --auto flag, set to AUTO mode
			if autoFlag {
				setAutoMode()
			} else {
				if speedFlag > 0 {
					changeFanSpeed(speedFlag)
				}
				if maxSpeedFlag > 0 {
					setMaxFanSpeed(maxSpeedFlag)
				}
				if minSpeedFlag > 0 {
					setMinFanSpeed(minSpeedFlag)
				}
			}			

			// Print out system information
			fmt.Println(printDate())
			fanMode, speed := getFanspeed()

			// Create and display the table
			printTable(fanMode, speed)
		},
	}

	// Define flags and their shorthand
	rootCmd.Flags().IntVarP(&maxSpeedFlag, "max-fan-speed", "m", 0, "Set the max fan speed. Default value 80")
	rootCmd.Flags().IntVarP(&minSpeedFlag, "min-fan-speed", "n", 0, "Set the min fan speed. Default value 40")
	rootCmd.Flags().IntVarP(&speedFlag, "fan", "f", 0, "Set the fan speed.")
	rootCmd.Flags().BoolVarP(&autoFlag, "auto", "a", false, "Set fan control to AUTO mode.")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printTable(fanMode, speed string) {
	fanMode, speed = getFanspeed()

	// Create a new table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleDefault)
	t.AppendHeader(table.Row{"GPU NAME", "FAMILY CODE NAME", "CODE NAME", "GPU CHIPSET"})
	t.AppendRow(table.Row{getGpuName(), getFamilyName(getCodename()), getCodename(), getChipset()})
	t.AppendSeparator()

	// Add second section
	t.AppendRow(table.Row{"TEMPERATURE", "BUS ID", "FAN STATUS", "FAN SPEED"})
	t.AppendSeparator()
	t.AppendRow(table.Row{getTemp(), getBusID(), fanMode, speed})
	t.AppendSeparator()

	// Define column configurations to center the text in each column
	t.SetColumnConfigs([]table.ColumnConfig{
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
	t.Render()
}

func main() {
}