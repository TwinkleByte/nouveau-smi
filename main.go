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

type NvidiaInfo struct {
    chipset  string
    family   string
}

var nvidiaFamilies = map[string]NvidiaInfo{
    "NV04": {"", "NV04 family (Fahrenheit)"},
    "NV05": {"", "NV04 family (Fahrenheit)"},
    "NV0A": {"", "NV04 family (Fahrenheit)"},
    "NV10": {"", "NV10 family (Celsius)"},
    "NV11": {"", "NV10 family (Celsius)"},
    "NV15": {"", "NV10 family (Celsius)"},
    "NV17": {"", "NV10 family (Celsius)"},
    "NV18": {"", "NV10 family (Celsius)"},
    "NV1A": {"", "NV10 family (Celsius)"},
    "NV1F": {"", "NV10 family (Celsius)"},
    "NV19": {"", "NV10 family (Celsius)"},
    "NV20": {"", "NV20 family (Kelvin)"},
    "NV25": {"", "NV20 family (Kelvin)"},
    "NV28": {"", "NV20 family (Kelvin)"},
    "NV2A": {"", "NV20 family (Kelvin)"},
    "NV30": {"", "NV30 family (Rankine)"},
    "NV31": {"", "NV30 family (Rankine)"},
    "NV34": {"", "NV30 family (Rankine)"},
    "NV35": {"", "NV30 family (Rankine)"},
    "NV36": {"", "NV30 family (Rankine)"},
    "NV37": {"", "NV30 family (Rankine)"},
    "NV38": {"", "NV30 family (Rankine)"},
    "NV39": {"", "NV30 family (Rankine)"},
    "NV40": {"", "NV40 family (Curie)"},
    "NV41": {"", "NV40 family (Curie)"},
    "NV42": {"", "NV40 family (Curie)"},
    "NV43": {"", "NV40 family (Curie)"},
    "NV44": {"", "NV40 family (Curie)"},
    "NV46": {"G72", "NV40 family (Curie)"},
    "NV47": {"G70", "NV40 family (Curie)"},
    "NV49": {"G71", "NV40 family (Curie)"},
    "NV4A": {"NV44A", "NV40 family (Curie)"},
    "NV4B": {"G73", "NV40 family (Curie)"},
    "NV4C": {"MCP61", "NV40 family (Curie)"},
    "NV4E": {"C51", "NV40 family (Curie)"},
    "NV63": {"MCP73", "NV40 family (Curie)"},
    "NV67": {"MCP67", "NV40 family (Curie)"},
    "NV68": {"MCP68", "NV40 family (Curie)"},
    "NV50": {"G80", "NV50 family (Tesla)"},
    "NV84": {"G84", "NV50 family (Tesla)"},
    "NV86": {"G86", "NV50 family (Tesla)"},
    "NV92": {"G92", "NV50 family (Tesla)"},
    "NV94": {"G94", "NV50 family (Tesla)"},
    "NV96": {"G96", "NV50 family (Tesla)"},
    "NV98": {"G98", "NV50 family (Tesla)"},
    "NVA0": {"GT200", "NV50 family (Tesla)"},
    "NVA3": {"GT215", "NV50 family (Tesla)"},
    "NVA5": {"GT216", "NV50 family (Tesla)"},
    "NVA8": {"GT218", "NV50 family (Tesla)"},
    "NVAA": {"MCP77/MCP78", "NV50 family (Tesla)"},
    "NVAC": {"MCP79/MCP7A", "NV50 family (Tesla)"},
    "NVAF": {"MCP89", "NV50 family (Tesla)"},
    "NVC0": {"GF100", "NVC0 family (Fermi)"},
    "NVC1": {"GF108", "NVC0 family (Fermi)"},
    "NVC3": {"GF106", "NVC0 family (Fermi)"},
    "NVC4": {"GF104", "NVC0 family (Fermi)"},
    "NVC8": {"GF110", "NVC0 family (Fermi)"},
    "NVCE": {"GF114", "NVC0 family (Fermi)"},
    "NVCF": {"GF116", "NVC0 family (Fermi)"},
    "NVD7": {"GF117", "NVC0 family (Fermi)"},
    "NVD9": {"GF119", "NVC0 family (Fermi)"},
    "NVE4": {"GK104", "NVE0 family (Kepler)"},
    "NVE7": {"GK107", "NVE0 family (Kepler)"},
    "NVE6": {"GK106", "NVE0 family (Kepler)"},
    "NVF0": {"GK110", "NVE0 family (Kepler)"},
    "NVF1": {"GK110B", "NVE0 family (Kepler)"},
    "NV106": {"GK208B", "NVE0 family (Kepler)"},
    "NV108": {"GK208", "NVE0 family (Kepler)"},
    "NVEA": {"GK20A", "NVE0 family (Kepler)"},
    "NV110": {"GM107", "NV110 family (Maxwell)"},
    "NV117": {"GM108", "NV110 family (Maxwell)"},
    "NV118": {"GM200", "NV110 family (Maxwell)"},
    "NV120": {"GM204", "NV110 family (Maxwell)"},
    "NV124": {"GM206", "NV110 family (Maxwell)"},
    "NV126": {"GM20B", "NV110 family (Maxwell)"},
    "NV12B": {"GM20B", "NV110 family (Maxwell)"},
    "NV130": {"GP102", "NV130 family (Pascal)"},
    "NV132": {"GP104", "NV130 family (Pascal)"},
    "NV134": {"GP106", "NV130 family (Pascal)"},
    "NV136": {"GP107", "NV130 family (Pascal)"},
    "NV137": {"GP108", "NV130 family (Pascal)"},
    "NV138": {"GP108", "NV130 family (Pascal)"},
    "NV140": {"GV100", "NV140 family (Volta)"},
    "NV160": {"TU102", "NV160 family (Turing)"},
    "NV162": {"TU104", "NV160 family (Turing)"},
    "NV164": {"TU106", "NV160 family (Turing)"},
    "NV166": {"TU116", "NV160 family (Turing)"},
    "NV168": {"TU117", "NV160 family (Turing)"},
    "NV167": {"TU116", "NV160 family (Turing)"},
    "NV170": {"GA102", "NV170 family (Ampere)"},
    "NV172": {"GA104", "NV170 family (Ampere)"},
    "NV174": {"GA106", "NV170 family (Ampere)"},
    "NV176": {"GA107", "NV170 family (Ampere)"},
    "NV177": {"GA107", "NV170 family (Ampere)"},
    "NV190": {"AD102", "NV190 family (Ada Lovelace)"},
    "NV192": {"AD103", "NV190 family (Ada Lovelace)"},
    "NV193": {"AD104", "NV190 family (Ada Lovelace)"},
    "NV194": {"AD106", "NV190 family (Ada Lovelace)"},
    "NV196": {"AD106", "NV190 family (Ada Lovelace)"},
    "NV197": {"AD107", "NV190 family (Ada Lovelace)"},
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
		if strings.Contains(line, "Chipset:") {
			chipset := strings.TrimSpace(strings.Split(line, "Chipset:")[1])
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

	outputStr := out.String()

	re2 := regexp.MustCompile(`(?m)NVIDIA\s+Corporation\s+([A-Za-z0-9]+)\s+\[(.*?)\]`)
	matches := re2.FindAllStringSubmatch(outputStr, -1)

	if len(matches) == 0 {
		return "", "", fmt.Errorf("No NVIDIA GPU found")
	}

	chipsetModel := matches[0][1]
	gpuName := matches[0][2]

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
    info, found := nvidiaFamilies[codename]
	if !found {
        return "Unknown Family"
	}
    return info.family
}

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
			continue
		}

		if strings.Contains(driverLink, driverName) {
			return dir, nil
		}
	}

	return "", fmt.Errorf("could not find DRM device directory for %s driver", driverName)
}

