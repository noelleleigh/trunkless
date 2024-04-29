package cutup

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	srcDir  = "/home/vilmibm/pg_plaintext/files"
	tgtDir  = "/home/vilmibm/pg_plaintext/cutup"
	workers = 10
)

// TODO configurable src/tgt dir
// TODO generalize so it's not gutenberg specific

func worker(paths <-chan string, sources chan<- string) {
	// TODO generalize to n character phrase markers, write new function
	phraseMarkers := map[rune]bool{
		';': true,
		',': true,
		':': true,
		'.': true,
		'?': true,
		'!': true,
		//'(':  true,
		')': true,
		//'{':  true,
		'}': true,
		//'[':  true,
		']': true,
		//'\'': true,
		//'"':  true,
		//'“':  true,
		'”': true,
		'=': true,
		'`': true,
		'-': true,
		'|': true,
		'>': true,
	}

	for p := range paths {
		f, err := os.Open(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open '%s': %s\n", p, err.Error())
		}
		s := bufio.NewScanner(f)

		phraseBuff := []byte{}
		written := 0
		inHeader := true
		title := ""
		sourceid := ""

		var of *os.File
		var cleaned string
		var ok bool
		var asStr string
		var text string
		var prefix string

		for s.Scan() {
			text = strings.TrimSpace(s.Text())
			if strings.HasPrefix(text, "*** START") {
				title, _ = strings.CutPrefix(text, "*** START OF THE PROJECT GUTENBERG")
				title, _ = strings.CutPrefix(title, " EBOOK")
				title = strings.Map(rep, title)
				title = strings.TrimSpace(title)
				inHeader = false
				continue
			}
			if inHeader {
				continue
			}
			if strings.HasPrefix(text, "*** END") {
				break
			}
			if title == "" {
				fmt.Fprintf(os.Stderr, "got to cutup phase with no title: '%s'", p)
				break
			}
			if sourceid == "" {
				sourceid = fmt.Sprintf("%x", sha1.Sum([]byte(title)))[0:6]
				prefix = sourceid + "\t"
				of, err = os.Create(path.Join(tgtDir, sourceid))
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not open '%s' for writing: %s", sourceid, err.Error())
					break
				}
			}
			for i, r := range text {
				if ok = phraseMarkers[r]; ok {
					if len(phraseBuff) >= 10 {
						cleaned = clean(phraseBuff)
						if len(cleaned) > 0 {
							fmt.Fprintln(of, prefix+cleaned)
							written++
						}
					}
					phraseBuff = []byte{}
				} else if v := conjPrep(phraseBuff, r); v > 0 {
					// TODO erase or keep? starting with erase.
					phraseBuff = phraseBuff[0 : len(phraseBuff)-v]
					// TODO this pasta is copied
					if len(phraseBuff) >= 10 {
						cleaned = clean(phraseBuff)
						if len(cleaned) > 0 {
							fmt.Fprintln(of, prefix+cleaned)
							written++
						}
					}
					phraseBuff = []byte{}
				} else {
					asStr = string(phraseBuff)
					if r == ' ' && strings.HasSuffix(asStr, " ") {
						continue
					}
					if i == 0 && len(phraseBuff) > 0 && phraseBuff[len(phraseBuff)-1] != ' ' && r != ' ' {
						phraseBuff = append(phraseBuff, byte(' '))
					}
					phraseBuff = append(phraseBuff, byte(r))
				}
			}
		}
		of.Close()
		if written == 0 {
			// there are a bunch of empty books in gutenberg :( these are text files
			// that just have start and end markers with nothing in between. nothing
			// i can do about it.
			fmt.Fprintf(os.Stderr, "WARN: no content found in '%s' '%s'\n", sourceid, p)
		}
		sources <- fmt.Sprintf("%s\t%s", sourceid, title)
	}
}

func CutupFiles() error {
	err := os.Mkdir(tgtDir, 0770)
	if err != nil {
		return err
	}

	dir, err := os.Open(srcDir)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", srcDir, err)
	}
	entries, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", srcDir, err)
	}

	paths := make(chan string, len(entries))
	sources := make(chan string, len(entries))

	for x := 0; x < workers; x++ {
		go worker(paths, sources)
	}

	for _, e := range entries {
		paths <- path.Join(srcDir, e)
	}
	close(paths)

	ixFile, err := os.Create(path.Join(tgtDir, "_title_index.tsv"))
	if err != nil {
		return fmt.Errorf("could not open index file: %w", err)
	}
	defer ixFile.Close()

	for i := 0; i < len(entries); i++ {
		l := <-sources
		fmt.Printf("%d/%d\r", i+1, len(entries))
		fmt.Fprintln(ixFile, l)
	}
	close(sources)

	return nil
}

