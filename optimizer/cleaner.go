package optimizer

import (
	"goHtmlBuilder/css"
)

type quickSelectors struct {
	data map[string]bool
}

func (qs *quickSelectors) init(selectors []Selector) {
	qs.data = make(map[string]bool)
	for _, selector := range selectors {
		qs.data[selector.Value] = true

	}
}
func (qs *quickSelectors) isContain(selector string) bool {
	return qs.data[selector]
}

func RemoveUnusedSelectors(cssContent css.Stylesheet, usedSelectors []Selector) css.Stylesheet {
	var qs quickSelectors
	qs.init(usedSelectors)

	var newCss css.Stylesheet
	for _, rule := range cssContent.Rules {
		if rule.Kind == css.AtRule {
			if rule.EmbedsRules() {
				newRule := clearRuleSubrules(rule, qs)
				if len(newRule.Rules) != 0 {
					newCss.Rules = append(newCss.Rules, rule)
				}
				continue
			}

			newCss.Rules = append(newCss.Rules, rule)
			continue
		}
		if fs := filterSelectors(rule.Selectors, qs); fs != nil {
			rule.Selectors = fs
			newCss.Rules = append(newCss.Rules, rule)
		}
	}
	return newCss
}

func clearRuleSubrules(rule *css.Rule, usedSelectors quickSelectors) *css.Rule {
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

func filterSelectors(sls []string, mustContain quickSelectors) []string {
	var filtered []string
	for _, sl := range sls {
		tokens := tokenizeCSSSelector(sl)
		for _, token := range tokens {
			if mustContain.isContain(token) {
				filtered = append(filtered, sl)
			}
		}
	}

	return filtered
}
