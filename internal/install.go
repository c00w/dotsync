package internal

import (
	"log"
	"os"
	"path/filepath"
)

func Install(curdir, home string) {
	index, err := OpenIndex(curdir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range index.ListFiles() {
		path := index.Get(file)
		err := os.MkdirAll(filepath.Join(home, filepath.Dir(path)), os.ModePerm)
		if err != nil {
			log.Fatal("error mkdir", err)
		}
		MoveFile(filepath.Join(curdir, file), filepath.Join(home, path))
	}
}
