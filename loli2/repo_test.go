package loli2

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visig9/lolikit-go/fstest"
)

func isRepoFS(path string) bool {
	cf := filepath.Join(path, ".lolinote")

	return fstest.IsDir(cf)
}

// getSampleDirs create some tmpdir for test due to those dir may be changed.
func getSampleDirs(t *testing.T) (dirs map[string]string, recycle func()) {
	createFS := func(path string) {
		kcd := filepath.Join(path, ".lolinote", "lolikit")
		if err := os.MkdirAll(kcd, 0777); err != nil {
			t.Fatal("can not create tmpdir")
		}
	}

	tmpdir, err := ioutil.TempDir("", "test-repo")
	if err != nil {
		t.Fatal("can not create tmpdir")
	}
	recycle = func() { os.RemoveAll(tmpdir) }
	dirs = make(map[string]string)

	// exist dir
	dirs["ed"] = filepath.Join(tmpdir, "exists-dir")
	os.Mkdir(dirs["ed"], 0777)

	// not exist dir
	dirs["ned"] = filepath.Join(tmpdir, "not-exists-dir")

	// lolinote dir
	dirs["ld"] = filepath.Join(tmpdir, "loli-dir")
	createFS(dirs["ld"])

	// lolinote's sub dir
	dirs["lsd"] = filepath.Join(dirs["ld"], "sub")
	os.Mkdir(dirs["lsd"], 0777)

	return
}

func TestRepo(t *testing.T) {
	d, recycle := getSampleDirs(t)
	defer recycle()

	d2, recycle := getSampleDirs(t)
	defer recycle()

	cases := []struct {
		p      string // input path
		create bool   // create?
		werr   bool   // wanted return an error?
	}{
		{d["ed"], false, true},
		{d["ned"], false, true},
		{d["ld"], false, false},
		{d["lsd"], false, false},

		{d2["ed"], true, false},
		{d2["ned"], true, false},
		{d2["ld"], true, false},
		{d2["lsd"], true, false},
	}

	for _, c := range cases {
		r, err := NewRepo(c.p, c.create)

		if c.werr {
			assert.NotNil(t, err)
		} else {
			if assert.Nil(t, err) {
				assert.DirExists(
					t,
					filepath.Join(
						r.Path(), ".lolinote",
					),
					"repo not init correctly",
				)
			}
		}
	}
}

func TestRepoDir(t *testing.T) {
	cases := []struct {
		p   string // assign path
		wdp string // want repo.Dir().Path()
	}{
		{"testdata/repo", "testdata/repo"},
		{"testdata/repo/dir", "testdata/repo"},
		{"testdata/repo/dir/sub-dir", "testdata/repo"},
		{"testdata/repo/.lolinote", "testdata/repo"},
		{"testdata/repo/complex-note", "testdata/repo"},
		{"testdata/repo/sub-repo", "testdata/repo/sub-repo"},
	}

	for _, c := range cases {
		r, _ := NewRepo(c.p, false)

		assert.Equal(
			t, c.wdp, r.Dir().Path(),
			"repo.Dir().Path() not as expected",
		)
	}
}
