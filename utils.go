package hush

import (
	"strings"
)

func buildFieldName(prefix, fieldName, separator string) string {
	if prefix == "" {
		return fieldName
	}
	return prefix + separator + fieldName
}

func defaultMaskFunc(value string) string {
	runes := []rune(value)
	length := len(runes)
	if length <= 8 {
		return strings.Repeat("*", length)
	}
	return string(runes[:1]) + strings.Repeat("*", length-2) + string(runes[length-1:])
}
