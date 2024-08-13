package main

import (
	"context"
	"fmt"

	"github.com/tlmanz/hush"
)

type Address struct {
	Street  string
	City    string
	Country string
}

type Account struct {
	ID       int
	Balance  float64 `hush:"mask"`
	Currency string
}

type ComplexUser struct {
	Name     string `hush:"mask"`
	Email    string `hush:"hide"`
	Age      int
	Address  Address `hush:"hide"`
	Accounts []Account
	Metadata map[string]string
	IsActive bool
}

func main() {
	user := ComplexUser{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Age:   28,
		Address: Address{
			Street:  "123 Main St",
			City:    "Anytown",
			Country: "USA",
		},
		Accounts: []Account{
			{ID: 1, Balance: 1000.50, Currency: "USD"},
			{ID: 2, Balance: 500.75, Currency: "EUR"},
		},
		Metadata: map[string]string{
			"lastLogin": "2023-04-01",
			"role":      "admin",
		},
		IsActive: true,
	}

	husher := hush.NewHush()

	result, err := husher.Hush(context.Background(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println("Complex Struct Example:")
	for _, field := range result {
		fmt.Printf("%s: %s\n", field[0], field[1])
	}
}
