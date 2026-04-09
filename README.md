# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Features

- 🔍 Continuously monitors TCP/UDP ports on your system
- 🚨 Alerts when unexpected ports open or close
- 📋 Whitelist known services to reduce noise
- 🪶 Minimal resource footprint
- 📊 Export monitoring data in JSON format

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch
go build -o portwatch
```

## Usage

Start monitoring with default settings:

```bash
portwatch start
```

Monitor with a custom whitelist:

```bash
portwatch start --whitelist 80,443,22,3000
```

Check current open ports:

```bash
portwatch scan
```

View monitoring logs:

```bash
portwatch logs
```

## Configuration

Create a `portwatch.yaml` in `~/.config/portwatch/`:

```yaml
interval: 5s
whitelist:
  - 22    # SSH
  - 80    # HTTP
  - 443   # HTTPS
alert_on_close: false
```

## License

MIT License - see [LICENSE](LICENSE) for details.