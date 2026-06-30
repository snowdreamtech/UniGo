// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(licenseCmd)
}

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Manage copyright license headers",
	Long:  "Manage copyright license headers in source files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("license command is a stub in this template.")
		return nil
	},
}
