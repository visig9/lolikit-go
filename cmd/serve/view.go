package serve

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	packr "github.com/gobuffalo/packr/v2"
	"gitlab.com/visig/lolikit-go/loli2"
)

var tplbox = packr.New("templates", "./templates")

func readfile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func getTemplate() *template.Template {
	tpl := template.New("").Funcs(
		template.FuncMap{
			"readfile": readfile,
		},
	)

	err := tplbox.Walk(func(path string, f packr.File) error {
		name := filepath.Base(path)

		if !strings.HasPrefix(name, ".") {
			tpltext, err := tplbox.FindString(path)
			if err != nil {
				panic(err)
			}
			tpl = template.Must(tpl.New(name).Parse(tpltext))
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return tpl
}

type viewHandler struct {
	repo loli2.Repo
	tpl  *template.Template
}

func newViewHandler(repo loli2.Repo) *viewHandler {
	return &viewHandler{
		repo: repo,
		tpl:  getTemplate(),
	}
}

func (h *viewHandler) getEntry(urlpath string) (loli2.Entry, error) {
	fspath := filepath.Join(
		h.repo.Path(),
		filepath.FromSlash(urlpath),
	)

	if visible, _ := isVisible(h.repo.Path(), fspath); !visible {
		return nil, fmt.Errorf("%v are invisible", fspath)
	}

	return h.repo.Dir().Get(fspath)
}

func (h *viewHandler) viewContent(w http.ResponseWriter, n loli2.Note) {
	relpath, err := filepath.Rel(h.repo.Path(), n.Path())
	if err != nil {
		panic(err)
	}
	reldir := filepath.Dir(relpath)

	data := struct {
		Note   loli2.Note
		Reldir string
	}{
		Note:   n,
		Reldir: reldir,
	}
	h.tpl.ExecuteTemplate(w, "view_content.html", data)
}

func (h *viewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	entry, err := h.getEntry(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if n, ok := entry.(loli2.Note); ok {
		h.viewContent(w, n)
	}
}
