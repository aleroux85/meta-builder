package refmap

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type ExecOp struct {
	Key     string
	Act     string
	Cmd     []string
	TimeOut int
	Rsp     chan ExecRsp
}

type ExecRsp struct {
	More           bool
	Key            string
	StdOut, StdErr string
	Err            error
}

func (o ExecOp) handle(execs map[string]command) {
	if o.Act != "register" && o.Act != "execute" {
		o.Rsp <- ExecRsp{}
		return
	}

	if o.Act == "execute" {
		var key string
		var exc command

		if len(execs) == 0 {
			o.Rsp <- ExecRsp{
				More: len(execs) > 0,
			}
			return
		}

		for key, exc = range execs {
			break
		}

		fmt.Println("executing command", key, exc.Cmd)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(exc.TimeOut)*time.Millisecond)
		defer cancel()

		cmd := exec.CommandContext(ctx, exc.Cmd[0], exc.Cmd[1:]...)

		var stdOut, stdErr bytes.Buffer
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr
		err := cmd.Run()

		o.Rsp <- ExecRsp{
			More:   len(execs) > 1,
			Key:    key,
			StdOut: stdOut.String(),
			StdErr: stdErr.String(),
			Err:    err,
		}

		delete(execs, key)
		return
	}

	if _, found := execs[o.Key]; !found {
		timeOut := 1000
		if o.TimeOut > 0 {
			timeOut = o.TimeOut
		}

		execs[o.Key] = command{
			Cmd:     o.Cmd,
			TimeOut: timeOut,
		}
		fmt.Println("registered command", o.Key, o.Cmd)
	}
	o.Rsp <- ExecRsp{}
}

func (r RefMap) Register(key string, cmd []string, timeOutOpt ...int) {
	timeOut := 0
	if len(timeOutOpt) > 0 {
		timeOut = timeOutOpt[0]
	}

	register := &ExecOp{
		Key:     key,
		Act:     "register",
		Cmd:     cmd,
		TimeOut: timeOut,
		Rsp:     make(chan ExecRsp),
	}
	r.Exec <- register
	<-register.Rsp
}

func (r RefMap) Execute() ExecRsp {
	execute := &ExecOp{
		Act: "execute",
		Rsp: make(chan ExecRsp),
	}
	r.Exec <- execute
	return <-execute.Rsp
}
