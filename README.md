[![CI](https://github.com/tlmanz/hush/actions/workflows/ci.yml/badge.svg)](https://github.com/tlmanz/hush/actions/workflows/ci.yml)
[![CodeQL](https://github.com/tlmanz/hush/actions/workflows/codequality.yml/badge.svg)](https://github.com/tlmanz/hush/actions/workflows/codequality.yml)
[![Coverage Status](https://coveralls.io/repos/github/tlmanz/hush/badge.svg)](https://coveralls.io/github/tlmanz/hush)
![Open Issues](https://img.shields.io/github/issues/tlmanz/hush)
[![Go Report Card](https://goreportcard.com/badge/github.com/tlmanz/hush)](https://goreportcard.com/report/github.com/tlmanz/hush)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tlmanz/hush)

# Hush ![hush_logo](https://mannapperuma.com/files/hush.png)

Hush is a Go package that provides a flexible and efficient way to process and mask sensitive data in structs and strings.

## Features

- Process structs and strings etc to mask or hide sensitive information
- Customizable field separators for nested structures
- Support for custom masking functions
- Concurrent processing of struct fields for improved performance
- Option to include or exclude private fields
- Context-aware processing with cancellation support
- Consistent handling of maps and slices

## Installation

To install Hush, use `go get`:

```
go get github.com/tlmanz/hush
```

## Usage

Here's a basic example of how to use Hush:

```go
package main

import (
	"context"
	"fmt"

	"github.com/tlmanz/hush"
)

type User struct {
	Name     string
	Password string `hush:"hide"`
	Age      int    `hush:"mask"`
	Email    string `hush:"mask"`
}

func main() {
	user := User{
		Name:     "John",
		Password: "secret123",
		Age:      30,
		Email:    "john@example.com",
	}

	husher := hush.NewHush()

	result, err := husher.Hush(context.Background(), 10, "TESTFIELD", hush.TagHide)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nString Usage Example (With Prefix):")
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
```

This will output:
```
String Usage Example (With Prefix):
TESTFIELD: HIDDEN

String Usage Example:
HIDDEN

Struct Usage Example:
Age: **
Email: jo************om
Name: John
Password: HIDDEN
```

## Configuration

Hush provides several options to customize its behavior:

- `WithSeparator(sep string)`: Set a custom separator for nested field names (default is ".")
- `WithMaskFunc(f func(string) string)`: Set a custom masking function
- `WithPrivateFields(include bool)`: Include or exclude private fields in the output

There are also options we can use specific to Non Composite types like strings, maps, slices, etc.

- `prefix string`: Set a prefix for the field name
- `maskType hush.HushType (hush.TagMask or hush.TagHide)`: Set the type of masking to be applied. By default it will return the value as is.

Examples:

```go
result, err := husher.Hush(context.Background(), data,
    hush.WithSeparator("_"),
    hush.WithPrivateFields(true),
    hush.WithMaskFunc(func(s string) string {
        return "CUSTOM_MASKED"
    }),
)
```

```go
result, err := husher.Hush(context.Background(), "johndoe@mail.com", "EMAIL", hush.TagMask)
```

## Private Fields

By default, Hush doesn't process private (unexported) fields. You can include private fields in the output by using the `WithPrivateFields` option:

```go
result, err := husher.Hush(context.Background(), "",
    hush.WithPrivateFields(true),
)
```

This will include private fields in the output, applying the same masking rules as public fields.

## Notes

- Map keys are sorted alphabetically in the output for consistent results
- Slices and arrays are processed with index-based field names

## Examples

Check out the `examples` folder for more detailed usage examples:

- `basic_usage.go`: Demonstrates basic usage with a simple struct
- `custom_options.go`: Shows how to use custom options like separators and masking functions
- `complex_struct.go`: Illustrates handling of complex structs with nested fields, slices, and maps
- `custom_options_table.go`: Shows how to use custom options like separators and masking functions and display the result in a table
- `private_fields.go`: Shows how to include private fields in the output

To run an example:

```
go run examples/basic_usage.go
go run examples/custom_options.go
go run examples/complex_struct.go
go run examples/custom_options_table.go
go run examples/private_fields.go
```

## License

Hush is released under the MIT License. See the LICENSE file for details.