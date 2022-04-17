package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const gExt = ".ghtml"

const includeConst = "@include"

const staticDirPath = "static/"

func main(){
	f, err := os.OpenFile("error.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)



	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == gExt {
			err = buildHtml("./" + f.Name())
			if err!=nil{
				log.Fatal(err)
			}
		}

	}

	err = copyDir(staticDirPath, "dist/static")
	if err!=nil{
		log.Println(err)
	}

}


func buildHtml(fpath string)error{
	fileStrs,err:= readAllFile(fpath)
	if err!=nil{
		return err
	}

	var resultStrs []string

	for _,s := range fileStrs{
		if strings.Contains(s, includeConst ){
			s = strings.ReplaceAll(s, includeConst, "")
			path := removeSpaceAndTab(s)

			htmlStrs, err:= readAllFile(path)
			if err!=nil{
				return err
			}


			resultStrs = append(resultStrs, htmlStrs...)
			continue

		}

		resultStrs = append(resultStrs, s)
	}

	fname:= getFileNameOnly(fpath)

	err = os.MkdirAll("./dist/", os.ModeDir)
	if err!=nil{
		return err
	}


	err = write2FileLineByLine(resultStrs, "./dist/" + fname + ".html")

	return err

}


func readAllFile(filepath string)([]string, error){
	bytesRead, err := ioutil.ReadFile(filepath)
	if err!=nil{
		return nil, err
	}
	file_content := string(bytesRead)
	file_content = removeHtmlComment(file_content)
	lines := strings.Split(file_content, "\n")
	return lines, nil
}

func write2FileLineByLine(lines []string, filePath string)error{

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	defer file.Close()

	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range lines {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()

	return nil
}

func removeExtraSpace(str string)string{
	// the character class \s matches a space, tab, new line, carriage return or form feed, and + says “one or more of those”.
	space := regexp.MustCompile(`\s{2,}`)
	return space.ReplaceAllString(str, " ")
}

func removeHtmlComment(str string)string  {


	for {
		startIndex:= strings.Index(str, "<!--")
		endIndex:= strings.Index(str, "-->")
		if startIndex == -1 || endIndex == -1 {
			break
		}

		str = str[:startIndex] + str[endIndex + 3:]
	}

	return str
}

func removeSpaceAndTab(str string)string{
	str = strings.ReplaceAll(str, " ", "")
	return strings.ReplaceAll(str, "	", "")
}


func getFileNameOnly(fpath string)string{
		fileName:= path.Base(fpath)
		if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
			return fileName[:pos]
		}
		return fileName
}


func copy(src, dst string) (int64, error) {
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


/*copy dir with all file inside*/
func copyDir(source, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(strings.ReplaceAll(path, "\\", "/"), source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.MkdirAll(filepath.Join(destination, relPath), 0755)
		} else {
			_,err = copy(filepath.Join(source, relPath),filepath.Join(destination, relPath) )
			return  err
		}
	})
	return err
}