func getBusID() string {
	drmDevicePath, err := findDrmDevicePath("nouveau")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	ueventFilePath := filepath.Join(drmDevicePath, "uevent")
	ueventData, err := os.ReadFile(ueventFilePath)
	if err != nil {
		fmt.Println("Error reading uevent file:", err)
		return ""
	}

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
			continue
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
	var rootCmd = &cobra.Command{
		Use:   "nouveau-smi",
		Short: "CLI Tool for Monitoring Nvidia GPU Using Nouveau Driver",
		Long: ` Simple Fast CLI Tool for Monitoring Nvidia GPU Using Nouveau Driver Written in Go `,
		Run: func(cmd *cobra.Command, args []string) {
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

			fmt.Println(printDate())
			fanMode, speed := getFanspeed()

			printTable(fanMode, speed)
		},
	}

	rootCmd.Flags().IntVarP(&maxSpeedFlag, "max-fan-speed", "m", 0, "Set the max fan speed. Default value 80")
	rootCmd.Flags().IntVarP(&minSpeedFlag, "min-fan-speed", "n", 0, "Set the min fan speed. Default value 40")
	rootCmd.Flags().IntVarP(&speedFlag, "fan", "f", 0, "Set the fan speed.")
	rootCmd.Flags().BoolVarP(&autoFlag, "auto", "a", false, "Set fan control to AUTO mode.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printTable(fanMode, speed string) {
	fanMode, speed = getFanspeed()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleDefault)
	t.AppendHeader(table.Row{"GPU NAME", "FAMILY CODE NAME", "CODE NAME", "GPU CHIPSET"})
	t.AppendRow(table.Row{getGpuName(), getFamilyName(getCodename()), getCodename(), getChipset()})
	t.AppendSeparator()

	t.AppendRow(table.Row{"TEMPERATURE", "BUS ID", "FAN STATUS", "FAN SPEED"})
	t.AppendSeparator()
	t.AppendRow(table.Row{getTemp(), getBusID(), fanMode, speed})
	t.AppendSeparator()

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

	t.Render()
}

func main() {
}
