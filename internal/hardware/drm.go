// internal/hardware/drm.go
package hardware

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func FindDrmDevicePath(driverName string) (string, error) {
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