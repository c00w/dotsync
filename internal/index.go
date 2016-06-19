package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Index struct {
	fd      *os.File
	entries map[string]string
}

func OpenIndex(dir string) (*Index, error) {
	index, err := os.OpenFile(filepath.Join(dir, "INDEX"), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	i := &Index{index, make(map[string]string)}
	i.readin()
	return i, nil
}

func (i *Index) readin() {
	contents, err := ioutil.ReadAll(i.fd)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range strings.Split(string(contents), "\n") {
		if len(e) == 0 {
			continue
		}
		i.entries[filepath.Base(e)] = e
	}
}

func (i *Index) write() {
	err := i.fd.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range i.entries {
		_, err := i.fd.WriteString(v + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (i *Index) ListFiles() []string {
	names := make([]string, 0, len(i.entries))
	for name, _ := range i.entries {
		names = append(names, name)
	}
	return names
}

func (i *Index) Close() {
	i.fd.Close()
}

func (i *Index) Update(srcpath string) {
	file := filepath.Base(srcpath)
	if i.entries[file] != srcpath {
		i.entries[file] = srcpath
		i.write()
	}
}

func (i *Index) Get(file string) string {
	return i.entries[file]
}
