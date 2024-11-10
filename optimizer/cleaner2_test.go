package optimizer

import (
	"goHtmlBuilder/css"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// TODO complete the test to check the correct operation on the bootstrap
func TestRemoveUnusedSelectors2(t *testing.T) {

	ParsedCss, err := css.Parse(bootstrap5css)
	if err != nil {
		t.Fatal("Failed to parse css", err, testCss)
	}

	type args struct {
		usedSelectors []Selector
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				usedSelectors: []Selector{{
					Value: "h1",
					SType: selectorTag,
				}, {
					Value: "h2",
					SType: selectorTag,
				}},
			},
			want: `@charset "UTF-8";
					h2, h1 {
					  margin-top: 0;
					  margin-bottom: 0.5rem;
					  font-weight: 500;
					  line-height: 1.2;
					  color: var(--bs-heading-color);
					}
					h1 {
					  font-size: calc(1.375rem + 1.5vw);
					}
					@media (min-width: 1200px) {
					  h1 {
						font-size: 2.5rem;
					  }
					}
					h2 {
					  font-size: calc(1.325rem + 0.9vw);
					}
					@media (min-width: 1200px) {
					  h2 {
						font-size: 2rem;
					  }
					}
				`,
		}, {
			name: "test2",
			args: args{
				usedSelectors: []Selector{{
					Value: "p",
					SType: selectorTag,
				}, {
					Value: "a",
					SType: selectorTag,
				}},
			},
			want: `@charset "UTF-8"; p { margin-top: 0; margin-bottom: 1rem; } a { color: rgba(var(--bs-link-color-rgb), var(--bs-link-opacity, 1)); text-decoration: underline; }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveUnusedSelectors(*ParsedCss, tt.args.usedSelectors); !reflect.DeepEqual(got, tt.want) {
				clearGot := cleanString(got.String())
				clearWant := cleanString(tt.want)
				if clearGot != clearWant {
					t.Errorf("RemoveUnusedSelectors():\n %v, \n\n want:\n %v", clearGot, clearWant)
				}
			}
		})
	}

}

func cleanString(s string) string {
	// Replace all tabs with a single space
	s = strings.ReplaceAll(s, "\t", " ")

	// Use regex to replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	// Trim any leading or trailing spaces
	s = strings.TrimSpace(s)

	return s
}
