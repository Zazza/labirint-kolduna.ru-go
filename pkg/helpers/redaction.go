package helpers

import "regexp"

var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`"password"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"secret"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"api_key"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"authorization"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"jwt"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"access_token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"refresh_token"\s*:\s*"[^"]*"`),
}

var sensitiveFieldNames = map[string]bool{
	"password":      true,
	"token":         true,
	"secret":        true,
	"api_key":       true,
	"authorization": true,
	"jwt":           true,
	"access_token":  true,
	"refresh_token": true,
}

func RedactSensitiveData(input string) string {
	for fieldName := range sensitiveFieldNames {
		pattern := regexp.MustCompile(`"` + regexp.QuoteMeta(fieldName) + `"\s*:\s*"[^"]*"`)
		input = pattern.ReplaceAllString(input, `"`+fieldName+`":"[REDACTED]"`)
	}
	return input
}

func RedactPassword(input string) string {
	passwordPattern := regexp.MustCompile(`"password"\s*:\s*"[^"]*"`)
	return passwordPattern.ReplaceAllString(input, `"password":"[REDACTED]"`)
}

func RedactToken(input string) string {
	tokenPattern := regexp.MustCompile(`"token"\s*:\s*"[^"]*"`)
	return tokenPattern.ReplaceAllString(input, `"token":"[REDACTED]"`)
}
