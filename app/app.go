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

	inputContent := a.cfg.Input
	if _, err := os.Stat(a.cfg.Input); err == nil {
		content, err := utils.ReadFile(a.cfg.Input)
		if err != nil {
			return fmt.Errorf("error reading input file: %w", err)
		}
		inputContent = string(content)
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

	if a.cfg.OutputRandomBase != "" {
		return a.handleOutputBaseRandom(inputContent)
	}

	if a.cfg.OutputOrderBase != "" {
		return a.handleOutputBaseSeq(inputContent)
	}

	if a.cfg.OutputUtf8Base != "" {
		return a.handleOutputUtf8Base(inputContent, a.cfg.InputBaseUtf8)
	}

	inputStr, err := a.processInputStr()
	if err != nil {
		return err
	}

	outputStr, err := a.processOutputStr()
	if err != nil {
		return err
	}

	if a.cfg.Output != "" && a.cfg.Input == a.cfg.Output {
		return fmt.Errorf("output file cannot be the same as input file")
	}

	inputConverter, err := a.createConverter(inputContent, inputStr, a.cfg.InputBaseStr, a.cfg.InputBaseNum, a.cfg.InputBaseUtf8, "input")
	if err != nil {
		return err
	}

	outputConverter, err := a.createConverter(inputContent, outputStr, a.cfg.OutputBaseStr, a.cfg.OutputBaseNum, a.cfg.OutputBaseUtf8, "output")
	if err != nil {
		return err
	}

	decimalValue := inputConverter.XToTenParallel(inputContent)
	result := outputConverter.TenToXParallel(decimalValue)

	return a.writeOutput(result)
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

func (a *App) createConverter(inputContent string, baseStr string, strParam string, baseNum int, baseUtf8 int, kind string) (*radix.Radix, error) {
	count := 0
	if strParam != "" {
		count++
	}
	if baseNum != 0 {
		count++
	}
	if baseUtf8 != 0 {
		count++
	}
	if count != 1 {
		if kind == "input" {
			return nil, fmt.Errorf("must use exactly one of -is, -ib, -iu")
		}
		return nil, fmt.Errorf("must use exactly one of -os, -ob, -ou")
	}

	if strParam != "" {
		return radix.NewRadixByString(baseStr), nil
	}

	if baseNum >= 2 && baseNum <= 62 {
		return radix.NewRadixByBit(baseNum), nil
	}

	if baseUtf8 != 0 {
		var N int
		if baseUtf8 < 2 {
			runes := []rune(inputContent)
			maxVal := 0
			for _, r := range runes {
				if int(r) > maxVal {
					maxVal = int(r)
				}
			}
			N = maxVal + 1
		} else {
			N = baseUtf8
		}
		// var result []rune
		// for i := 0; i < N; i++ {
		// 	result = append(result, rune(i))
		// }
		return radix.NewCharacterRadix(N), nil
		// return radix.NewRadixByString(string(result)), nil
	}

	return nil, fmt.Errorf("invalid base configuration")
}

func (a *App) writeOutput(result string) error {
	if a.cfg.Output == "" {
		fmt.Println(result)
		return nil
	}

	err := utils.WriteFile(a.cfg.Output, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}
	fmt.Printf("Result written to %s\n", a.cfg.Output)
	return nil
}
