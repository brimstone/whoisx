# whoisx

whoisx is a CLI tool and Go library that looks up owner information for IPv4 and IPv6 addresses by querying the appropriate RIR whois server (ARIN, RIPE, APNIC, AFRINIC, LACNIC) and recursing to the correct registry.

## Installation

Pre-built binaries are available on the [GitHub Releases page](https://github.com/brimstone/whoisx/releases).

Download the archive for your platform, extract, and place the `whoisx` binary on your `PATH`.

## Usage

Reads IP addresses from stdin, one per line. Output is JSON on stdout.

```console
$ whoisx <<EOF
8.8.8.8
2001:4860:4860::8888
EOF
```

### Flags

| Flag | Description |
|------|-------------|
| `-j`, `--json` | Output in JSON |
| `--debug` | Enable debug logging |

### Go Library

```go
import "github.com/brimstone/whoisx"

info, err := whoisx.IPLookup("8.8.8.8")
```

## Credits

This package is a fork from [digineo/ripego](https://github.com/digineo/ripego), which in turn is a fork from [c1982/ripego](https://github.com/c1982/ripego).
