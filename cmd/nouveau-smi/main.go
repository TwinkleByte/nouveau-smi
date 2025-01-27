
package main

import (
	"fmt"
	"os"
	"time"
	"github.com/spf13/cobra"
	"github.com/TwinkleByte/nouveau-smi/internal/fancontrol"
	"github.com/TwinkleByte/nouveau-smi/internal/table"
)

func printDate() string {
	return time.Now().Format("Mon Jan 2 15:04:05 2006")
}

var (
	maxSpeedFlag int
	minSpeedFlag int
	speedFlag    int
	autoFlag     bool
	version = "v1.1.0"

)

func main() {
	var rootCmd = &cobra.Command{
	  Use:   "nouveau-smi",
	  Short: "CLI Tool for Monitoring Nvidia GPU Using Nouveau Driver",
	  Long:  `Simple Fast CLI Tool for Monitoring Nvidia GPU Using Nouveau Driver Written in Go`,
	  Run: func(cmd *cobra.Command, args []string) {
			if autoFlag {
				fancontrol.SetAutoMode()
			} else {
				if speedFlag > 0 {
					fancontrol.ChangeFanSpeed(speedFlag)
				}
				if maxSpeedFlag > 0 {
					fancontrol.SetMaxFanSpeed(maxSpeedFlag)
				}
				if minSpeedFlag > 0 {
					fancontrol.SetMinFanSpeed(minSpeedFlag)
				}
			}

			fmt.Println(printDate())
			fanMode, speed := fancontrol.GetFanspeed()

			table.PrintTable(fanMode, speed)
		},
	}

	rootCmd.Version = version
    rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")
    rootCmd.Flags().BoolP("version", "v", false, "Print version information")
	rootCmd.Flags().IntVarP(&maxSpeedFlag, "max-fan-speed", "m", 0, "Set the max fan speed. Default value 80")
	rootCmd.Flags().IntVarP(&minSpeedFlag, "min-fan-speed", "n", 0, "Set the min fan speed. Default value 40")
	rootCmd.Flags().IntVarP(&speedFlag, "fan", "f", 0, "Set the fan speed.")
	rootCmd.Flags().BoolVarP(&autoFlag, "auto", "a", false, "Set fan control to AUTO mode.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
