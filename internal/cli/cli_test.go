package cli

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/gen2brain/avif"
)

// createTestAVIF creates a simple AVIF image file for testing
func createTestAVIF(t *testing.T, path string) {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, red)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test AVIF file: %v", err)
	}
	defer file.Close()

	if err := avif.Encode(file, img); err != nil {
		t.Fatalf("failed to encode test AVIF: %v", err)
	}
}

// setupTestDir creates a temporary directory for tests
func setupTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "avif2png-cli-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return dir
}

// ==================== ParseFlags Tests ====================

func TestParseFlags_ValidInput(t *testing.T) {
	args := []string{"image.avif"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.InputPath != "image.avif" {
		t.Errorf("expected InputPath 'image.avif', got: %s", config.InputPath)
	}
	if config.OutputDir != DefaultOutputDir {
		t.Errorf("expected OutputDir '%s', got: %s", DefaultOutputDir, config.OutputDir)
	}
	if config.Verbose != false {
		t.Error("expected Verbose to be false")
	}
}

func TestParseFlags_WithOutputFlag(t *testing.T) {
	args := []string{"-o", "./custom-output", "image.avif"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.OutputDir != "./custom-output" {
		t.Errorf("expected OutputDir './custom-output', got: %s", config.OutputDir)
	}
}

func TestParseFlags_WithOutputFlagLong(t *testing.T) {
	args := []string{"--output", "./custom-output", "image.avif"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.OutputDir != "./custom-output" {
		t.Errorf("expected OutputDir './custom-output', got: %s", config.OutputDir)
	}
}

func TestParseFlags_WithVerboseFlag(t *testing.T) {
	args := []string{"-v", "image.avif"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.Verbose != true {
		t.Error("expected Verbose to be true")
	}
}

func TestParseFlags_WithVerboseFlagLong(t *testing.T) {
	args := []string{"--verbose", "image.avif"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.Verbose != true {
		t.Error("expected Verbose to be true")
	}
}

func TestParseFlags_WithRecursiveFlag(t *testing.T) {
	args := []string{"-r", "my-images/"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.Recursive != true {
		t.Error("expected Recursive to be true")
	}
}

func TestParseFlags_WithRecursiveFlagLong(t *testing.T) {
	args := []string{"--recursive", "my-images/"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.Recursive != true {
		t.Error("expected Recursive to be true")
	}
}

func TestParseFlags_WithAllFlags(t *testing.T) {
	args := []string{"-v", "-r", "-o", "./out", "my-images/"}

	config, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.InputPath != "my-images/" {
		t.Errorf("expected InputPath 'my-images/', got: %s", config.InputPath)
	}
	if config.OutputDir != "./out" {
		t.Errorf("expected OutputDir './out', got: %s", config.OutputDir)
	}
	if config.Recursive != true {
		t.Error("expected Recursive to be true")
	}
	if config.Verbose != true {
		t.Error("expected Verbose to be true")
	}
}

func TestParseFlags_NoArguments(t *testing.T) {
	args := []string{}

	_, err := ParseFlags(args)

	if err == nil {
		t.Fatal("expected error for no arguments, got nil")
	}
}

func TestParseFlags_TooManyArguments(t *testing.T) {
	args := []string{"image1.avif", "image2.avif"}

	_, err := ParseFlags(args)

	if err == nil {
		t.Fatal("expected error for too many arguments, got nil")
	}
}

func TestParseFlags_InvalidFlag(t *testing.T) {
	args := []string{"--invalid-flag", "image.avif"}

	_, err := ParseFlags(args)

	if err == nil {
		t.Fatal("expected error for invalid flag, got nil")
	}
}

// ==================== ValidateInputPath Tests ====================

func TestValidateInputPath_ValidFile(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	createTestAVIF(t, inputPath)

	isDir, err := ValidateInputPath(inputPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if isDir {
		t.Error("expected isDir to be false for file")
	}
}

func TestValidateInputPath_ValidDirectory(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	isDir, err := ValidateInputPath(testDir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !isDir {
		t.Error("expected isDir to be true for directory")
	}
}

func TestValidateInputPath_NonExistent(t *testing.T) {
	_, err := ValidateInputPath("/nonexistent/path")

	if err == nil {
		t.Fatal("expected error for non-existent path, got nil")
	}
}

func TestValidateInputPath_WrongExtension(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.png")
	if err := os.WriteFile(inputPath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err := ValidateInputPath(inputPath)

	if err == nil {
		t.Fatal("expected error for wrong file extension, got nil")
	}
}

// ==================== ValidateInputFile Tests ====================

func TestValidateInputFile_ValidFile(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	createTestAVIF(t, inputPath)

	err := ValidateInputFile(inputPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateInputFile_NonExistent(t *testing.T) {
	err := ValidateInputFile("/nonexistent/path/image.avif")

	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestValidateInputFile_WrongExtension(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.png")
	if err := os.WriteFile(inputPath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := ValidateInputFile(inputPath)

	if err == nil {
		t.Fatal("expected error for wrong extension, got nil")
	}
}

func TestValidateInputFile_UppercaseExtension(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.AVIF")
	if err := os.WriteFile(inputPath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := ValidateInputFile(inputPath)
	if err != nil {
		t.Fatalf("expected no error for uppercase extension, got: %v", err)
	}
}

// ==================== Run Tests ====================

func TestRun_Success(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.avif")
	outputDir := filepath.Join(testDir, "output")

	createTestAVIF(t, inputPath)

	config := &Config{
		InputPath: inputPath,
		OutputDir: outputDir,
		Recursive: false,
		Verbose:   false,
	}

	err := Run(config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRun_InvalidInputFile(t *testing.T) {
	config := &Config{
		InputPath: "/nonexistent/image.avif",
		OutputDir: "./output",
		Recursive: false,
		Verbose:   false,
	}

	err := Run(config)

	if err == nil {
		t.Fatal("expected error for invalid input, got nil")
	}
}

func TestRun_WrongExtension(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	inputPath := filepath.Join(testDir, "test.jpg")
	if err := os.WriteFile(inputPath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	config := &Config{
		InputPath: inputPath,
		OutputDir: filepath.Join(testDir, "output"),
		Verbose:   false,
	}

	err := Run(config)

	if err == nil {
		t.Fatal("expected error for wrong extension, got nil")
	}
}

func TestRun_DirectoryConversion(t *testing.T) {
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

	config := &Config{
		InputPath: inputDir,
		OutputDir: outputDir,
		Recursive: false,
		Verbose:   false,
	}

	err := Run(config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify output files exist
	if _, err := os.Stat(filepath.Join(outputDir, "image1.png")); os.IsNotExist(err) {
		t.Error("expected image1.png to exist")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "image2.png")); os.IsNotExist(err) {
		t.Error("expected image2.png to exist")
	}
}

func TestRun_RecursiveDirectoryConversion(t *testing.T) {
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

	config := &Config{
		InputPath: inputDir,
		OutputDir: outputDir,
		Recursive: true,
		Verbose:   false,
	}

	err := Run(config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify both files were converted
	if _, err := os.Stat(filepath.Join(outputDir, "image1.png")); os.IsNotExist(err) {
		t.Error("expected image1.png to exist")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "image2.png")); os.IsNotExist(err) {
		t.Error("expected image2.png to exist (flattened)")
	}
}
