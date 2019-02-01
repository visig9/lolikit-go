package loli2

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/stretchr/testify/assert"
	"github.com/visig9/lolikit-go/fstest"
	"github.com/visig9/lolikit-go/loli2/simplerepo"
)

func TestNotesByMod(t *testing.T) {
	buildExpected := func(
		snp, cnp []string,
		cnip map[string]string,
	) []interface{} {
		getMtime := func(path string) time.Time {
			if fstest.IsDir(path) {
				path = cnip[path]
			}

			fi, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}

			return fi.ModTime()
		}

		list := arraylist.New()
		for _, v := range snp {
			list.Add(v)
		}
		for _, v := range cnp {
			list.Add(v)
		}

		list.Sort(func(a, b interface{}) int {
			ta := getMtime(a.(string))
			tb := getMtime(b.(string))

			return -int(ta.Sub(tb))
		})

		return list.Values()
	}

	toPaths := func(i int, v interface{}) interface{} {
		return v.(Note).Path()
	}

	sr := simplerepo.New()
	defer sr.Close()

	expected := buildExpected(sr.SNPaths, sr.CNPaths, sr.CNIPaths)

	r, err := NewRepo(sr.Path, false)
	if assert.Nil(t, err) {
		actual := r.NotesByMod().Map(toPaths).Values()
		assert.Equal(t, expected, actual)
	}
}

func TestGetNoteRel(t *testing.T) {
	cases := []struct {
		p    string   // path
		t    []string // terms
		wc   bool     // with content
		wrel float64  // want relevance
	}{
		{
			"testdata/repo/note-with-content.md",
			[]string{"note"},
			false,
			4.0 / 17.0,
		},
		{
			"testdata/repo/note-with-content.md",
			[]string{"note"},
			true,
			4.0/17.0 + 4.0/11.0,
		},
	}

	for _, c := range cases {
		actual := getNoteRel(&SimpleNote{path: c.p}, c.t, c.wc)
		assert.Equal(t, c.wrel, actual, fmt.Sprint(c))
	}
}

func TestNotesByRel(t *testing.T) {
	cases := []struct {
		ts  []string      // terms
		tts *hashset.Set  // text types
		wnl []interface{} // want note list
	}{
		{
			[]string{"a"}, hashset.New(),
			[]interface{}{
				&SimpleNote{path: "testdata/rel-repo/apple.md"},
				&SimpleNote{path: "testdata/rel-repo/orange.txt"},
			},
		},
		{
			[]string{"apple"}, hashset.New(),
			[]interface{}{
				&SimpleNote{path: "testdata/rel-repo/apple.md"},
			},
		},
		{
			[]string{"orange"}, hashset.New(),
			[]interface{}{
				&SimpleNote{path: "testdata/rel-repo/orange.txt"},
			},
		},
		{
			[]string{"apple"}, hashset.New("txt", "md"),
			[]interface{}{
				&SimpleNote{path: "testdata/rel-repo/apple.md"},
				&SimpleNote{path: "testdata/rel-repo/orange.txt"},
			},
		},
		{
			[]string{"orange"}, hashset.New("txt", "md"),
			[]interface{}{
				&SimpleNote{path: "testdata/rel-repo/orange.txt"},
			},
		},
	}

	r, err := NewRepo("testdata/rel-repo", false)
	if assert.Nil(t, err) {
		for _, c := range cases {
			actual := r.NotesByRel(c.ts, c.tts)
			assert.Equal(t, c.wnl, actual.Values())
		}
	}
}
