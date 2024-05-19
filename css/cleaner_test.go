package css

import (
	"reflect"
	"testing"
)

const testCss = `
h6, .h6, h5, .h5, h4, .h4, h3, .h3, h2, .h2, h1, .h1 {
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-weight: 500;
  line-height: 1.2;
}

h1, .h1 {
  font-size: calc(1.375rem + 1.5vw);
}

h2, .h2 {
  font-size: calc(1.325rem + 0.9vw);
}

h3, .h3 {
  font-size: calc(1.3rem + 0.6vw);
}


h4, .h4 {
  font-size: calc(1.275rem + 0.3vw);
}


h5, .h5 {
  font-size: 1.25rem;
}

h6, .h6 {
  font-size: 1rem;
}

p {
  margin-top: 0;
  margin-bottom: 1rem;
}

@media (min-width: 1200px) {
  h1, .h1 {
    font-size: 2.5rem;
  }
}


`

func TestRemoveUnusedSelectors(t *testing.T) {

	ParsedCss, err := Parse(testCss)
	if err != nil {
		t.Fatal("Failed to parse css", err, testCss)
	}

	type args struct {
		css           Stylesheet
		usedSelectors []string
	}
	tests := []struct {
		name string
		args args
		want Stylesheet
	}{
		{
			name: "test1",
			args: args{
				css:           *ParsedCss,
				usedSelectors: []string{"h1"},
			},
			want: Stylesheet{
				Rules: []*Rule{
					{
						Prelude:   "h6, .h6, h5, .h5, h4, .h4, h3, .h3, h2, .h2, h1, .h1",
						Selectors: []string{"h6", ".h6", "h5", ".h5", "h4", ".h4", "h3", ".h3", "h2", ".h2", "h1", ".h1"},
						Declarations: []*Declaration{
							{
								Property: "margin-top",
								Value:    "0",
							}, {
								Property: "margin-bottom",
								Value:    "0.5rem",
							}, {
								Property: "font-weight",
								Value:    "500",
							}, {
								Property: "line-height",
								Value:    "1.2",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},

					{
						Prelude:   "h1, .h1",
						Selectors: []string{"h1", ".h1"},
						Declarations: []*Declaration{
							{
								Property: "font-size",
								Value:    "calc(1.375rem + 1.5vw)",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},

					{
						Kind:    AtRule,
						Name:    "@media",
						Prelude: "(min-width: 1200px)",

						Rules: []*Rule{
							{
								Kind:      QualifiedRule,
								Prelude:   "h1, .h1",
								Selectors: []string{"h1", ".h1"},
								Declarations: []*Declaration{
									{
										Property: "font-size",
										Value:    "2.5rem",
									},
								},
								EmbedLevel: 1,
							},
						},
						EmbedLevel: 0,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := RemoveUnusedSelectors(tt.args.css, tt.args.usedSelectors); !reflect.DeepEqual(got, tt.want) {
				for i, g := range got.Rules {
					if !reflect.DeepEqual(g, tt.want.Rules[i]) {
						t.Errorf("not equal \n\n %+v\n want \n %+v\n", *g, *tt.want.Rules[i])
						t.Errorf("rule \n\n %+v\n want \n %+v\n", *g.Rules[0], *tt.want.Rules[i].Rules[0])
					}
				}
				t.Errorf("RemoveUnusedSelectors() = %v, want %v", got, tt.want)
			}
		})
	}
}
