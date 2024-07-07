package optimizer

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetAllClasses(t *testing.T) {
	type args struct {
		htmlCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{"<p class='hello'><span class='test test2'> world</span> </p>"},
			want:    []string{"hello", "test", "test2"},
			wantErr: false,
		}, {
			name:    "test2",
			args:    args{"<p class='hello'><span class='test test2'> world</span> </p> <p class='outside'>text</p>"},
			want:    []string{"hello", "test", "test2", "outside"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllClasses(tt.args.htmlCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllClasses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotClasses []string
			for _, s := range got {
				gotClasses = append(gotClasses, s.Value)
			}
			if !reflect.DeepEqual(gotClasses, tt.want) {
				t.Errorf("GetAllClasses() got = [%s], want [%s]", strings.Join(gotClasses, `','`), strings.Join(tt.want, `','`))
			}
		})
	}
}

func TestGetAllIds(t *testing.T) {
	type args struct {
		htmlCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{"<p id='hello'><span id='test'> world</span> </p>"},
			want:    []string{"hello", "test"},
			wantErr: false,
		}, {
			name:    "test2",
			args:    args{"<p id='hello'><span id='test'> world</span> </p> <p id='outside'>text</p>"},
			want:    []string{"hello", "test", "outside"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllIds(tt.args.htmlCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var gotIds []string
			for _, s := range got {
				gotIds = append(gotIds, s.Value)
			}
			if !reflect.DeepEqual(gotIds, tt.want) {
				t.Errorf("GetAllIds() got = %v, want %v", gotIds, tt.want)
			}
		})
	}
}
func TestGetAllSelectors(t *testing.T) {
	tests := []struct {
		name           string
		htmlCode       string
		expectedOutput []Selector
		expectError    bool
	}{
		{
			name:     "HTML with both class and id attributes",
			htmlCode: `<div class="test class1" id="unique"></div>`,
			expectedOutput: []Selector{
				{Value: "test", SType: selectorClass},
				{Value: "class1", SType: selectorClass},
				{Value: "unique", SType: selectorId},
			},
			expectError: false,
		},
		{
			name:     "HTML with only class attribute",
			htmlCode: `<span class="onlyclass"></span>`,
			expectedOutput: []Selector{
				{Value: "onlyclass", SType: selectorClass},
			},
			expectError: false,
		},
		{
			name:     "HTML with only id attribute",
			htmlCode: `<p id="onlyid"></p>`,
			expectedOutput: []Selector{
				{Value: "onlyid", SType: selectorId},
			},
			expectError: false,
		},
		{
			name:           "HTML with no class or id attributes",
			htmlCode:       `<div></div>`,
			expectedOutput: []Selector{},
			expectError:    false,
		},
		{
			name:     "HTML with multiple class and id attributes",
			htmlCode: `<div class="multi class" id="test"></div>`,
			expectedOutput: []Selector{
				{Value: "multi", SType: selectorClass},
				{Value: "class", SType: selectorClass},
				{Value: "test", SType: selectorId},
			},
			expectError: false,
		},
		{
			name:     "HTML with mixed class and id attributes",
			htmlCode: `<div class="mixed" id="test"></div>`,
			expectedOutput: []Selector{
				{Value: "mixed", SType: selectorClass},
				{Value: "test", SType: selectorId},
			},
			expectError: false,
		},
		{
			name:           "Complete HTML document with no class or id attributes",
			htmlCode:       `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
			expectedOutput: []Selector{},
			expectError:    false,
		},
		{
			name:     "Nested HTML with multiple class and id attributes",
			htmlCode: `<div class="class1 class2" id="id1"><span class="class3" id="id2"></span></div>`,
			expectedOutput: []Selector{
				{Value: "class1", SType: selectorClass},
				{Value: "class2", SType: selectorClass},
				{Value: "id1", SType: selectorId},
				{Value: "class3", SType: selectorClass},
				{Value: "id2", SType: selectorId},
			},
			expectError: false,
		},
		{
			name:           "Empty HTML input",
			htmlCode:       ``,
			expectedOutput: []Selector{},
			expectError:    false,
		},
		{
			name:           "HTML with spaces in class and id attributes",
			htmlCode:       `<div class="   " id="   "></div>`,
			expectedOutput: []Selector{},
			expectError:    false,
		},
		{
			name:     "HTML with extra spaces in class and id attributes",
			htmlCode: `<div class="class1    class2" id="id1   "></div>`,
			expectedOutput: []Selector{
				{Value: "class1", SType: selectorClass},
				{Value: "class2", SType: selectorClass},
				{Value: "id1", SType: selectorId},
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := GetAllSelectors(test.htmlCode)
			if test.expectError && err == nil {
				t.Errorf("Expected error but got none for input: %s", test.htmlCode)
			}
			if !test.expectError && err != nil {
				t.Errorf("Did not expect error but got: %v for input: %s", err, test.htmlCode)
			}
			if !test.expectError && !compareSelectors(output, test.expectedOutput) {
				t.Errorf("Expected output %v but got %v for input: %s", test.expectedOutput, output, test.htmlCode)
			}
		})
	}
}

func compareSelectors(a, b []Selector) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetCSSFileNamesFromHtml(t *testing.T) {
	tests := []struct {
		html       string
		expected   []string
		shouldFail bool
	}{
		{
			html: `
            <!DOCTYPE html>
            <html>
            <head>
                <link rel="stylesheet" href="styles.css">
                <link rel="stylesheet" href="theme.css">
            </head>
            <body>
                <h1>Hello, World!</h1>
            </body>
            </html>
            `,
			expected: []string{"styles.css", "theme.css"},
		},
		{
			html: `
            <!DOCTYPE html>
            <html>
            <head>
                <link rel="alternate" href="feed.xml">
            </head>
            <body>
                <h1>Hello, World!</h1>
            </body>
            </html>
            `,
			expected: []string{},
		},
		{
			html: `
            <!DOCTYPE html>
            <html>
            <head>
                <link rel="stylesheet" href="">
            </head>
            <body>
                <h1>Hello, World!</h1>
            </body>
            </html>
            `,
			expected: []string{},
		},
		{
			html:       "<html><head><link rel=\"stylesheet\" href=\"styles.css\"></head><body><h1>Hello</h1></body></html>",
			expected:   []string{"styles.css"},
			shouldFail: false,
		},
		{
			html:       "<html><head><link rel=\"stylesheet\" href=\"\"></head><body><h1>Hello</h1></body></html>",
			expected:   []string{},
			shouldFail: false,
		},
		{
			html:       "<html><head><link rel=\"stylesheet\"></head><body><h1>Hello</h1></body></html>",
			expected:   []string{},
			shouldFail: false,
		},
		{
			html:       "<html><head><link rel=\"alternate\" href=\"feed.xml\"></head><body><h1>Hello</h1></body></html>",
			expected:   []string{},
			shouldFail: false,
		},
		{
			html: `
            <!DOCTYPE html>
            <html>
            <head>
                <link rel="stylesheet" href="main.css">
                <link rel="stylesheet" href="theme.css">
                <link rel="icon" href="favicon.ico">
            </head>
            <body>
                <h1>Hello, World!</h1>
            </body>
            </html>
            `,
			expected: []string{"main.css", "theme.css"},
		},
		{
			html: `
            <!DOCTYPE html>
            <html>
            <head>
                <link rel="stylesheet" href="css/styles1.css">
                <link rel="stylesheet" href="css/styles2.css">
                <link rel="stylesheet" href="css/styles3.css">
            </head>
            <body>
                <h1>Hello, World!</h1>
            </body>
            </html>
            `,
			expected: []string{"styles1.css", "styles2.css", "styles3.css"},
		},
		{
			html:       "",
			expected:   []string{},
			shouldFail: false,
		},
	}

	for _, test := range tests {
		r := strings.NewReader(test.html)
		cssFiles, err := GetCSSFileNamesFromHtml(r)
		if test.shouldFail {
			if err == nil {
				t.Errorf("Expected an error for input: %s", test.html)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect an error for input: %s, got: %v", test.html, err)
			}
			if len(cssFiles) != len(test.expected) {
				t.Errorf("Expected %d CSS files, got %d for input: %s", len(test.expected), len(cssFiles), test.html)
			}
			for i, css := range test.expected {
				if cssFiles[i].fileName != css {
					t.Errorf("Expected CSS file %s, got %s for input: %s", css, cssFiles[i], test.html)
				}
			}
		}
	}
}

// TestGetFileName tests the GetFileName function with different full path inputs.
func TestGetFileName(t *testing.T) {
	tests := []struct {
		fullPath string
		expected string
	}{
		{
			fullPath: "/home/user/documents/report.pdf",
			expected: "report.pdf",
		},
		{
			fullPath: "C:\\Users\\user\\documents\\report.pdf",
			expected: "report.pdf",
		},
		{
			fullPath: "/var/log/syslog",
			expected: "syslog",
		},
		{
			fullPath: "report.pdf",
			expected: "report.pdf",
		},
		{
			fullPath: "/home/user/music/song.mp3",
			expected: "song.mp3",
		},
		{
			fullPath: "C:\\path\\to\\file\\image.png",
			expected: "image.png",
		},
		{
			fullPath: "/relative/path/to/file.txt",
			expected: "file.txt",
		},
		{
			fullPath: "filename.ext",
			expected: "filename.ext",
		},
		{
			fullPath: "",
			expected: "",
		},
		{
			fullPath: "/",
			expected: "",
		},
	}

	for i, test := range tests {
		result := getFileName(test.fullPath)
		if result != test.expected {
			t.Errorf("Test case %d: GetFileName(%q) = %q; want %q", i+1, test.fullPath, result, test.expected)
		}
	}
}
