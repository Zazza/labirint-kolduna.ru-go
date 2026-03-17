package helpers

import (
	"testing"
)

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "script tag",
			input:    "<script>alert('XSS')</script>",
			expected: "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
		},
		{
			name:     "img tag onerror",
			input:    "<img src=x onerror=alert('XSS')>",
			expected: "&lt;img src=x onerror=alert(&#39;XSS&#39;)&gt;",
		},
		{
			name:     "svg xss",
			input:    "<svg onload=alert('XSS')>",
			expected: "&lt;svg onload=alert(&#39;XSS&#39;)&gt;",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "unicode emojis",
			input:    "🏆 Успешно",
			expected: "🏆 Успешно",
		},
		{
			name:     "html entity",
			input:    "&lt;",
			expected: "&amp;lt;",
		},
		{
			name:     "mixed content",
			input:    "Hello <b>World</b> & <script>alert('XSS')</script>",
			expected: "Hello &lt;b&gt;World&lt;/b&gt; &amp; &lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHTML(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapeHTMLArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "simple array",
			input:    []string{"Hello", "World"},
			expected: []string{"Hello", "World"},
		},
		{
			name:     "xss array",
			input:    []string{"<script>alert('XSS')</script>", "<img src=x onerror=alert('XSS')>"},
			expected: []string{"&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;", "&lt;img src=x onerror=alert(&#39;XSS&#39;)&gt;"},
		},
		{
			name:     "empty array",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "nil array",
			input:    nil,
			expected: []string(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHTMLArray(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("EscapeHTMLArray() length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("EscapeHTMLArray()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestEscapeHTML_Identity(t *testing.T) {
	input := "Hello World"
	escaped := EscapeHTML(input)
	reescaped := EscapeHTML(escaped)

	if reescaped != escaped {
		t.Errorf("Double escape changed output: %q -> %q -> %q", input, escaped, reescaped)
	}
}
