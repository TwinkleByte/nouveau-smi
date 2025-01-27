// internal/table/table.go
package table

import (
	"os"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/TwinkleByte/nouveau-smi/internal/gpuinfo"
	"github.com/TwinkleByte/nouveau-smi/internal/temperature"
)

func PrintTable(fanMode, speed string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleDefault)

	// GPU Information
	t.AppendHeader(table.Row{"GPU NAME", "FAMILY CODE NAME", "CODE NAME", "GPU CHIPSET"})
	t.AppendRow(table.Row{
		gpuinfo.GetGpuName(),
		gpuinfo.GetFamilyName(gpuinfo.GetCodename()),
		gpuinfo.GetCodename(),
		gpuinfo.GetChipset(),
	})
	t.AppendSeparator()

	// Temperature and Fan Information
	t.AppendRow(table.Row{"TEMPERATURE", "BUS ID", "FAN STATUS", "FAN SPEED"})
	t.AppendSeparator()
	t.AppendRow(table.Row{
		temperature.GetTemp(),
		gpuinfo.GetBusID(),
		fanMode,
		speed,
	})
	t.AppendSeparator()

	// Table Formatting
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "GPU NAME",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "FAMILY CODE NAME",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "CODE NAME",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "GPU CHIPSET",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
	})

	t.Render()
}