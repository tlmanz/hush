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
	length := len(value)
	if length <= 4 {
		return strings.Repeat("*", length)
	}
	return value[:2] + strings.Repeat("*", length-4) + value[length-2:]
}
