package loli2

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gitlab.com/visig/lolikit-go/loli2/simplerepo"
)

func getMTimeForTest(path string) time.Time {
	if fi, err := os.Stat(path); err != nil {
		panic(err)
	} else {
		return fi.ModTime()
	}
}

func TestIsSimpleNote(t *testing.T) {
	sr := simplerepo.New()
	defer sr.Close()

	for _, path := range sr.SNPaths {
		assert.Equal(
			t, true, IsSimpleNote(path),
			"IsSimpleNote() on simple note path",
		)
	}

	var notSNP []string
	notSNP = append(notSNP, sr.XPaths...)
	notSNP = append(notSNP, sr.CNPaths...)
	notSNP = append(notSNP, sr.DPaths...)

	for _, path := range notSNP {
		assert.Equal(
			t, false, IsSimpleNote(path),
			"IsSimpleNote() on not simple note path",
		)
	}
}

func TestSimpleNote(t *testing.T) {
	cases := []struct {
		p   string // input path
		wp  string // want Path()
		wt  string // want Title()
		wcp string // want ContentPath()
		wct string // want ContentType()
	}{
		{
			"testdata/repo/simple-note.txt",
			"testdata/repo/simple-note.txt",
			"simple-note",
			"testdata/repo/simple-note.txt",
			"txt",
		},
		{
			"testdata/repo/simple-note.2.md",
			"testdata/repo/simple-note.2.md",
			"simple-note.2",
			"testdata/repo/simple-note.2.md",
			"md",
		},
	}

	for _, c := range cases {
		n := &SimpleNote{path: c.p}
		assert.Equal(t, c.wp, n.Path())
		assert.Equal(t, c.wt, n.Title())
		assert.Equal(t, c.wcp, n.ContentPath())
		assert.Equal(t, c.wct, n.ContentType())
		assert.Equal(t, getMTimeForTest(c.wcp), n.MTime())

		njd := noteJSONData{}
		json.Unmarshal(n.JSON(), &njd)
		assert.Equal(t, c.wp, njd.Path)
		assert.Equal(t, c.wt, njd.Title)
		assert.Equal(t, c.wcp, njd.ContentPath)
		assert.Equal(t, c.wct, njd.ContentType)
		assert.Equal(t, getMTimeForTest(c.wcp), njd.ModTime)
	}
}

func TestIsComplexNote(t *testing.T) {
	sr := simplerepo.New()
	defer sr.Close()

	for _, path := range sr.CNPaths {
		assert.Equal(
			t, true, IsComplexNote(path),
			"IsComplexNote() on complex note path",
		)
	}

	var notCNP []string
	notCNP = append(notCNP, sr.XPaths...)
	notCNP = append(notCNP, sr.SNPaths...)
	notCNP = append(notCNP, sr.DPaths...)

	for _, path := range notCNP {
		assert.Equal(
			t, false, IsComplexNote(path),
			"IsComplexNote() on not complex note path",
		)
	}
}

func TestComplexNote(t *testing.T) {
	cases := []struct {
		p   string // input path
		wp  string // want Path()
		wt  string // want Title()
		wcp string // want ContentPath()
		wct string // want ContentType()
	}{
		{
			"testdata/repo/complex-note",
			"testdata/repo/complex-note",
			"complex-note",
			"testdata/repo/complex-note/index.txt",
			"txt",
		},
		{
			"testdata/repo/complex-note-2",
			"testdata/repo/complex-note-2",
			"complex-note-2",
			"testdata/repo/complex-note-2/index.dot",
			"dot",
		},
		{
			"testdata/repo/complex-note.3",
			"testdata/repo/complex-note.3",
			"complex-note.3",
			"testdata/repo/complex-note.3/index.html",
			"html",
		},
	}

	for _, c := range cases {
		n := ComplexNote{path: c.p}
		assert.Equal(t, c.wp, n.Path())
		assert.Equal(t, c.wt, n.Title())
		assert.Equal(t, c.wcp, n.ContentPath())
		assert.Equal(t, c.wct, n.ContentType())
		assert.Equal(t, getMTimeForTest(c.wcp), n.MTime())

		njd := noteJSONData{}
		json.Unmarshal(n.JSON(), &njd)
		assert.Equal(t, c.wp, njd.Path)
		assert.Equal(t, c.wt, njd.Title)
		assert.Equal(t, c.wcp, njd.ContentPath)
		assert.Equal(t, c.wct, njd.ContentType)
		assert.Equal(t, getMTimeForTest(c.wcp), njd.ModTime)
	}
}
