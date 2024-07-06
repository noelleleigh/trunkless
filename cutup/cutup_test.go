package cutup

import "testing"

func Test_shouldBreak(t *testing.T) {
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
		{
			name:     "phrase marker",
			args:     args{[]byte("whither good"), ';'},
			expected: 0,
		},
		// TODO test phrasemarkers
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			result := shouldBreak(c.args.buff, c.args.r)
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
			arg:      "abcdefghijklmnopqrstuvwxyz",
			expected: 100.0,
		},
		{
			arg:      "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
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

func test_shouldSkipLine(t *testing.T) {
	cases := []struct {
		name     string
		arg      string
		expected bool
	}{
		{
			name: "blank",
			arg:  "",
		},
		{
			name: "lol",
			arg:  "lol",
		},
		{
			name:     "head",
			arg:      "head",
			expected: true,
		},
		{
			name:     "HEAD",
			arg:      "HEAD",
			expected: true,
		},
		{
			name:     "/HEAD",
			arg:      "/HEAD",
			expected: false,
		},
		{
			name:     "/head",
			arg:      "/head",
			expected: false,
		},
		{
			name:     "style",
			arg:      "style",
			expected: true,
		},
		{
			name:     "STYLE",
			arg:      "STYLE",
			expected: true,
		},
		{
			name:     "/STYLE",
			arg:      "/STYLE",
			expected: false,
		},
		{
			name:     "/style",
			arg:      "/style",
			expected: false,
		},
		{
			name:     "script",
			arg:      "script",
			expected: true,
		},
		{
			name:     "SCRIPT",
			arg:      "SCRIPT",
			expected: true,
		},
		{
			name:     "/SCRIPT",
			arg:      "/SCRIPT",
			expected: false,
		},
		{
			name:     "/script",
			arg:      "/script",
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.arg, func(t *testing.T) {
			result := shouldSkipLine(c.arg)
			if result != c.expected {
				t.Errorf("got '%v', expected '%v'", result, c.expected)
			}
		})
	}
}
