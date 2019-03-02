package refmap

import (
	"fmt"
	"path/filepath"
)

type WriteOp struct {
	src string
	dst string
	val RefVal
}

func (o WriteOp) handle(location string, refs map[string]*RefLink) {
	source := filepath.Join(location, o.src)

	if _, found := refs[source]; !found {
		refs[source] = NewRefLink()
		fmt.Printf("adding %s\n", source)
	} else {
		if refs[source].Change == DataStable {
			refs[source].Change = DataFlagged
		}
	}

	if _, found := refs[source].Files[o.dst]; !found {
		refs[source].Files[o.dst] = o.val
		refs[source].Files[o.dst].SetChange(DataAdded)
		fmt.Printf("adding %s\t-> %s\n", source, o.dst)
	} else {
		if refs[source].Files[o.dst].SetChange() == DataAdded {
			return
		}
		if refs[source].Files[o.dst].GetHash() == o.val.GetHash() {
			o.val.SetChange(DataFlagged)
		} else {
			o.val.SetChange(DataUpdated)
			fmt.Printf("updating %s\t-> %s\n", source, o.dst)
		}
		refs[source].Files[o.dst] = o.val
	}
}

func (r RefMap) Write(src, dst string, val RefVal) {
	write := &WriteOp{
		src: src,
		dst: dst,
		val: val,
	}
	r.Writes <- write
}
