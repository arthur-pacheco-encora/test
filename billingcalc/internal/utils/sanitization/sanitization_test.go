//go:build !selectTest || unitTest

package sanitization_test

import (
	"testing"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Only alphabets", "ArthurPacheco", "ArthurPacheco"},
		{"Alphabets with symbols", "Arthur$%Pacheco", "ArthurPacheco"},
		{"Only numbers", "123456", "123456"},
		{"Numbers with symbols", "123$%^456", "123456"},
		{"Mixed characters", "ArthurPacheco789!", "ArthurPacheco789"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitization.SanitizeName(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeName(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSanitizeFloatString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Only numbers", "123", "123"},
		{"Float number", "123.456", "123.456"},
		{"Numbers with symbols", "123.45^&*6", "123.456"},
		{"Only symbols", "$%^&", "0"},
		{"Empty string", "", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitization.SanitizeFloatString(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeFloatString(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
