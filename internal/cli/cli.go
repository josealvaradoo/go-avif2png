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
	Recursive bool
	Verbose   bool
}

// ParseFlags parses command line arguments and returns a Config
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("avif2png", flag.ContinueOnError)

	outputDir := fs.String("output", DefaultOutputDir, "Output directory for converted PNG files")
	fs.StringVar(outputDir, "o", DefaultOutputDir, "Output directory (shorthand)")

	recursive := fs.Bool("recursive", false, "Recursively process subdirectories")
	fs.BoolVar(recursive, "r", false, "Recursively process subdirectories (shorthand)")

	verbose := fs.Bool("verbose", false, "Enable verbose output")
	fs.BoolVar(verbose, "v", false, "Enable verbose output (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "🖼️  AVIF to PNG Converter\n\n")
		fmt.Fprintf(os.Stderr, "Usage: avif2png [options] <input.avif or directory>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Convert single file\n")
		fmt.Fprintf(os.Stderr, "  avif2png image.avif\n")
		fmt.Fprintf(os.Stderr, "  avif2png -o ./converted image.avif\n\n")
		fmt.Fprintf(os.Stderr, "  # Convert directory\n")
		fmt.Fprintf(os.Stderr, "  avif2png my-images/\n")
		fmt.Fprintf(os.Stderr, "  avif2png -r my-images/\n")
		fmt.Fprintf(os.Stderr, "  avif2png -r -o ./converted my-images/\n\n")
		fmt.Fprintf(os.Stderr, "  # Verbose mode\n")
		fmt.Fprintf(os.Stderr, "  avif2png --verbose image.avif\n")
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remainingArgs := fs.Args()
	if len(remainingArgs) != 1 {
		return nil, errors.New("exactly one input file or directory is required")
	}

	return &Config{
		InputPath: remainingArgs[0],
		OutputDir: *outputDir,
		Recursive: *recursive,
		Verbose:   *verbose,
	}, nil
}

// ValidateInputPath validates that the input path exists and is either a valid file or directory
// Returns true if the path is a directory, false if it's a file
func ValidateInputPath(path string) (isDir bool, err error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("input path does not exist: %s", path)
	}
	if err != nil {
		return false, fmt.Errorf("failed to access input path: %w", err)
	}

	// If it's a directory, just verify it's readable
	if info.IsDir() {
		return true, nil
	}

	// If it's a file, check extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".avif" {
		return false, fmt.Errorf("input file must have .avif extension, got: %s", ext)
	}

	return false, nil
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

// runSingleFileConversion handles conversion of a single AVIF file
func runSingleFileConversion(config *Config) error {
	return converter.AVIFToPNG(config.InputPath, config.OutputDir, config.Verbose)
}

// runDirectoryConversion handles conversion of all AVIF files in a directory
func runDirectoryConversion(config *Config) error {
	result, err := converter.ConvertDirectory(config.InputPath, config.OutputDir, config.Recursive, config.Verbose)
	if err != nil {
		return err
	}

	// Print summary for non-verbose mode
	if !config.Verbose && result.TotalFiles > 0 {
		if result.Failed > 0 || result.Skipped > 0 {
			fmt.Printf("✅ Converted %d/%d files", result.Successful, result.TotalFiles)
			if result.Skipped > 0 {
				fmt.Printf(" (%d skipped - already exist)", result.Skipped)
			}
			if result.Failed > 0 {
				fmt.Printf(" (%d failed)", result.Failed)
			}
			fmt.Println()
		} else {
			fmt.Printf("✅ Converted %d file(s)\n", result.Successful)
		}
	}

	// Print verbose summary
	if config.Verbose && result.TotalFiles > 0 {
		fmt.Printf("\n📊 Summary: %d successful, %d skipped, %d failed\n",
			result.Successful, result.Skipped, result.Failed)
	}

	// Print error details
	if len(result.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "\n❌ Failed conversions:\n")
		for _, fileErr := range result.Errors {
			fmt.Fprintf(os.Stderr, "  - %s: %v\n", filepath.Base(fileErr.FilePath), fileErr.Error)
		}
		return fmt.Errorf("completed with %d error(s)", len(result.Errors))
	}

	// If no files were found
	if result.TotalFiles == 0 {
		fmt.Println("⚠️  No AVIF files found in directory")
	}

	return nil
}

// Run executes the main application logic
func Run(config *Config) error {
	isDir, err := ValidateInputPath(config.InputPath)
	if err != nil {
		return err
	}

	if isDir {
		return runDirectoryConversion(config)
	}
	return runSingleFileConversion(config)
}
