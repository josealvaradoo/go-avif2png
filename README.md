# 🖼️ avif2png

A fast and simple CLI tool to convert AVIF images to PNG format, written in Go.

## Features

- ✅ Convert AVIF to PNG format
- 📁 Bulk directory conversion (with optional recursive mode)
- 🛡️ Overwrite protection (automatically skips existing files)
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

### Single File Conversion

```bash
# Basic usage (outputs to ./output/)
avif2png image.avif

# Custom output directory
avif2png -o ./converted image.avif
avif2png --output ./converted image.avif

# Verbose mode
avif2png -v image.avif
avif2png --verbose image.avif
```

### Bulk Directory Conversion

```bash
# Convert all AVIF files in a directory (non-recursive)
avif2png my-images/
avif2png -o ./converted my-images/

# Recursive mode (includes subdirectories)
avif2png -r my-images/
avif2png --recursive my-images/

# Combine flags
avif2png -r -v -o ./converted my-images/
```

### Output Structure

When converting directories, all PNG files are saved directly to the output directory with a flattened structure:

```
input/
  ├── photo1.avif
  └── subfolder/
      └── photo2.avif

# After: avif2png -r input/ -o output/

output/
  ├── photo1.png
  └── photo2.png  (flattened, not in subfolder)
```

## Options

| Flag          | Short | Description                         | Default    |
| ------------- | ----- | ----------------------------------- | ---------- |
| `--output`    | `-o`  | Output directory                    | `./output` |
| `--recursive` | `-r`  | Recursively process subdirectories  | `false`    |
| `--verbose`   | `-v`  | Enable verbose output               | `false`    |

## Behavior

- **Overwrite Protection**: Existing PNG files are automatically skipped (not overwritten)
- **Hidden Files**: Files starting with `.` are ignored
- **Case Insensitive**: Accepts both `.avif` and `.AVIF` extensions
- **Flattened Output**: Directory conversion outputs all files to a single directory (no subdirectories)

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
