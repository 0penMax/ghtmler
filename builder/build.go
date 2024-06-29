package builder

import (
	"bufio"
	"errors"
	"fmt"
	"goHtmlBuilder/css"
	"goHtmlBuilder/minify"
	"goHtmlBuilder/optimizer"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const includeConst = "@include"

const staticDirPath = "static/"

const LIVE_RELOAD_FOLDER = "./liveReload/"

type GhtmlFile struct {
	filename     string
	content      []string
	cssFilesPath []optimizer.CssFile
	isLiveReload bool
	minifyParams minify.Params
}

// TODO test this func
// TODO add minify
func (g *GhtmlFile) save() error {
	err := writeLines2File("./dist/"+g.filename+".html", g.content)
	if err != nil {
		return err
	}

	selectors, err := optimizer.GetAllSelectors(strings.Join(g.content, ""))

	for _, cssFile := range g.cssFilesPath {

		cssContent, err := readAllFile(cssFile.GetContentPath())
		if err != nil {
			return err
		}

		styles, err := css.Parse(strings.Join(cssContent, ""))
		if err != nil {
			return err
		}
		OptimizedStyles := css.RemoveUnusedSelectors(*styles, selectors)

		err = saveToFile(cssFile.GetSavePath(), OptimizedStyles.String())
		if err != nil {
			return err
		}

	}

	if g.isLiveReload {
		err = writeLines2File(LIVE_RELOAD_FOLDER+g.filename+".html", injectLiveReloadScript(g.content))
		if err != nil {
			return err
		}

	}

	return nil
}

func BuildGthmlFile(file string, isLiveReload bool, minifyParams minify.Params) (GhtmlFile, error) {
	content, err := buildHtml(file)
	if err != nil {
		return GhtmlFile{}, errors.New(fmt.Sprintf("error build from file %s: %s", file, err.Error()))
	}
	fname := getFileNameOnly(file)

	r := strings.NewReader(strings.Join(content, ""))

	cssPaths, err := optimizer.GetCSSFileNamesFromHtml(r)

	ghtmlFile := GhtmlFile{
		filename:     fname,
		content:      content,
		cssFilesPath: cssPaths,
		isLiveReload: isLiveReload,
		minifyParams: minifyParams,
	}

	return ghtmlFile, nil

}

func Build(ghmlFiles []string, isLiveReload bool, minifyParam minify.Params) error {

	for _, f := range ghmlFiles {

		ghtmlFile, err := BuildGthmlFile(f, isLiveReload, minifyParam)
		if err != nil {
			return err
		}

		err = ghtmlFile.save()
		if err != nil {
			return err
		}

	}

	err := copyDir(staticDirPath, "dist/static")
	if err != nil {
		return err
	}

	return nil
}

func buildHtml(fpath string) ([]string, error) {
	fileStrs, err := readAllFile(fpath)
	if err != nil {
		return nil, err
	}

	var resultStrs []string

	for _, s := range fileStrs {
		if strings.Contains(s, includeConst) {
			s = strings.ReplaceAll(s, includeConst, "")
			includePath := removeSpaceAndTab(s)

			htmlStrs, err := readAllFile(includePath)
			if err != nil {
				return nil, err
			}

			resultStrs = append(resultStrs, htmlStrs...)
			continue

		}

		resultStrs = append(resultStrs, s)
	}

	return resultStrs, nil

}

// saveToFile saves the provided string to the specified file.
// If the file already exists, it will be overwritten.
func saveToFile(filename, data string) error {
	// Open the file with write permissions. Create it if it doesn't exist.
	// The file permissions are set to 0644, meaning read and write for the owner, and read-only for others.
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func writeLines2File(fpath string, content []string) error {
	if err := os.MkdirAll(filepath.Dir(fpath), 0770); err != nil {
		return err
	}
	file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	return write2FileLineByLine(file, content)
}

func readAllFile(filepath string) ([]string, error) {
	bytesRead, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	file_content := string(bytesRead)
	file_content = removeHtmlComment(file_content)
	lines := strings.Split(file_content, "\n")
	return lines, nil
}

func write2FileLineByLine(file *os.File, lines []string) error {
	datawriter := bufio.NewWriter(file)

	for _, data := range lines {
		_, _ = datawriter.WriteString(data + "\n")
	}

	return datawriter.Flush()
}

func removeExtraSpace(str string) string {
	// the character class \s matches a space, tab, new line, carriage return or form feed, and + says “one or more of those”.
	space := regexp.MustCompile(`\s{2,}`)
	return space.ReplaceAllString(str, " ")
}

func removeHtmlComment(str string) string {

	for {
		startIndex := strings.Index(str, "<!--")
		endIndex := strings.Index(str, "-->")
		if startIndex == -1 || endIndex == -1 {
			break
		}

		str = str[:startIndex] + str[endIndex+3:]
	}

	return str
}

func removeSpaceAndTab(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, " ", "")
	return strings.ReplaceAll(str, "	", "")
}

func getFileNameOnly(fpath string) string {
	fileName := path.Base(fpath)
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func copyAll(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0770); err != nil {
		return 0, err
	}
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// dir with all file inside
func copyDir(source, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(strings.ReplaceAll(path, "\\", "/"), source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.MkdirAll(filepath.Join(destination, relPath), 0755)
		} else {
			_, err = copyAll(filepath.Join(source, relPath), filepath.Join(destination, relPath))
			return err
		}
	})
	return err
}
