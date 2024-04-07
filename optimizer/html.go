package optimizer

import (
	"golang.org/x/net/html"
	"strings"
)

func GetAllClasses(htmlCode string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlCode))
	if err != nil {
		return nil, err
	}

	var classes []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, el := range n.Attr {
				if el.Key == "class" {
					classes = append(classes, strings.Split(el.Val, " ")...)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return classes, nil
}
