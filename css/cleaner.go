package css

import (
	"goHtmlBuilder/optimizer"
	"slices"
)

func RemoveUnusedSelectors(css Stylesheet, usedSelectors []optimizer.Selector) Stylesheet {
	sl := getStringSlice(usedSelectors)

	var newCss Stylesheet
	for _, rule := range css.Rules {
		if rule.Kind == AtRule {
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

func getStringSlice(sl []optimizer.Selector) []string {
	var strs []string
	for _, selector := range sl {
		strs = append(strs, selector.Value)
	}

	return strs
}
