package loli2

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets"
	"github.com/visig9/tf/textrel"
)

func (r *Repo) notes() *arraylist.List {
	notes := arraylist.New()

	walkfn := func(entry Entry) {
		if note, isNote := entry.(Note); isNote {
			notes.Add(note)
		}
	}

	r.Dir().Walk(walkfn)

	return notes
}

// NotesByMod return a slice contain all sub-notes ordering
// by modtime (desc).
func (r *Repo) NotesByMod() *arraylist.List {
	MTimeComparator := func(a, b interface{}) int {
		ta := a.(Note).MTime()
		tb := b.(Note).MTime()

		return -int(ta.Sub(tb))
	}

	notes := r.notes()
	notes.Sort(MTimeComparator)

	return notes
}

// getNoteRel calculate the relevance between terms and title.
//
// if withContent is ture, add extra relevance between terms and
// it's content.
func getNoteRel(n Note, terms []string, withContent bool) float64 {
	rel := textrel.ByTermsCI(n.Title(), terms)

	if withContent {
		crel, err := textrel.FileByTerms(
			n.ContentPath(), terms, textrel.CaseInsensitive,
		)
		if err != nil { // file not found should not happen
			panic(err)
		}

		rel += crel
	}

	return rel
}

// NotesByRel return a slice contain all sub-notes ordering
// by relevance.
func (r *Repo) NotesByRel(terms []string, textTypes sets.Set) *arraylist.List {
	type relNote struct {
		rel  float64
		note Note
	}

	rns := r.notes().Map(func(i int, v interface{}) interface{} {
		n := v.(Note)
		withContent := textTypes.Contains(n.ContentType())

		return relNote{
			rel:  getNoteRel(n, terms, withContent),
			note: n,
		}
	})

	filtedRNs := rns.Select(func(i int, v interface{}) bool {
		rn := v.(relNote)
		if rn.rel != 0 {
			return true
		}

		return false
	})

	filtedRNs.Sort(func(a, b interface{}) int {
		switch delta := a.(relNote).rel - b.(relNote).rel; {
		case delta > 0:
			return -1
		case delta < 0:
			return 1
		}

		return 0
	})

	return filtedRNs.Map(func(i int, v interface{}) interface{} {
		rn := v.(relNote)

		return rn.note
	})
}
