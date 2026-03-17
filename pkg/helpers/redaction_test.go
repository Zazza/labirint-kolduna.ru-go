package helpers

import (
	"testing"
)

func TestRedactSensitiveData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "password field",
			input:    `{"username":"user","password":"secret123"}`,
			expected: `{"username":"user","password":"[REDACTED]"}`,
		},
		{
			name:     "token field",
			input:    `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
			expected: `{"token":"[REDACTED]"}`,
		},
		{
			name:     "secret field",
			input:    `{"secret":"my-secret-key"}`,
			expected: `{"secret":"[REDACTED]"}`,
		},
		{
			name:     "api_key field",
			input:    `{"api_key":"sk-1234567890"}`,
			expected: `{"api_key":"[REDACTED]"}`,
		},
		{
			name:     "authorization field",
			input:    `{"authorization":"Bearer token123"}`,
			expected: `{"authorization":"[REDACTED]"}`,
		},
		{
			name:     "jwt field",
			input:    `{"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
			expected: `{"jwt":"[REDACTED]"}`,
		},
		{
			name:     "access_token field",
			input:    `{"access_token":"token123"}`,
			expected: `{"access_token":"[REDACTED]"}`,
		},
		{
			name:     "refresh_token field",
			input:    `{"refresh_token":"refresh123"}`,
			expected: `{"refresh_token":"[REDACTED]"}`,
		},
		{
			name:     "multiple sensitive fields",
			input:    `{"password":"pass123","token":"token123","username":"user"}`,
			expected: `{"password":"[REDACTED]","token":"[REDACTED]","username":"user"}`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no sensitive data",
			input:    `{"username":"user","email":"test@example.com"}`,
			expected: `{"username":"user","email":"test@example.com"}`,
		},
		{
			name:     "malformed JSON",
			input:    `{"password":"pass123"{"invalid"}`,
			expected: `{"password":"[REDACTED]"{"invalid"}`,
		},
		{
			name:     "password with spaces",
			input:    `{"password" : "secret"}`,
			expected: `{"password":"[REDACTED]"}`,
		},
		{
			name:     "password with tabs",
			input:    `{"password"\t:\t"secret"}`,
			expected: `{"password":"[REDACTED]"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactSensitiveData(tt.input)
			if result != tt.expected {
				t.Errorf("RedactSensitiveData(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRedactPassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic password",
			input:    `{"password":"secret123"}`,
			expected: `{"password":"[REDACTED]"}`,
		},
		{
			name:     "password with special chars",
			input:    `{"password":"p@ssw0rd!"}`,
			expected: `{"password":"[REDACTED]"}`,
		},
		{
			name:     "empty password",
			input:    `{"password":""}`,
			expected: `{"password":"[REDACTED]"}`,
		},
		{
			name:     "no password field",
			input:    `{"username":"user"}`,
			expected: `{"username":"user"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactPassword(tt.input)
			if result != tt.expected {
				t.Errorf("RedactPassword(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRedactToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic token",
			input:    `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
			expected: `{"token":"[REDACTED]"}`,
		},
		{
			name:     "empty token",
			input:    `{"token":""}`,
			expected: `{"token":"[REDACTED]"}`,
		},
		{
			name:     "no token field",
			input:    `{"username":"user"}`,
			expected: `{"username":"user"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactToken(tt.input)
			if result != tt.expected {
				t.Errorf("RedactToken(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRedactSensitiveData_RealWorld(t *testing.T) {
	input := `{
		"user": {
			"username": "testuser",
			"password": "MySecurePassword123!",
			"email": "test@example.com"
		},
		"auth": {
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			"refresh_token": "another-token-here"
		}
	}`

	result := RedactSensitiveData(input)

	// Verify redaction
	if contains(result, `"password":"MySecurePassword123!"`) {
		t.Error("Password was not redacted")
	}
	if contains(result, `"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`) {
		t.Error("Token was not redacted")
	}
	if contains(result, `"refresh_token":"another-token-here"`) {
		t.Error("Refresh token was not redacted")
	}

	// Verify non-sensitive data is preserved (note: spacing may be normalized)
	if !contains(result, `"username"`) || !contains(result, `testuser`) {
		t.Error("Username was incorrectly redacted")
	}
	if !contains(result, `"email"`) || !contains(result, `test@example.com`) {
		t.Error("Email was incorrectly redacted")
	}

	// Verify redaction actually happened
	if !contains(result, `"password":"[REDACTED]"`) {
		t.Error("Password redaction marker not found")
	}
	if !contains(result, `"token":"[REDACTED]"`) {
		t.Error("Token redaction marker not found")
	}
	if !contains(result, `"refresh_token":"[REDACTED]"`) {
		t.Error("Refresh token redaction marker not found")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
