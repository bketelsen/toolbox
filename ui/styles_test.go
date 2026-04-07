package ui_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bketelsen/toolbox/ui"
)

func TestBold(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple text",
			input:    "Hello",
			expected: "Hello",
		},
		{
			name:     "text with spaces",
			input:    "Hello World",
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.Bold(tt.input)
			// Bold may add formatting codes when in a TTY
			// so we just verify it contains the text
			assert.True(t, strings.Contains(result, tt.input) || result == tt.expected)
		})
	}
}

func TestRed(t *testing.T) {
	t.Parallel()

	input := "Error text"
	result := ui.Red(input)

	// Red may add formatting codes, so we check it contains the input
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestYellow(t *testing.T) {
	t.Parallel()

	input := "Warning text"
	result := ui.Yellow(input)

	// Yellow may add formatting codes
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestGreen(t *testing.T) {
	t.Parallel()

	input := "Success text"
	result := ui.Green(input)

	// Green may add formatting codes
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestTimestamp(t *testing.T) {
	t.Parallel()

	now := time.Now()
	result := ui.Timestamp(now)

	// Should contain the formatted time
	assert.NotEmpty(t, result)
}

func TestKeyword(t *testing.T) {
	t.Parallel()

	input := "keyword"
	result := ui.Keyword(input)

	// Should contain the keyword
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestPlaceholder(t *testing.T) {
	t.Parallel()

	input := "<placeholder>"
	result := ui.Placeholder(input)

	// Should contain the placeholder
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestWrap(t *testing.T) {
	t.Parallel()

	input := "This is a long text that might need to be wrapped"
	result := ui.Wrap(input)

	// Should contain the input text
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestCode(t *testing.T) {
	t.Parallel()

	input := "fmt.Println(\"hello\")"
	result := ui.Code(input)

	// Should contain the code
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestField(t *testing.T) {
	t.Parallel()

	input := "field_name"
	result := ui.Field(input)

	// Should contain the field name
	assert.True(t, strings.Contains(result, input) || result == input)
}

func TestKeyValuePair(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		value    string
		contains string
	}{
		{
			name:     "simple pair",
			key:      "name",
			value:    "John",
			contains: "name",
		},
		{
			name:     "with colon",
			key:      "id",
			value:    "123",
			contains: ":",
		},
		{
			name:     "empty value",
			key:      "status",
			value:    "",
			contains: "status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.KeyValuePair(tt.key, tt.value)
			assert.Contains(t, result, tt.contains)
		})
	}
}

func TestDefaultStylesInitialized(t *testing.T) {
	t.Parallel()

	// DefaultStyles should be initialized
	styles := ui.DefaultStyles
	assert.NotNil(t, styles)

	// Verify all fields are set
	assert.NotNil(t, styles.Code)
	assert.NotNil(t, styles.DateTimeStamp)
	assert.NotNil(t, styles.Error)
	assert.NotNil(t, styles.Field)
	assert.NotNil(t, styles.Hyperlink)
	assert.NotNil(t, styles.Keyword)
	assert.NotNil(t, styles.Placeholder)
	assert.NotNil(t, styles.Prompt)
	assert.NotNil(t, styles.FocusedPrompt)
	assert.NotNil(t, styles.Fuchsia)
	assert.NotNil(t, styles.Warn)
	assert.NotNil(t, styles.Success)
	assert.NotNil(t, styles.Info)
	assert.NotNil(t, styles.Debug)
	assert.NotNil(t, styles.Wrap)
}
