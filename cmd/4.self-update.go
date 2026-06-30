// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update to the latest version",
	Long:  "Update the application to the latest available version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("self-update is not implemented in this template.")
		return nil
	},
}
