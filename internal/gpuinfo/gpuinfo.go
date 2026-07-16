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
    // On multi-GPU laptops the default GL context is the integrated GPU (no NV codename); query the NVIDIA card directly.
    cmd := exec.Command("glxinfo")
    cmd.Env = append(os.Environ(), "DRI_PRIME=1")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println("Error:", err)
        return ""
    }
    if err := cmd.Start(); err != nil {
        fmt.Println("Error starting glxinfo command:", err)
        return ""
    }
    defer stdout.Close()

    scanner := bufio.NewScanner(stdout)
    var codename string

    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "OpenGL renderer") {
            // Example line: "OpenGL renderer string: NV106"
            parts := strings.Split(line, ":")
            if len(parts) < 2 {
                continue
            }
            renderer := strings.TrimSpace(parts[1])
            // Look for the "NV" substring in the renderer info
            if idx := strings.Index(renderer, "NV"); idx != -1 {
                codename = renderer[idx:]
                codename = strings.Fields(codename)[0] // extract the first word (e.g., "NV106")
                break
            }
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading glxinfo output:", err)
    }

    if codename == "" {
        return ""
    }

    return codename
}
