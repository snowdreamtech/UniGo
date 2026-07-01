// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	if rootCmd != nil {
		rootCmd.AddCommand(generateCmd)
	}
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate integration files",
	Long:  "Generate configuration files for CI/CD, pre-commit hooks, or IDE integrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("Running generate command...")
		
		content := `# UniGo GitHub Actions Workflow
name: Go

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:

  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.24'

    - name: Build
      run: go build -v ./...

    - name: Test
      run: go test -v ./...
`
		
		dir := filepath.Join(".github", "workflows")
		fileName := filepath.Join(dir, "go.yml")
		
		if _, err := os.Stat(fileName); err == nil {
			pterm.Warning.Printf("File %s already exists. Skipping.\n", fileName)
			return nil
		}
		
		if err := os.MkdirAll(dir, 0755); err != nil {
			pterm.Error.Printf("Failed to create directory %s: %v\n", dir, err)
			return err
		}
		
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			pterm.Error.Printf("Failed to generate %s: %v\n", fileName, err)
			return err
		}
		
		pterm.Success.Printf("Successfully generated %s\n", fileName)
		return nil
	},
}
