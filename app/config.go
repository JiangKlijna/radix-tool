package app

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Input            string
	InputBaseNum     int
	InputBaseStr     string
	Output           string
	OutputBaseNum    int
	OutputBaseStr    string
	OutputRandomBase string
	OutputOrderBase  string
	Help             bool
}

func ParseFlags(args []string) (*Config, error) {
	cfg := &Config{}

	flagSet := flag.NewFlagSet("radix-tool", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	flagSet.StringVar(&cfg.Input, "i", "", "Input value (number or file path)")
	flagSet.StringVar(&cfg.Input, "input", "", "Input value (number or file path)")

	flagSet.IntVar(&cfg.InputBaseNum, "in", 0, "Input base number (2-62)")
	flagSet.IntVar(&cfg.InputBaseNum, "input-base-num", 0, "Input base number (2-62)")

	flagSet.StringVar(&cfg.InputBaseStr, "is", "", "Input base characters")
	flagSet.StringVar(&cfg.InputBaseStr, "input-base-str", "", "Input base characters")

	flagSet.StringVar(&cfg.Output, "o", "", "Output file path (if empty, prints to cmd)")
	flagSet.StringVar(&cfg.Output, "output", "", "Output file path (if empty, prints to cmd)")

	flagSet.IntVar(&cfg.OutputBaseNum, "on", 0, "Output base number (2-62, defaults to input base)")
	flagSet.IntVar(&cfg.OutputBaseNum, "output-base-num", 0, "Output base number (2-62, defaults to input base)")

	flagSet.StringVar(&cfg.OutputBaseStr, "os", "", "Output base characters")
	flagSet.StringVar(&cfg.OutputBaseStr, "output-base-str", "", "Output base characters")

	flagSet.StringVar(&cfg.OutputRandomBase, "orb", "", "Randomly output characters from -is (empty means to stderr)")
	flagSet.StringVar(&cfg.OutputRandomBase, "output-random-base", "", "Randomly output characters from -is (empty means to stderr)")

	flagSet.StringVar(&cfg.OutputOrderBase, "oob", "", "Sequence output characters from -is (empty means to stderr)")
	flagSet.StringVar(&cfg.OutputOrderBase, "output-order-base", "", "Sequence output characters from -is (empty means to stderr)")

	flagSet.BoolVar(&cfg.Help, "help", false, "Show help")
	flagSet.BoolVar(&cfg.Help, "h", false, "Show help")

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ShowHelp() {
	fmt.Println(`radix-tool - A command-line tool for numeral base conversions

USAGE:
  radix-tool [FLAGS] [OPTIONS]

FLAGS:
  -h, --help              Show this help message

OPTIONS:
  -i, --input VALUE        Input value (number or file path) [required]
  -in, --input-base-num N  Input base number (2-62)
  -is, --input-base-str STR  Input base characters
  -o, --output FILE        Output file path (prints to stdout if omitted)
  -on, --output-base-num N Output base number (2-62, defaults to input base) 
  -os, --output-base-str STR Output base characters
  -orb, --output-random-base  Randomly output characters from -is (empty means to stderr)
  -oob, --output-order-base  Sequence output characters from -is (empty means to stderr)

EXAMPLES:
  # Convert decimal 255 to hexadecimal
  radix-tool -i "255" -in 10 -on 16
  
  # Convert binary string to decimal
  radix-tool -i "1010" -in 2 -on 10
  
  # Convert decimal 255 to binary
  radix-tool --input "255" --input-base-num 10 --output-base-num 2`)
}
