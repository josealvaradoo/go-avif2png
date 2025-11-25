package converter

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/gen2brain/avif"
)

// createTestAVIF creates a simple AVIF image file for testing
func createTestAVIF(t *testing.T, path string) {
	t.Helper()

	// Create a simple 10x10 red image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, red)
		}
	}

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test AVIF file: %v", err)
	}
	defer file.Close()

	// Encode as AVIF
	if err := avif.Encode(file, img); err != nil {
		t.Fatalf("failed to encode test AVIF: %v", err)
	}
}

// setupTestDir creates a temporary directory for tests
func setupTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "avif2png-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return dir
}

// ==================== AVIFToPNG Tests ====================

func TestAVIFToPNG_Success(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	outputDir := filepath.Join(testDir, "output")

	createTestAVIF(t, inputPath)

	err := AVIFToPNG(inputPath, outputDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	outputPath := filepath.Join(outputDir, "test.png")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("expected output PNG file to exist")
	}
}

func TestAVIFToPNG_SuccessVerbose(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	outputDir := filepath.Join(testDir, "output")

	createTestAVIF(t, inputPath)

	err := AVIFToPNG(inputPath, outputDir, true)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	outputPath := filepath.Join(outputDir, "test.png")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("expected output PNG file to exist")
	}
}

func TestAVIFToPNG_NonExistentInput(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "nonexistent.avif")
	outputDir := filepath.Join(testDir, "output")

	err := AVIFToPNG(inputPath, outputDir, false)

	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestAVIFToPNG_InvalidAVIF(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "invalid.avif")
	outputDir := filepath.Join(testDir, "output")

	if err := os.WriteFile(inputPath, []byte("not a valid avif file"), 0644); err != nil {
		t.Fatalf("failed to create invalid test file: %v", err)
	}

	err := AVIFToPNG(inputPath, outputDir, false)

	if err == nil {
		t.Fatal("expected error for invalid AVIF file, got nil")
	}
}

func TestAVIFToPNG_OutputAlreadyExists(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	outputDir := filepath.Join(testDir, "output")
	outputPath := filepath.Join(outputDir, "test.png")

	createTestAVIF(t, inputPath)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	if err := os.WriteFile(outputPath, []byte("existing file"), 0644); err != nil {
		t.Fatalf("failed to create existing output file: %v", err)
	}

	err := AVIFToPNG(inputPath, outputDir, false)

	if err == nil {
		t.Fatal("expected error for existing output file, got nil")
	}
}

func TestAVIFToPNG_CreatesOutputDirectory(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	outputDir := filepath.Join(testDir, "nested", "deep", "output")

	createTestAVIF(t, inputPath)

	err := AVIFToPNG(inputPath, outputDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Fatal("expected nested output directory to be created")
	}
}

func TestAVIFToPNG_FileNameWithSpaces(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "my test image.avif")
	outputDir := filepath.Join(testDir, "output")

	createTestAVIF(t, inputPath)

	err := AVIFToPNG(inputPath, outputDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	outputPath := filepath.Join(outputDir, "my test image.png")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("expected output PNG file with spaces in name to exist")
	}
}

func TestAVIFToPNG_FileNameWithSpecialChars(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test_image-2024.avif")
	outputDir := filepath.Join(testDir, "output")

	createTestAVIF(t, inputPath)

	err := AVIFToPNG(inputPath, outputDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	outputPath := filepath.Join(outputDir, "test_image-2024.png")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("expected output PNG file with special chars to exist")
	}
}

func TestAVIFToPNG_OutputPNGIsValid(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	outputDir := filepath.Join(testDir, "output")

	createTestAVIF(t, inputPath)

	err := AVIFToPNG(inputPath, outputDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify the output PNG is a valid image
	outputPath := filepath.Join(outputDir, "test.png")
	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("failed to open output PNG: %v", err)
	}
	defer file.Close()

	_, err = png.Decode(file)
	if err != nil {
		t.Fatalf("output PNG is not a valid PNG image: %v", err)
	}
}
