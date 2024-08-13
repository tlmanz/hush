package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/tlmanz/hush"
)

type PersonConfig struct {
	Email     string `hush:"mask"`
	SecretKey string
	Debug     bool
}

// MaskString masks the matched portions of the input string based on a pattern.
func maskEmail(input string) string {
	pattern := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b` // Regular expression for matching email addresses
	re := regexp.MustCompile(pattern)
	mask := "*"
	return re.ReplaceAllStringFunc(input, func(match string) string {
		if len(match) <= 2 {
			return match
		}
		return match[:2] + strings.Repeat(mask, len(match)-2)
	})
}

func main() {
	config := PersonConfig{
		Email:     "This is a test for masking tlmannapperuma@gmail.com email address",
		SecretKey: "verysecret",
		Debug:     true,
	}

	husher := hush.NewHush()

	result, err := husher.Hush(context.Background(), config,
		hush.WithSeparator("."),
		hush.WithMaskFunc(maskEmail),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Custom Options Example:")
	for _, field := range result {
		fmt.Printf("%s: %s\n", field[0], field[1])
	}
}
