// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDirectorySize(t *testing.T) {
	tmpDir := t.TempDir()

	size, err := CalculateDirectorySize(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), size)

	// Create a dummy file
	f1 := filepath.Join(tmpDir, "dummy.txt")
	err = os.WriteFile(f1, []byte("hello"), 0644)
	assert.NoError(t, err)

	size, err = CalculateDirectorySize(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), size)

	// Subdir
	subDir := filepath.Join(tmpDir, "sub")
	err = os.Mkdir(subDir, 0755)
	assert.NoError(t, err)

	f2 := filepath.Join(subDir, "dummy2.txt")
	err = os.WriteFile(f2, []byte("world"), 0644)
	assert.NoError(t, err)

	size, err = CalculateDirectorySize(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), size)
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{"Zero", 0, "0 B"},
		{"Bytes", 500, "500 B"},
		{"Kilobytes", 1024, "1.0 KB"},
		{"Megabytes", 1024 * 1024, "1.0 MB"},
		{"Megabytes Decimals", 1500 * 1024, "1.5 MB"},
		{"Gigabytes", 1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.size)
			assert.Equal(t, tt.expected, result)
		})
	}
}
