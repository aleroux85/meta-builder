package builder

import (
	"path/filepath"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/aleroux85/utils"
	"github.com/pkg/errors"
)

type Project interface {
	Load(string, *refmap.RefMap) error
	LoadSecrets(string)
}

type Config struct {
	Project     Project
	source      string
	destination string
	configFile  string
	links       *refmap.RefMap
	force       bool
	tmplMon     *utils.Monitor
	cnfgMon     *utils.Monitor
	err         error
}

func NewConfig(s ...string) *Config {
	c := new(Config)

	if len(s) > 0 {
		c.source = s[0]
	}

	if len(s) > 1 {
		c.destination = s[1]
	}

	c.links = refmap.NewRefMap()
	c.links.Start()
	c.links.Set("location", c.source)

	c.tmplMon = new(utils.Monitor)
	c.tmplMon.Error = &c.err

	c.cnfgMon = new(utils.Monitor)
	c.cnfgMon.Error = &c.err
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
		c.source = dst[0]
	}

	return c.destination
}

func (c *Config) Force() bool {
	return c.force
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

func (c *Config) NewProject(err *error) Project {
	p := new(ProjectDefault)
	p.Error = err
	return p
}

func (c *Config) Load(cnf ...string) {
	if c.err != nil {
		return
	}

	c.Project = c.NewProject(&c.err)

	if len(cnf) > 0 {
		c.configFile = cnf[0]
	}

	err := c.Project.Load(c.configFile, c.links)
	if err != nil {
		c.err = errors.Wrap(err, "loading configuration file")
		return
	}

	c.Project.LoadSecrets(filepath.Join(c.destination, "passwords.json"))
	c.links.Assess()
	// c.Project.Blackboard = c.source
}

func (c *Config) Sync() {
	if c.err != nil {
		return
	}

	sync := &refmap.SyncOp{
		Mon: c.tmplMon,
		Err: make(chan error),
	}
	c.links.Sync <- sync
	c.err = <-sync.Err
}

func (c *Config) Build(force bool) {
	if c.err != nil {
		return
	}

	c.force = force

	changed := &refmap.ChangedOp{make(chan *refmap.RefLink)}
	c.links.Changed <- changed
	for ref := range changed.Refs {
		for _, val := range ref.Files {
			val.Build(c)
		}
	}

	if c.err != nil {
		c.err = errors.Wrap(c.err, "building")
	}
}
