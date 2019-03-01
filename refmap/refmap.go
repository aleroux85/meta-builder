package refmap

import (
	"fmt"
	"path/filepath"

	"github.com/aleroux85/utils"
)

type RefMap struct {
	Reads   chan *ReadOp
	Writes  chan *WriteOp
	Sets    chan *SetOp
	Sync    chan *SyncOp
	Changed chan *ChangedOp
}

type RefVal interface {
	Build(config)
	SetChange(...uint) uint
	GetHash() string
}

type ReadOp struct {
	Src string
	Rsp chan *RefLink
}

func (r RefMap) DoRead(src string) *RefLink {
	read := &ReadOp{
		Src: src,
		Rsp: make(chan *RefLink),
	}
	r.Reads <- read
	return <-read.Rsp
}

func (r RefMap) Set(key, val string) {
	setter := &SetOp{
		Key: key,
		Val: val,
		Err: make(chan error),
	}
	r.Sets <- setter
	<-setter.Err
}

func (r RefMap) Finish() {
	setter := &SetOp{
		Key: "finish",
		Err: make(chan error),
	}
	r.Sets <- setter
	<-setter.Err
}

func (r RefMap) Assess() {
	setter := &SetOp{
		Key: "assess",
		Err: make(chan error),
	}
	r.Sets <- setter
	<-setter.Err
}

type ChangedOp struct {
	Refs chan *RefLink
}

func (r RefMap) ChangedRefs() chan *RefLink {
	changed := &ChangedOp{
		Refs: make(chan *RefLink),
	}
	r.Changed <- changed
	return changed.Refs
}

type WriteOp struct {
	src string
	dst string
	val RefVal
	err chan error
}

type SetOp struct {
	Key string
	Val string
	Err chan error
}

type SyncOp struct {
	Mon *utils.Monitor
	Err chan error
}

func NewRefMap() *RefMap {
	m := new(RefMap)
	m.Reads = make(chan *ReadOp)
	m.Writes = make(chan *WriteOp)
	m.Sets = make(chan *SetOp)
	m.Sync = make(chan *SyncOp)
	m.Changed = make(chan *ChangedOp)
	return m
}

func (m *RefMap) Start() {
	go func() {
		refs := make(map[string]*RefLink)
		sourceStartLocation := ""

		for {
			select {
			case read := <-m.Reads:
				read.Rsp <- refs[read.Src]
			case write := <-m.Writes:
				source := filepath.Join(sourceStartLocation, write.src)
				if _, found := refs[source]; !found {
					refs[source] = NewRefLink()
					fmt.Printf("adding %s\n", source)
				} else {
					if refs[source].Change == DataStable {
						refs[source].Change = DataFlagged
					}
				}
				if _, found := refs[source].Files[write.dst]; !found {
					refs[source].Files[write.dst] = write.val
					refs[source].Files[write.dst].SetChange(DataAdded)
					fmt.Printf("adding %s\t-> %s\n", source, write.dst)
				} else {
					if refs[source].Files[write.dst].GetHash() == write.val.GetHash() {
						write.val.SetChange(DataFlagged)
					} else {
						write.val.SetChange(DataUpdated)
						fmt.Printf("updating %s\t-> %s\n", source, write.dst)
					}
					refs[source].Files[write.dst] = write.val
				}
				write.err <- nil
			case setter := <-m.Sets:
				if setter.Key == "sourceStartLocation" {
					sourceStartLocation = setter.Val
					setter.Err <- nil
					continue
				}
				if setter.Key == "assess" {
					for _, ref := range refs {
						if ref.Change == DataStable {
							ref.Change = DataRemove
							continue
						}
						for _, file := range ref.Files {
							if file.SetChange() == DataStable {
								file.SetChange(DataRemove)
							}
						}
					}
					setter.Err <- nil
					continue
				}
				if setter.Key == "finish" {
					for src, ref := range refs {
						if ref.Change == DataRemove {
							fmt.Println("removing", src, "from refmap")
							delete(refs, src)
							continue
						}
						if ref.Change != DataFlagged {
							fmt.Println("setting", src, ref.Change, "stable")
						}
						ref.Change = DataStable
						for dst, file := range ref.Files {
							if file.SetChange() == DataRemove {
								fmt.Println("removing", src, ":", dst, "from refmap")
								delete(ref.Files, dst)
								continue
							}
							if ref.Change != DataFlagged {
								fmt.Println("setting", src, ":", dst, "stable")
							}
							file.SetChange(DataStable)
						}
					}
					setter.Err <- nil
					continue
				}
				var value uint
				if setter.Val == "update" {
					value = DataUpdated
				} else {
					setter.Err <- fmt.Errorf("unknown value")
					continue
				}
				if ref, found := refs[setter.Key]; found {
					ref.Change = value
					fmt.Println("setting", setter.Key, setter.Val)
					setter.Err <- nil
				}
			case sync := <-m.Sync:
				if sync.Mon.Watcher == nil {
					sync.Err <- fmt.Errorf("monitor has nil watcher")
					continue
				}
				for source, ref := range refs {
					if ref.Change == DataAdded {
						sync.Mon.Watcher.Add(source)
						fmt.Println("watching", source)
					} else if ref.Change == DataRemove {
						sync.Mon.Watcher.Remove(source)
						fmt.Println("un-watching", source)
					}
				}
				sync.Err <- nil
			case changed := <-m.Changed:
				for src, ref := range refs {
					if ref.Change == DataUpdated || ref.Change == DataAdded {
						fmt.Println("returning changed", src)
						changed.Refs <- ref
					}
				}
				close(changed.Refs)
			}
		}
	}()
}

type RefLink struct {
	Files  map[string]RefVal
	Change uint
}

func NewRefLink() *RefLink {
	r := new(RefLink)
	r.Files = make(map[string]RefVal)
	r.Change = DataAdded
	return r
}
