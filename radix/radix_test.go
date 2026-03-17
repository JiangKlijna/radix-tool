package radix

import (
	"math/big"
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
