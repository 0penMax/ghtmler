package optimizer

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

func (qs *quickSelectors) filterSelectors(selectors []string) []string {
	var filtered []string
	for _, sl := range selectors {
		tokens := tokenizeCSSSelector(sl)
		for _, token := range tokens {
			if qs.isContain(token) {
				filtered = append(filtered, sl)
				break
			}
		}
	}

	return filtered
}
