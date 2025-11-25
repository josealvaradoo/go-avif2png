package converter

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/gen2brain/avif"
)

// ErrFileExists is returned when an output file already exists
var ErrFileExists = errors.New("output file already exists")

// FileError represents an error that occurred while processing a specific file
type FileError struct {
	FilePath string
	Error    error
}

// ConversionResult holds the results of a bulk conversion operation
type ConversionResult struct {
	TotalFiles int
	Successful int
	Skipped    int
	Failed     int
	Errors     []FileError
}

// collectAVIFFiles scans a directory for AVIF files
// If recursive is true, it scans subdirectories as well
// Hidden files (starting with '.') are skipped
func collectAVIFFiles(rootDir string, recursive bool) ([]string, error) {
	var avifFiles []string

	if recursive {
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Skip hidden files
			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			// Check for .avif extension (case-insensitive)
			if strings.ToLower(filepath.Ext(info.Name())) == ".avif" {
				avifFiles = append(avifFiles, path)
			}

			return nil
		})
		return avifFiles, err
	}

	// Non-recursive: only scan immediate directory
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Check for .avif extension (case-insensitive)
		if strings.ToLower(filepath.Ext(entry.Name())) == ".avif" {
			avifFiles = append(avifFiles, filepath.Join(rootDir, entry.Name()))
		}
	}

	return avifFiles, nil
}

// ConvertDirectory converts all AVIF files in a directory to PNG format
// It returns a ConversionResult with statistics about the operation
func ConvertDirectory(inputDir, outputDir string, recursive, verbose bool) (*ConversionResult, error) {
	// Collect all AVIF files
	avifFiles, err := collectAVIFFiles(inputDir, recursive)
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	result := &ConversionResult{
		TotalFiles: len(avifFiles),
		Errors:     []FileError{},
	}

	// If no files found, return early
	if result.TotalFiles == 0 {
		return result, nil
	}

	if verbose {
		recursiveMsg := ""
		if recursive {
			recursiveMsg = " (recursive)"
		}
		fmt.Printf("📂 Processing directory: %s%s\n", inputDir, recursiveMsg)
		fmt.Printf("📊 Found %d AVIF file(s)\n", result.TotalFiles)
	}

	// Process each file
	for i, filePath := range avifFiles {
		if verbose {
			fmt.Printf("  [%d/%d] Converting %s... ", i+1, result.TotalFiles, filepath.Base(filePath))
		}

		err := AVIFToPNG(filePath, outputDir, false)

		if err != nil {
			if errors.Is(err, ErrFileExists) {
				// File already exists, skip it
				result.Skipped++
				if verbose {
					fmt.Println("⚠️  Skipped (already exists)")
				}
			} else {
				// Actual error occurred
				result.Failed++
				result.Errors = append(result.Errors, FileError{
					FilePath: filePath,
					Error:    err,
				})
				if verbose {
					fmt.Printf("❌ Failed: %v\n", err)
				}
			}
		} else {
			result.Successful++
			if verbose {
				fmt.Println("✅")
			}
		}
	}

	return result, nil
}

// AVIFToPNG converts an AVIF file to PNG format
func AVIFToPNG(inputPath, outputDir string, verbose bool) error {
	// Open the input AVIF file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	if verbose {
		fmt.Printf("📂 Reading: %s\n", inputPath)
	}

	// Decode the AVIF image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("failed to decode AVIF image: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate output file path
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outputPath := filepath.Join(outputDir, baseName+".png")

	// Check if output file already exists (overwrite protection)
	if _, err := os.Stat(outputPath); err == nil {
		return ErrFileExists
	}

	// Create the output PNG file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode and write PNG
	if err := png.Encode(outputFile, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	if verbose {
		fmt.Printf("✅ Saved: %s\n", outputPath)
	}

	return nil
}
