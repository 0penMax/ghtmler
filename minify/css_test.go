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
		// -- new test cases below --
		{
			input:    ``, // empty input
			expected: ``,
		},
		{
			input:    `/* only comment should be removed */`,
			expected: ``,
		},
		{
			input: `
				ul li, ol li {
					margin: 0px 0px 5px 10px;  /* mixed zero values */
					padding: 0;
				}
				`,
			expected: `ul li,ol li{margin:0 0 5px 10px;padding:0}`,
		},
		{
			input: `
					.btn-primary {
						background: #FF0000;
						border: 1px solid #000;
						opacity: .5;
					}
					`,
			expected: `.btn-primary{background:red;border:1px solid #000;opacity:.5}`,
		},
		{
			input: `
					@media screen and (max-width: 600px) {
					  .responsive {
						display: block;
						/* comment inside media */
						width: 100%;
					  }
					}
					`,
			expected: `@media screen and (max-width:600px){.responsive{display:block;width:100%}}`,
		},
		{
			input: `
					div::before, div::after {
						content: "";
						display: table;
						clear: both;
					}
					`,
			expected: `div::before,div::after{content:"";display:table;clear:both}`,
		},
		{
			input: `
					a:hover {
						color: #00ff00;
						text-decoration:    underline  ;
					}
					`,
			expected: `a:hover{color:#0f0;text-decoration:underline}`,
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
