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
			if rule.EmbedsRules() {
				newRule := clearRuleSubrules(rule, sl)
				if len(newRule.Rules) != 0 {
					newCss.Rules = append(newCss.Rules, rule)
				}
				continue
			}

			newCss.Rules = append(newCss.Rules, rule)
			continue
		}
		if isContain(sl, rule.Selectors) {
			newCss.Rules = append(newCss.Rules, rule)
		}
	}
	return newCss
}

func clearRuleSubrules(rule *css.Rule, usedSelectors []string) *css.Rule {
	embRules := rule.Rules
	rule.Rules = nil

	for _, erule := range embRules {
		if isContain(usedSelectors, erule.Selectors) {
			rule.Rules = append(rule.Rules, erule)
		}
	}

	return rule
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
