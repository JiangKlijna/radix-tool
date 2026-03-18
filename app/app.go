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

	if a.cfg.OutputRandomBase != "" && a.cfg.OutputOrderBase != "" {
		return fmt.Errorf("cannot use both -orb and -oob at the same time")
	}

	if a.cfg.OutputRandomBase != "" {
		return a.handleOutputBaseRandom(inputContent)
	}

	if a.cfg.OutputOrderBase != "" {
		return a.handleOutputBaseSeq(inputContent)
	}

	if a.cfg.InputBaseStr == "" && a.cfg.InputBaseNum == 0 {
		return fmt.Errorf("must use either -is or -ib")
	}

	if a.cfg.InputBaseStr != "" && a.cfg.InputBaseNum != 0 {
		return fmt.Errorf("cannot use both -is and -ib at the same time, please use either one")
	}

	if a.cfg.OutputBaseStr != "" && a.cfg.OutputBaseNum != 0 {
		return fmt.Errorf("cannot use both -os and -ob at the same time, please use either one")
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

	inputConverter, err := a.createConverter(inputStr, a.cfg.InputBaseStr, a.cfg.InputBaseNum)
	if err != nil {
		return err
	}

	outputConverter, err := a.createConverter(outputStr, a.cfg.OutputBaseStr, a.cfg.OutputBaseNum)
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

func (a *App) createConverter(str string, strParam string, base int) (*radix.Radix, error) {
	if strParam != "" {
		return radix.NewRadixByString(str), nil
	}

	if base >= 2 && base <= 62 {
		return radix.NewRadixByBit(base), nil
	}

	return nil, fmt.Errorf("base must be between 2 and 62, got: %d", base)
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
