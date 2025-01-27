// internal/hardware/hwmon.go
package hardware

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func FindHwmonPath(driverName string) (string, error) {
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