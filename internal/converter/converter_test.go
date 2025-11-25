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

// ==================== collectAVIFFiles Tests ====================

func TestCollectAVIFFiles_SingleDirectory(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Create test files
	createTestAVIF(t, filepath.Join(testDir, "image1.avif"))
	createTestAVIF(t, filepath.Join(testDir, "image2.avif"))
	if err := os.WriteFile(filepath.Join(testDir, "other.png"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	files, err := collectAVIFFiles(testDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 AVIF files, got: %d", len(files))
	}
}

func TestCollectAVIFFiles_RecursiveDirectory(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Create test files in root
	createTestAVIF(t, filepath.Join(testDir, "image1.avif"))

	// Create subdirectory with files
	subDir := filepath.Join(testDir, "subfolder")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}
	createTestAVIF(t, filepath.Join(subDir, "image2.avif"))
	createTestAVIF(t, filepath.Join(subDir, "image3.avif"))

	// Non-recursive should find only 1 file
	files, err := collectAVIFFiles(testDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("non-recursive: expected 1 AVIF file, got: %d", len(files))
	}

	// Recursive should find all 3 files
	files, err = collectAVIFFiles(testDir, true)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("recursive: expected 3 AVIF files, got: %d", len(files))
	}
}

func TestCollectAVIFFiles_SkipsHiddenFiles(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Create visible file
	createTestAVIF(t, filepath.Join(testDir, "image.avif"))

	// Create hidden file
	createTestAVIF(t, filepath.Join(testDir, ".hidden.avif"))

	files, err := collectAVIFFiles(testDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 AVIF file (hidden should be skipped), got: %d", len(files))
	}
}

func TestCollectAVIFFiles_MixedFileTypes(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Create various file types
	createTestAVIF(t, filepath.Join(testDir, "image1.avif"))
	createTestAVIF(t, filepath.Join(testDir, "image2.AVIF")) // uppercase extension
	if err := os.WriteFile(filepath.Join(testDir, "photo.jpg"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "doc.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	files, err := collectAVIFFiles(testDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 AVIF files, got: %d", len(files))
	}
}

func TestCollectAVIFFiles_EmptyDirectory(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	files, err := collectAVIFFiles(testDir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("expected 0 AVIF files in empty directory, got: %d", len(files))
	}
}

func TestCollectAVIFFiles_NonExistentDirectory(t *testing.T) {
	_, err := collectAVIFFiles("/nonexistent/directory", false)

	if err == nil {
		t.Fatal("expected error for non-existent directory, got nil")
	}
}

// ==================== ConvertDirectory Tests ====================

func TestConvertDirectory_Success(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputDir := filepath.Join(testDir, "input")
	outputDir := filepath.Join(testDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("failed to create input dir: %v", err)
	}

	// Create test files
	createTestAVIF(t, filepath.Join(inputDir, "image1.avif"))
	createTestAVIF(t, filepath.Join(inputDir, "image2.avif"))

	result, err := ConvertDirectory(inputDir, outputDir, false, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.TotalFiles != 2 {
		t.Errorf("expected 2 total files, got: %d", result.TotalFiles)
	}
	if result.Successful != 2 {
		t.Errorf("expected 2 successful conversions, got: %d", result.Successful)
	}
	if result.Skipped != 0 {
		t.Errorf("expected 0 skipped files, got: %d", result.Skipped)
	}
	if result.Failed != 0 {
		t.Errorf("expected 0 failed conversions, got: %d", result.Failed)
	}

	// Verify output files exist
	if _, err := os.Stat(filepath.Join(outputDir, "image1.png")); os.IsNotExist(err) {
		t.Error("expected image1.png to exist")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "image2.png")); os.IsNotExist(err) {
		t.Error("expected image2.png to exist")
	}
}

func TestConvertDirectory_SkipsExistingFiles(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputDir := filepath.Join(testDir, "input")
	outputDir := filepath.Join(testDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("failed to create input dir: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create test files
	createTestAVIF(t, filepath.Join(inputDir, "image1.avif"))
	createTestAVIF(t, filepath.Join(inputDir, "image2.avif"))

	// Create existing output file
	if err := os.WriteFile(filepath.Join(outputDir, "image1.png"), []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	result, err := ConvertDirectory(inputDir, outputDir, false, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Successful != 1 {
		t.Errorf("expected 1 successful conversion, got: %d", result.Successful)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped file, got: %d", result.Skipped)
	}
	if result.Failed != 0 {
		t.Errorf("expected 0 failed conversions, got: %d", result.Failed)
	}
}

func TestConvertDirectory_EmptyDirectory(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputDir := filepath.Join(testDir, "input")
	outputDir := filepath.Join(testDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("failed to create input dir: %v", err)
	}

	result, err := ConvertDirectory(inputDir, outputDir, false, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.TotalFiles != 0 {
		t.Errorf("expected 0 total files, got: %d", result.TotalFiles)
	}
}

func TestConvertDirectory_RecursiveMode(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputDir := filepath.Join(testDir, "input")
	outputDir := filepath.Join(testDir, "output")
	subDir := filepath.Join(inputDir, "subfolder")

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	// Create test files
	createTestAVIF(t, filepath.Join(inputDir, "image1.avif"))
	createTestAVIF(t, filepath.Join(subDir, "image2.avif"))

	result, err := ConvertDirectory(inputDir, outputDir, true, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.TotalFiles != 2 {
		t.Errorf("expected 2 total files (recursive), got: %d", result.TotalFiles)
	}
	if result.Successful != 2 {
		t.Errorf("expected 2 successful conversions, got: %d", result.Successful)
	}
}

func TestConvertDirectory_FlattenStructure(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputDir := filepath.Join(testDir, "input")
	outputDir := filepath.Join(testDir, "output")
	subDir := filepath.Join(inputDir, "subfolder")

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	// Create test file in subdirectory
	createTestAVIF(t, filepath.Join(subDir, "nested.avif"))

	result, err := ConvertDirectory(inputDir, outputDir, true, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Successful != 1 {
		t.Errorf("expected 1 successful conversion, got: %d", result.Successful)
	}

	// Verify output is flattened (directly in output dir, not in subfolder)
	if _, err := os.Stat(filepath.Join(outputDir, "nested.png")); os.IsNotExist(err) {
		t.Error("expected nested.png to exist in flat output directory")
	}

	// Verify it's NOT in a subdirectory
	if _, err := os.Stat(filepath.Join(outputDir, "subfolder", "nested.png")); !os.IsNotExist(err) {
		t.Error("expected output to be flattened, not preserve directory structure")
	}
}

func TestConvertDirectory_PartialFailure(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputDir := filepath.Join(testDir, "input")
	outputDir := filepath.Join(testDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("failed to create input dir: %v", err)
	}

	// Create one valid AVIF
	createTestAVIF(t, filepath.Join(inputDir, "valid.avif"))

	// Create one invalid AVIF
	if err := os.WriteFile(filepath.Join(inputDir, "invalid.avif"), []byte("not a valid avif"), 0644); err != nil {
		t.Fatalf("failed to create invalid file: %v", err)
	}

	result, err := ConvertDirectory(inputDir, outputDir, false, false)
	if err != nil {
		t.Fatalf("expected no error from ConvertDirectory, got: %v", err)
	}

	if result.Successful != 1 {
		t.Errorf("expected 1 successful conversion, got: %d", result.Successful)
	}
	if result.Failed != 1 {
		t.Errorf("expected 1 failed conversion, got: %d", result.Failed)
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error in result, got: %d", len(result.Errors))
	}
}
