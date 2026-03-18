package radix

import (
	"math/big"
	"strings"
	"testing"
)

func TestTenStrToX(t *testing.T) {
	// Test decimal to radix conversion
	radixConverter := NewRadixByBit(16)
	result := radixConverter.TenStrToX("255")
	expected := "ff"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTenToX(t *testing.T) {
	// Test big.Int to hex conversion
	radixConverter := NewRadixByBit(16)
	bigNum := big.NewInt(255)
	result := radixConverter.TenToX(bigNum)
	expected := "ff"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestXToTenStr(t *testing.T) {
	// Test hex to decimal conversion
	radixConverter := NewRadixByBit(16)
	result := radixConverter.XToTenStr("ff")
	expected := "255"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestBinaryConversion(t *testing.T) {
	// Test binary to decimal
	binaryConverter := NewRadixByBit(2)
	result := binaryConverter.XToTenStr("1010")
	expected := "10"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test decimal to binary
	result = binaryConverter.TenStrToX("10")
	expected = "1010"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestBase62Conversion(t *testing.T) {
	// Test base 62 conversion
	base62Converter := NewRadixByBit(62)
	testCases := []struct {
		decimal string
		base62  string
	}{
		{"0", "0"},
		{"1", "1"},
		{"9", "9"},
		{"10", "a"},
		{"35", "z"},
		{"36", "A"},
		{"61", "Z"},
		{"62", "10"},
		{"123", "1Z"},
		{"12345", "3d7"},
	}

	for _, tc := range testCases {
		// Test decimal to base62
		result := base62Converter.TenStrToX(tc.decimal)
		if result != tc.base62 {
			t.Errorf("Converting decimal %s to base62, expected '%s', got '%s'", tc.decimal, tc.base62, result)
		}

		// Test base62 to decimal
		result = base62Converter.XToTenStr(tc.base62)
		if result != tc.decimal {
			t.Errorf("Converting base62 '%s' to decimal, expected '%s', got '%s'", tc.base62, tc.decimal, result)
		}
	}
}

func TestCustomAlphabet(t *testing.T) {
	// Test custom alphabet conversion
	customAlphabet := "0123456789ABCDEF"
	radixConverter := NewRadixByString(customAlphabet)

	// Test that it behaves like standard hexadecimal
	result := radixConverter.XToTenStr("FF") // FF in our alphabet is 15*16 + 15 = 255
	expected := "255"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test reverse operation
	result = radixConverter.TenStrToX("255")
	expected = "FF"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestNewCharacterRadix(t *testing.T) {
	// Test characterRadix implementation
	characterRadix := NewCharacterRadix(10000) // Big value from original source code
	number := big.NewInt(1234567)

	converted := characterRadix.TenToX(number)
	backToNum := characterRadix.XToTen(converted)

	if backToNum.Cmp(number) != 0 {
		t.Errorf("Round trip conversion failed. Original: %s, After roundtrip: %s",
			number.String(), backToNum.String())
	}
}

// TestTenToXParallel tests the parallel version using existing functionality to guarantee correctness
func TestTenToXParallel(t *testing.T) {
	testCases := []int{2, 16, 36, 62}

	for _, base := range testCases {
		converter := NewRadixByBit(base)

		// Test for small values (won't trigger parallel path, just correctness check)
		smallValue := big.NewInt(255)
		sequentialResult := converter.TenToX(smallValue)
		parallelResult := converter.TenToXParallel(smallValue)
		if sequentialResult != parallelResult {
			t.Errorf("Base %d: Sequential (%s) and parallel (%s) results differ for 255",
				base, sequentialResult, parallelResult)
		}

		// For parallel testing: need BitLen() >= 1024 to trigger parallel code path
		// Use 2^1500 which has 1501 bits > 1024 threshold
		largeValue := new(big.Int).Exp(big.NewInt(2), big.NewInt(1500), nil)
		sequentialResult = converter.TenToX(largeValue)
		parallelResult = converter.TenToXParallel(largeValue)
		if sequentialResult != parallelResult {
			t.Errorf("Base %d: Sequential and parallel results differ for large input (2^1500)", base)
		}

		// Extra large value to ensure multiple levels of parallel recursion
		extraLargeValue := new(big.Int).Exp(big.NewInt(2), big.NewInt(5000), nil)
		sequentialResult = converter.TenToX(extraLargeValue)
		parallelResult = converter.TenToXParallel(extraLargeValue)
		if sequentialResult != parallelResult {
			t.Errorf("Base %d: Sequential and parallel results differ for extra large input (2^5000)", base)
		}
	}
}

// TestXToTenParallel tests the parallel version using existing functionality to guarantee correctness
func TestXToTenParallel(t *testing.T) {
	testCases := []struct {
		base int
		str  string
	}{
		{2, "11111111"},
		{16, "ff"},
		{36, "zz"},
		{62, "Zz"},
	}

	for _, tc := range testCases {
		converter := NewRadixByBit(tc.base)

		// Test simple cases (short strings won't trigger parallel path, just correctness check)
		testStr := tc.str
		sequentialResult := converter.XToTen(testStr).String()
		parallelResult := converter.XToTenParallel(testStr).String()
		if sequentialResult != parallelResult {
			t.Errorf("Base %d: Sequential (%s) and parallel (%s) results differ for '%s'",
				tc.base, sequentialResult, parallelResult, testStr)
		}

		// For testing with longer strings: need len >= 256 to trigger parallel code path
		// Use 1000 chars > 256 threshold
		longStr := strings.Repeat(string(converter.Radixer.GetRuneByInt(int64(tc.base-1))), 1000)
		sequentialResult = converter.XToTen(longStr).String()
		parallelResult = converter.XToTenParallel(longStr).String()
		if sequentialResult != parallelResult {
			t.Errorf("Base %d: Sequential and parallel results differ for long input (1000 chars)", tc.base)
		}

		// Extra long string to ensure multiple levels of parallel recursion
		extraLongStr := strings.Repeat(string(converter.Radixer.GetRuneByInt(int64(tc.base-1))), 5000)
		sequentialResult = converter.XToTen(extraLongStr).String()
		parallelResult = converter.XToTenParallel(extraLongStr).String()
		if sequentialResult != parallelResult {
			t.Errorf("Base %d: Sequential and parallel results differ for extra long input (5000 chars)", tc.base)
		}
	}
}

// TestRoundTripConsistency verifies round trip conversions work correctly
func TestRoundTripConsistency(t *testing.T) {
	testCases := []struct {
		base  int
		value string
	}{
		{2, "1010"},
		{8, "777"},
		{10, "12345"},
		{16, "abcdef"},
		{36, "zyzabc"},
		{62, "AZaz09"},
	}

	for _, tc := range testCases {
		converter := NewRadixByBit(tc.base)

		// From base-X to base-10 to base-X
		numInDec := converter.XToTenStr(tc.value)
		newValue := converter.TenStrToX(numInDec)
		if newValue != tc.value {
			t.Errorf("Base %d: Round trip failed: %s -> %s -> %s", tc.base, tc.value, numInDec, newValue)
		}

		// Test with parallel methods too
		numInDecParallel := converter.XToTenParallel(tc.value).String()
		if numInDec != numInDecParallel {
			t.Errorf("Base %d: Sequential (%s) and parallel (%s) results differ for %s",
				tc.base, numInDec, numInDecParallel, tc.value)
		}

		// Test parallel-to-parallel round trip using big.Int
		bigVal := new(big.Int)
		bigVal.SetString(numInDec, 10)
		valueFromParallel := converter.TenToXParallel(bigVal)
		if valueFromParallel != tc.value {
			t.Errorf("Base %d: Parallel round trip failed: %s -> %s -> %s", tc.base, tc.value, numInDec, valueFromParallel)
		}
	}
}

// TestLargeValues tests with large numbers that will definitely exercise parallel algorithms
func TestLargeValues(t *testing.T) {
	converter := NewRadixByBit(16)

	// Large number that will definitely trigger the parallel algorithm (> 1024 bits)
	hugeValue := new(big.Int)
	hugeValue.Exp(big.NewInt(2), big.NewInt(2048), nil) // 2^2048 has 2049 bits

	seqResult := converter.TenToX(hugeValue)
	parResult := converter.TenToXParallel(hugeValue)

	if seqResult != parResult {
		t.Errorf("Sequential (%s) and parallel (%s) results differ for large value", seqResult, parResult)
	}

	if len(seqResult) < 100 { // Result should be substantial since it's 2^2048
		t.Errorf("Unexpectedly short result for 2^2048: %s", seqResult)
	}
}

// Benchmark-style test for parallel methods
func TestParallelMethodsConsistencyWithLongInputs(t *testing.T) {
	bases := []int{2, 8, 10, 16, 36, 62}

	for _, base := range bases {
		converter := NewRadixByBit(base)

		// Create a large string/input that will engage parallel code paths

		// Large numeric value - should trigger parallel processing in TenToXParallel
		hugeNum := new(big.Int)
		hugeNum.SetString(strings.Repeat("9", 1000), 10) // 1000-digit number
		sequentialTenToX := converter.TenToX(hugeNum)
		parallelTenToX := converter.TenToXParallel(hugeNum)
		if sequentialTenToX != parallelTenToX {
			t.Errorf("Base %d: TenToX vs TenToXParallel differ for large input", base)
		}

		// Large base-X string - should trigger parallel processing in XToTenParallel
		largestdigit := converter.Radixer.GetRuneByInt(int64(base - 1))
		longXString := strings.Repeat(string(largestdigit), 1000) // Very long base-X string
		sequentialXToTen := converter.XToTen(longXString).String()
		parallelXToTen := converter.XToTenParallel(longXString).String()
		if sequentialXToTen != parallelXToTen {
			t.Errorf("Base %d: XToTen vs XToTenParallel differ for long input", base)
		}
	}
}
