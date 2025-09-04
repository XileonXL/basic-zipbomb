# Go Zip Bomb - Build and Usage Instructions

## Overview

This Go program generates highly compressed ZIP files that expand to significantly larger sizes when decompressed, demonstrating zip bomb functionality for educational purposes.

## Features

- Creates nested ZIP archives with multiple copies per level
- Achieves compression ratios of 10x or higher (often much higher)
- Configurable parameters for size, nesting levels, and copies
- Clean, idiomatic Go code using only standard library
- Command-line interface with comprehensive options

## Requirements

- Go 1.19 or later
- Standard system with sufficient disk space

## Building the Program

### Step 1: Initialize Go Module (if not already done)
```bash
go mod init zipbomb
```

### Step 2: Compile the Program
```bash
go build -o zipbomb main.go
```

This will create an executable named `zipbomb` in the current directory.

## Usage

### Command-Line Arguments

```bash
./zipbomb [options]
```

**Required Parameters:**
- `-size int`: Size of the initial file in MB (required)

**Optional Parameters:**
- `-levels int`: Number of nested ZIP levels to create (default: 3)
- `-copies int`: Number of copies to include in each level (default: 3)
- `-output string`: Output filename for the zip bomb (default: "zipbomb.zip")

### Examples

#### Basic Usage (1MB initial, 3 levels, 3 copies each)
```bash
./zipbomb -size 1
```
This creates a zip bomb with:
- 1MB initial file
- 3 nested levels
- 3 copies per level
- Expected expansion: 1 × 3³ = 27MB
- Actual compression ratio: ~20,000x

#### Custom Configuration
```bash
./zipbomb -size 2 -levels 4 -copies 2 -output mybomb.zip
```
This creates a zip bomb with:
- 2MB initial file
- 4 nested levels
- 2 copies per level
- Expected expansion: 2 × 2⁴ = 32MB

#### High Compression Example
```bash
./zipbomb -size 5 -levels 5 -copies 4
```
This creates a zip bomb with:
- 5MB initial file
- 5 nested levels
- 4 copies per level
- Expected expansion: 5 × 4⁵ = 5,120MB (5GB)

### Verification Commands

#### Check Compressed Size
```bash
ls -lh zipbomb.zip
```

#### View ZIP Structure (be careful with large files)
```bash
unzip -l zipbomb.zip
```

#### Calculate Actual Compression Ratio
The program automatically displays the compression ratio upon completion.

## Implementation Details

### Algorithm
1. **Initial File Creation**: Creates a file filled with null bytes of specified size
2. **First Compression**: Compresses the initial file into a ZIP archive
3. **Nested Levels**: For each level:
   - Creates multiple copies of the previous ZIP file
   - Compresses all copies into a new ZIP archive
   - Removes the previous level's file
4. **Final Output**: Renames the final ZIP to the specified output name

### Compression Mechanics
- Uses null bytes (`\0`) for maximum compressibility
- ZIP compression algorithm efficiently compresses repetitive data
- Nested structure amplifies the compression effect exponentially
- Each level multiplies the effective size by the number of copies

### Safety Measures
- Cleans up temporary files automatically
- Uses efficient memory handling for large files
- Provides clear error messages for invalid parameters
- Shows expected vs actual compression ratios

## Performance Comparison with Python Reference

Testing with identical parameters (`-size 1 -levels 3 -copies 3`):

| Implementation | Compressed Size | Compression Ratio |
|----------------|-----------------|-------------------|
| Go             | 1.37 KB         | ~20,194x         |
| Python         | 1.11 KB         | ~24,892x         |

Both implementations achieve excellent compression ratios well above the 10x target.

## Safety and Legal Considerations

**⚠️ Important Notice:**
- This tool is for educational and research purposes only
- Zip bombs can consume excessive system resources when decompressed
- Never use maliciously or deploy on systems without proper authorization
- Be cautious when testing - large zip bombs can fill disk space rapidly
- Always test in controlled environments with adequate resources

## Troubleshooting

### Common Issues

1. **"go: command not found"**
   - Install Go: `apt install golang-go` (Debian/Ubuntu) or download from golang.org

2. **"build cache is required" error**
   - Set environment variables: `export HOME=/tmp && export GOCACHE=/tmp/.cache`

3. **Large file creation fails**
   - Check available disk space
   - Reduce initial size or nesting levels

4. **Memory issues with large bombs**
   - The program is designed to be memory-efficient
   - For extremely large bombs, monitor system resources

### Optimization Tips
- Start with small parameters and gradually increase
- Monitor disk space during creation
- Use compression ratios to gauge effectiveness before creating larger bombs

## Technical Notes

- Uses Go's `archive/zip` package for ZIP operations
- Implements custom `zeroReader` for efficient null byte generation
- Temporary files are automatically cleaned up
- Supports files larger than 4GB using ZIP64 format
- Cross-platform compatibility (Windows, Linux, macOS)