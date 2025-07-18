package agent

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-def/internal/model"
)

// Claude implements the Claude-specific conversion logic
type Claude struct{}

// ID returns the unique identifier for Claude agent
func (c *Claude) ID() string {
	return "claude"
}

// Name returns the display name for Claude agent
func (c *Claude) Name() string {
	return "Claude"
}

// FormatFile converts a path to Claude's file reference format
func (c *Claude) FormatFile(path string) string {
	return "@" + path
}

// FormatMCP formats an MCP command for Claude agent
func (c *Claude) FormatMCP(agent, command string, args ...string) string {
	return formatMCP(agent, command, args...)
}

// FormatMemory processes a memory context for Claude agent
func (c *Claude) FormatMemory(content string) (string, error) {
	// Process the memory context for Claude
	// For now, we're simply returning the content as is
	// Any Claude-specific formatting rules could be applied here
	return content, nil
}

// FormatCommand processes command definitions for Claude agent
func (c *Claude) FormatCommand(commands []model.Command) (string, error) {
	// Format commands according to Claude's requirements
	var formattedCommands []string

	for _, cmd := range commands {
		var formattedCmd string

		// Include frontmatter if Claude-specific attributes exist
		hasFrontmatter := cmd.Claude.Description != "" || cmd.Claude.AllowedTools != ""

		if hasFrontmatter {
			formattedCmd = "---\n"

			if cmd.Claude.Description != "" {
				formattedCmd += fmt.Sprintf("description: %s\n", cmd.Claude.Description)
			}

			if cmd.Claude.AllowedTools != "" {
				formattedCmd += fmt.Sprintf("allowed-tools: %s\n", cmd.Claude.AllowedTools)
			}

			formattedCmd += "---\n\n"
		}

		// Add the main content
		formattedCmd += cmd.Content

		formattedCommands = append(formattedCommands, formattedCmd)
	}

	return strings.Join(formattedCommands, "\n\n---\n\n"), nil
}

// DefaultMemoryPath determines the output path for Claude agent memory files
func (c *Claude) DefaultMemoryPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	if userScope {
		return filepath.Join(outputBaseDir, ".claude", "CLAUDE.md"), nil
	}
	return filepath.Join(outputBaseDir, "CLAUDE.md"), nil
}

// DefaultCommandPath determines the output path for Claude agent command files
func (c *Claude) DefaultCommandPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	ext := ".md"
	name := strings.TrimSuffix(fileName, ext)

	return filepath.Join(outputBaseDir, ".claude", "commands", name+ext), nil
}

// ShouldConcatenate determines if content should be concatenated for Claude agent
func (c *Claude) ShouldConcatenate(taskType string) bool {
	return taskType == "memory"
}
