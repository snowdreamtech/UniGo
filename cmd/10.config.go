// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/snowdreamtech/unigo/internal/config"
	"github.com/snowdreamtech/unigo/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	if rootCmd != nil {
		rootCmd.AddCommand(configCmd)
		configCmd.AddCommand(configGetCmd)
		configCmd.AddCommand(configSetCmd)
		configCmd.AddCommand(configListCmd)
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Query or modify the UniGo configuration settings stored in unigo.toml.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return configListCmd.RunE(cmd, args)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		tableData := [][]string{
			{"KEY", "VALUE"},
			{"debug", strconv.FormatBool(cfg.Debug)},
		}

		return pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.ToLower(args[0])
		logger.Debug("Getting config value", "key", key)

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		switch key {
		case "debug":
			fmt.Println(cfg.Debug)
		default:
			pterm.Warning.Printf("Key '%s' not found in configuration\n", key)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.ToLower(args[0])
		val := args[1]
		logger.Debug("Setting config value", "key", key, "value", val)

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		switch key {
		case "debug":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid boolean value for debug: %w", err)
			}
			cfg.Debug = b
		default:
			return fmt.Errorf("unsupported configuration key '%s'", key)
		}

		if err := cfg.Save(); err != nil {
			pterm.Error.Printf("Failed to save configuration: %v\n", err)
			return err
		}

		pterm.Success.Printf("Set '%s' to '%s'\n", key, val)
		return nil
	},
}
