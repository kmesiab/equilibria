package encoding

import (
	"testing"
)

// TestIsStandardGSMChar tests if standard GSM characters are correctly identified.
func TestIsStandardGSMChar(t *testing.T) {
	// Include a sample of standard GSM characters
	standardChars := "Aa@£$ΓΩΣ0123\n\r"
	for _, char := range standardChars {
		if !isStandardGSMChar(char) {
			t.Errorf("Expected true for standard GSM character %q, got false", char)
		}
	}

	// Include a character not in the standard GSM set
	if isStandardGSMChar('€') {
		t.Error("Expected false for non-standard GSM character '€', got true")
	}
}

// TestIsExtendedGSMChar tests if extended GSM characters are correctly identified.
func TestIsExtendedGSMChar(t *testing.T) {
	// Include a sample of extended GSM characters
	extendedChars := "{}\\[~]|€^"
	for _, char := range extendedChars {
		if !isExtendedGSMChar(char) {
			t.Errorf("Expected true for extended GSM character %q, got false", char)
		}
	}

	// Include a character that is not in the extended GSM set
	if isExtendedGSMChar('@') {
		t.Error("Expected false for non-extended GSM character '@', got true")
	}
}

// TestIsGSMEncoded tests if strings are correctly identified as GSM encoded or not.
func TestIsGSMEncoded(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{"Hello, World!", true},
		{"Standard GSM: Aa@£$ΓΩΣ0123", true},
		{"Extended GSM: {}\\[~]|€^", true},
		{"Mixed: ΩΣ€@A", true},
		{"Non-GSM: 😊", false},
		{"Empty string: ", true},
		{"Chinese characters: 汉字", false},
		{"Plain English text", true},
		{"Text with newlines\nand carriage returns\r", true},
		{"1234567890", true},
		{"Symbols !\"#$%&'()*+,-./:;<=>?@[\\]^_{|}~", true},
		{"`", false},
		{"Extended €{}[~]|^\\ but still GSM", true},
		{"Contains emoji 😊", false},
		{"汉字 - Chinese characters", false},
		{"Mixed content 😊 with GSM", false},
		{"Just € extended symbols", true},
		{"Empty string", true},
		{"Single char: @", true},
		{"Single extended char: €", true},
		{"Single non-GSM char: 😊", false},
	}

	for _, c := range cases {
		got := IsGSMEncoded(c.text)
		if got != c.want {
			t.Errorf("IsGSMEncoded(%q) == %t, want %t", c.text, got, c.want)
		}
	}
}

// TestIsGSMEncodedEdgeCases tests edge cases for GSM encoding identification.
func TestIsGSMEncodedEdgeCases(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{"", true},         // Empty string case
		{"€€€€€", true},    // String of only extended characters
		{"\n\r\n\r", true}, // String of only control characters
		{"@@@@@", true},    // String of repeated standard characters
		{"😊😊😊", false},     // String of only emojis
	}

	for _, c := range cases {
		got := IsGSMEncoded(c.text)
		if got != c.want {
			t.Errorf("IsGSMEncodedEdgeCases(%q) == %t, want %t", c.text, got, c.want)
		}
	}
}
