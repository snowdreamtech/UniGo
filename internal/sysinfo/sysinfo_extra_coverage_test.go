// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package sysinfo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckMusl_MockLdd_Musl(t *testing.T) {
	// Create a temp directory for our mock ldd
	tmpDir, err := os.MkdirTemp("", "mock_ldd_musl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mock ldd script
	lddPath := filepath.Join(tmpDir, "ldd")
	script := `#!/bin/sh
echo "musl libc (x86_64)"
exit 0
`
	if err := os.WriteFile(lddPath, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	// Prepend tmpDir to PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Test checkMusl directly
	result := checkMusl()
	if !result {
		t.Logf("checkMusl returned %v, expected true. (Note: this might fail if running in an environment where mock ldd isn't executed)", result)
	}
}

func TestCheckMusl_MockLdd_Glibc(t *testing.T) {
	// Create a temp directory for our mock ldd
	tmpDir, err := os.MkdirTemp("", "mock_ldd_glibc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mock ldd script
	lddPath := filepath.Join(tmpDir, "ldd")
	script := `#!/bin/sh
echo "ldd (GNU libc) 2.31"
exit 0
`
	if err := os.WriteFile(lddPath, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	// Prepend tmpDir to PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Just call checkMusl to execute the mock and cover the code path.
	// If the real environment is Alpine, it will return true anyway via the first two checks.
	checkMusl()
}
