package optimizer

import (
	"golang.org/x/net/html"
	"slices"
	"strings"
)

type SelectorType string

const (
	selectorClass SelectorType = "class"
	selectorId    SelectorType = "id"
)

type Selector struct {
	Value string
	SType SelectorType
}

func GetAllSelectors(htmlCode string) ([]Selector, error) {
	doc, err := html.Parse(strings.NewReader(htmlCode))
	if err != nil {
		return nil, err
	}

	var selectors []Selector

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, el := range n.Attr {
				if el.Key == "class" {
					for _, class := range strings.Split(el.Val, " ") {
						s := strings.ReplaceAll(class, " ", "")
						if s != "" {
							selectors = append(selectors, Selector{
								Value: s,
								SType: selectorClass,
							})
						}

					}
				}

				if el.Key == "id" {
					s := strings.ReplaceAll(el.Val, " ", "")
					if s != "" {
						selectors = append(selectors, Selector{
							Value: s,
							SType: selectorId,
						})
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return selectors, nil
}

func GetAllClasses(htmlCode string) ([]Selector, error) {
	allSelectors, err := GetAllSelectors(htmlCode)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(allSelectors, func(selector Selector) bool {
		return selector.SType != selectorClass
	}), nil
}

func GetAllIds(htmlCode string) ([]Selector, error) {
	allSelectors, err := GetAllSelectors(htmlCode)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(allSelectors, func(selector Selector) bool {
		return selector.SType != selectorId
	}), nil
}
