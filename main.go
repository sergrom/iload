package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const NumberOfThreadsDefault = 5

func main() {
	var NumThreads int
	var InputFile string
	var OutputDir string
	var Verbose bool

	flag.StringVar(&InputFile, "f", "", "Input file with urls you want to download. Each url must be in separate line.")
	flag.StringVar(&OutputDir, "d", "", "Output directory to which you want to save downloaded files.")
	flag.IntVar(&NumThreads, "th", NumberOfThreadsDefault,"Number of threads.")
	flag.BoolVar(&Verbose, "v", false, "Verbose")
	flag.Parse()

	OutputDir = strings.TrimRight(OutputDir, "/")

	if InputFile == "" || OutputDir == "" || NumThreads < 1 {
		flag.Usage()
		return
	}

	var dl DownloadList

	// 1. Открываем файл
	err := dl.Open(InputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer dl.Close()

	// 2. Проверить директорию
	if _, err1 := os.Stat(OutputDir); err1 != nil {
		log.Fatal(err1)
	}

	c := make(chan int, NumThreads)

	// 3. Качаем файлы
	for i:=0; i < NumThreads; i++ {
		go func(i int) {
			for {
				d, done := dl.getNext()
				if done == true {
					break
				}

				size, fname, e := saveFile(d.url, OutputDir + "/" + d.filename)
				if Verbose {
					if e == nil {
						fmt.Printf("[%d] Ok: url \"%s\" saved as \"%s\" (%s)\n", i , d.url, fname, getPrettySize(size))
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

	for i:=0; i < NumThreads; i++ {
		<-c
	}
}

func saveFile(url string, outputFile string) (int64, string, error) {
	response, err1 := http.Get(url)
	if err1 != nil {
		return 0, "", err1
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return 0, "", errors.New(string(response.Status))
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
	re1,_ := regexp.Compile("^(.+)_([0-9]{1,2})(\\..+)$")
	re2,_ := regexp.Compile("^(.+)(\\..+)$")

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

	nn, _:= strconv.Atoi(partNumber)
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
