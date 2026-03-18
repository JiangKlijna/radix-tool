package app

import (
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedInput string
		expectedIB    int
		expectedOB    int
		expectHelp    bool
	}{
		{
			name:          "basic conversion",
			args:          []string{"-i", "255", "-in", "10", "-on", "16"},
			expectedInput: "255",
			expectedIB:    10,
			expectedOB:    16,
			expectHelp:    false,
		},
		{
			name:          "help flag short",
			args:          []string{"-h"},
			expectedInput: "",
			expectedIB:    0,
			expectedOB:    0,
			expectHelp:    true,
		},
		{
			name:          "help flag long",
			args:          []string{"--help"},
			expectedInput: "",
			expectedIB:    0,
			expectedOB:    0,
			expectHelp:    true,
		},
		{
			name:          "long flags",
			args:          []string{"--input", "1010", "--input-base-num", "2", "--output-base-num", "10"},
			expectedInput: "1010",
			expectedIB:    2,
			expectedOB:    10,
			expectHelp:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ParseFlags(tt.args)
			if err != nil {
				t.Fatalf("ParseFlags() error = %v", err)
			}
			if cfg.Input != tt.expectedInput {
				t.Errorf("ParseFlags() Input = %v, want %v", cfg.Input, tt.expectedInput)
			}
			if cfg.InputBaseNum != tt.expectedIB {
				t.Errorf("ParseFlags() InputBaseNum = %v, want %v", cfg.InputBaseNum, tt.expectedIB)
			}
			if cfg.OutputBaseNum != tt.expectedOB {
				t.Errorf("ParseFlags() OutputBaseNum = %v, want %v", cfg.OutputBaseNum, tt.expectedOB)
			}
			if cfg.Help != tt.expectHelp {
				t.Errorf("ParseFlags() Help = %v, want %v", cfg.Help, tt.expectHelp)
			}
		})
	}
}

func TestAppRun_Help(t *testing.T) {
	cfg := &Config{Help: true}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() with help flag should return nil, got %v", err)
	}
}

func TestAppRun_NoInput(t *testing.T) {
	cfg := &Config{Input: ""}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() with no input should return nil, got %v", err)
	}
}

func TestAppRun_MissingInputBaseNum(t *testing.T) {
	cfg := &Config{
		Input:        "255",
		InputBaseNum: 0,
		InputBaseStr: "",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when missing input base")
	}
}

func TestAppRun_InvalidInputBaseNum(t *testing.T) {
	cfg := &Config{
		Input:        "255",
		InputBaseNum: 1,
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error for invalid input base")
	}
}

func TestAppRun_InvalidOutputBaseNum(t *testing.T) {
	cfg := &Config{
		Input:         "255",
		InputBaseNum:  10,
		OutputBaseNum: 100,
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error for invalid output base")
	}
}

func TestAppRun_ConflictingInputOptions(t *testing.T) {
	cfg := &Config{
		Input:        "255",
		InputBaseNum: 10,
		InputBaseStr: "0123456789",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when both -is and -ib are set")
	}
}

func TestAppRun_ConflictingOutputOptions(t *testing.T) {
	cfg := &Config{
		Input:         "255",
		InputBaseNum:  10,
		OutputBaseNum: 16,
		OutputBaseStr: "0123456789abcdef",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when both -os and -ob are set")
	}
}

func TestAppRun_ConflictingRandomAndSeq(t *testing.T) {
	cfg := &Config{
		Input:            "abc",
		OutputRandomBase: "random.txt",
		OutputOrderBase:  "seq.txt",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when both -obr and -obs are set")
	}
}

func TestAppRun_SameInputOutputFile(t *testing.T) {
	cfg := &Config{
		Input:        "test.txt",
		InputBaseNum: 10,
		Output:       "test.txt",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when input and output files are the same")
	}
}

func TestAppRun_DecimalToHex(t *testing.T) {
	cfg := &Config{
		Input:         "255",
		InputBaseNum:  10,
		OutputBaseNum: 16,
	}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
}

func TestAppRun_BinaryToDecimal(t *testing.T) {
	cfg := &Config{
		Input:         "1010",
		InputBaseNum:  2,
		OutputBaseNum: 10,
	}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
}

func TestAppRun_CustomCharSet(t *testing.T) {
	cfg := &Config{
		Input:         "FF",
		InputBaseStr:  "0123456789ABCDEF",
		OutputBaseStr: "0123456789",
	}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
}

func TestAppRun_InputStrWithDuplicates(t *testing.T) {
	cfg := &Config{
		Input:        "test",
		InputBaseStr: "001122",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when input string has duplicate characters")
	}
}

func TestAppRun_OutputStrWithDuplicates(t *testing.T) {
	cfg := &Config{
		Input:         "test",
		InputBaseNum:  10,
		OutputBaseStr: "001122",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when output string has duplicate characters")
	}
}
