
# Golang Package Examples

This repository contains a collection of example programs demonstrating the usage of various features in our Golang package. Each example showcases different functionalities, such as masking sensitive information, handling complex structures, and customizing options.

## Examples

### 1. Basic Usage

**File:** `basic_usage.go`

**Description:** Demonstrates basic usage of the package, including string handling with and without a prefix.

**Output:**

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

### 2. Complex Struct Handling

**File:** `complex_struct.go`

**Description:** Shows how to manage complex structures with nested fields, including masking sensitive data within structs.

**Output:**

```
Complex Struct Example:
Accounts[0].Balance: 10**.5
Accounts[0].Currency: USD
Accounts[0].ID: 1
Accounts[1].Balance: 50**75
Accounts[1].Currency: EUR
Accounts[1].ID: 2
Address.City: Anytown
Address.Country: USA
Address.Street: 123 Main St
Age: 28
Email: HIDDEN
IsActive: true
Metadata[lastLogin]: 2023-04-01
Metadata[role]: admin
Name: Al*********on
```

### 3. Custom Options

**File:** `custom_options.go`

**Description:** Illustrates how to use custom options in the package to configure behavior, such as masking specific fields.

**Output:**

```
Custom Options Example:
APIKey: ------------
Debug: true
SecretKey: HIDDEN
```

### 4. Custom Options with Table Output

**File:** `custom_options_table.go`

**Description:** Similar to the previous example but displays the custom options in a tabular format.

**Output:**

```
Custom Options Example:
+-----------+--------------+
|   FIELD   |    VALUE     |
+-----------+--------------+
| APIKey    | ------------ |
| Debug     | true         |
| SecretKey | HIDDEN       |
+-----------+--------------+
```

### 5. Private Fields Handling

**File:** `private_fields.go`

**Description:** Demonstrates how the package handles private fields within structs, with and without specific configurations.

**Output:**

```
Without private fields:
Address: 123 Main St
Name: Jo****oe

With private fields:
Address: 123 Main St
Name: Jo****oe
age: **
```

### 6. Custom Regex Function

**File:** `custom_regex_function.go`

**Description:** Shows how to use a custom regular expression function to handle specific field patterns within the package.

**Output:**

```
Custom Options Example:
Debug: true
Email: This is a test for masking tl********************** email address
SecretKey: verysecret
```

## How to Run the Examples

To run any of these examples, simply execute the corresponding Go file using the following command:

```bash
go run <example_file.go>
```

For example, to run the `basic_usage.go` file, you would use:

```bash
go run basic_usage.go
```

## Contributing

Contributions to improve the examples or add new ones are welcome! Please submit a pull request or open an issue if you encounter any problems.
