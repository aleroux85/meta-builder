package builder

import (
	"github.com/aleroux85/meta-builder/refmap"
	"github.com/aleroux85/utils"
)

type project interface {
	Load(string, *refmap.RefMap) error
}

type config struct {
	project     *project
	source      string
	destination string
	configFile  string
	links       *refmap.RefMap
	force       bool
	tmplMon     *utils.Monitor
	cnfgMon     *utils.Monitor
	err         error
}

func NewConfig(s ...string) *config {
	c := new(config)

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
