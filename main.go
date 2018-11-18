package main

import (
	"bufio"
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
	"sync"
)

const NumberOfThreadsDefault = 5

func main() {
	var MaxWorkers int
	var InputFile string
	var OutputDir string
	var Verbose bool

	flag.StringVar(&InputFile, "f", "", "Input file with urls you want to download. Each url must be in separate line.")
	flag.StringVar(&OutputDir, "d", "", "Output directory to which you want to save images.")
	flag.IntVar(&MaxWorkers, "th", NumberOfThreadsDefault,"Number of threads.")
	flag.BoolVar(&Verbose, "v", false, "Verbose")
	flag.Parse()

	OutputDir = strings.TrimRight(OutputDir, "/")

	if InputFile == "" || MaxWorkers < 1 {
		flag.Usage()
		return
	}

	// 1. Файл в урлами
	file, err := os.Open(InputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 2. Проверить директорию
	if _, err1 := os.Stat(OutputDir); err1 != nil {
		log.Fatal(err1)
	}

	// 3. Получаем урлы
	var downloads []Download

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url, filename := parseRow(scanner.Text())
		if url != "" {
			downloads = append(downloads, Download{url, filename})
		}
	}

	dl := DownloadList{downloads, sync.Mutex{}, 0}
	dl.fillOutputFileNames()

	c := make(chan int, MaxWorkers)

	// 4. Качаем картинки
	for i:=0; i < MaxWorkers; i++ {
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
			c <- 0
		}(i)
	}

	for i:=0; i < MaxWorkers; i++ {
		<-c
	}
}

func parseRow(fileRow string) (string, string) {
	var url, filename string

	if len(fileRow) == 0 {
		return url, filename
	}

	splitRow := strings.Split(fileRow, "|")
	url = splitRow[0]

	if len(splitRow) > 1 {
		filename = splitRow[1]
	}

	return url, filename
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
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(size)/1024/1024)
	default:
		return fmt.Sprintf("%.1f GB", float64(size)/1024/1024/1024)
	}
}