package refmap

import "fmt"

type ChangedOp struct {
	Refs chan *RefLink
}

func (o ChangedOp) handle(refs map[string]*RefLink) {
	for src, ref := range refs {
		if ref.Change == DataUpdated || ref.Change == DataAdded {
			fmt.Println("returning changed", src)
			o.Refs <- ref
		}
	}
	close(o.Refs)
}

func (r RefMap) ChangedRefsChan() chan *RefLink {
	changed := &ChangedOp{
		Refs: make(chan *RefLink),
	}
	r.Changed <- changed
	return changed.Refs
}

func (r RefMap) ChangedRefs() (refs []*RefLink) {
	for ref := range r.ChangedRefsChan() {
		refs = append(refs, ref)
	}
	return refs
}
