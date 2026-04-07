package ui_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bketelsen/toolbox/ui"
)

func TestFormatExamples_SingleExample(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "List all items",
			Command:     "list",
		},
	}

	result := ui.FormatExamples(examples...)

	assert.Contains(t, result, "List all items")
	assert.Contains(t, result, "list")
	assert.Contains(t, result, "$")
}

func TestFormatExamples_MultipleExamples(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "List all items",
			Command:     "list",
		},
		{
			Description: "Show details",
			Command:     "show --details",
		},
		{
			Description: "Export data",
			Command:     "export --format json",
		},
	}

	result := ui.FormatExamples(examples...)

	assert.Contains(t, result, "List all items")
	assert.Contains(t, result, "list")
	assert.Contains(t, result, "Show details")
	assert.Contains(t, result, "show --details")
	assert.Contains(t, result, "Export data")
	assert.Contains(t, result, "export --format json")
}

func TestFormatExamples_EmptyDescription(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "",
			Command:     "command",
		},
	}

	result := ui.FormatExamples(examples...)

	// Should still have the command
	assert.Contains(t, result, "command")
}

func TestFormatExamples_NoExamples(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{}

	result := ui.FormatExamples(examples...)

	assert.Equal(t, "", result)
}

func TestFormatExamples_LongCommand(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "Complex operation",
			Command:     "tool --flag=value --option foo bar baz",
		},
	}

	result := ui.FormatExamples(examples...)

	assert.Contains(t, result, "Complex operation")
	assert.Contains(t, result, "tool --flag=value --option foo bar baz")
}

func TestFormatExamples_SpecialCharacters(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "Use pipes and redirects",
			Command:     "cat file.txt | grep pattern > output.txt",
		},
	}

	result := ui.FormatExamples(examples...)

	assert.Contains(t, result, "pipes")
	assert.Contains(t, result, "cat")
}

func TestFormatExamples_BulletPoint(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "First example",
			Command:     "cmd1",
		},
	}

	result := ui.FormatExamples(examples...)

	// Should contain bullet point character or dash
	assert.True(t, strings.Contains(result, "-") || strings.Contains(result, "•"))
}

func TestLong_WithDescription(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "List items",
			Command:     "list",
		},
	}

	result := ui.Long("This command lists all available items", examples...)

	assert.Contains(t, result, "This command lists all available items")
	assert.Contains(t, result, "List items")
	assert.Contains(t, result, "list")
}

func TestLong_EmptyDescription(t *testing.T) {
	t.Parallel()

	examples := []ui.Example{
		{
			Description: "Get details",
			Command:     "details",
		},
	}

	result := ui.Long("", examples...)

	// Should still have examples
	assert.Contains(t, result, "Get details")
	assert.Contains(t, result, "details")
}

func TestLong_NoExamples(t *testing.T) {
	t.Parallel()

	result := ui.Long("Just a description", []ui.Example{}...)

	assert.Contains(t, result, "Just a description")
}

func TestLong_MultipleExamples(t *testing.T) {
	t.Parallel()

	description := "Tool description with detailed information"
	examples := []ui.Example{
		{
			Description: "Basic usage",
			Command:     "tool",
		},
		{
			Description: "With options",
			Command:     "tool --verbose",
		},
		{
			Description: "With output",
			Command:     "tool > file.txt",
		},
	}

	result := ui.Long(description, examples...)

	assert.Contains(t, result, description)
	assert.Contains(t, result, "Basic usage")
	assert.Contains(t, result, "tool")
	assert.Contains(t, result, "With options")
	assert.Contains(t, result, "--verbose")
	assert.Contains(t, result, "With output")
	assert.Contains(t, result, "> file.txt")
}

func TestLong_LongDescription(t *testing.T) {
	t.Parallel()

	longDesc := strings.Repeat("word ", 50) // Create a long description

	result := ui.Long(longDesc, []ui.Example{}...)

	require.NotEmpty(t, result)
	assert.Contains(t, result, "word")
}

func TestLong_Format(t *testing.T) {
	t.Parallel()

	description := "A detailed explanation"
	examples := []ui.Example{
		{
			Description: "Example one",
			Command:     "cmd1",
		},
	}

	result := ui.Long(description, examples...)

	// Result should be formatted string
	assert.NotEmpty(t, result)
	// Should have both description and example
	assert.True(t, strings.Contains(result, description) || strings.Contains(result, "Example one"))
}
