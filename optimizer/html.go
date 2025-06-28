package optimizer

import (
	"io"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type SelectorType string

const (
	selectorClass SelectorType = "class"
	selectorId    SelectorType = "id"
	selectorTag   SelectorType = "tag"
)

const staticCssPath = "static/css"
const distStaticCssPath = "dist/static/css"

type Selector struct {
	Value string
	SType SelectorType
}

func getFileName(fullPath string) string {
	if fullPath == "" || fullPath == "/" || fullPath == "\\" {
		return ""
	}
	fullPath = strings.ReplaceAll(fullPath, "\\\\", "/")
	fullPath = strings.ReplaceAll(fullPath, "\\", "/")
	return filepath.Base(fullPath)
}

//TODO think, maybe i can combine CssFile and JsFile in one universal struct

func GetCssAndJsFileNamesFromHtml(r io.Reader) ([]CssFile, []JsFile, error) {
	var functions []func(token html.Token) injectedFile
	functions = append(functions, getCSSFileNames, GetJSFileNames)
	files, err := getInjectedFilesFromHtml(r, functions)
	if err != nil {
		return nil, nil, err
	}
	var cssFiles []CssFile
	var jsFiles []JsFile
	for _, file := range files {
		switch file.fType {
		case fTypeCss:
			cssFiles = append(cssFiles, CssFile{fileName: file.fileName})
		case fTypeJs:
			jsFiles = append(jsFiles, JsFile{fileName: file.fileName})
		}
	}

	return cssFiles, jsFiles, nil
}

type JsFile struct {
	fileName string
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
