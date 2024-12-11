# Nouveau-SMI

Note: This project is still under active development, and significant changes are expected. The current version does not follow any formal versioning, and changes may be pushed frequently.

`nouveau-smi` is a tool for monitoring NVIDIA GPUs using the Nouveau driver. It provides real-time information about the GPU’s status, such as temperature, DRAM usage, and fan settings.

### Prerequisites
- **Git**: Required to clone the repository
- **Go**: Version 1.23.4 or later is required. Verify by running `go version`.
- **Nouveau Driver**: The Nouveau driver for NVIDIA GPUs must be installed and active.
- **Go Modules**:
  - `go-pretty` v6.6.4 or later by jedib0t
  - `cobra` v1.8.1 or later by spf13
- **mesa-utils**: Required for DRAM info (note: doesn't work in TTY unless you have a display running).

### Example output
```
Tue Dec 10 00:38:44 2024
+----------------+----------------------+------------+-------------+
|    GPU NAME    |   FAMILY CODE NAME   |  CODE NAME | GPU CHIPSET |
+----------------+----------------------+------------+-------------+
| GeForce GT 710 | NVE0 family (Kepler) |    NV106   |    GK208B   |
+----------------+----------------------+------------+-------------+
|   TEMPERATURE  |         DRAM         | FAN STATUS |  FAN SPEED  |
+----------------+----------------------+------------+-------------+
|     46.0°C     |   389MiB / 1945MiB   |    AUTO    |      44     |
+----------------+----------------------+------------+-------------+
```
### Clone the Repository and Build the Tool
```bash
git clone https://github.com/TwinkleByte/nouveau-smi.git
cd nouveau-smi
go build -o nouveau-smi nouveau-smi.go
sudo install -m 755 nouveau-smi /usr/local/bin/
nouveau-smi
```
### Uninstall
```
sudo rm /usr/local/bin/nouveau-smi
```
### Monitor Nouveau GPU status every second
```
watch -n1 --no-title nouveau-smi
```
### Usage
```
nouveau-smi allows you to control the fan speed of your NVIDIA GPU
and view system information such as temperature, fan status, etc.

Usage:
  nouveau-smi [flags]

Flags:
  -a, --auto      Set fan control to AUTO mode.
  -f, --fan int   Set the fan speed (range: 40 to 80).
  -h, --help      help for nouveau-smi
```
