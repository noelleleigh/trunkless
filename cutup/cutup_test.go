package cutup

import "testing"

func Test_conjPrep(t *testing.T) {
	type args struct {
		buff []byte
		r    rune
	}
	cs := []struct {
		name     string
		args     args
		expected int
	}{
		{
			name:     "empty buffer",
			args:     args{[]byte(""), ' '},
			expected: -1,
		},
		{
			name:     "not a space yet",
			args:     args{[]byte("saccharine juice from"), 'x'},
			expected: -1,
		},
		{
			name:     "from",
			args:     args{[]byte("i will eat from"), ' '},
			expected: 4,
		},
		{
			name:     "no preceding space",
			args:     args{[]byte("wakkabarblurpfrom"), ' '},
			expected: -1,
		},
		{
			name:     "however",
			args:     args{[]byte("my eyes are hollow, however"), ' '},
			expected: 7,
		},
		{
			name:     "at",
			args:     args{[]byte("there will be no more joy at"), ' '},
			expected: 2,
		},
		{
			name:     "but",
			args:     args{[]byte("i buried him, but"), ' '},
			expected: 3,
		},
		{
			name:     "yet",
			args:     args{[]byte("the echoes quited yet"), ' '},
			expected: 3,
		},
		{
			name:     "though",
			args:     args{[]byte("my eyes were closed though"), ' '},
			expected: 6,
		},
		{
			name:     "and",
			args:     args{[]byte("i raised the torch and"), ' '},
			expected: 3,
		},
		{
			name:     "to",
			args:     args{[]byte("thousands more to"), ' '},
			expected: 2,
		},
		{
			name:     "on",
			args:     args{[]byte("bringing rain down on"), ' '},
			expected: 2,
		},
		{
			name:     "no match",
			args:     args{[]byte("i raised the torch"), ' '},
			expected: -1,
		},
		{
			name:     "or",
			args:     args{[]byte("whether good or"), ' '},
			expected: 2,
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			result := conjPrep(c.args.buff, c.args.r)
			if result != c.expected {
				t.Errorf("got '%v', expected '%v'", result, c.expected)
			}
		})
	}
}

func Test_isAlpha(t *testing.T) {
	cs := []struct {
		arg      rune
		expected bool
	}{
		{arg: 'a', expected: true},
		{arg: 'b', expected: true},
		{arg: 'c', expected: true},
		{arg: 'd', expected: true},
		{arg: 'e', expected: true},
		{arg: 'f', expected: true},
		{arg: 'g', expected: true},
		{arg: 'h', expected: true},
		{arg: 'i', expected: true},
		{arg: 'j', expected: true},
		{arg: 'k', expected: true},
		{arg: 'l', expected: true},
		{arg: 'm', expected: true},
		{arg: 'n', expected: true},
		{arg: 'o', expected: true},
		{arg: 'p', expected: true},
		{arg: 'q', expected: true},
		{arg: 'r', expected: true},
		{arg: 's', expected: true},
		{arg: 't', expected: true},
		{arg: 'u', expected: true},
		{arg: 'v', expected: true},
		{arg: 'w', expected: true},
		{arg: 'x', expected: true},
		{arg: 'y', expected: true},
		{arg: 'z', expected: true},
		{arg: '1'},
		{arg: '2'},
		{arg: '3'},
		{arg: '\''},
		{arg: '"'},
		{arg: '#'},
		{arg: '%'},
	}

	for _, c := range cs {
		t.Run(string(c.arg), func(t *testing.T) {
			result := isAlpha(c.arg)
			if result != c.expected {
				t.Errorf("got '%v', expected '%v'", result, c.expected)
			}
		})
	}
}

func Test_alphaPercent(t *testing.T) {
	cs := []struct {
		arg      string
		expected float64
	}{
		{
			arg:      "abcd",
			expected: 100.0,
		},
		{
			arg:      "a1b2c3d4",
			expected: 50.0,
		},
		{
			arg:      "--------",
			expected: 0.0,
		},
	}

	for _, c := range cs {
		t.Run(c.arg, func(t *testing.T) {
			result := alphaPercent(c.arg)
			if result != c.expected {
				t.Errorf("got '%v', expected '%v'", result, c.expected)
			}
		})
	}
}

func Test_clean(t *testing.T) {
	cs := []struct {
		name     string
		arg      string
		expected string
	}{
		{
			name:     "all whitespace rejected",
			arg:      "    		 ",
			expected: "",
		},
		{
			name:     "trimmed",
			arg:      " cats eat fish    ",
			expected: "cats eat fish",
		},
		{
			name:     "dquotes removed",
			arg:      "cats \"eat\" fish",
			expected: "cats eat fish",
		},
		{
			name:     "leading quote removed",
			arg:      "'cats eat fish",
			expected: "cats eat fish",
		},
		{
			name:     "leading double quote removed",
			arg:      "\"cats eat fish",
			expected: "cats eat fish",
		},
		{
			name:     "lowered",
			arg:      "Cats Eat Fish",
			expected: "cats eat fish",
		},
		{
			name:     "dumb quote replaced",
			arg:      "cat’s eaten fish",
			expected: "cat's eaten fish",
		},
		{
			name:     "rejects low alphabetic content",
			arg:      "----- --- -a- ---a-dsbbca---asd--",
			expected: "",
		},
	}

	for _, c := range cs {
		t.Run(c.arg, func(t *testing.T) {
			result := clean([]byte(c.arg))
			if result != c.expected {
				t.Errorf("got '%v', expected '%v'", result, c.expected)
			}
		})
	}
}
