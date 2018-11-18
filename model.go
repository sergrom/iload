package main

import (
	"path/filepath"
	"sync"
)

type Download struct {
	url      string
	filename string
}

type DownloadList struct {
	downloads    []Download
	mux          sync.Mutex
	currentIndex int
}

func (dl *DownloadList) getNext() (Download, bool) {
	var download Download
	var done bool

	dl.mux.Lock()

	if dl.currentIndex >= len(dl.downloads) {
		download, done = Download{}, true
	} else {
		download, done = dl.downloads[dl.currentIndex], false
		dl.currentIndex++
	}

	dl.mux.Unlock()

	return download, done
}

func (dl *DownloadList) fillOutputFileNames() {
	for i, d := range dl.downloads {
		if d.filename == "" {
			dl.downloads[i].filename = filepath.Base(d.url)
		}
	}
}