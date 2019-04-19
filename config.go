package builder

import (
	"fmt"
	"path/filepath"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/pkg/errors"
)

type Config struct {
	project     ProjectLoader
	source      string
	destination string
	force       bool
	metaFile    string
	refMap      *refmap.RefMap
	// tmplMon     *utils.Monitor
	// cnfgMon     *utils.Monitor
	err error
}

func NewConfig(l ...string) *Config {
	c := new(Config)

	if len(l) > 0 {
		c.source = l[0]
	}

	if len(l) > 1 {
		c.destination = l[1]
	}

	c.refMap = refmap.NewRefMap()
	c.refMap.Start()
	c.refMap.Set("location", c.source)

	// c.tmplMon = new(utils.Monitor)
	// c.tmplMon.Error = &c.err

	// c.cnfgMon = new(utils.Monitor)
	// c.cnfgMon.Error = &c.err
	return c
}

func (c *Config) Source(src ...string) string {
	if len(src) > 0 {
		c.source = src[0]
	}
	return c.source
}

func (c *Config) Destination(dst ...string) string {
	if len(dst) > 0 {
		c.destination = dst[0]
	}
	return c.destination
}

func (c *Config) Force(f ...bool) bool {
	if len(f) > 0 {
		c.force = f[0]
		return c.force
	}

	if c.force {
		c.force = false
		return true
	}
	return false
}

func (c *Config) Error(err ...error) error {
	if c.err != nil {
		return c.err
	}

	if len(err) > 0 {
		c.err = err[0]
	}

	return c.err
}

func (c Config) RegisterCmd(name string, cmd []string, timeOutOpt ...int) {
	timeOut := 0
	if len(timeOutOpt) > 0 {
		timeOut = timeOutOpt[0]
	}

	fmt.Println("RegisterCmd", name, cmd, timeOut)
	c.refMap.Register(name, cmd, timeOut)
}

func (c *Config) Load(p ProjectLoader, mf ...string) {
	if c.err != nil {
		return
	}

	c.project = p

	if len(mf) > 0 {
		c.metaFile = mf[0]
	}

	c.project.Load(c.metaFile)
	if c.err != nil {
		c.err = errors.Wrap(c.err, "loading configuration file")
		return
	}
	c.project.Process(c.refMap)
	if c.err != nil {
		c.err = errors.Wrap(c.err, "processing configuration file")
		return
	}

	c.project.LoadSecrets(filepath.Join(c.destination, "passwords.json"))
	c.refMap.Assess()
}

func (c *Config) BuildAll(force bool) {
	if c.err != nil {
		return
	}

	c.Force(force)

	for _, ref := range c.refMap.ChangedRefs() {
		for _, val := range ref.Files {
			val.Build(c)
		}
	}

	rsp := c.rExec()
	for rsp.More && c.err == nil {
		rsp = c.rExec()
	}

	if c.err != nil {
		c.err = errors.Wrap(c.err, "building")
	}
}

func (c *Config) rExec() refmap.ExecRsp {
	if c.err != nil {
		return refmap.ExecRsp{}
	}
	rsp := c.refMap.Execute()
	if rsp.Err != nil {
		c.err = rsp.Err
	}
	return rsp
}
