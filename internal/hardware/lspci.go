// internal/hardware/lspci.go
package hardware

import (
    "bytes"
    "fmt"
    "os/exec"
    "regexp"
)

func ParseLspciOutput() (string, string, error) {
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