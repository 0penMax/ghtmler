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
		if fs := filterSelectors(rule.Selectors, sl); fs != nil {
			rule.Selectors = fs
			newCss.Rules = append(newCss.Rules, rule)
		}
	}
	return newCss
}

func clearRuleSubrules(rule *css.Rule, usedSelectors []string) *css.Rule {
	embRules := rule.Rules
	rule.Rules = nil

	for _, erule := range embRules {
		if fs := filterSelectors(erule.Selectors, usedSelectors); fs != nil {
			erule.Selectors = fs
			rule.Rules = append(rule.Rules, erule)
		}
	}

	rule.Selectors = filterSelectors(rule.Selectors, usedSelectors)

	return rule
}

func filterSelectors(slice1, mustContain []string) []string {
	// Create a map to hold the values of mustContain for faster lookup
	set := make(map[string]struct{})
	for _, item := range mustContain {
		set[item] = struct{}{}
	}

	// Filter slice1
	var filtered []string
	for _, item := range slice1 {
		if _, exists := set[item]; exists {
			filtered = append(filtered, item)
		}
	}

	return filtered
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
