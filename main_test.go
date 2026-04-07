package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"radix-tool/app"
	"radix-tool/radix"
	"radix-tool/utils"
)

var testDir string

func setupTestDir(t *testing.T) {
	testDir = filepath.Join("temp", "testfiles")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	textContent := "Hello World!"
	if err := os.WriteFile(filepath.Join(testDir, "input_text.txt"), []byte(textContent), 0644); err != nil {
		t.Fatalf("Failed to create input_text.txt: %v", err)
	}

	binaryData := []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100}
	if err := os.WriteFile(filepath.Join(testDir, "input_binary.bin"), binaryData, 0644); err != nil {
		t.Fatalf("Failed to create input_binary.bin: %v", err)
	}

	base64Like := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	if err := os.WriteFile(filepath.Join(testDir, "base_custom.txt"), []byte(base64Like), 0644); err != nil {
		t.Fatalf("Failed to create base_custom.txt: %v", err)
	}
}

func verifyFiles(t *testing.T, originalPath, restoredPath, testName string) {
	origBytes, err := os.ReadFile(originalPath)
	if err != nil {
		t.Errorf("[%s] Failed to read original file: %v", testName, err)
		return
	}

	restoredBytes, err := os.ReadFile(restoredPath)
	if err != nil {
		t.Errorf("[%s] Failed to read restored file: %v", testName, err)
		return
	}

	if len(origBytes) != len(restoredBytes) {
		t.Errorf("[%s] File lengths differ: original=%d, restored=%d", testName, len(origBytes), len(restoredBytes))
		return
	}

	if !bytes.Equal(origBytes, restoredBytes) {
		t.Errorf("[%s] Files content differs", testName)
		return
	}

	t.Logf("[%s] Success: files identical (%d bytes)", testName, len(origBytes))
}

func TestTextFileToHexRoundTrip(t *testing.T) {
	setupTestDir(t)

	cfg1 := &app.Config{
		Input:         filepath.Join(testDir, "input_text.txt"),
		InputBaseByte: true,
		OutputBaseNum: 16,
		Output:        filepath.Join(testDir, "test1_hex.txt"),
	}
	app1 := app.New(cfg1)
	if err := app1.Run(); err != nil {
		t.Fatalf("Step 1 failed: %v", err)
	}

	cfg2 := &app.Config{
		Input:          filepath.Join(testDir, "test1_hex.txt"),
		InputBaseNum:   16,
		OutputBaseByte: true,
		Output:         filepath.Join(testDir, "test1_restored.txt"),
	}
	app2 := app.New(cfg2)
	if err := app2.Run(); err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}

	verifyFiles(t, filepath.Join(testDir, "input_text.txt"), filepath.Join(testDir, "test1_restored.txt"), "Test1")
}

func TestBinaryFileToBase62RoundTrip(t *testing.T) {
	setupTestDir(t)

	cfg1 := &app.Config{
		Input:         filepath.Join(testDir, "input_binary.bin"),
		InputBaseByte: true,
		OutputBaseNum: 62,
		Output:        filepath.Join(testDir, "test2_base62.txt"),
	}
	app1 := app.New(cfg1)
	if err := app1.Run(); err != nil {
		t.Fatalf("Step 1 failed: %v", err)
	}

	cfg2 := &app.Config{
		Input:          filepath.Join(testDir, "test2_base62.txt"),
		InputBaseNum:   62,
		OutputBaseByte: true,
		Output:         filepath.Join(testDir, "test2_restored.bin"),
	}
	app2 := app.New(cfg2)
	if err := app2.Run(); err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}

	verifyFiles(t, filepath.Join(testDir, "input_binary.bin"), filepath.Join(testDir, "test2_restored.bin"), "Test2")
}

func TestBinaryFileToCustomBaseRoundTrip(t *testing.T) {
	setupTestDir(t)

	cfg1 := &app.Config{
		Input:         filepath.Join(testDir, "input_binary.bin"),
		InputBaseByte: true,
		OutputBaseStr: filepath.Join(testDir, "base_custom.txt"),
		Output:        filepath.Join(testDir, "test3_custom.txt"),
	}
	app1 := app.New(cfg1)
	if err := app1.Run(); err != nil {
		t.Fatalf("Step 1 failed: %v", err)
	}

	cfg2 := &app.Config{
		Input:          filepath.Join(testDir, "test3_custom.txt"),
		InputBaseStr:   filepath.Join(testDir, "base_custom.txt"),
		OutputBaseByte: true,
		Output:         filepath.Join(testDir, "test3_restored.bin"),
	}
	app2 := app.New(cfg2)
	if err := app2.Run(); err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}

	verifyFiles(t, filepath.Join(testDir, "input_binary.bin"), filepath.Join(testDir, "test3_restored.bin"), "Test3")
}

