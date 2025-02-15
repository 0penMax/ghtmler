package optimizer

import (
	"reflect"
	"testing"
)

func TestTokenizeCSSSelector(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    ".row>*",
			expected: []string{".row", ">", "*"},
		},
		{
			input:    ".row > *",
			expected: []string{".row", ">", "*"},
		},
		{
			input:    `div#main.content[data-type="example"]:hover > span::before`,
			expected: []string{"div", "#main", ".content[data-type=\"example\"]", ">", "span"},
		},
		{
			input:    "a.link:visited, a.link:hover",
			expected: []string{"a", ".link", ",", "a", ".link"},
		},
		{
			input:    "ul li:nth-child(2) > a",
			expected: []string{"ul", "li", ">", "a"},
		},
		{
			input:    "button:active",
			expected: []string{"button"},
		},
		{
			input:    "p::first-line",
			expected: []string{"p"},
		},
		{
			input:    "section.content:focus-within div.item::after",
			expected: []string{"section", ".content", "div", ".item"},
		},
	}

	for _, tc := range testCases {
		result := tokenizeCSSSelector(tc.input)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("For input %q,\nexpected tokens %v,\nbut got       %v", tc.input, tc.expected, result)
		}
	}
}
