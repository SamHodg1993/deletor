package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func FormatSize(bytes int64) string {
	const (
		KB = 1 << 10 // 1024
		MB = 1 << 20 // 1024 * 1024
		GB = 1 << 30 // 1024 * 1024 * 1024
		TB = 1 << 40 // 1024 * 1024 * 1024 * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func ExpandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, path[1:])
}

func ToBytes(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(strings.ToLower(sizeStr))

	var unitIndex int
	for unitIndex = 0; unitIndex < len(sizeStr); unitIndex++ {
		if sizeStr[unitIndex] < '0' || sizeStr[unitIndex] > '9' {
			if sizeStr[unitIndex] != '.' {
				break
			}
		}
	}

	numStr := sizeStr[:unitIndex]
	unit := sizeStr[unitIndex:]

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %v", err)
	}

	var multiplier float64
	switch unit {
	case "b":
		multiplier = 1
	case "kb":
		multiplier = 1024
	case "mb":
		multiplier = 1024 * 1024
	case "gb":
		multiplier = 1024 * 1024 * 1024
	case "tb":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown unit of measurement: %s", unit)
	}

	bytes := num * multiplier
	return int64(bytes), nil
}

func ToBytesOrDefault(sizeStr string) int64 {
	size, err := ToBytes(sizeStr)
	if err != nil {
		return 0 // Default to 0 if conversion fails
	}
	return size
}

func LogDeletionToFile(files map[string]string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	const (
		DELETION_FILE_NAME = "deletor.log"
	)
	var deletionLogs string
	deletionTimestamp := time.Now().Format("2006-01-02 15:04:05")
	for path, size := range files {
		deletionLogs += fmt.Sprintf("[%s] %s | %s\n", deletionTimestamp, path, size)
	}

	file, err := os.OpenFile(DELETION_FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(yellow("Error:"), "Failed to open deleted files")
		return
	}
	_, err = file.WriteString(deletionLogs)
	if err != nil {
		fmt.Println(yellow("Error:"), "Failed to save deleted files")
		return
	}
	defer file.Close()
}
