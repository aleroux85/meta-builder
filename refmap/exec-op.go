package refmap

import (
	"fmt"
	"os/exec"
)

type ExecOp struct {
	Key string
	Cmd *exec.Cmd
	Err chan error
}

func (o ExecOp) handle(execs map[string]*exec.Cmd) {
	if _, found := execs[o.Key]; !found {
		execs[o.Key] = o.Cmd
		fmt.Println("registring command", o.Key, o.Cmd)
		o.Err <- nil
	}
}

func (r RefMap) Register(key string, cmd *exec.Cmd) {
	pusher := &ExecOp{
		Key: key,
		Cmd: cmd,
		Err: make(chan error),
	}
	r.Exec <- pusher
	<-pusher.Err
}
