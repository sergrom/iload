package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Download struct {
	url      string
	filename string
}

type DownloadList struct {
	file         *os.File
	scanner      *bufio.Scanner
	downloads    []Download
	mux          sync.Mutex
	currentIndex int
}

func (dl *DownloadList) Open(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	dl.file = f
	dl.scanner = bufio.NewScanner(dl.file)
	return nil
}

func (dl *DownloadList) Close() {
	dl.file.Close()
}

func (dl *DownloadList) getNext() (Download, bool) {
	download, done := Download{}, true

	dl.mux.Lock()

	for dl.scanner.Scan() {
		url, filename := parseRow(dl.scanner.Text())
		if url != "" {
			if filename == "" {
				filename = filepath.Base(url)
			}

			download, done = Download{url, filename}, false
			break
		}
	}

	dl.mux.Unlock()

	return download, done
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