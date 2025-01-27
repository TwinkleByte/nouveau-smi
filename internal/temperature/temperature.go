// internal/temperature/temperature.go
package temperature

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"nouveau-smi/internal/hardware"
)

func GetTemp() string {
	hwmonPath, err := hardware.FindHwmonPath("nouveau")
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