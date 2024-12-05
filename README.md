# Nouveau-SMI
`nouveau-smi` is a tool made by me for monitoring NVIDIA GPUs with the Nouveau driver. It provides essential GPU information, such as temperature, DRAM usage, fan status, and more.

## Installation Guide

### Prerequisites
- Go: Version 1.23.4 or later.
- **Nouveau Driver**: The tool requires the Nouveau driver for NVIDIA GPUs. Make sure it is installed and active on your system.

### Clone the Repository and Build the tool
```bash
git clone https://github.com/TwinkleByte/nouveau-smi.git
cd nouveau-smi
go mod init nouveau-smi
go get github.com/olekukonko/tablewriter
go build -o nouveau-smi nouveau-smi.go
./nouveau-smi
```

You can control your GPU's fan speed using the following commands:  
- `sudo nouveau -f 40` (sets fan speed to 40%)  
- `sudo nouveau -auto` (enables automatic fan speed control)  
