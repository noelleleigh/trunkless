package cutup

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/vilmibm/trunkless/db"
)

type CutupOpts struct {
	SrcDir           string
	CutupDir         string
	NumWorkers       int
	Flavor           string
	headerEndCheck   func(string) bool
	footerBeginCheck func(string) bool
}

func defaultHeaderEndCheck(string) bool   { return true }
func defaultFooterBeginCheck(string) bool { return false }

func gutenbergHeaderEndCheck(s string) bool {
	return strings.HasPrefix(s, "*** START")
}

func gutenbergFooterBeginCheck(s string) bool {
	return strings.HasPrefix(s, "*** END")
}

func extractGutenbergTitle(s string) string {
	title, _ := strings.CutPrefix(s, "*** START OF THE PROJECT GUTENBERG")
	title, _ = strings.CutPrefix(title, " EBOOK")
	return strings.TrimSpace(strings.Map(rep, title))
}

func Cutup(opts CutupOpts) error {
	switch opts.Flavor {
	case "gutenberg":
		opts.headerEndCheck = gutenbergHeaderEndCheck
		opts.footerBeginCheck = gutenbergFooterBeginCheck
	default:
		opts.headerEndCheck = defaultHeaderEndCheck
		opts.footerBeginCheck = defaultFooterBeginCheck
	}

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
		go worker(opts, paths, sources)
	}

	for _, e := range entries {
		paths <- path.Join(opts.SrcDir, e)
	}
	close(paths)

	ixPath := path.Join(opts.CutupDir, "_title_index.tsv")
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

func worker(opts CutupOpts, paths <-chan string, sources chan<- string) {
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
		var asStr string
		var text string
		var prefix string

		// geocities
		var inTag bool
		var tagSkip bool
		tagBuff := []byte{}

		for s.Scan() {
			text = strings.TrimSpace(s.Text())
			if inHeader && opts.headerEndCheck(text) {
				inHeader = false
				if opts.Flavor == "gutenberg" {
					title = extractGutenbergTitle(text)
					continue
				} else {
					base := path.Base(p)
					title = strings.TrimSuffix(base, filepath.Ext(base))
				}
			}
			if inHeader {
				continue
			}
			if opts.footerBeginCheck(text) {
				break
			}
			if title == "" {
				fmt.Fprintf(os.Stderr, "got to cutup phase with no title: '%s'", p)
				break
			}
			if sourceid == "" {
				sourceid = db.StrToID(title)
				fmt.Fprintln(os.Stderr, sourceid, p)
				prefix = sourceid + "\t"
				of, err = os.Create(path.Join(opts.CutupDir, sourceid))
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not open '%s' for writing: %s", sourceid, err.Error())
					break
				}
			}
			for i, r := range text {
				if opts.Flavor == "geocities" {
					if r == '<' {
						inTag = true
						continue
					} else if r == '>' {
						tagSkip = shouldSkipLine(string(tagBuff))
						inTag = false
						tagBuff = []byte{}
					}
					if inTag {
						tagBuff = append(tagBuff, byte(r))
						continue
					}
					if tagSkip {
						continue
					}
				}
				if v := shouldBreak(phraseBuff, r); v >= 0 {
					if len(phraseBuff) > 0 {
						phraseBuff = phraseBuff[0 : len(phraseBuff)-v]
					}
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

			// in gfaqs I got a few no content files; they have all of their content
			// on one line with a bunch of special characters. it's a pathological
			// case and i'm shocked more doesn't break but somehow in this printf
			// sourceid renders as '' and it doesn't end up in the title index.
			fmt.Fprintf(os.Stderr, "WARN: no content found in '%s' '%s'\n", sourceid, p)
		}
		// I would leave empty stuff out of sources but we need the number of
		// things written to sources to match the number of initial entires. this
		// means some sources in the DB that will never be used but that's fine.
		sources <- fmt.Sprintf("%s\t%s", sourceid, title)
	}
}

var phraseMarkers = map[rune]bool{
	';': true,
	',': true,
	':': true,
	'.': true,
	'?': true,
	'!': true,
	')': true,
	'}': true,
	']': true,
	'”': true,
	'=': true,
	'`': true,
	'-': true,
	'|': true,
	'>': true,
}

var suffices = []string{"from", "at", "but", "however", "yet", "though", "and", "to", "on", "or"}

const maxSuffixLen = 8 // magic number based on longest suffix

func shouldBreak(phraseBuff []byte, r rune) int {
	if ok := phraseMarkers[r]; ok {
		return 0
	}

	if r != ' ' {
		return -1
	}

	offset := len(phraseBuff) - maxSuffixLen
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

func alphaPercent(s string) float64 {
	total := 0.0
	alpha := 0.0

	for _, r := range s {
		total++
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
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

var ignoreTags = []string{
	"head",
	"script",
	"style",
}

func shouldSkipLine(tagBuff string) bool {
	var s string
	for _, t := range ignoreTags {
		s = strings.ToLower(tagBuff)
		if strings.Contains(s, "/"+t) {
			return false
		}
		if strings.Contains(s, t) {
			return true
		}
	}
	return false
}
