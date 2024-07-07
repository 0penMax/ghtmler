package optimizer

import (
	"goHtmlBuilder/css"
	"goHtmlBuilder/utils"
	"path/filepath"
	"strings"
)

type CssFile struct {
	fileName string
}

func (f CssFile) GetContentPath() string {
	return filepath.Join(staticCssPath, f.fileName)
}

func (f CssFile) GetSavePath() string {
	return filepath.Join(distStaticCssPath, f.fileName)
}

func (f CssFile) GetOptimizedContent(usedSelectors []Selector) (string, error) {
	cssContent, err := utils.ReadAllFile(f.GetContentPath())
	if err != nil {
		return "", err
	}

	styles, err := css.Parse(strings.Join(cssContent, ""))
	if err != nil {
		return "", err
	}
	OptimizedStyles := RemoveUnusedSelectors(*styles, usedSelectors)

	return OptimizedStyles.String(), nil
}
