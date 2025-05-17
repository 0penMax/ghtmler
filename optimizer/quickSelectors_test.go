package optimizer

import (
	"goHtmlBuilder/css"
	"reflect"
	"testing"
)

func TestFilterSelectors(t *testing.T) {
	tests := []struct {
		name         string
		dataKeys     []string       // which tokens are “allowed”
		input        []css.Selector // selectors to filter
		wantFiltered []css.Selector // expected output
	}{
		{
			name:     "match one selector",
			dataKeys: []string{"foo"},
			input: []css.Selector{
				"foo bar",
				"baz",
			},
			wantFiltered: []css.Selector{
				"foo bar",
			},
		},
		{
			name:     "match multiple selectors",
			dataKeys: []string{"a", "b"},
			input: []css.Selector{
				"a",
				"x, b, y",
				"c",
			},
			wantFiltered: []css.Selector{
				"a",
				"x, b, y",
			},
		},
		{
			name:     "no selectors match",
			dataKeys: []string{"zz"},
			input: []css.Selector{
				"foo",
				"bar",
			},
			wantFiltered: nil,
		},
		{
			name:         "empty input slice",
			dataKeys:     []string{"anything"},
			input:        nil,
			wantFiltered: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// build quickSelectors with the given dataKeys
			qs := &quickSelectors{data: make(map[string]bool)}
			for _, k := range tc.dataKeys {
				qs.data[k] = true
			}

			got := qs.filterSelectors(tc.input)
			if !reflect.DeepEqual(got, tc.wantFiltered) {
				t.Errorf("filterSelectors() = %v, want %v", got, tc.wantFiltered)
			}
		})
	}
}
