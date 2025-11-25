package cli

import (
	"avif2png/internal/converter"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultOutputDir = "./output"
)

// Config holds the CLI configuration
type Config struct {
	InputPath string
	OutputDir string
	Verbose   bool
}

// ParseFlags parses command line arguments and returns a Config
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("avif2png", flag.ContinueOnError)

	outputDir := fs.String("output", DefaultOutputDir, "Output directory for converted PNG files")
	fs.StringVar(outputDir, "o", DefaultOutputDir, "Output directory (shorthand)")

	verbose := fs.Bool("verbose", false, "Enable verbose output")
	fs.BoolVar(verbose, "v", false, "Enable verbose output (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "🖼️  AVIF to PNG Converter\n\n")
		fmt.Fprintf(os.Stderr, "Usage: avif2png [options] <input.avif>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  avif2png image.avif\n")
		fmt.Fprintf(os.Stderr, "  avif2png -o ./converted image.avif\n")
		fmt.Fprintf(os.Stderr, "  avif2png --verbose image.avif\n")
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remainingArgs := fs.Args()
	if len(remainingArgs) != 1 {
		return nil, errors.New("exactly one input file is required")
	}

	return &Config{
		InputPath: remainingArgs[0],
		OutputDir: *outputDir,
		Verbose:   *verbose,
	}, nil
}

// ValidateInputFile validates that the input file exists and has .avif extension
func ValidateInputFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", path)
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".avif" {
		return fmt.Errorf("input file must have .avif extension, got: %s", ext)
	}

	return nil
}

// Run executes the main application logic
func Run(config *Config) error {
	if err := ValidateInputFile(config.InputPath); err != nil {
		return err
	}

	return converter.AVIFToPNG(config.InputPath, config.OutputDir, config.Verbose)
}
