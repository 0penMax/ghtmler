package optimizer

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// TODO complete the test to check the correct operation on the bootstrap
func TestRemoveUnusedSelectors2(t *testing.T) {

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

		{
			name: "test9_single_class_selector",
			args: args{
				usedSelectors: []Selector{{
					Value: ".display-1",
					SType: selectorClass,
				}},
			},
			want: `@charset "UTF-8";
				.display-1 {
				  font-size: calc(1.625rem + 4.5vw);
				  font-weight: 300;
				  line-height: 1.2;
				}
				@media (min-width: 1200px) {
				  .display-1 {
				    font-size: 5rem;
				  }
				}`,
		},

		//TODO add to cleaner option to correct work with special character like >, <, ::, * and etc
		// {
		// 	name: "test11_utility_classes",
		// 	args: args{
		// 		usedSelectors: []Selector{{
		// 			Value: ".row",
		// 			SType: selectorClass,
		// 		}},
		// 	},
		// 	want: `@charset "UTF-8";
		// 		.row {
		// 		  --bs-gutter-x: 1.5rem;
		// 		  --bs-gutter-y: 0;
		// 		  display: flex;
		// 		  flex-wrap: wrap;
		// 		  margin-top: calc(-1 * var(--bs-gutter-y));
		// 		  margin-right: calc(-0.5 * var(--bs-gutter-x));
		// 		  margin-left: calc(-0.5 * var(--bs-gutter-x));
		// 		}
		// 		.row > * {
		// 		  flex-shrink: 0;
		// 		  width: 100%;
		// 		  max-width: 100%;
		// 		  padding-right: calc(var(--bs-gutter-x) * 0.5);
		// 		  padding-left: calc(var(--bs-gutter-x) * 0.5);
		// 		  margin-top: var(--bs-gutter-y);
		// 		}`,
		// },
		{
			name: "test12_multiple_selectors_mixed",
			args: args{
				usedSelectors: []Selector{{
					Value: "h3",
					SType: selectorTag,
				}, {
					Value: "blockquote",
					SType: selectorTag,
				}},
			},
			want: ` @charset "UTF-8"; h3 { margin-top: 0; margin-bottom: 0.5rem; font-weight: 500; line-height: 1.2; color: var(--bs-heading-color); } h3 { font-size: calc(1.3rem + 0.6vw); } @media (min-width: 1200px) { h3 { font-size: 1.75rem; } } blockquote { margin: 0 0 1rem; }`,
		},
		{
			name: "test13_unused_selectors",
			args: args{
				usedSelectors: []Selector{{
					Value: ".nonexistent-class",
					SType: selectorClass,
				}},
			},
			want: `@charset "UTF-8";`, // No matching selectors, empty result aside from charset
		},
		{
			name: "test14_selector_with_media_query",
			args: args{
				usedSelectors: []Selector{{
					Value: ".col-sm",
					SType: selectorClass,
				}},
			},
			want: `@charset "UTF-8";
				@media (min-width: 576px) {
				  .col-sm {
				    flex: 1 0 0%;
				  }
				}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveUnusedSelectors(getParsedCss4Test(bootstrap5css), tt.args.usedSelectors); !reflect.DeepEqual(got, tt.want) {
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
