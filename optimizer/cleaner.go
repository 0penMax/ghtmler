package optimizer

import (
	"goHtmlBuilder/css"
	"slices"
)

func RemoveUnusedSelectors(cssContent css.Stylesheet, usedSelectors []Selector) css.Stylesheet {
	sl := getStringSlice(usedSelectors)

	var newCss css.Stylesheet
	for _, rule := range cssContent.Rules {
		if rule.Kind == css.AtRule {
			newCss.Rules = append(newCss.Rules, rule)
			continue
		}
		if isContain(sl, rule.Selectors) {
			newCss.Rules = append(newCss.Rules, rule)
		}
	}
	return newCss
}

func isContain(what []string, where []string) bool {
	for _, s := range what {
		if c := slices.Contains(where, s); c {
			return true
		}
	}

	return false
}

func getStringSlice(sl []Selector) []string {
	var strs []string
	for _, selector := range sl {
		strs = append(strs, selector.Value)
	}

	return strs
}