func conjPrep(phraseBuff []byte, r rune) int {
	if r != ' ' {
		return -1
	}

	suffices := []string{"from", "at", "but", "however", "yet", "though", "and", "to", "on", "or"}
	maxLen := 8 // TODO magic number based on longest suffix
	offset := len(phraseBuff) - maxLen
	if offset < 0 {
		offset = 0
	}
	end := string(phraseBuff[offset:])
	for _, s := range suffices {
		if strings.HasSuffix(end, " "+s) {
			return len(s)
		}
	}
	return -1
}

func isAlpha(r rune) bool {
	// TODO use rune numerical ranges for this
	switch strings.ToLower(string(r)) {
	case "a":
		return true
	case "b":
		return true
	case "c":
		return true
	case "d":
		return true
	case "e":
		return true
	case "f":
		return true
	case "g":
		return true
	case "h":
		return true
	case "i":
		return true
	case "j":
		return true
	case "k":
		return true
	case "l":
		return true
	case "m":
		return true
	case "n":
		return true
	case "o":
		return true
	case "p":
		return true
	case "q":
		return true
	case "r":
		return true
	case "s":
		return true
	case "t":
		return true
	case "u":
		return true
	case "v":
		return true
	case "w":
		return true
	case "x":
		return true
	case "y":
		return true
	case "z":
		return true
	}

	return false
}

func alphaPercent(s string) float64 {
	total := 0.0
	alpha := 0.0

	for _, r := range s {
		total++
		if isAlpha(r) {
			alpha++
		}
	}

	return 100 * (alpha / total)
}

func rep(r rune) (s rune) {
	s = r
	switch s {
	case '’':
		return '\''
	case '“':
		return '"'
	case '”':
		return '"'
	case '"':
		return -1
	case '(':
		return -1
	case '[':
		return -1
	case '{':
		return -1
	case '<':
		return -1
	case '_':
		return -1
	case '*':
		return -1
	case '\r':
		return -1
	case '\t':
		return -1
	case '\n': // should not need this but stray \n ending up in output...
		return -1
	case 0x1c:
		return -1
	case 0x19:
		return -1
	case 0x01:
		return -1
	case 0x0f:
		return -1
	case 0x00:
		return -1
	case 0xb0:
		return -1
	case 0x1b:
		return -1
	case '\\':
		return '/'
	}
	return
}

func clean(bs []byte) string {
	s := strings.ToLower(
		strings.TrimSpace(
			strings.TrimRight(
				strings.TrimLeft(
					strings.Map(rep, strings.ToValidUTF8(string(bs), "")), "'\""), "'\"")))

	if alphaPercent(s) < 50.0 {
		return ""
	}

	return s
}

type CutupOpts struct {
	SrcDir     string
	CutupDir   string
	NumWorkers int
}

func Cutup(opts CutupOpts) error {
	err := os.Mkdir(opts.CutupDir, 0775)
	if err != nil {
		return fmt.Errorf("could not make '%s': %w", opts.CutupDir, err)
	}

	src, err := os.Open(opts.SrcDir)
	if err != nil {
		return fmt.Errorf("could not open '%s' for reading: %w", opts.SrcDir, err)
	}

	entries, err := src.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("could not read '%s': %w", opts.SrcDir, err)
	}

	paths := make(chan string, len(entries))
	sources := make(chan string, len(entries))

	for x := 0; x < opts.NumWorkers; x++ {
		go worker(paths, sources)
	}

	for _, e := range entries {
		paths <- path.Join(opts.SrcDir, e)
	}
	close(paths)

	ixPath := path.Join(opts.CutupDir, "_title_index.csv")
	ixFile, err := os.Create(ixPath)
	if err != nil {
		return fmt.Errorf("could not open '%s': %w", ixPath, err)
	}
	defer ixFile.Close()

	for i := 0; i < len(entries); i++ {
		l := <-sources
		fmt.Printf("%d/%d\r", i+1, len(entries))
		fmt.Fprintln(ixFile, l)
	}
	close(sources)

	return nil
}
