# Getting Started

## 🔧 Installation

#### macOS / Linux (Homebrew)

```bash
brew tap jaxxstorm/tap
brew install tscli          # upgrades via ‘brew upgrade’
```

#### Windows (Scoop)

```powershell
scoop bucket add jaxxstorm https://github.com/jaxxstorm/scoop-bucket.git
scoop install tscli
```

#### Nix

```bash
nix shell github:jaxxstorm/tscli
```

#### Manual download

Pre-built archives for **macOS, Linux, Windows (x86-64 / arm64)** are published on every release:

```bash
# install the newest tscli (Linux/macOS, amd64/arm64)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
esac

curl -sSL "$(curl -sSL \
  https://api.github.com/repos/jaxxstorm/tscli/releases/latest \
  | grep -oE "https.*tscli_.*_${OS}_${ARCH}\.tar\.gz" \
  | head -n1)" \
| sudo tar -xz -C /usr/local/bin tscli
```

#### Go install (always builds from HEAD)

```bash
go install github.com/jaxxstorm/tscli@latest
```

After any method, confirm:

```bash
tscli --version
```

## First call

```bash
TAILSCALE_API_KEY=tskey-xxx \
TAILSCALE_TAILNET=example.com \
tscli list devices
```

## Global flags

- `--api-key`, `-k`: Tailscale API key
- `--tailnet`, `-n`: tailnet name (defaults to `-`)
- `--output`, `-o`: `json`, `yaml`, `human`, or `pretty`
- `--debug`, `-d`: dump raw HTTP traffic

