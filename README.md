# radix-tool

A powerful command-line utility for numeral base conversions, supporting arbitrary bases (2-62) with both standard and custom character sets.

## Overview

`radix-tool` is a versatile tool for converting numbers between different numeral systems. It supports:
- Standard numeral bases (2-62): binary, octal, decimal, hexadecimal, etc.
- Custom character sets for specialized encoding schemes
- Both string and file input/output operations
- Comprehensive command-line interface

## Installation

```bash
# Clone and build
git clone <repository-url>
cd radix-tool
go build .

# Or install directly
go install radix-tool@latest
```

## Usage

```
radix-tool [FLAGS] [OPTIONS]
```

### Flags
```
-h, --help              Show this help message
```

### Options
```
-i, --input VALUE        Input value (number or file path) [required]
-ib, --input-base-num N  Input base number (2-62) [default: 10]
-is, --input-base-str STR  Input base characters [default: "0123456789"]
-o, --output FILE        Output file path (prints to stdout if omitted)
-ob, --output-base-num N Output base number (2-62, defaults to input base)
-os, --output-base-str STR Output base characters (defaults to input string)
```

### Examples

```bash
# Convert decimal 255 to hexadecimal  
radix-tool -i "255" -ib 10 -ob 16
# Output: "ff"

# Convert binary string to decimal
radix-tool -i "1010" -ib 2 -ob 10
# Output: "10"

# Convert from file input 
echo "ff" > input.txt
radix-tool -i input.txt -ib 16 -ob 10 -o output.txt

# Convert decimal 255 to binary
radix-tool --input "255" --input-base-num 10 --output-base-num 2
# Output: "11111111"

# Use custom character sets
radix-tool -i "FF" -is "0123456789ABCDEF" -ob 10
# Output: "255" (treats FF as hexadecimal in that custom alphabet)
```

## Build 

```bash
go build .
```

## Testing

```bash
go test ./...
```

## Features

- **Flexible Base Conversion**: Supports all standard bases from binary (base 2) to hexatrigesimal (base 62) and custom alphabets
- **Custom Character Sets**: Ability to define your own character mapping for non-standard numeral systems  
- **File Operations**: Both read from and write to files for batch processing
- **Efficient Implementation**: Leverages Go's `math/big` for accurate computations with large numbers

## Architecture

- `radix/` - Core radix conversion implementation with customizable base systems
- `main.go` - Command-line interface parsing and execution

## Contributing

Feel free to submit issues and enhancement requests through the standard channels.

## License

This project is distributed under the terms described in the LICENSE file.
