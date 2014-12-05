package internal

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type collator []string

// Shoves all files ending in rc down the channel.
func (c *collator) Walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("Error walking: ", filepath.Join(path, info.Name()), err)
		return nil
	}
	if info.IsDir() {
		return nil
	}
	if info.Name() == "xmonad.hs" {
		*c = append(*c, path)
		return nil
	}
	if strings.HasSuffix(info.Name(), "rc") {
		*c = append(*c, path)
	}
	return nil
}

// FlatWalk walks a directory without recursing.
func FlatWalk(path string, walker filepath.WalkFunc) error {
	homefd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer homefd.Close()

	files, err := homefd.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range files {
		err := walker(filepath.Join(path, fi.Name()), fi, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// Says whether a and b are different, note that if b does not exist the error will be hidden
func DiffFile(a, b string) (bool, error) {
	afd, err := os.Open(a)
	if err != nil {
		return true, err
	}
	defer afd.Close()
	bfd, err := os.Open(b)
	// Ignore destination does not exist errors
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer bfd.Close()
	ahash := sha512.New()
	bhash := sha512.New()
	_, err = io.Copy(ahash, afd)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(bhash, bfd)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.Equal(ahash.Sum(nil), bhash.Sum(nil)), nil
}

// FilterDiff checks that the file has changed from the current directory, and output all which have.
func FilterDiff(in []string, dir string) []string {
	out := make([]string, 0)
	for _, f := range in {
		same, err := DiffFile(f, filepath.Join(dir, filepath.Base(f)))
		if err != nil {
			log.Fatal(err)
		}
		if !same {
			out = append(out, f)
		}
	}
	return out
}

// Return y or no, and whether it is permanent
func Prompt(file string) (bool, bool) {
	for {
		fmt.Print(file, " Copy? Y/N/A(all):")
		input := ""
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		switch input {
		case "Y":
			return true, false
		case "N":
			return false, false
		case "A":
			return true, true
		default:
			fmt.Println("I don't understand that")
		}
	}
}

// Prompt promts the user if they want to move the file, offering Y/N/(All)
func PromptFilter(in []string) []string {
	out := make([]string, 0, len(in))
	all := false
	for _, f := range in {
		if all {
			out = append(out, f)
			continue
		}
		yes := false
		yes, all = Prompt(f)
		if yes {
			out = append(out, f)
		}
	}
	return out
}

func MoveFile(srcpath, curdir string) {
	dstpath := filepath.Join(curdir, filepath.Base(srcpath))
	dstfd, err := os.Create(dstpath)
	if err != nil {
		log.Fatal(err)
	}
	defer dstfd.Close()
	srcfd, err := os.Open(srcpath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcfd.Close()
	_, err = io.Copy(dstfd, srcfd)
	if err != nil {
		log.Fatal(err)
	}
}

func striphome(srcpath, homedir string) string {
	if strings.HasPrefix(srcpath, homedir+string(os.PathSeparator)) == false {
		log.Fatal("Src does not have HOME inside of it")
	}
	srcpath = srcpath[len(homedir)+1:]
	return srcpath
}

// Copy each file into the current direction.
func OverWrite(in []string, homedir, curdir string) {
	i, err := OpenIndex(curdir)
	if err != nil {
		log.Fatal(err)
	}
	defer i.Close()
	for _, srcpath := range in {
		MoveFile(srcpath, curdir)
		i.Update(striphome(srcpath, homedir))
	}
}

// Save saves files ending in rc from home to curdir. It prompts for All/Some/None before saving.
func Save(home, curdir string) {
	c := collator(make([]string, 0))

	// Gather rc files
	fmt.Println("Finding changed files")

	FlatWalk(home, c.Walker)

	// walk .config
	config := filepath.Join(home, ".config")
	filepath.Walk(config, c.Walker)

	// Remove unchanged
	files := FilterDiff(c, curdir)

	if len(files) == 0 {
		fmt.Println("No files found")
		return
	}

	fmt.Println("Files changed:")
	for _, name := range files {
		fmt.Println(name)
	}

	// Prompt about changed files
	files = PromptFilter(files)
	OverWrite(files, home, curdir)
}
