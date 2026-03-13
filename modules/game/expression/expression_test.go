package expression

import (
	"testing"
)

func TestRun_ValidExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected any
	}{
		{
			name:     "Simple addition",
			expr:     "5 + 3",
			expected: 8,
		},
		{
			name:     "Simple subtraction",
			expr:     "10 - 4",
			expected: 6,
		},
		{
			name:     "Multiplication",
			expr:     "6 * 7",
			expected: 42,
		},
		{
			name:     "Division",
			expr:     "20 / 4",
			expected: 5,
		},
		{
			name:     "Complex expression",
			expr:     "(10 + 5) * 2 - 3",
			expected: 27,
		},
		{
			name:     "Boolean true",
			expr:     "5 == 5",
			expected: true,
		},
		{
			name:     "Boolean false",
			expr:     "5 == 6",
			expected: false,
		},
		{
			name:     "Greater than",
			expr:     "10 > 5",
			expected: true,
		},
		{
			name:     "Less than",
			expr:     "5 < 10",
			expected: true,
		},
		{
			name:     "Greater or equal",
			expr:     "10 >= 10",
			expected: true,
		},
		{
			name:     "Less or equal",
			expr:     "5 <= 5",
			expected: true,
		},
		{
			name:     "Not equal",
			expr:     "5 != 6",
			expected: true,
		},
		{
			name:     "And operator",
			expr:     "true && true",
			expected: true,
		},
		{
			name:     "Or operator",
			expr:     "true || false",
			expected: true,
		},
		{
			name:     "Modulo",
			expr:     "10 % 3",
			expected: 1,
		},
		{
			name:     "Power",
			expr:     "2 ** 3",
			expected: 8,
		},
		{
			name:     "Negative number",
			expr:     "-5",
			expected: -5,
		},
		{
			name:     "Float division",
			expr:     "7 / 2",
			expected: 3.5,
		},
		{
			name:     "Float result",
			expr:     "3.5 * 2",
			expected: 7.0,
		},
		{
			name:     "Complex boolean",
			expr:     "(5 > 3) && (10 < 20)",
			expected: true,
		},
		{
			name:     "Complex boolean",
			expr:     "(!true && !false) && (!false && ! false)",
			expected: false,
		},
		{
			name:     "Parentheses priority",
			expr:     "2 + 3 * 4",
			expected: 14,
		},
		{
			name:     "Parentheses override",
			expr:     "(2 + 3) * 4",
			expected: 20,
		},
		{
			name:     "Zero",
			expr:     "0",
			expected: 0,
		},
		{
			name:     "Large number",
			expr:     "1000000",
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Run(tt.expr)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Compare values handling different types
			switch expected := tt.expected.(type) {
			case int:
				if resultInt, ok := result.(int); ok {
					if resultInt != expected {
						t.Errorf("Expected %d, got %d", expected, resultInt)
					}
				} else if resultFloat, ok := result.(float64); ok {
					if int(resultFloat) != expected {
						t.Errorf("Expected %d, got %f", expected, resultFloat)
					}
				} else {
					t.Errorf("Expected int, got %T", result)
				}
			case float64:
				if resultFloat, ok := result.(float64); ok {
					if resultFloat != expected {
						t.Errorf("Expected %f, got %f", expected, resultFloat)
					}
				} else {
					t.Errorf("Expected float64, got %T", result)
				}
			case bool:
				if resultBool, ok := result.(bool); ok {
					if resultBool != expected {
						t.Errorf("Expected %t, got %t", expected, resultBool)
					}
				} else {
					t.Errorf("Expected bool, got %T", result)
				}
			}
		})
	}
}

