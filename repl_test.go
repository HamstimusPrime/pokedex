package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "HELLO wOrlD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  CAULdron wOrlD BuBBLe",
			expected: []string{"cauldron", "world", "bubble"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("slice length do not match. Expected %v, but got %v", len(c.expected), len(actual))
			return
		}
		for i := range actual {
			word := actual[i]
			if word != c.expected[i] {
				t.Errorf("words do not match. Expected %s, but got %s", c.expected[i], word)
				return
			}
		}
	}

}

