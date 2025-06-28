package optimizer

import (
	"golang.org/x/net/html"
	"io"
)

type injectedFile struct {
	fileName string
	fType    string //like js or css
}

const fTypeCss = "css"
const fTypeJs = "js"
const fTypeNone = "none"

func getInjectedFilesFromHtml(r io.Reader, linkExtractF []func(token html.Token) injectedFile) ([]injectedFile, error) {
	var injectedFiles []injectedFile
	tokenizer := html.NewTokenizer(r)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return injectedFiles, nil
			}
			return nil, err
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			for _, f := range linkExtractF {
				inj := f(token)
				if inj.fType != fTypeNone {
					injectedFiles = append(injectedFiles, inj)
				}
			}
		}
	}
}

func getCSSFileNames(token html.Token) injectedFile {
	if token.Data == "link" {
		var isStylesheet bool
		var href string
		for _, attr := range token.Attr {
			if attr.Key == "rel" && attr.Val == "stylesheet" {
				isStylesheet = true
			}
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		if isStylesheet && href != "" {

			return injectedFile{fileName: getFileName(href), fType: fTypeCss}
		}
	}
	return injectedFile{fType: fTypeNone}
}

// TODO write test
func GetJSFileNames(token html.Token) injectedFile {
	if token.Data == "script" {
		var src string
		for _, attr := range token.Attr {
			if attr.Key == "src" {
				src = attr.Val
			}
		}
		if src != "" {
			return injectedFile{fileName: getFileName(src), fType: fTypeJs}
		}
	}
	return injectedFile{fType: fTypeNone}
}
