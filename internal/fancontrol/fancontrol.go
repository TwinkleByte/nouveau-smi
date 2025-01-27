// internal/fancontrol/fancontrol.go
package fancontrol

import (
	"fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
	"nouveau-smi/internal/hardware"
)


func GetFanspeed() (string, string) {
	hwmonPath, err := hardware.FindHwmonPath("nouveau")
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

func SetMaxFanSpeed(speed int) {
	if speed < 10 || speed > 100 {
		log.Fatalf("Error: Invalid max fan speed. It must be between 10 and 100.\n")
	}

	hwmonPath, err := hardware.FindHwmonPath("nouveau")
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

func SetMinFanSpeed(speed int) {
	if speed < 10 || speed > 100 {
		log.Fatalf("Error: Invalid min fan speed. It must be between 10 and 100.\n")
	}

	hwmonPath, err := hardware.FindHwmonPath("nouveau")
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

func ChangeFanSpeed(speed int) {
	hwmonPath, err := hardware.FindHwmonPath("nouveau")
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

func SetAutoMode() {
	hwmonPath, err := hardware.FindHwmonPath("nouveau")
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