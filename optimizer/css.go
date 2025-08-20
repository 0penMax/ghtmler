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

func (f *CssFile) Process(usedSelectors []Selector, params minify.Params) error {
	if params.IsOptiMiniCss() || params.IsOptimizeCss {
		err := f.optimize(usedSelectors)
		if err != nil {
			return err
		}
	}
	if params.IsOptiMiniCss() || params.IsMinifyCss {
		err := f.minimize()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *CssFile) Save() error {
	return utils.SaveToFile(f.getSavePath(), f.content)
}

func (f *CssFile) getContentPath() string {
	return filepath.Join(staticCssPath, f.fileName)
}

func (f *CssFile) getSavePath() string {
	return filepath.Join(distStaticCssPath, f.fileName)
}

func (f *CssFile) getContent() (string, error) {
	if f.content == "" {
		cssContent, err := utils.ReadAllFile(f.getContentPath())
		if err != nil {
			return "", err
		}

		f.content = strings.Join(cssContent, "")
	}

	return f.content, nil
}

func (f *CssFile) optimize(usedSelectors []Selector) error {
	cssContent, err := f.getContent()
	if err != nil {
		return err
	}

	styles, err := css.Parse(cssContent)
	if err != nil {
		return err
	}
	OptimizedStyles := RemoveUnusedSelectors(*styles, usedSelectors)
	f.content = OptimizedStyles.String()

	return nil
}

func (f *CssFile) minimize() error {
	c, err := f.getContent()
	if err != nil {
		return err
	}

	c, err = minify.MinifyCSS(c)
	if err != nil {
		return err
	}

	f.content = c

	return nil
}
