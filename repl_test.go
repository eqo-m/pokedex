package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "A B C",
			expected: []string{"A", "B", "C"},
		},
		{
			input:    "Hell8no",
			expected: []string{"Hell8no"},
		},
		{
			input:    "Hello World",
			expected: []string{"Hello", "World"},
		}}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected length of %d, got %d", len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if expectedWord != word {
				t.Errorf("%s expected to be %s", word, expectedWord)
			}
		}
	}

}
