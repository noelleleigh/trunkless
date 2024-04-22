package cutup

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

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

func Cutup(ins io.Reader) {
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

	// I want to experiment with treating prepositions and conjunctions as phrase
	// markers.

	// to do this i would need to check the phraseBuff when I check phraseMarkers and then split accordingly

	s := bufio.NewScanner(ins)
	phraseBuff := []byte{}
	printed := false
	for s.Scan() {
		text := strings.TrimSpace(s.Text())
		for i, r := range text {
			if ok := phraseMarkers[r]; ok {
				if len(phraseBuff) >= 10 {
					cleaned := clean(phraseBuff)
					if len(cleaned) > 0 {
						fmt.Println(cleaned)
						printed = true
					}
				}
				if !printed {
					//fmt.Fprintf(os.Stderr, "SKIP: %s\n", string(phraseBuff))
				}
				printed = false
				phraseBuff = []byte{}
			} else if v := conjPrep(phraseBuff, r); v > 0 {
				// TODO erase or keep? starting with erase.
				phraseBuff = phraseBuff[0 : len(phraseBuff)-v]
				// TODO this pasta is copied
				if len(phraseBuff) >= 10 {
					cleaned := clean(phraseBuff)
					if len(cleaned) > 0 {
						fmt.Println(cleaned)
						printed = true
					}
				}
				if !printed {
					//fmt.Fprintf(os.Stderr, "SKIP: %s\n", string(phraseBuff))
				}
				printed = false
				phraseBuff = []byte{}
			} else {
				asStr := string(phraseBuff)
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
}

func isAlpha(r rune) bool {
	alphaChars := map[rune]bool{
		'a': true,
		'b': true,
		'c': true,
		'd': true,
		'e': true,
		'f': true,
		'g': true,
		'h': true,
		'i': true,
		'j': true,
		'k': true,
		'l': true,
		'm': true,
		'n': true,
		'o': true,
		'p': true,
		'q': true,
		'r': true,
		's': true,
		't': true,
		'u': true,
		'v': true,
		'w': true,
		'x': true,
		'y': true,
		'z': true,
	}
	lookup := strings.ToLower(string(r))
	return alphaChars[rune(lookup[0])]
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

func clean(bs []byte) string {
	s := string(bs)
	s = strings.ReplaceAll(s, "’", "'")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.TrimLeft(s, "'\"")
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if alphaPercent(s) < 50.0 {
		return ""
	}

	return s
}
