package optimizer

import (
	"goHtmlBuilder/css"
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

	ParsedCss, err := css.Parse(testCss)
	if err != nil {
		t.Fatal("Failed to parse css", err, testCss)
	}

	type args struct {
		css           css.Stylesheet
		usedSelectors []Selector
	}
	tests := []struct {
		name string
		args args
		want css.Stylesheet
	}{
		{
			name: "test1",
			args: args{
				css: *ParsedCss,
				usedSelectors: []Selector{{
					Value: "h1",
					SType: selectorTag,
				}},
			},
			want: css.Stylesheet{
				Rules: []*css.Rule{
					{
						Prelude:   "h6, .h6, h5, .h5, h4, .h4, h3, .h3, h2, .h2, h1, .h1",
						Selectors: []string{"h6", ".h6", "h5", ".h5", "h4", ".h4", "h3", ".h3", "h2", ".h2", "h1", ".h1"},
						Declarations: []*css.Declaration{
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
						Declarations: []*css.Declaration{
							{
								Property: "font-size",
								Value:    "calc(1.375rem + 1.5vw)",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},
					{
						Kind:    css.AtRule,
						Name:    "@media",
						Prelude: "(min-width: 1200px)",
						Rules: []*css.Rule{
							{
								Kind:      css.QualifiedRule,
								Prelude:   "h1, .h1",
								Selectors: []string{"h1", ".h1"},
								Declarations: []*css.Declaration{
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
		{
			name: "Single Class Selector",
			args: args{
				css: *ParsedCss,
				usedSelectors: []Selector{{
					Value: ".h1",
					SType: selectorClass,
				}},
			},
			want: css.Stylesheet{
				Rules: []*css.Rule{
					{
						Prelude:   "h6, .h6, h5, .h5, h4, .h4, h3, .h3, h2, .h2, h1, .h1",
						Selectors: []string{"h6", ".h6", "h5", ".h5", "h4", ".h4", "h3", ".h3", "h2", ".h2", "h1", ".h1"},
						Declarations: []*css.Declaration{
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
						Declarations: []*css.Declaration{
							{
								Property: "font-size",
								Value:    "calc(1.375rem + 1.5vw)",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},
					{
						Kind:    css.AtRule,
						Name:    "@media",
						Prelude: "(min-width: 1200px)",
						Rules: []*css.Rule{
							{
								Kind:      css.QualifiedRule,
								Prelude:   "h1, .h1",
								Selectors: []string{"h1", ".h1"},
								Declarations: []*css.Declaration{
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
		{
			name: "Multiple Tag Selectors",
			args: args{
				css: *ParsedCss,
				usedSelectors: []Selector{
					{
						Value: "h1",
						SType: selectorTag,
					},
					{
						Value: "h2",
						SType: selectorTag,
					},
					{
						Value: "h3",
						SType: selectorTag,
					},
				},
			},
			want: css.Stylesheet{
				Rules: []*css.Rule{
					{
						Prelude:   "h6, .h6, h5, .h5, h4, .h4, h3, .h3, h2, .h2, h1, .h1",
						Selectors: []string{"h6", ".h6", "h5", ".h5", "h4", ".h4", "h3", ".h3", "h2", ".h2", "h1", ".h1"},
						Declarations: []*css.Declaration{
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
						Declarations: []*css.Declaration{
							{
								Property: "font-size",
								Value:    "calc(1.375rem + 1.5vw)",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},
					{
						Prelude:   "h2, .h2",
						Selectors: []string{"h2", ".h2"},
						Declarations: []*css.Declaration{
							{
								Property: "font-size",
								Value:    "calc(1.325rem + 0.9vw)",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},
					{
						Prelude:   "h3, .h3",
						Selectors: []string{"h3", ".h3"},
						Declarations: []*css.Declaration{
							{
								Property: "font-size",
								Value:    "calc(1.3rem + 0.6vw)",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},
					{
						Kind:    css.AtRule,
						Name:    "@media",
						Prelude: "(min-width: 1200px)",
						Rules: []*css.Rule{
							{
								Kind:      css.QualifiedRule,
								Prelude:   "h1, .h1",
								Selectors: []string{"h1", ".h1"},
								Declarations: []*css.Declaration{
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
		{
			name: "Unused Selectors",
			args: args{
				css: *ParsedCss,
				usedSelectors: []Selector{{
					Value: "h1",
					SType: selectorTag,
				}},
			},
			want: css.Stylesheet{
				Rules: []*css.Rule{
					{
						Prelude:   "h6, .h6, h5, .h5, h4, .h4, h3, .h3, h2, .h2, h1, .h1",
						Selectors: []string{"h6", ".h6", "h5", ".h5", "h4", ".h4", "h3", ".h3", "h2", ".h2", "h1", ".h1"},
						Declarations: []*css.Declaration{
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
						Declarations: []*css.Declaration{
							{
								Property: "font-size",
								Value:    "calc(1.375rem + 1.5vw)",
							},
						},
						Rules:      nil,
						EmbedLevel: 0,
					},
					{
						Kind:    css.AtRule,
						Name:    "@media",
						Prelude: "(min-width: 1200px)",
						Rules: []*css.Rule{
							{
								Kind:      css.QualifiedRule,
								Prelude:   "h1, .h1",
								Selectors: []string{"h1", ".h1"},
								Declarations: []*css.Declaration{
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
		{
			name: "No Used Selectors",
			args: args{
				css:           *ParsedCss,
				usedSelectors: []Selector{},
			},
			want: css.Stylesheet{
				Rules: []*css.Rule{{
					Kind:    css.AtRule,
					Name:    "@media",
					Prelude: "(min-width: 1200px)",
					Rules: []*css.Rule{
						{
							Kind:      css.QualifiedRule,
							Prelude:   "h1, .h1",
							Selectors: []string{"h1", ".h1"},
							Declarations: []*css.Declaration{
								{
									Property: "font-size",
									Value:    "2.5rem",
								},
							},
							EmbedLevel: 1,
						},
					},
					EmbedLevel: 0,
				}},
			},
		},
		{
			name: "All Selectors Used",
			args: args{
				css: *ParsedCss,
				usedSelectors: []Selector{
					{
						Value: "h1",
						SType: selectorTag,
					},
					{
						Value: "h2",
						SType: selectorTag,
					},
					{
						Value: "h3",
						SType: selectorTag,
					},
					{
						Value: "h4",
						SType: selectorTag,
					},
					{
						Value: "h5",
						SType: selectorTag,
					},
					{
						Value: "h6",
						SType: selectorTag,
					},
					{
						Value: ".h1",
						SType: selectorClass,
					},
					{
						Value: ".h2",
						SType: selectorClass,
					},
					{
						Value: ".h3",
						SType: selectorClass,
					},
					{
						Value: ".h4",
						SType: selectorClass,
					},
					{
						Value: ".h5",
						SType: selectorClass,
					},
					{
						Value: ".h6",
						SType: selectorClass,
					},
					{
						Value: "p",
						SType: selectorTag,
					},
				},
			},
			want: *ParsedCss,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveUnusedSelectors(tt.args.css, tt.args.usedSelectors); !reflect.DeepEqual(got, tt.want) {
				for i, g := range got.Rules {
					if i >= len(tt.want.Rules) {
						t.Errorf("unexpected rule: %+v", *g)
						continue
					}
					if !reflect.DeepEqual(g, tt.want.Rules[i]) {
						t.Errorf("not equal \n\n %+v\n want \n %+v\n", *g, *tt.want.Rules[i])
						if len(g.Rules) > 0 && len(tt.want.Rules[i].Rules) > 0 {
							t.Errorf("rule \n\n %+v\n want \n %+v\n", *g.Rules[0], *tt.want.Rules[i].Rules[0])
						}
					}
				}
				if len(got.Rules) < len(tt.want.Rules) {
					t.Errorf("missing rules: %+v", tt.want.Rules[len(got.Rules):])
				}
				t.Errorf("RemoveUnusedSelectors():\n %v, \n\n want:\n %v", got, tt.want)
			}
		})
	}
}
