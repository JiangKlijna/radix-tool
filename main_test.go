package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Helper function to create temporary files
func createTestFile(content string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func TestBinaryToDecimal(t *testing.T) {
	cmd := exec.Command("./radix-tool.exe", "-i", "1010", "-ib", "2", "-ob", "10")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}
	result := strings.TrimSpace(string(output))

	if result != "10" {
		t.Errorf("Expected 10, got %s", result)
	}
}

func TestDecimalToBinary(t *testing.T) {
	cmd := exec.Command("./radix-tool.exe", "--input", "255", "--input-base-num", "10", "--output-base-num", "2")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}
	result := strings.TrimSpace(string(output))

	if result != "11111111" {
		t.Errorf("Expected 11111111, got %s", result)
	}
}

func TestDecimalToHex(t *testing.T) {
	cmd := exec.Command("./radix-tool.exe", "-i", "255", "-ib", "10", "-ob", "16")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}
	result := strings.TrimSpace(string(output))

	if result != "ff" {
		t.Errorf("Expected ff, got %s", result)
	}
}

func TestHexToDecimal(t *testing.T) {
	cmd := exec.Command("./radix-tool.exe", "--input", "ff", "--input-base-num", "16", "--output-base-num", "10")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}
	result := strings.TrimSpace(string(output))

	if result != "255" {
		t.Errorf("Expected 255, got %s", result)
	}
}

func TestFileInput(t *testing.T) {
	fileName, err := createTestFile("ff")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName)

	cmd := exec.Command("./radix-tool.exe", "-i", fileName, "-ib", "16", "-ob", "10")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}
	result := strings.TrimSpace(string(output))

	if result != "255" {
		t.Errorf("Expected 255 from file input, got %s", result)
	}
}

func TestFileOutput(t *testing.T) {
	// Test writing to a file
	outFileName := "test_output.txt"
	defer os.Remove(outFileName)

	cmd := exec.Command("./radix-tool.exe", "-i", "ff", "-ib", "16", "-ob", "10", "-o", outFileName)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	content, err := ioutil.ReadFile(outFileName)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	result := strings.TrimSpace(string(content))

	if result != "255" {
		t.Errorf("Expected 255 in output file, got %s", result)
	}
}