func TestRun_InvalidExpressions(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		expectError bool
	}{
		{
			name:        "Unmatched parenthesis",
			expr:        "(5 + 3",
			expectError: true,
		},
		{
			name:        "Invalid operator",
			expr:        "5 @ 3",
			expectError: true,
		},
		{
			name:        "Empty expression",
			expr:        "",
			expectError: true,
		},
		{
			name:        "Only operators",
			expr:        "+++",
			expectError: true,
		},
		{
			name:        "Invalid syntax",
			expr:        "5 5 5",
			expectError: true,
		},
		{
			name:        "Undefined variable",
			expr:        "x + 5",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Run(tt.expr)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestRunAndReturnRoundUint_ValidExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected uint
	}{
		{
			name:     "Uint result",
			expr:     "10 + 5",
			expected: 15,
		},
		{
			name:     "Float64 round down",
			expr:     "7.3",
			expected: 7,
		},
		{
			name:     "Float64 round up",
			expr:     "7.8",
			expected: 8,
		},
		{
			name:     "Float64 exact half",
			expr:     "7.5",
			expected: 8,
		},
		{
			name:     "Complex expression with float",
			expr:     "(10.5 + 5.3) * 2",
			expected: 32,
		},
		{
			name:     "Zero",
			expr:     "0",
			expected: 0,
		},
		{
			name:     "Large number",
			expr:     "999999.9",
			expected: 1000000,
		},
		{
			name:     "Small decimal",
			expr:     "0.1",
			expected: 0,
		},
		{
			name:     "0.5 rounds to 1",
			expr:     "0.5",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunAndReturnRoundUint(tt.expr)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestRunAndReturnRoundUint_InvalidExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			name: "Invalid syntax",
			expr: "5 +",
		},
		{
			name: "Empty",
			expr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RunAndReturnRoundUint(tt.expr)

			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestRunAndReturnBoolean_ValidExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{
			name:     "True equality",
			expr:     "5 == 5",
			expected: true,
		},
		{
			name:     "False equality",
			expr:     "5 == 6",
			expected: false,
		},
		{
			name:     "True greater than",
			expr:     "10 > 5",
			expected: true,
		},
		{
			name:     "False greater than",
			expr:     "5 > 10",
			expected: false,
		},
		{
			name:     "True less than",
			expr:     "3 < 7",
			expected: true,
		},
		{
			name:     "False less than",
			expr:     "7 < 3",
			expected: false,
		},
		{
			name:     "True greater or equal",
			expr:     "5 >= 5",
			expected: true,
		},
		{
			name:     "True greater or equal (strict)",
			expr:     "6 >= 5",
			expected: true,
		},
		{
			name:     "False greater or equal",
			expr:     "4 >= 5",
			expected: false,
		},
		{
			name:     "True less or equal",
			expr:     "5 <= 5",
			expected: true,
		},
		{
			name:     "True less or equal (strict)",
			expr:     "4 <= 5",
			expected: true,
		},
		{
			name:     "False less or equal",
			expr:     "6 <= 5",
			expected: false,
		},
		{
			name:     "True not equal",
			expr:     "5 != 6",
			expected: true,
		},
		{
			name:     "False not equal",
			expr:     "5 != 5",
			expected: false,
		},
		{
			name:     "True AND true",
			expr:     "true && true",
			expected: true,
		},
		{
			name:     "True AND false",
			expr:     "true && false",
			expected: false,
		},
		{
			name:     "False AND true",
			expr:     "false && true",
			expected: false,
		},
		{
			name:     "False AND false",
			expr:     "false && false",
			expected: false,
		},
		{
			name:     "True OR true",
			expr:     "true || true",
			expected: true,
		},
		{
			name:     "True OR false",
			expr:     "true || false",
			expected: true,
		},
		{
			name:     "False OR true",
			expr:     "false || true",
			expected: true,
		},
		{
			name:     "False OR false",
			expr:     "false || false",
			expected: false,
		},
		{
			name:     "Complex boolean AND",
			expr:     "(5 > 3) && (10 < 20)",
			expected: true,
		},
		{
			name:     "Complex boolean OR",
			expr:     "(5 > 10) || (10 < 20)",
			expected: true,
		},
		{
			name:     "Complex boolean mixed",
			expr:     "(5 > 3) && (10 > 20) || (3 < 5)",
			expected: true,
		},
		{
			name:     "Negation",
			expr:     "!true",
			expected: false,
		},
		{
			name:     "Double negation",
			expr:     "!!true",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunAndReturnBoolean(tt.expr)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestRunAndReturnBoolean_InvalidExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			name: "Non-boolean result",
			expr: "5 + 3",
		},
		{
			name: "String result",
			expr: "'hello'",
		},
		{
			name: "Invalid syntax",
			expr: "5 >",
		},
		{
			name: "Empty",
			expr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RunAndReturnBoolean(tt.expr)

			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestRunAndReturnBoolean_GameExpressions(t *testing.T) {
	// Test typical game expressions
	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{
			name:     "Dice check exact",
			expr:     "5 == 5",
			expected: true,
		},
		{
			name:     "Dice check greater or equal",
			expr:     "5 >= 5",
			expected: true,
		},
		{
			name:     "Dice check fail",
			expr:     "4 >= 5",
			expected: false,
		},
		{
			name:     "Two dice sum check",
			expr:     "9 >= 8",
			expected: true,
		},
		{
			name:     "Two dice sum fail",
			expr:     "6 >= 8",
			expected: false,
		},
		{
			name:     "Health check",
			expr:     "10 > 5",
			expected: true,
		},
		{
			name:     "Gold check",
			expr:     "100 >= 50",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunAndReturnBoolean(tt.expr)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestRunAndReturnRoundUint_GameExpressions(t *testing.T) {
	// Test typical game expressions
	tests := []struct {
		name     string
		expr     string
		expected uint
	}{
		{
			name:     "Health increase",
			expr:     "10 + 5",
			expected: 15,
		},
		{
			name:     "Health decrease",
			expr:     "10 - 5",
			expected: 5,
		},
		{
			name:     "Health multiplication",
			expr:     "10 * 2",
			expected: 20,
		},
		{
			name:     "Gold calculation",
			expr:     "100 + 50",
			expected: 150,
		},
		{
			name:     "Weapon count increase",
			expr:     "3 + 1",
			expected: 4,
		},
		{
			name:     "Weapon count decrease",
			expr:     "3 - 1",
			expected: 2,
		},
		{
			name:     "Float dice damage",
			expr:     "5.5",
			expected: 6,
		},
		{
			name:     "Complex health calculation",
			expr:     "(10 + 5) * 2 - 3",
			expected: 27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunAndReturnRoundUint(tt.expr)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
