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
		fmt.Printf("adding %s (status %s)\n", source, statusText[refs[source].Change])
	} else {
		if refs[source].Change == DataStable {
			refs[source].Change = DataFlagged
		}
	}

	if _, found := refs[source].Files[o.dst]; !found {
		refs[source].Files[o.dst] = o.val
		refs[source].Files[o.dst].Change(DataAdded)
		fmt.Printf("adding %s -> %s (status %s)\n", source, o.dst, statusText[refs[source].Files[o.dst].Change()])
	} else {
		if refs[source].Files[o.dst].Change() == DataAdded {
			o.val.Change(DataAdded)
			fmt.Printf("WARNING - duplicate %s -> %s added, over-writing previous entry\n", source, o.dst)
		} else {
			if refs[source].Files[o.dst].Hash() == o.val.Hash() {
				o.val.Change(DataFlagged)
			} else {
				o.val.Change(DataUpdated)
				fmt.Printf("updating %s -> %s to status %s\n", source, o.dst, statusText[refs[source].Files[o.dst].Change()])
			}
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

var statusText map[uint8]string = map[uint8]string{
	DataStable:  "DataStable",
	DataFlagged: "DataFlagged",
	DataUpdated: "DataUpdated",
	DataAdded:   "DataAdded",
	DataRemove:  "DataRemove",
}
