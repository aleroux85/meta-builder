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
	Op      string
	Cmd     []string
	Deps    []string
	TimeOut int
	Rsp     chan ExecRsp
}

type ExecRsp struct {
	Key            string
	StdOut, StdErr string
	Err            error
}

type action struct {
	Cmd     []string
	Deps    []string
	TimeOut int
}

func (act action) exec(stdOut, stdErr *bytes.Buffer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(act.TimeOut)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, act.Cmd[0], act.Cmd[1:]...)

	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	return cmd.Run()
}

func (o ExecOp) handle(actions map[string]action) {
	if o.Op != "register" && o.Op != "execute" {
		o.Rsp <- ExecRsp{}
		return
	}

	if o.Op == "execute" {
		if len(actions) == 0 {
			close(o.Rsp)
			return
		}

		for key := range actions {
			err := o.execute(actions, key)
			if err != nil {
				for key := range actions {
					delete(actions, key)
				}
				break
			}
		}

		close(o.Rsp)
		return
	}

	if _, found := actions[o.Key]; !found {
		timeOut := 1000
		if o.TimeOut > 0 {
			timeOut = o.TimeOut
		}

		actions[o.Key] = action{
			Cmd:     o.Cmd,
			Deps:    o.Deps,
			TimeOut: timeOut,
		}
		fmt.Println("registered command", o.Key, o.Cmd)
	}
	o.Rsp <- ExecRsp{}
}

func (o ExecOp) execute(actions map[string]action, key string) error {
	var stdOut, stdErr bytes.Buffer

	action, found := actions[key]
	if !found {
		return nil
	}

	for _, dep := range action.Deps {
		err := o.execute(actions, dep)
		if err != nil {
			return err
		}
	}

	fmt.Println("executing command", key, action.Cmd)
	err := action.exec(&stdOut, &stdErr)
	if err != nil {
		o.Rsp <- ExecRsp{
			Key:    key,
			StdOut: stdOut.String(),
			StdErr: stdErr.String(),
			Err:    err,
		}
		return err
	}

	o.Rsp <- ExecRsp{
		Key:    key,
		StdOut: stdOut.String(),
		StdErr: stdErr.String(),
		Err:    nil,
	}

	delete(actions, key)
	return nil
}

func (r RefMap) Register(key string, cmd, deps []string, timeOutOpt ...int) {
	timeOut := 0
	if len(timeOutOpt) > 0 {
		timeOut = timeOutOpt[0]
	}

	register := &ExecOp{
		Key:     key,
		Op:      "register",
		Cmd:     cmd,
		Deps:    deps,
		TimeOut: timeOut,
		Rsp:     make(chan ExecRsp),
	}
	r.Exec <- register
	<-register.Rsp
}

func (r RefMap) Execute() []ExecRsp {
	execute := &ExecOp{
		Op:  "execute",
		Rsp: make(chan ExecRsp),
	}
	r.Exec <- execute
	var responses []ExecRsp
	for response := range execute.Rsp {
		responses = append(responses, response)
	}

	return responses
}
