package optimizer

import (
	"goHtmlBuilder/css"
	"goHtmlBuilder/minify"
	"goHtmlBuilder/utils"
	"path/filepath"
	"strings"
)

type CssFile struct {
	fileName string
	content  string
}

func (f CssFile) GetContentPath() string {
	return filepath.Join(staticCssPath, f.fileName)
}

func (f CssFile) GetSavePath() string {
	return filepath.Join(distStaticCssPath, f.fileName)
}

func (f CssFile) GetContent() (string, error) {
	if f.content == "" {
		cssContent, err := utils.ReadAllFile(f.GetContentPath())
		if err != nil {
			return "", err
		}

		f.content = strings.Join(cssContent, "")
	}

	return f.content, nil
}

func (f CssFile) GetOptimizedContent(usedSelectors []Selector) (string, error) {
	cssContent, err := f.GetContent()
	if err != nil {
		return "", err
	}

	styles, err := css.Parse(cssContent)
	if err != nil {
		return "", err
	}
	OptimizedStyles := RemoveUnusedSelectors(*styles, usedSelectors)

	return OptimizedStyles.String(), nil
}

func (f CssFile) GetOptimizedAndMinifiedContent(usedSelectors []Selector) (string, error) {
	c, err := f.GetOptimizedContent(usedSelectors)
	if err != nil {
		return "", err
	}

	return minify.MinifyCSS(c)
}

func (f CssFile) SaveOptimizedContent(usedSelectors []Selector) error {
	c, err := f.GetOptimizedContent(usedSelectors)
	if err != nil {
		return err
	}

	return utils.SaveToFile(f.GetSavePath(), c)
}

func (f CssFile) SaveOptimizedAndMinifiedContent(usedSelectors []Selector) error {
	c, err := f.GetOptimizedAndMinifiedContent(usedSelectors)
	if err != nil {
		return err
	}

	return utils.SaveToFile(f.GetSavePath(), c)
}
