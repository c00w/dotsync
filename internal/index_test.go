package internal

import (
    "testing"
    "path/filepath"
    "log"
    "io/ioutil"
    "os"
)

func createIndex(i string) (*Index, string) {
    path := filepath.Join(os.TempDir(), "INDEX")
    fd, err := os.Create(path)
    if err != nil {
        log.Fatal(err)
    }
    defer fd.Close()
    _, err = fd.WriteString(i)
    if err != nil {
        log.Fatal(err)
    }
    I, err := OpenIndex(os.TempDir())
    if err != nil {
        log.Fatal(err)
    }
    return I, path
}

func readindex(path string) string {
    fd, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer fd.Close()
    r, err := ioutil.ReadAll(fd)
    if err != nil {
        log.Fatal(err)
    }
    return string(r)
}
func TestIndex(t *testing.T) {
    I, fp := createIndex("")
    defer os.Remove(fp)
    I.Update("foo")
    I.Close()
    if readindex(fp) != "foo\n" {
        t.Errorf("Got %q != %q", readindex(fp), "foo\n")
    }
    if I.Get("foo")!= "foo" {
        t.Errorf("Got %q != %q", I.Get("foo"), "foo")
    }

}
