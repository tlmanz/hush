package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/tlmanz/hush"
)

type Config struct {
	APIKey    string `hush:"mask"`
	SecretKey string `hush:"hide"`
	Debug     bool
}

func customMask(s string) string {
	return strings.Repeat("-", len(s))
}

func main() {
	config := Config{
		APIKey:    "abcdef123456",
		SecretKey: "verysecret",
		Debug:     true,
	}

	husher := hush.NewHush()

	result, err := husher.Hush(context.Background(), config,
		hush.WithSeparator("_"),
		hush.WithMaskFunc(customMask),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Custom Options Example:")
	for _, field := range result {
		fmt.Printf("%s: %s\n", field[0], field[1])
	}
}
