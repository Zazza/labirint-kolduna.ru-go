package helpers

import "html/template"

func EscapeHTML(input string) string {
	return template.HTMLEscapeString(input)
}

func EscapeHTMLArray(inputs []string) []string {
	escaped := make([]string, len(inputs))
	for i, input := range inputs {
		escaped[i] = EscapeHTML(input)
	}
	return escaped
}