func TestDecimalToHexRoundTrip(t *testing.T) {
	original := "255"

	converter10 := radix.NewRadixByBit(10)
	converter16 := radix.NewRadixByBit(16)

	runes10 := []rune(original)
	decimal := converter10.XToTen(runes10)
	hexRunes := converter16.TenToX(decimal)
	decimal2 := converter16.XToTen(hexRunes)
	resultRunes := converter10.TenToX(decimal2)
	result := string(resultRunes)

	if original != result {
		t.Errorf("[Test4] Expected %s, got %s", original, result)
	} else {
		t.Logf("[Test4] Success: %s -> %s -> %s", original, string(hexRunes), result)
	}
}

func TestDecimalToBase62RoundTrip(t *testing.T) {
	original := "999999999999"

	converter10 := radix.NewRadixByBit(10)
	converter62 := radix.NewRadixByBit(62)

	runes10 := []rune(original)
	decimal := converter10.XToTen(runes10)
	base62Runes := converter62.TenToX(decimal)
	decimal2 := converter62.XToTen(base62Runes)
	resultRunes := converter10.TenToX(decimal2)
	result := string(resultRunes)

	if original != result {
		t.Errorf("[Test5] Expected %s, got %s", original, result)
	} else {
		t.Logf("[Test5] Success: %s -> %s -> %s", original, string(base62Runes), result)
	}
}

func TestDecimalToBase36RoundTrip(t *testing.T) {
	original := "12345678"

	converter10 := radix.NewRadixByBit(10)
	converter36 := radix.NewRadixByBit(36)

	runes10 := []rune(original)
	decimal := converter10.XToTen(runes10)
	base36Runes := converter36.TenToX(decimal)
	decimal2 := converter36.XToTen(base36Runes)
	resultRunes := converter10.TenToX(decimal2)
	result := string(resultRunes)

	if original != result {
		t.Errorf("[Test6] Expected %s, got %s", original, result)
	} else {
		t.Logf("[Test6] Success: %s -> %s -> %s", original, string(base36Runes), result)
	}
}

func TestBinaryToDecimalRoundTrip(t *testing.T) {
	original := "10101010"

	converter2 := radix.NewRadixByBit(2)
	converter10 := radix.NewRadixByBit(10)

	runes2 := []rune(original)
	decimal := converter2.XToTen(runes2)
	decRunes := converter10.TenToX(decimal)
	decimal2 := converter10.XToTen(decRunes)
	resultRunes := converter2.TenToX(decimal2)
	result := string(resultRunes)

	if original != result {
		t.Errorf("[Test7] Expected %s, got %s", original, result)
	} else {
		t.Logf("[Test7] Success: %s -> %s -> %s", original, string(decRunes), result)
	}
}

func TestOctalToDecimalRoundTrip(t *testing.T) {
	original := "777"

	converter8 := radix.NewRadixByBit(8)
	converter10 := radix.NewRadixByBit(10)

	runes8 := []rune(original)
	decimal := converter8.XToTen(runes8)
	decRunes := converter10.TenToX(decimal)
	decimal2 := converter10.XToTen(decRunes)
	resultRunes := converter8.TenToX(decimal2)
	result := string(resultRunes)

	if original != result {
		t.Errorf("[Test8] Expected %s, got %s", original, result)
	} else {
		t.Logf("[Test8] Success: %s -> %s -> %s", original, string(decRunes), result)
	}
}

func TestBytesToRunesToBytes(t *testing.T) {
	original := []byte{0, 1, 2, 127, 128, 255, 72, 101, 108, 108, 111}

	runes := utils.BytesToRunes(original)
	restored, err := utils.RunesToBytes(runes)

	if err != nil {
		t.Errorf("[UtilsTest] RunesToBytes error: %v", err)
		return
	}

	if !bytes.Equal(original, restored) {
		t.Errorf("[UtilsTest] Bytes differ")
		return
	}

	t.Logf("[UtilsTest] Success: %d bytes round-trip", len(original))
}

func TestBytesToRunesToBytesWithOutOfRange(t *testing.T) {
	runes := []rune{0, 127, 255, 256, 300}

	_, err := utils.RunesToBytes(runes)

	if err == nil {
		t.Errorf("[UtilsTest] Expected error for out-of-range runes, got nil")
		return
	}

	t.Logf("[UtilsTest] Success: correctly detected out-of-range rune")
}

func cleanupTestDir() {
	if testDir != "" {
		os.RemoveAll(testDir)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	cleanupTestDir()
	os.Exit(code)
}
