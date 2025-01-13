# Nouveau-SMI
`nouveau-smi` is a tool for monitoring NVIDIA GPUs using the Nouveau driver. It provides real-time information about the GPU’s status, such as temperature, and fan settings.

### Prerequisites
- **Git**: Required to clone the repository
- **Go**: Version 1.23.4 or later is required. Verify by running `go version`.
- **Nouveau Driver**: The Nouveau driver for NVIDIA GPUs must be installed and active.
- **Go Modules**:
  - `go-pretty` v6.6.4 or later by jedib0t
  - `cobra` v1.8.1 or later by spf13

### Example output
```
Wed Dec 11 16:30:33 2024
+----------------+----------------------+------------+-------------+
|    GPU NAME    |   FAMILY CODE NAME   |  CODE NAME | GPU CHIPSET |
+----------------+----------------------+------------+-------------+
| GeForce GT 710 | NVE0 family (Kepler) |    NV106   |    GK208B   |
+----------------+----------------------+------------+-------------+
|   TEMPERATURE  |        BUS ID        | FAN STATUS |  FAN SPEED  |
+----------------+----------------------+------------+-------------+
|     45.0°C     |     0000:01:00.0     |    AUTO    |      44     |
+----------------+----------------------+------------+-------------+
```
### Clone the Repository and Build the Tool:
```bash
git clone https://github.com/TwinkleByte/nouveau-smi.git
cd nouveau-smi
go build -o nouveau-smi main.go
```
### Put the binary to $PATH:
```bash
sudo install -m 755 nouveau-smi /usr/local/bin/
nouveau-smi
```
### Uninstall
```
sudo rm /usr/local/bin/nouveau-smi
```
### Monitor Nouveau GPU status every second:
```
watch -n1 --no-title nouveau-smi
```
### Usage:
```
 Simple Fast CLI Tool for Monitoring Nvidia GPU Using Nouveau Driver Written in Go

Usage:
  nouveau-smi [flags]

Flags:
  -a, --auto                Set fan control to AUTO mode.
  -f, --fan int             Set the fan speed.
  -h, --help                help for nouveau-smi
  -m, --max-fan-speed int   Set the max fan speed. Default value 80
  -n, --min-fan-speed int   Set the min fan speed. Default value 40
```
