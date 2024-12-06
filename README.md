# Nouveau-SMI
`nouveau-smi` is a tool made by me for monitoring NVIDIA GPUs with the Nouveau driver. It provides essential GPU information, such as temperature, DRAM usage, fan status, and more.

## Installation Guide

### Prerequisites
- Git: To clone the repo
- Go: Version 1.23.4 or later.
- **Nouveau Driver**: The tool requires the Nouveau driver for NVIDIA GPUs. Make sure it is installed and active on your system.
- Go module named "tablewriter v0.0.5" by olekukonko and go-runewidth v0.0.9 by mattn
- mesa-utils

### Clone the Repository and Build the tool
```bash
git clone https://github.com/TwinkleByte/nouveau-smi.git
cd nouveau-smi
go build -o nouveau-smi nouveau-smi.go
sudo mv nouveau-smi /usr/local/bin/
nouveau-smi
```
### Usage
```
Usage of nouveau-smi:
Options:
  -a         Set fan control to AUTO mode
  --auto     Set fan control to AUTO mode
  -f         Set the fan speed (40 to 80)
  --fan      Set the fan speed (40 to 80)
```
