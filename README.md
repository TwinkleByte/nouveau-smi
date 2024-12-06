# Nouveau-SMI

`nouveau-smi` is a tool for monitoring NVIDIA GPUs using the Nouveau driver. It provides real-time information about the GPU’s status, such as temperature, DRAM usage, and fan settings.

### Prerequisites
- **Git**: Required to clone the repository
- **Go**: Version 1.23.4 or later is required. Verify by running `go version`.
- **Nouveau Driver**: The Nouveau driver for NVIDIA GPUs must be installed and active.
- **Go Modules**:
  - `tablewriter` v0.0.5 by olekukonko
- **mesa-utils**: Required for DRAM info (note: doesn't work in TTY).

### Example output
```
Sat Dec 7 02:36:45 2024
+----------------+----------------------+------------+-------------+
|    GPU NAME    |   FAMILY CODE NAME   |  CODE NAME | GPU CHIPSET |
+----------------+----------------------+------------+-------------+
| GeForce GT 710 | NVE0 family (Kepler) |    NV106   |    GK208B   |
+----------------+----------------------+------------+-------------+
| ·············· | ···················· |  ········· | ··········· |
+----------------+----------------------+------------+-------------+
|   TEMPERATURE  |         DRAM         | FAN STATUS |  FAN SPEED  |
+----------------+----------------------+------------+-------------+
|     47.0°C     |   389MiB / 1945MiB   |    AUTO    |      46     |
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
Usage of nouveau-smi:
Options:
  -a, --auto   Set fan control to AUTO mode.
  -f, --fan    Set the fan speed (range: 40 to 80).
```
