package builder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const includeConst = "@include"

const staticDirPath = "static/"

func Build(ghmlFiles []string) error {
	for _, f := range ghmlFiles {
		err := buildHtml("./" + f)
		if err != nil {
			return errors.New(fmt.Sprintf("error build from file %s: %s", f, err.Error()))
		}
	}

	err := copyDir(staticDirPath, "dist/static")
	if err != nil {
		return err
	}

	return nil
}

func buildHtml(fpath string) error {
	fileStrs, err := readAllFile(fpath)
	if err != nil {
		return err
	}

	var resultStrs []string

	for _, s := range fileStrs {
		if strings.Contains(s, includeConst) {
			s = strings.ReplaceAll(s, includeConst, "")
			includePath := removeSpaceAndTab(s)

			htmlStrs, err := readAllFile(includePath)
			if err != nil {
				return err
			}

			resultStrs = append(resultStrs, htmlStrs...)
			continue

		}

		resultStrs = append(resultStrs, s)
	}

	fname := getFileNameOnly(fpath)

	err = os.MkdirAll("./dist/", os.ModeDir)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("./dist/"+fname+".html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	defer file.Close()

	if err != nil {
		return err
	}

	err = write2FileLineByLine(file, resultStrs)

	return err

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
