package app

import (
	"fmt"
	"os"

	"radix-tool/radix"
	"radix-tool/utils"
)

type App struct {
	cfg *Config
}

func New(cfg *Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() error {
	if a.cfg.Help {
		ShowHelp()
		return nil
	}

	if a.cfg.Input == "" {
		ShowHelp()
		return nil
	}

	count := 0
	if a.cfg.OutputRandomBase != "" {
		count++
	}
	if a.cfg.OutputOrderBase != "" {
		count++
	}
	if a.cfg.OutputUtf8Base != "" {
		count++
	}
	if count > 1 {
		return fmt.Errorf("can only use one of -orb, -oob, -oub at a time")
	}

	var inputRunes []rune
	var isFile bool

	if _, err := os.Stat(a.cfg.Input); err == nil {
		isFile = true
		content, err := utils.ReadFile(a.cfg.Input)
		if err != nil {
			return fmt.Errorf("error reading input file: %w", err)
		}
		if a.cfg.InputBaseByte {
			inputRunes = utils.BytesToRunes(content)
		} else {
			inputRunes = []rune(string(content))
		}
	} else {
		if a.cfg.InputBaseByte {
			return fmt.Errorf("-ib requires input to be a file path")
		}
		inputRunes = []rune(a.cfg.Input)
	}

	if a.cfg.OutputRandomBase != "" {
		return a.handleOutputBaseRandom(string(inputRunes))
	}

	if a.cfg.OutputOrderBase != "" {
		return a.handleOutputBaseSeq(string(inputRunes))
	}

	if a.cfg.OutputUtf8Base != "" {
		return a.handleOutputUtf8Base(string(inputRunes), a.cfg.InputBaseUtf8)
	}

	inputStr, err := a.processInputStr()
	if err != nil {
		return err
	}

	outputStr, err := a.processOutputStr()
	if err != nil {
		return err
	}

	if a.cfg.Output != "" && isFile && a.cfg.Input == a.cfg.Output {
		return fmt.Errorf("output file cannot be the same as input file")
	}

	inputConverter, err := a.createInputConverter(inputRunes, inputStr)
	if err != nil {
		return err
	}

	outputConverter, err := a.createOutputConverter(outputStr)
	if err != nil {
		return err
	}

	decimalValue := inputConverter.XToTenParallel(inputRunes)
	resultRunes := outputConverter.TenToXParallel(decimalValue)

	return a.writeOutputRunes(resultRunes)
}

func (a *App) processInputStr() (string, error) {
	if a.cfg.InputBaseStr == "" {
		return "", nil
	}

	inputStr := a.cfg.InputBaseStr
	if _, err := os.Stat(a.cfg.InputBaseStr); err == nil {
		content, err := utils.ReadFile(a.cfg.InputBaseStr)
		if err != nil {
			return "", fmt.Errorf("error reading -is file: %w", err)
		}
		inputStr = string(content)
	}

	if utils.HasDuplicateChars(inputStr) {
		return "", fmt.Errorf("input string contains duplicate characters")
	}

	return inputStr, nil
}

func (a *App) processOutputStr() (string, error) {
	if a.cfg.OutputBaseStr == "" {
		return "", nil
	}

	outputStr := a.cfg.OutputBaseStr
	if _, err := os.Stat(a.cfg.OutputBaseStr); err == nil {
		content, err := utils.ReadFile(a.cfg.OutputBaseStr)
		if err != nil {
			return "", fmt.Errorf("error reading -os file: %w", err)
		}
		outputStr = string(content)
	}

	if utils.HasDuplicateChars(outputStr) {
		return "", fmt.Errorf("output string contains duplicate characters")
	}

	return outputStr, nil
}

func (a *App) handleOutputBaseRandom(inputContent string) error {
	result := utils.ShuffleString(utils.RemoveDuplicates(inputContent))

	err := utils.WriteFile(a.cfg.OutputRandomBase, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("error writing random output file: %w", err)
	}
	fmt.Printf("Random output written to %s\n", a.cfg.OutputRandomBase)
	return nil
}

func (a *App) handleOutputBaseSeq(inputContent string) error {
	result := utils.RemoveDuplicates(inputContent)

	err := utils.WriteFile(a.cfg.OutputOrderBase, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("error writing sequence output file: %w", err)
	}
	fmt.Printf("Sequence output written to %s\n", a.cfg.OutputOrderBase)
	return nil
}

func (a *App) handleOutputUtf8Base(inputContent string, oub int) error {
	var N int
	if oub < 2 {
		runes := []rune(inputContent)
		maxVal := 0
		for _, r := range runes {
			if int(r) > maxVal {
				maxVal = int(r)
			}
		}
		N = maxVal + 1
	} else {
		N = oub
	}

	var result []rune
	for i := 0; i < N; i++ {
		result = append(result, rune(i))
	}

	err := utils.WriteFile(a.cfg.OutputUtf8Base, []byte(string(result)), 0644)
	if err != nil {
		return fmt.Errorf("error writing sequence output file: %w", err)
	}
	fmt.Printf("Utf8 %d output written to %s\n", N, a.cfg.OutputUtf8Base)
	return nil
}

func (a *App) createInputConverter(inputRunes []rune, baseStr string) (*radix.Radix, error) {
	count := 0
	if a.cfg.InputBaseStr != "" {
		count++
	}
	if a.cfg.InputBaseNum != 0 {
		count++
	}
	if a.cfg.InputBaseUtf8 != 0 {
		count++
	}
	if a.cfg.InputBaseByte {
		count++
	}
	if count != 1 {
		return nil, fmt.Errorf("must use exactly one of -is, -in, -iu, -ib")
	}

	if a.cfg.InputBaseByte {
		return radix.NewCharacterRadix(256), nil
	}

	if a.cfg.InputBaseStr != "" {
		return radix.NewRadixByString(baseStr), nil
	}

	if a.cfg.InputBaseNum >= 2 && a.cfg.InputBaseNum <= 62 {
		return radix.NewRadixByBit(a.cfg.InputBaseNum), nil
	}

	if a.cfg.InputBaseUtf8 != 0 {
		var N int
		if a.cfg.InputBaseUtf8 < 2 {
			maxVal := 0
			for _, r := range inputRunes {
				if int(r) > maxVal {
					maxVal = int(r)
				}
			}
			N = maxVal + 1
		} else {
			N = a.cfg.InputBaseUtf8
		}
		return radix.NewCharacterRadix(N), nil
	}

	return nil, fmt.Errorf("invalid input base configuration")
}

func (a *App) createOutputConverter(baseStr string) (*radix.Radix, error) {
	count := 0
	if a.cfg.OutputBaseStr != "" {
		count++
	}
	if a.cfg.OutputBaseNum != 0 {
		count++
	}
	if a.cfg.OutputBaseUtf8 != 0 {
		count++
	}
	if a.cfg.OutputBaseByte {
		count++
	}
	if count != 1 {
		return nil, fmt.Errorf("must use exactly one of -os, -on, -ou, -ob")
	}

	if a.cfg.OutputBaseByte {
		return radix.NewCharacterRadix(256), nil
	}

	if a.cfg.OutputBaseStr != "" {
		return radix.NewRadixByString(baseStr), nil
	}

	if a.cfg.OutputBaseNum >= 2 && a.cfg.OutputBaseNum <= 62 {
		return radix.NewRadixByBit(a.cfg.OutputBaseNum), nil
	}

	if a.cfg.OutputBaseUtf8 != 0 {
		N := a.cfg.OutputBaseUtf8
		if N < 2 {
			N = 256
		}
		return radix.NewCharacterRadix(N), nil
	}

	return nil, fmt.Errorf("invalid output base configuration")
}

func (a *App) writeOutputRunes(result []rune) error {
	var data []byte
	var err error
	if a.cfg.OutputBaseByte {
		data, err = utils.RunesToBytes(result)
		if err != nil {
			return fmt.Errorf("error converting runes to bytes: %w", err)
		}
	} else {
		data = []byte(string(result))
	}

	if a.cfg.Output == "" {
		if a.cfg.OutputBaseByte {
			return fmt.Errorf("-ob requires -o to specify output file")
		}
		fmt.Println(string(result))
		return nil
	}

	err = utils.WriteFile(a.cfg.Output, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}
	fmt.Printf("Result written to %s (%d bytes)\n", a.cfg.Output, len(data))
	return nil
}
