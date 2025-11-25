# 🖼️ avif2png

A fast and simple CLI tool to convert AVIF images to PNG format, written in Go.

## Features

- ✅ Convert AVIF to PNG format
- 📁 Custom output directory support
- 🛡️ Overwrite protection
- 📝 Verbose mode for detailed output
- ⚡ Fast and lightweight

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/avif2png.git
cd avif2png

# Build
make build

# Or install to $GOPATH/bin
make install
```

### Using Go

```bash
go install github.com/yourusername/avif2png/cmd/avif2png@latest
```

## Usage

```bash
# Basic usage (outputs to ./output/)
avif2png image.avif

# Custom output directory
avif2png -o ./converted image.avif
avif2png --output ./converted image.avif

# Verbose mode
avif2png -v image.avif
avif2png --verbose image.avif

# Combine flags
avif2png -v -o ./my-folder image.avif
```

## Options

| Flag        | Short | Description           | Default    |
| ----------- | ----- | --------------------- | ---------- |
| `--output`  | `-o`  | Output directory      | `./output` |
| `--verbose` | `-v`  | Enable verbose output | `false`    |

## Development

### Prerequisites

- Go 1.21 or higher

### Build

```bash
make build
```

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Project Structure

```
avif2png/
├── cmd/
│   └── avif2png/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── cli.go
│   │   └── cli_test.go
│   └── converter/
│       ├── converter.go
│       └── converter_test.go
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

## License

MIT License - see [LICENSE](LICENSE) for details.
