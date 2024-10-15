# Go JSON Parser

A simple JSON parser implemented in Go, featuring a lexer and a parser for JSON validation and tokenization.

## Features

- Lexer for tokenizing JSON input
- Parser for validating JSON structure
- Supports all JSON data types: objects, arrays, strings, numbers, booleans, and null
- Handles nested structures
- Proper error reporting for invalid JSON

## Package Structure

The project consists of two main files:

1. `lexer.go`: Contains the lexer implementation for tokenizing JSON input.
2. `parser.go`: Implements the parser for validating JSON structure.

## Usage

To use this JSON parser in your Go project:

1. Import the package:

```go
import "path/to/jsonparser"
```

2. Create a new parser instance with your JSON input:

```go
input := []byte(`{"key": "value", "array": [1, 2, 3]}`)
parser, err := jsonparser.NewParser(input)
if err != nil {
    // Handle error
}
```

3. Parse the JSON:

```go
valid, err := parser.Parse()
if err != nil {
    // Handle parsing error
}
if valid {
    fmt.Println("JSON is valid")
} else {
    fmt.Println("JSON is invalid")
}
```

## Lexer

The lexer (`lexer.go`) tokenizes the input JSON string into a series of tokens. It supports all JSON token types, including:

- Braces and brackets
- Colons and commas
- Strings (with proper escape sequence handling)
- Numbers
- Booleans
- Null
- Whitespace (ignored)

## Parser

The parser (`parser.go`) validates the JSON structure using the tokens provided by the lexer. It checks for:

- Balanced braces and brackets
- Proper object and array structures
- Valid value types
- Correct use of commas and colons

## Error Handling

The parser provides detailed error messages for various JSON structure issues, including:

- Unbalanced braces or brackets
- Unexpected tokens
- Invalid number formats
- Improper use of commas
- Incomplete JSON structures

## Limitations

- The parser focuses on validation rather than data extraction.
- It does not create a structured representation of the JSON data (e.g., as a map or struct).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License
