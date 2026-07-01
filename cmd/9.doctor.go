// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cmd

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/pterm/pterm"
	"github.com/snowdreamtech/unigo/internal/pkg/env"
	"github.com/spf13/cobra"
)

func init() {
	if rootCmd != nil {
		rootCmd.AddCommand(doctorCmd)
	}
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health and diagnose issues",
	Long:  "Run a series of checks to ensure the application environment is healthy.",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("Running doctor checks...")
		
		pterm.DefaultSection.Println("UniGo System Health Check")
		
		var checksPassed int
		var totalChecks int
		
		// 1. Check OS/Arch
		totalChecks++
		pterm.Success.Printf("System: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		checksPassed++
		
		// 2. Check Cache Dir
		totalChecks++
		cacheDir := env.GetCacheDir()
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			pterm.Info.Printf("Cache dir: %s (Not created yet)\n", cacheDir)
			checksPassed++ // Not an error if it simply hasn't been created
		} else if err != nil {
			pterm.Error.Printf("Cache dir: %s (%v)\n", cacheDir, err)
		} else {
			pterm.Success.Printf("Cache dir: %s (OK)\n", cacheDir)
			checksPassed++
		}
		
		// 3. Check CWD Permissions
		totalChecks++
		if wd, err := os.Getwd(); err == nil {
			if f, err := os.CreateTemp(wd, ".unigo-doctor-*"); err == nil {
				f.Close()
				os.Remove(f.Name())
				pterm.Success.Printf("Current dir: %s (Writable)\n", wd)
				checksPassed++
			} else {
				pterm.Warning.Printf("Current dir: %s (Read-only or restricted: %v)\n", wd, err)
				checksPassed++ // Warning, but doesn't completely fail health check
			}
		} else {
			pterm.Error.Printf("Current dir: Cannot determine (%v)\n", err)
		}
		
		pterm.Println()
		if checksPassed == totalChecks {
			pterm.Success.Println("Status: All checks passed! 🚀")
		} else {
			pterm.Warning.Printf("Status: %d/%d checks passed.\n", checksPassed, totalChecks)
		}
		
		return nil
	},
}
