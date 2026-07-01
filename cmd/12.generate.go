// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"os"

	"github.com/snowdreamtech/unigo/internal/cli/output"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	generateOutput     string
	generateManpageDir string
)

func init() {
	generateGithubActionCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "write to file instead of stdout")
	generateGitlabCiCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "write to file instead of stdout")
	generateDockerfileCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "write to file instead of stdout")
	generatePreCommitCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "write to file instead of stdout")
	generateManpageCmd.Flags().StringVarP(&generateManpageDir, "dir", "d", "", "directory to export manpages to")

	generateCmd.AddCommand(generateGithubActionCmd)
	generateCmd.AddCommand(generateGitlabCiCmd)
	generateCmd.AddCommand(generateDockerfileCmd)
	generateCmd.AddCommand(generatePreCommitCmd)
	generateCmd.AddCommand(generateShellAliasCmd)
	generateCmd.AddCommand(generateManpageCmd)
	if rootCmd != nil {
		rootCmd.AddCommand(generateCmd)
	}
}

// generateCmd is the root of the generate sub-command group.
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate integration files (GitHub Actions, pre-commit hooks, etc.)",
	Aliases: []string{"gen"},
	Long: `Generate integration files for common tooling.

Sub-commands:
  github-action   Generate a GitHub Actions workflow step
  gitlab-ci       Generate a GitLab CI script snippet
  dockerfile      Generate a Dockerfile snippet for UniGo
  pre-commit      Generate a .pre-commit-hooks.yaml snippet
  shell-alias     Print shell alias definitions
  manpage         Generate manpages

Examples:
  unigo generate github-action
  unigo generate gitlab-ci
  unigo generate dockerfile
  unigo generate pre-commit --output .pre-commit-hooks.yaml
  unigo generate shell-alias >> ~/.zshrc
  unigo generate manpage -d ./manpages`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// ─── github-action ────────────────────────────────────────────────────────────

var generateGithubActionCmd = &cobra.Command{
	Use:   "github-action",
	Short: "Generate a GitHub Actions workflow snippet for UniGo",
	Args:  cobra.NoArgs,
	RunE:  runGenerateGithubAction,
}

const githubActionTemplate = `# Add this step to your GitHub Actions workflow to install UniGo:
- name: Install UniGo
  uses: snowdreamtech/setup-unigo@v1
  # Or install manually:
  # run: curl -fsSL https://github.com/snowdreamtech/unigo/raw/main/install.sh | sh

- name: Install tools
  run: unigo install

- name: Verify lock file
  run: unigo lock --check
`

func runGenerateGithubAction(cmd *cobra.Command, args []string) error {
	return writeOrPrint(generateOutput, githubActionTemplate)
}

// ─── gitlab-ci ────────────────────────────────────────────────────────────────

var generateGitlabCiCmd = &cobra.Command{
	Use:   "gitlab-ci",
	Short: "Generate a GitLab CI script snippet for UniGo",
	Args:  cobra.NoArgs,
	RunE:  runGenerateGitlabCi,
}

const gitlabCiTemplate = `# Add this to your .gitlab-ci.yml to install UniGo:
.unigo-setup:
  before_script:
    - curl -fsSL https://github.com/snowdreamtech/unigo/raw/main/install.sh | sh
    - export PATH="$HOME/.local/share/unigo/shims:$PATH"
    - unigo install
    - unigo lock --check
`

func runGenerateGitlabCi(cmd *cobra.Command, args []string) error {
	return writeOrPrint(generateOutput, gitlabCiTemplate)
}

// ─── dockerfile ───────────────────────────────────────────────────────────────

var generateDockerfileCmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "Generate a Dockerfile snippet for UniGo",
	Args:  cobra.NoArgs,
	RunE:  runGenerateDockerfile,
}

const dockerfileTemplate = `# Add UniGo to your Dockerfile:
RUN curl -fsSL https://github.com/snowdreamtech/unigo/raw/main/install.sh | sh
ENV PATH="/root/.local/share/unigo/shims:$PATH"

# Copy configuration and install tools
COPY .unigo.toml ./
# COPY unigo.lock ./
RUN unigo install
`

func runGenerateDockerfile(cmd *cobra.Command, args []string) error {
	return writeOrPrint(generateOutput, dockerfileTemplate)
}

// ─── pre-commit ───────────────────────────────────────────────────────────────

var generatePreCommitCmd = &cobra.Command{
	Use:   "pre-commit",
	Short: "Generate a pre-commit hook snippet for UniGo",
	Args:  cobra.NoArgs,
	RunE:  runGeneratePreCommit,
}

const preCommitTemplate = `# Add to .pre-commit-config.yaml:
repos:
  - repo: local
    hooks:
      - id: unigo-lock-check
        name: UniGo lock file check
        language: system
        entry: unigo lock --check
        pass_filenames: false
        always_run: true
`

func runGeneratePreCommit(cmd *cobra.Command, args []string) error {
	return writeOrPrint(generateOutput, preCommitTemplate)
}

// ─── shell-alias ──────────────────────────────────────────────────────────────

var generateShellAliasCmd = &cobra.Command{
	Use:   "shell-alias",
	Short: "Print shell alias definitions for common UniGo commands",
	Args:  cobra.NoArgs,
	RunE:  runGenerateShellAlias,
}

const shellAliasTemplate = `# UniGo shell aliases — add to ~/.bashrc or ~/.zshrc:
alias u='unigo'
alias ui='unigo install'
alias ul='unigo list'
alias uu='unigo outdated'
alias ub='unigo backends'
`

func runGenerateShellAlias(cmd *cobra.Command, args []string) error {
	return writeOrPrint("", shellAliasTemplate) // always stdout for aliases
}

// ─── manpage ──────────────────────────────────────────────────────────────────

var generateManpageCmd = &cobra.Command{
	Use:   "manpage",
	Short: "Generate manpages for UniGo",
	Args:  cobra.NoArgs,
	RunE:  runGenerateManpage,
}

func runGenerateManpage(cmd *cobra.Command, args []string) error {
	if generateManpageDir == "" {
		return fmt.Errorf("directory is required, use -d/--dir flag")
	}

	if err := os.MkdirAll(generateManpageDir, os.ModePerm); err != nil {
		return err
	}

	header := &doc.GenManHeader{
		Title:   "UNIGO",
		Section: "1",
		Source:  "Auto generated by spf13/cobra",
	}

	if err := doc.GenManTree(rootCmd, header, generateManpageDir); err != nil {
		return err
	}

	formatter := output.NewFormatter(output.FormatterOptions{
		Format:  getOutputFormat(),
		NoColor: false,
		Writer:  os.Stdout,
		Quiet:   quiet,
	})
	formatter.Success(fmt.Sprintf("Manpages successfully generated in %s", generateManpageDir), nil)
	return nil
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func writeOrPrint(path, content string) error {
	formatter := output.NewFormatter(output.FormatterOptions{
		Format:  getOutputFormat(),
		NoColor: false,
		Writer:  os.Stdout,
		Quiet:   quiet,
	})
	if path != "" {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			formatter.Error(fmt.Sprintf("Failed to write %s: %v", path, err))
			return err
		}
		formatter.Success(fmt.Sprintf("Written to %s", path), nil)
		return nil
	}
	fmt.Print(content)
	return nil
}
