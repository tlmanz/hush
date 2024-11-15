package main

import (
	"context"
	"fmt"

	"github.com/tlmanz/hush"
)

type User struct {
	Name     string
	Password string   `hush:"hide"`
	Ages     []int    `hush:"mask"`
	Emails   []string `hush:"mask"`
}

func main() {
	user := User{
		Name:     "John",
		Password: "secret123",
		Ages:     []int{301, 312, 323},
		Emails:   []string{"john@example.com", "john2@example.com"},
	}

	husher := hush.NewHush()

	result, err := husher.Hush(context.Background(), 10, "TESTFIELD", hush.TagHide)
	if err != nil {
		panic(err)
	}

	fmt.Println("String Usage Example (With Prefix):")
	for _, field := range result {
		fmt.Printf("%s: %s\n", field[0], field[1])
	}

	result, err = husher.Hush(context.Background(), 10, hush.TagHide)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nString Usage Example:")
	for _, field := range result {
		fmt.Printf("%s\n", field[0])
	}

	result, err = husher.Hush(context.Background(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nStruct Usage Example:")
	for _, field := range result {
		fmt.Printf("%s: %s\n", field[0], field[1])
	}
}
