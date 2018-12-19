package loli2

import (
	"io/ioutil"
	"path/filepath"
)

// Dir represent a lolinote-folder structure in lolinote.
//
// Notice: not all fs-folder was a lolinote-folder. Some fs-folder was
// the "leaf" structure in lolinote data tree, such as ComplexNote,
// Noise dir, sub repo, etc.
type Dir struct {
	path string
}

// Path return dir's path.
func (d Dir) Path() string {
	return d.path
}

// allPaths return all path under this folder with string's order.
func (d Dir) allPaths() []string {
	files, err := ioutil.ReadDir(d.path)
	if err != nil { // file not found should not happen
		panic(err)
	}

	results := make([]string, 0, len(files))
	for _, file := range files {
		results = append(
			results,
			filepath.Join(d.path, file.Name()),
		)
	}

	return results
}

// All return all elem under this folder with string's order.
func (d Dir) All() []Entry {
	paths := d.allPaths()
	entries := make([]Entry, 0, len(paths))

	for _, p := range paths {
		entry, err := buildEntry(p)
		if err != nil {
			panic(err)
		}

		entries = append(entries, entry)
	}

	return entries
}

// Dirs return all Dir under this folder with string's order.
func (d Dir) Dirs() []Dir {
	entries := d.All()
	dirs := make([]Dir, 0, len(entries))

	for _, entry := range entries {
		d, ok := entry.(Dir)
		if ok {
			dirs = append(dirs, d)
		}
	}

	return dirs
}

// Walk can call the walkfn on all entries in all the sub-tree entries
// recursively. but not call walkfn on d itself.
func (d Dir) Walk(walkfn func(Entry)) {
	entries := d.All()
	for _, entry := range entries {
		if subDir, isdir := entry.(Dir); isdir {
			subDir.Walk(walkfn)
		}

		walkfn(entry)
	}
}
