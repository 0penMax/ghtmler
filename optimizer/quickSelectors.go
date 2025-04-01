package optimizer

import "goHtmlBuilder/css"

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

func (qs *quickSelectors) filterSelectors(selectors []css.Selector) []css.Selector {
	var filtered []css.Selector
	for _, sl := range selectors {
		for _, token := range sl.Tokenize() {
			if qs.isContain(token) {
				filtered = append(filtered, sl)
				break
			}
		}
	}

	return filtered
}
