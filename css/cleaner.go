package css

import (
	"slices"
)

func RemoveUnusedSelectors(css Stylesheet, usedSelectors []string) Stylesheet {
	var newCss Stylesheet
	for _, rule := range css.Rules {
		if rule.Kind == AtRule {
			newCss.Rules = append(newCss.Rules, rule)
			continue
		}
		if isContain(usedSelectors, rule.Selectors) {
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
