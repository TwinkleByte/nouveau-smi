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
### **Installation**

#### 1. Install the CLI tool:
```bash
go install github.com/TwinkleByte/nouveau-smi/cmd/nouveau-smi@latest
```

#### 2. Add the Go binary directory to your `PATH`:
The binary is installed to `$GOPATH/bin` (default: `~/go/bin`) or `$GOBIN` if set. Add the appropriate directory to your shell’s `PATH`:

**Bash**:
```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

**Zsh**:
```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

**Fish** (preferred method):
```fish
fish_add_path (go env GOPATH)/bin
```
*If you don’t have `fish_add_path` (older Fish versions):*
```fish
echo 'set -gx PATH $PATH (go env GOPATH)/bin' >> ~/.config/fish/config.fish
```

#### 3. Verify Installation:
```bash
nouveau-smi --version
```

---

### **Troubleshooting**
- **Command not found?**
  - Ensure `$GOPATH/bin` or `$GOBIN` is in your `PATH`:
    ```bash
    echo $PATH | grep "$(go env GOPATH)/bin"
    ```
  - If you’ve set `$GOBIN`, use that path instead of `$GOPATH/bin`.
  - Restart your shell or run `source ~/.bashrc`/`source ~/.zshrc`.

- **Avoid duplicates in `PATH`:**
  If you’ve already added the path, don’t run the command again. Check with:
  ```bash
  echo $PATH
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
