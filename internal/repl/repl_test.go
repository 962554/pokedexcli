package repl

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "pre-lc string",
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "empty",
			input:    " ",
			expected: []string{},
		},
		{
			name:     "single-word",
			input:    " hello ",
			expected: []string{"hello"},
		},
		{
			name:     "capatalised words",
			input:    " Hello World ",
			expected: []string{"hello", "world"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := cleanInput(tc.input)
			if len(got) != len(tc.expected) {
				t.Errorf("slice lengths don't match: expected: %d, got: %d", len(tc.expected), len(got))
			}
			for i := range got {
				word := got[i]
				expectedWord := tc.expected[i]
				if expectedWord != word {
					t.Errorf("input: %q, expected: %s, got: %s", tc.input, expectedWord, word)
				}
			}
		})
	}
}
