package app

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Input            string
	InputBase        int
	InputStr         string
	Output           string
	OutputBase       int
	OutputStr        string
	OutputBaseRandom string
	OutputBaseSeq    string
	Help             bool
}

func ParseFlags(args []string) (*Config, error) {
	cfg := &Config{}

	flagSet := flag.NewFlagSet("radix-tool", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	flagSet.StringVar(&cfg.Input, "i", "", "Input value (number or file path)")
	flagSet.StringVar(&cfg.Input, "input", "", "Input value (number or file path)")

	flagSet.IntVar(&cfg.InputBase, "ib", 0, "Input base number (2-62)")
	flagSet.IntVar(&cfg.InputBase, "input-base-num", 0, "Input base number (2-62)")

	flagSet.StringVar(&cfg.InputStr, "is", "", "Input base characters")

	flagSet.StringVar(&cfg.Output, "o", "", "Output file path (if empty, prints to cmd)")
	flagSet.StringVar(&cfg.Output, "output", "", "Output file path (if empty, prints to cmd)")

	flagSet.IntVar(&cfg.OutputBase, "ob", 0, "Output base number (2-62, defaults to input base)")
	flagSet.IntVar(&cfg.OutputBase, "output-base-num", 0, "Output base number (2-62, defaults to input base)")

	flagSet.StringVar(&cfg.OutputStr, "os", "", "Output base characters")

	flagSet.StringVar(&cfg.OutputBaseRandom, "obr", "", "Randomly output characters from -is (empty means to stderr)")
	flagSet.StringVar(&cfg.OutputBaseRandom, "output-base-random", "", "Randomly output characters from -is (empty means to stderr)")

	flagSet.StringVar(&cfg.OutputBaseSeq, "obs", "", "Sequence output characters from -is (empty means to stderr)")
	flagSet.StringVar(&cfg.OutputBaseSeq, "output-base-sequence", "", "Sequence output characters from -is (empty means to stderr)")

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
  -ib, --input-base-num N  Input base number (2-62)
  -is, --input-base-str STR  Input base characters
  -o, --output FILE        Output file path (prints to stdout if omitted)
  -ob, --output-base-num N Output base number (2-62, defaults to input base) 
  -os, --output-base-str STR Output base characters
  -obr, --output-base-random  Randomly output characters from -i input
  -obs, --output-base-sequence  Sequence output characters from -i input

EXAMPLES:
  # Convert decimal 255 to hexadecimal
  radix-tool -i "255" -ib 10 -ob 16
  
  # Convert binary string to decimal
  radix-tool -i "1010" -ib 2 -ob 10
  
  # Convert decimal 255 to binary
  radix-tool --input "255" --input-base-num 10 --output-base-num 2`)
}
