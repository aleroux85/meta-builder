package refmap

import "fmt"

type SetOp struct {
	Key string
	Val string
	Err chan error
}

func (o SetOp) handle(location *string, refs map[string]*RefLink) {
	if o.Key == "location" {
		*location = o.Val
		o.Err <- nil
		return
	}
	if o.Key == "assess" {
		assess(refs)
		o.Err <- nil
		return
	}
	if o.Key == "finish" {
		finish(refs)
		o.Err <- nil
		return
	}

	var value uint
	if o.Val == "update" {
		value = DataUpdated
	} else {
		o.Err <- fmt.Errorf("unknown value")
		return
	}
	if ref, found := refs[o.Key]; found {
		ref.Change = value
		fmt.Println("setting", o.Key, o.Val)
		o.Err <- nil
	}
}

func assess(refs map[string]*RefLink) {
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
}

func finish(refs map[string]*RefLink) {
	for src, ref := range refs {
		if ref.Change == DataRemove {
			fmt.Println("removing", src, "from refmap")
			delete(refs, src)
			continue
		}

		if ref.Change != DataFlagged && ref.Change != DataStable {
			fmt.Println("setting", src, ref.Change, "-> stable")
		}
		ref.Change = DataStable

		for dst, file := range ref.Files {
			if file.SetChange() == DataRemove {
				fmt.Println("removing", src, ":", dst, "from refmap")
				delete(ref.Files, dst)
				continue
			}
			if file.SetChange() != DataFlagged && file.SetChange() != DataStable {
				fmt.Println("setting", src, ":", dst, file.SetChange(), "-> stable")
			}
			file.SetChange(DataStable)
		}
	}
}

func (r RefMap) Set(key, val string) error {
	setter := &SetOp{
		Key: key,
		Val: val,
		Err: make(chan error),
	}
	r.Sets <- setter
	return <-setter.Err
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
