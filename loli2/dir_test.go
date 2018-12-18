package loli2

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/visig/lolikit-go/loli2/simplerepo"
)

func entriesToPaths(entries interface{}) []string {
	esv := reflect.ValueOf(entries)
	paths := make([]string, 0, esv.Len())

	for i := 0; i < esv.Len(); i++ {
		path := esv.Index(i).MethodByName("Path").Call(
			[]reflect.Value{},
		)[0].String()
		paths = append(paths, path)
	}

	return paths
}

func TestDirAllPaths(t *testing.T) {
	sr := simplerepo.New()
	defer sr.Close()

	expected := make([]string, 0)
	expected = append(expected, sr.SNPaths...)
	expected = append(expected, sr.CNPaths...)
	expected = append(expected, sr.XPaths...)
	expected = append(expected, sr.DPaths...)

	d := Dir{sr.Path}
	actual := d.allPaths()

	assert.ElementsMatch(t, expected, actual)
}

func TestDirAll(t *testing.T) {
	sr := simplerepo.New()
	defer sr.Close()

	expected := make([]string, 0)
	expected = append(expected, sr.SNPaths...)
	expected = append(expected, sr.CNPaths...)
	expected = append(expected, sr.XPaths...)
	expected = append(expected, sr.DPaths...)

	d := Dir{sr.Path}
	actual := entriesToPaths(d.All())
	assert.ElementsMatch(t, expected, actual)
}

func TestDirDirs(t *testing.T) {
	sr := simplerepo.New()
	defer sr.Close()

	expected := make([]string, 0)
	expected = append(expected, sr.DPaths...)

	d := Dir{sr.Path}
	actual := entriesToPaths(d.Dirs())
	assert.ElementsMatch(t, expected, actual)
}

func TestDirWalk(t *testing.T) {
	getWalkfn := func() (func(Entry), *[]string) {
		actual := make([]string, 0)

		return func(e Entry) {
			actual = append(actual, e.Path())
		}, &actual
	}

	sr := simplerepo.New()
	defer sr.Close()

	expected := make([]string, 0)
	expected = append(expected, sr.SNPaths...)
	expected = append(expected, sr.CNPaths...)
	expected = append(expected, sr.XPaths...)
	expected = append(expected, sr.DPaths...)

	walkfn, actualPtr := getWalkfn()

	d := Dir{sr.Path}
	d.Walk(walkfn)
	assert.ElementsMatch(t, expected, *actualPtr)
}
