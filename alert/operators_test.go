package alert

import (
	"testing"
)

func TestGreaterInt(t *testing.T) {
	tests := []struct {
		name     string
		first    int
		second   int
		expected bool
	}{
		{"First greater", 5, 3, true},
		{"Second greater", 3, 5, false},
		{"Equal", 4, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GreaterInt(tt.first, tt.second); got != tt.expected {
				t.Errorf("GreaterInt(%d, %d) = %v, want %v", tt.first, tt.second, got, tt.expected)
			}
		})
	}
}

func TestLessInt(t *testing.T) {
	tests := []struct {
		name     string
		first    int
		second   int
		expected bool
	}{
		{"First less", 3, 5, true},
		{"Second less", 5, 3, false},
		{"Equal", 4, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LessInt(tt.first, tt.second); got != tt.expected {
				t.Errorf("LessInt(%d, %d) = %v, want %v", tt.first, tt.second, got, tt.expected)
			}
		})
	}
}

func TestEqualsInt(t *testing.T) {
	tests := []struct {
		name     string
		first    int
		second   int
		expected bool
	}{
		{"Equal", 4, 4, true},
		{"Not equal (first greater)", 5, 3, false},
		{"Not equal (second greater)", 3, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualsInt(tt.first, tt.second); got != tt.expected {
				t.Errorf("EqualsInt(%d, %d) = %v, want %v", tt.first, tt.second, got, tt.expected)
			}
		})
	}
}
