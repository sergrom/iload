package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Loader ...
type Loader struct {
	inputFile       string
	outputDirectory string
}

// NewLoader ...
func NewLoader() *Loader {
	return &Loader{}
}

// WithInputFile ...
func (l *Loader) WithInputFile(inputFile string) *Loader {
	l.inputFile = inputFile
	return l
}

// WithOutputDirectory ...
func (l *Loader) WithOutputDirectory(outputDirectory string) *Loader {
	l.outputDirectory = outputDirectory
	return l
}

// Download ...
func (l *Loader) Download(threadsNum int, verbose bool) error {
	outputDir := strings.TrimRight(l.outputDirectory, "/")

	if threadsNum < 1 {
		return fmt.Errorf("wrong threads number value: %d", threadsNum)
	}

	var dl DownloadList

	// 1. Открываем файл
	err := dl.Open(l.inputFile)
	if err != nil {
		return err
	}
	defer dl.Close()

	// 2. Проверить директорию
	if _, err := os.Stat(outputDir); err != nil {
		return err
	}

	c := make(chan int, threadsNum)

	// 3. Качаем файлы
	for i := 0; i < threadsNum; i++ {
		go func(i int) {
			for {
				d, done := dl.getNext()
				if done == true {
					break
				}

				size, fName, e := saveFile(d.url, outputDir+"/"+d.filename)
				if verbose {
					if e == nil {
						fmt.Printf("[%d] Ok: url \"%s\" saved as \"%s\" (%s)\n", i, d.url, fName, getPrettySize(size))
					} else {
						fmt.Printf("[%d] Fail: %s. %s\n", i, d.url, e)
					}
				}
			}
			defer func() {
				c <- 0
			}()
		}(i)
	}

	for i := 0; i < threadsNum; i++ {
		<-c
	}

	return nil
}

func saveFile(url string, outputFile string) (int64, string, error) {
	response, err1 := http.Get(url)
	if err1 != nil {
		return 0, "", err1
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return 0, "", errors.New(response.Status)
	}

	finalOutputFile := getOutputFile(outputFile)

	file, err2 := os.Create(finalOutputFile)
	if err2 != nil {
		return 0, "", err2
	}
	defer file.Close()

	fSize, err3 := io.Copy(file, response.Body)
	if err3 != nil {
		return 0, "", err3
	}

	return fSize, finalOutputFile, nil
}

func getOutputFile(outputFile string) string {
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return outputFile
	}

	var dir, filename, partName, partNumber, partExt string
	re1, _ := regexp.Compile("^(.+)_([0-9]{1,2})(\\..+)$")
	re2, _ := regexp.Compile("^(.+)(\\..+)$")

	filename = filepath.Base(outputFile)
	dir = filepath.Dir(outputFile)

	parts := re1.FindStringSubmatch(filename)

	if len(parts) == 4 {
		partName, partNumber, partExt = parts[1], parts[2], parts[3]
	} else {
		parts = re2.FindStringSubmatch(filename)
		if len(parts) == 3 {
			partName, partNumber, partExt = parts[1], "0", parts[2]
		} else {
			partName, partNumber, partExt = parts[1], "0", ""
		}
	}

	nn, _ := strconv.Atoi(partNumber)
	fName := partName + "_" + partNumber + partExt
	for nn < 99 {
		nn++
		fName = partName + "_" + strconv.Itoa(nn) + partExt
		if _, err := os.Stat(dir + "/" + fName); os.IsNotExist(err) {
			return dir + "/" + fName
		}
	}

	return dir + "/" + fName
}

func getPrettySize(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2f MB", float64(size)/1024/1024)
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/1024/1024/1024)
	}
}
