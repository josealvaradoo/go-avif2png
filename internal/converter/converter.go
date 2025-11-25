package converter

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/gen2brain/avif"
)

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
		return fmt.Errorf("output file already exists: %s (skipping to prevent overwrite)", outputPath)
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
