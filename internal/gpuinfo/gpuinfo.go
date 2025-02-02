// internal/gpuinfo/gpuinfo.go
package gpuinfo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TwinkleByte/nouveau-smi/internal/hardware"
	"github.com/TwinkleByte/nouveaunvinfo"
)

func GetFamilyName(codename string) string {
	info, found := nouveaunvinfo.NvidiaFamilies[codename]
	if !found {
		return "Unknown Family"
	}
	return info.Family
}

func GetBusID() string {
	drmDevicePath, err := hardware.FindDrmDevicePath("nouveau")
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

func GetChipset() string {
	chipsetModel, _, err := hardware.ParseLspciOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return chipsetModel
}

func GetGpuName() string {
	_, gpuName, err := hardware.ParseLspciOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return gpuName
}

func GetCodename() string {
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