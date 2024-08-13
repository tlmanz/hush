package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tlmanz/hush"
)

type TableConfig struct {
	APIKey    string `hush:"mask"`
	SecretKey string `hush:"hide"`
	Debug     bool
}

func customMask(s string) string {
	return strings.Repeat("-", len(s))
}

func main() {
	config := TableConfig{
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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.AppendBulk(result)
	table.Render()
}
