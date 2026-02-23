package repl_test

import (
	"testing"

	"github.com/962554/pokedexcli/internal/repl"
)

func TestCleanInput(t *testing.T) {
	t.Parallel()

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

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got := repl.CleanInput(tcase.input)
			if len(got) != len(tcase.expected) {
				t.Errorf("slice lengths don't match: expected: %d, got: %d", len(tcase.expected), len(got))
			}

			for i := range got {
				word := got[i]

				expectedWord := tcase.expected[i]
				if expectedWord != word {
					t.Errorf("input: %q, expected: %s, got: %s", tcase.input, expectedWord, word)
				}
			}
		})
	}
}
