// Package simplerepo offer test simple for loli2 project.
package simplerepo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// for build some difference of file's modify time
var createTime = time.Now()

// SimpleRepo is a loli2 repo in filesystem used for testing.
// All attribute should not be modify.
type SimpleRepo struct {
	Path     string            // repo path
	XPaths   []string          // noise paths
	SNPaths  []string          // simple note paths
	CNPaths  []string          // complex note paths
	CNIPaths map[string]string // complex note index paths
	DPaths   []string          // dir paths
}

// Close remove the tmpdir.
func (sr *SimpleRepo) Close() {
	os.RemoveAll(sr.Path)
}

func createFile(path string) {
	if f, err := os.Create(path); err != nil {
		panic(err)
	} else {
		if errC := f.Close(); errC != nil {
			panic(err)
		}
	}

	createTime = createTime.Add(1 * time.Second)
	os.Chtimes(path, createTime, createTime)
}

func createDir(path string) {
	if err := os.MkdirAll(path, 0777); err != nil {
		panic(err)
	}
}

func createNoises(dir string) (noises []string) {
	// bad filename
	fnames := []string{".txt", ".txt.", "txt.", "txt"}
	for _, fname := range fnames {
		fpath := filepath.Join(dir, fname)
		createFile(fpath)
		noises = append(noises, fpath)
	}

	// bad dirname
	dnames := []string{".git", ".orz"}
	for _, dname := range dnames {
		dpath := filepath.Join(dir, dname)
		createDir(dpath)
		noises = append(noises, dpath)
	}

	// bad complex note entry
	badComplexPath := filepath.Join(dir, "not-complex-note")
	createDir(badComplexPath)
	badComplexFilenames := []string{"index.md", "index.txt"}
	for _, fname := range badComplexFilenames {
		fpath := filepath.Join(badComplexPath, fname)
		createFile(fpath)
	}
	noises = append(noises, badComplexPath)

	// sub repo
	subRepoDirNames := []string{"sub-repo-1", "sub-repo-2"}
	for _, srdn := range subRepoDirNames {
		srp := filepath.Join(dir, srdn)
		srcp := filepath.Join(srp, ".lolinote")
		createDir(srcp)
		noises = append(noises, srp)
	}

	return
}

func createDirs(dir string) (dirs []string) {
	dnames := []string{"dir", "dir."}
	for _, dname := range dnames {
		dpath := filepath.Join(dir, dname)
		createDir(dpath)
		dirs = append(dirs, dpath)
	}

	return
}

func createSimpleNotes(dir string) (sns []string) {
	fnames := []string{"simple-note-1.txt", "simple-note.2.txt"}

	for _, fname := range fnames {
		snp := filepath.Join(dir, fname)
		createFile(snp)
		sns = append(sns, snp)
	}

	return
}

func createComplexNotes(dir string) (cns []string, cnis map[string]string) {
	cnis = make(map[string]string)
	data := []struct {
		name     string
		internal []string
	}{
		{"complex-note-1", []string{"index.md", "index.1.txt"}},
		{"complex-note.2", []string{"index.html"}},
		{"complex-note.3.", []string{"index.txt"}},
	}

	for _, d := range data {
		dpath := filepath.Join(dir, d.name)
		createDir(dpath)
		for i, in := range d.internal {
			inpath := filepath.Join(dpath, in)
			if i == 0 {
				cnis[dpath] = inpath
			}
			createFile(inpath)
		}

		cns = append(cns, dpath)
	}

	return
}

// New return a SimpleRepo.
func New() (sr *SimpleRepo) {
	tmpdir, err := ioutil.TempDir("", "test-repo")
	if err != nil {
		panic("can not create tmpdir")
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.RemoveAll(tmpdir)
			panic("SimpleRepo create failed, corrupted data be removed.")
		}
	}()
	kcd := filepath.Join(tmpdir, ".lolinote")
	if err := os.MkdirAll(kcd, 0777); err != nil {
		panic("can not create .lolinote dir")
	}

	sr = new(SimpleRepo)
	sr.XPaths = append(sr.XPaths, kcd)

	sr.Path = tmpdir
	sr.XPaths = append(sr.XPaths, createNoises(tmpdir)...)
	sr.SNPaths = createSimpleNotes(tmpdir)
	sr.CNPaths, sr.CNIPaths = createComplexNotes(tmpdir)
	sr.DPaths = createDirs(tmpdir)

	return
}
