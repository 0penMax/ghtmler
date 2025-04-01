package optimizer

import (
	"goHtmlBuilder/css"
)

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
		if fs := qs.filterSelectors(rule.Selectors); fs != nil {
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
		if fs := usedSelectors.filterSelectors(erule.Selectors); fs != nil {
			erule.Selectors = fs
			rule.Rules = append(rule.Rules, erule)
		}
	}

	rule.Selectors = usedSelectors.filterSelectors(rule.Selectors)

	return rule
}
