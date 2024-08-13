package main

import (
	"context"
	"fmt"

	"github.com/tlmanz/hush"
)

type Person struct {
	Name    string `hush:"mask"`
	age     int    `hush:"mask"`
	Address string
}

func main() {
	person := Person{
		Name:    "John Doe",
		age:     30,
		Address: "123 Main St",
	}

	husher := hush.NewHush()

	// Without private fields (default behavior)
	resultWithoutPrivate, err := husher.Hush(context.Background(), person)
	if err != nil {
		panic(err)
	}

	fmt.Println("Without private fields:")
	for _, field := range resultWithoutPrivate {
		fmt.Printf("%s: %s\n", field[0], field[1])
	}

	fmt.Println()

	// With private fields
	resultWithPrivate, err := husher.Hush(context.Background(), person,
		hush.WithPrivateFields(true),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("With private fields:")
	for _, field := range resultWithPrivate {
		fmt.Printf("%s: %s\n", field[0], field[1])
	}
}
