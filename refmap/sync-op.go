package refmap

import (
	"fmt"
)

type SyncOp struct {
	Mon *Monitor
	Err chan error
}

func (o SyncOp) handle(refs map[string]*RefLink) {
	if o.Mon.Watcher == nil {
		o.Err <- fmt.Errorf("monitor has nil watcher")
		return
	}

	for source, ref := range refs {
		if ref.Change == DataAdded {
			o.Mon.Watcher.Add(source)
			fmt.Println("watching", source)
		} else if ref.Change == DataRemove {
			o.Mon.Watcher.Remove(source)
			fmt.Println("un-watching", source)
		}
	}
	o.Err <- nil
}

func (r RefMap) Sync(mon *Monitor) error {
	sync := &SyncOp{
		Mon: mon,
		Err: make(chan error),
	}
	r.sync <- sync
	return <-sync.Err
}
