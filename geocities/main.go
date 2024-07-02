package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	// TODO LOWERCASE
	// TODO NUMBERS
	src            = "/home/vilmibm/geocities/%s/geocities/YAHOOIDS"
	userDirPattern = "/home/vilmibm/geocities/*/geocities/YAHOOIDS/?/?/*"
	t              = "/home/vilmibm/gc"
)

func main() {
	userDirs := make(chan string)
	var wg sync.WaitGroup

	walkFn := func(s string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		isUserDir, err := filepath.Match(userDirPattern, s)
		if err != nil {
			return err
		}

		if isUserDir {
			userDirs <- s
		}

		return nil
	}

	go func() {
		dirs := []string{"UPPERCASE", "LOWERCASE", "NUMBERS"}
		for _, dir := range dirs {
			err := filepath.WalkDir(fmt.Sprintf(src, dir), walkFn)
			if err != nil {
				panic(err)
			}
		}
		close(userDirs)
	}()

	totalUserDirs := 0

	for ud := range userDirs {
		totalUserDirs++
		wg.Add(1)
		go processUserDir(&wg, ud)
	}

	wg.Wait()

	fmt.Printf("processed %d user dirs\n", totalUserDirs)
}

func processUserDir(wg *sync.WaitGroup, ud string) {
	defer wg.Done()

	var outFile *os.File
	walkFn := func(s string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if ignoreSuffix(s) {
			return nil
		}
		fmt.Printf("READING: %s\n", s)
		f, err := os.Open(s)
		if err != nil {
			return err
		}
		defer f.Close()
		if !wantSuffix(s) {
			// if it's not a text NOR an ignore suffix, try and guess if text

			// cribbed from godoc's source
			var buf [1024]byte
			n, err := f.Read(buf[0:])
			if err != nil {
				fmt.Printf("\terror reading file: %s\n", err.Error())
				return nil
			}

			if !IsText(buf[0:n]) {
				fmt.Printf("NOT TEXT: %s\n", s)
				return nil
			}
		}

		if outFile == nil {
			outFile, err = os.Create(filepath.Join(t, filepath.Base(ud)))
			if err != nil {
				panic(err)
			}
		}
		all, err := os.ReadFile(s)
		if err != nil {
			panic(err)
		}
		wrote, _ := outFile.Write(all)
		fmt.Printf("WROTE: %s (%d bytes)\n", s, wrote)

		return nil
	}

	err := filepath.WalkDir(ud, walkFn)
	if err != nil {
		panic(err)
	}
	if outFile != nil {
		outFile.Close()
	}
}

// from godoc source
func IsText(s []byte) bool {
	const max = 1024 // at least utf8.UTFMax
	if len(s) > max {
		s = s[0:max]
	}
	for i, c := range string(s) {
		if i+utf8.UTFMax > len(s) {
			// last char may be incomplete - ignore
			break
		}
		if c == 0xFFFD || c < ' ' && c != '\n' && c != '\t' && c != '\f' {
			// decoding error or control character - not a text file
			return false
		}
	}
	return true
}

var ignoreSuffices []string
var wantSuffices []string

func init() {
	ignoreSuffices = []string{
		"jpg",
		"jpeg",
		"gif",
		"js",
		"css",
		"mp3",
		"wav",
		"midi",
		"JPG",
		"JPEG",
		"GIF",
		"JS",
		"CSS",
		"MP3",
		"WAV",
		"MIDI",
	}
	wantSuffices = []string{
		"html",
		"htm",
		"txt",
		"HTML",
		"HTM",
		"TXT",
	}
}

func ignoreSuffix(p string) bool {
	for _, s := range ignoreSuffices {
		if strings.HasSuffix(p, s) {
			return true
		}
	}

	return false
}

func wantSuffix(p string) bool {
	for _, s := range wantSuffices {
		if strings.HasSuffix(p, s) {
			return true
		}
	}

	return false

}
