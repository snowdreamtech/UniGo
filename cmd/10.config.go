// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

func init() {
	if rootCmd != nil {
		rootCmd.AddCommand(configCmd)
		configCmd.AddCommand(configGetCmd)
		configCmd.AddCommand(configSetCmd)
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Query or modify the configuration settings.",
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		slog.Debug("Getting config value", "key", key)
		fmt.Printf("Placeholder: Getting value for '%s'\n", key)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		val := args[1]
		slog.Debug("Setting config value", "key", key, "value", val)
		fmt.Printf("Placeholder: Setting '%s' to '%s'\n", key, val)
		return nil
	},
}
