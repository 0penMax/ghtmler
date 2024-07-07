package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func WriteLines2File(fpath string, content []string) error {
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

func ReadAllFile(filepath string) ([]string, error) {
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
