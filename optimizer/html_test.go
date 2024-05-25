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
