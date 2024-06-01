package minify

import "testing"

func TestMinifyCSS(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `
body {
	background-color: white; /* This is a comment */
	color: black;
}
.container {
	width: 100%;
	padding: 10px;
	margin: 0 auto;
}
`,
			expected: `body{background-color:#fff;color:#000}.container{width:100%;padding:10px;margin:0 auto}`,
		},
		{
			input: `
/* Another comment */
p {
  font-size: 16px;
  line-height: 1.5;
}
`,
			expected: `p{font-size:16px;line-height:1.5}`,
		},
		{
			input: `
h1 {
	font-weight: bold;
}
a {
	text-decoration: none;
	color: blue;
}
`,
			expected: `h1{font-weight:700}a{text-decoration:none;color:blue}`,
		},
	}

	for _, test := range tests {
		result, err := MinifyCSS(test.input)
		if err != nil {
			t.Errorf("MinifyCSS(%q) returned error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("MinifyCSS(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
