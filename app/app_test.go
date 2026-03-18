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
			args:          []string{"-i", "255", "-ib", "10", "-ob", "16"},
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
			if cfg.InputBase != tt.expectedIB {
				t.Errorf("ParseFlags() InputBase = %v, want %v", cfg.InputBase, tt.expectedIB)
			}
			if cfg.OutputBase != tt.expectedOB {
				t.Errorf("ParseFlags() OutputBase = %v, want %v", cfg.OutputBase, tt.expectedOB)
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

func TestAppRun_MissingInputBase(t *testing.T) {
	cfg := &Config{
		Input:     "255",
		InputBase: 0,
		InputStr:  "",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when missing input base")
	}
}

func TestAppRun_InvalidInputBase(t *testing.T) {
	cfg := &Config{
		Input:     "255",
		InputBase: 1,
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error for invalid input base")
	}
}

func TestAppRun_InvalidOutputBase(t *testing.T) {
	cfg := &Config{
		Input:      "255",
		InputBase:  10,
		OutputBase: 100,
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error for invalid output base")
	}
}

func TestAppRun_ConflictingInputOptions(t *testing.T) {
	cfg := &Config{
		Input:     "255",
		InputBase: 10,
		InputStr:  "0123456789",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when both -is and -ib are set")
	}
}

func TestAppRun_ConflictingOutputOptions(t *testing.T) {
	cfg := &Config{
		Input:      "255",
		InputBase:  10,
		OutputBase: 16,
		OutputStr:  "0123456789abcdef",
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
		OutputBaseRandom: "random.txt",
		OutputBaseSeq:    "seq.txt",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when both -obr and -obs are set")
	}
}

func TestAppRun_SameInputOutputFile(t *testing.T) {
	cfg := &Config{
		Input:     "test.txt",
		InputBase: 10,
		Output:    "test.txt",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when input and output files are the same")
	}
}

func TestAppRun_DecimalToHex(t *testing.T) {
	cfg := &Config{
		Input:      "255",
		InputBase:  10,
		OutputBase: 16,
	}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
}

func TestAppRun_BinaryToDecimal(t *testing.T) {
	cfg := &Config{
		Input:      "1010",
		InputBase:  2,
		OutputBase: 10,
	}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
}

func TestAppRun_CustomCharSet(t *testing.T) {
	cfg := &Config{
		Input:     "FF",
		InputStr:  "0123456789ABCDEF",
		OutputStr: "0123456789",
	}
	app := New(cfg)

	err := app.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
}

func TestAppRun_InputStrWithDuplicates(t *testing.T) {
	cfg := &Config{
		Input:    "test",
		InputStr: "001122",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when input string has duplicate characters")
	}
}

func TestAppRun_OutputStrWithDuplicates(t *testing.T) {
	cfg := &Config{
		Input:     "test",
		InputBase: 10,
		OutputStr: "001122",
	}
	app := New(cfg)

	err := app.Run()
	if err == nil {
		t.Error("Run() should return error when output string has duplicate characters")
	}
}
