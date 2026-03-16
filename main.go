package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"radix-tool/radix"
	"radix-tool/utils"
)

func main() {
	var (
		input            string
		inputBase        int
		inputStr         string
		output           string
		outputBase       int
		outputStr        string
		outputBaseRandom string
		outputBaseSeq    string
		help             bool
	)

	// 定义简短参数名称
	flag.StringVar(&input, "i", "", "Input value (number or file path)")
	flag.StringVar(&input, "input", "", "Input value (number or file path)")

	flag.IntVar(&inputBase, "ib", 0, "Input base number (2-62)")
	flag.IntVar(&inputBase, "input-base-num", 0, "Input base number (2-62)")

	flag.StringVar(&inputStr, "is", "", "Input base characters")

	flag.StringVar(&output, "o", "", "Output file path (if empty, prints to cmd)")
	flag.StringVar(&output, "output", "", "Output file path (if empty, prints to cmd)")

	flag.IntVar(&outputBase, "ob", 0, "Output base number (2-62, defaults to input base)")
	flag.IntVar(&outputBase, "output-base-num", 0, "Output base number (2-62, defaults to input base)")

	flag.StringVar(&outputStr, "os", "", "Output base characters")

	flag.StringVar(&outputBaseRandom, "obr", "", "Randomly output characters from -is (empty means to stderr)")
	flag.StringVar(&outputBaseRandom, "output-base-random", "", "Randomly output characters from -is (empty means to stderr)")

	flag.StringVar(&outputBaseSeq, "obs", "", "Sequence output characters from -is (empty means to stderr)")
	flag.StringVar(&outputBaseSeq, "output-base-sequence", "", "Sequence output characters from -is (empty means to stderr)")

	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		showHelp()
		return
	}

	if input == "" {
		showHelp()
		return
	}

	// 处理输入-如果输入是文件路径，则读取文件内容
	inputContent := input
	if _, err := os.Stat(input); err == nil {
		// 是文件
		content, err := ioutil.ReadFile(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
			os.Exit(1)
		}
		inputContent = string(content)
	}

	if outputBaseRandom != "" && outputBaseSeq != "" {
		fmt.Fprintf(os.Stderr, "Cannot use both -obr and -obs at the same time\n")
		os.Exit(1)
	}

	// 处理 output-base-random
	if outputBaseRandom != "" {
		result := utils.ShuffleString(utils.RemoveDuplicates(inputContent))

		err := ioutil.WriteFile(outputBaseRandom, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing random output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", result)

		return
	}

	// 处理 output-base-sequence
	if outputBaseSeq != "" {
		result := utils.RemoveDuplicates(inputContent)

		err := ioutil.WriteFile(outputBaseSeq, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing sequence output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", result)

		return
	}

	if inputStr == "" && inputBase == 0 {
		fmt.Fprintf(os.Stderr, "Must use either -is or -ib\n")
		os.Exit(1)
	}

	if inputStr != "" && inputBase != 0 {
		fmt.Fprintf(os.Stderr, "Cannot use both -is and -ib at the same time, please use either one\n")
		os.Exit(1)
	}

	if outputStr != "" && outputBase != 0 {
		fmt.Fprintf(os.Stderr, "Cannot use both -os and -ob at the same time, please use either one\n")
		os.Exit(1)
	}

	if inputStr != "" {
		if _, err := os.Stat(inputStr); err == nil {
			// 是文件
			content, err := ioutil.ReadFile(inputStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading -is file: %v\n", err)
				os.Exit(1)
			}
			inputStr = string(content)
		}
		if utils.HasDuplicateChars(inputStr) {
			fmt.Fprintf(os.Stderr, "Input string contains duplicate characters\n")
			os.Exit(1)
		}
	}

	if outputStr != "" {
		if _, err := os.Stat(outputStr); err == nil {
			// 是文件
			content, err := ioutil.ReadFile(outputStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading -os file: %v\n", err)
				os.Exit(1)
			}
			outputStr = string(content)
		}
		if utils.HasDuplicateChars(outputStr) {
			fmt.Fprintf(os.Stderr, "Output string contains duplicate characters\n")
			os.Exit(1)
		}
	}

	// 如果指定输出文件，确保输入中包含文件扩展名避免覆盖
	if output != "" && input == output {
		fmt.Fprintf(os.Stderr, "Output file cannot be the same as input file\n")
		os.Exit(1)
	}

	var inputConverter, outputConverter *radix.Radix
	if inputStr != "" {
		// 使用自定义字符集
		inputConverter = radix.NewRadixByString(inputStr)
	} else {
		// 使用指定的基础数值
		if inputBase >= 2 && inputBase <= 62 {
			inputConverter = radix.NewRadixByBit(inputBase)
		} else {
			fmt.Fprintf(os.Stderr, "Input base must be between 2 and 62, got: %d\n", inputBase)
			os.Exit(1)
		}
	}

	if outputStr != "" {
		// 使用自定义字符集
		outputConverter = radix.NewRadixByString(outputStr)
	} else {
		// 使用指定的基础数值
		if outputBase >= 2 && outputBase <= 62 {
			outputConverter = radix.NewRadixByBit(outputBase)
		} else {
			fmt.Fprintf(os.Stderr, "Output base must be between 2 and 62, got: %d\n", outputBase)
			os.Exit(1)
		}
	}

	// 执行转换：首先将输入值转换为10进制big.Int，然后转换到目标格式
	decimalValue := inputConverter.XStrToTen(inputContent)
	result := outputConverter.TenToX(decimalValue)

	// 输出结果
	if output == "" {
		fmt.Println(result)
	} else {
		err := ioutil.WriteFile(output, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Result written to %s\n", output)
	}
}

func showHelp() {
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
