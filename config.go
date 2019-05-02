package builder

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/aleroux85/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

type Config struct {
	project     ProjectLoader
	source      string
	destination string
	force       bool
	watching    bool
	metaFile    string
	refMap      *refmap.RefMap
	tmplMon     *utils.Monitor
	cnfgMon     *utils.Monitor
	err         error
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
	}
	return c.force
}

func (c *Config) Watch(throttling time.Duration) {
	if c.err != nil {
		return
	}

	c.tmplMon = new(utils.Monitor)
	c.tmplMon.Error = &c.err

	c.cnfgMon = new(utils.Monitor)
	c.cnfgMon.Error = &c.err

	c.tmplMon.SetWatcher()
	err := c.refMap.Sync(c.tmplMon)
	if err != nil {
		c.err = err
	}

	c.cnfgMon.SetWatcher()
	c.cnfgMon.Watcher.Add(c.metaFile)

	c.watching = true
	go c.watch(throttling)
}

func (c Config) Watching() bool {
	return c.watching
}

func (c Config) StopWatching() {
	c.cnfgMon.Close()
	c.tmplMon.Close()
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

func (c *Config) Finish() {
	c.refMap.Finish()
}

func (c Config) RegisterCmd(name string, cmd, deps []string, timeoutOpt ...uint) {
	var timeout uint = 0
	if len(timeoutOpt) > 0 {
		timeout = timeoutOpt[0]
	}

	c.refMap.Register(name, cmd, deps, timeout)
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

	for _, act := range c.refMap.Execute() {
		fmt.Printf("%v\n%v", act.StdOut, act.StdErr)

		if act.Err != nil {
			c.err = act.Err
		}
	}

	if c.err != nil {
		c.err = errors.Wrap(c.err, "building")
	}
}

func (c *Config) watch(throttling time.Duration) {
	var tmplChange, cnfgChange bool

	for c.err == nil {
		select {
		case event, ok := <-c.cnfgMon.Watcher.Events:
			if !ok {
				goto stop
			}
			fmt.Println("fs event", event.Op, event.Name)
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Chmod == fsnotify.Chmod {
				cnfgChange = true
			}
		case err := <-c.cnfgMon.Watcher.Errors:
			c.err = err
		case event, ok := <-c.tmplMon.Watcher.Events:
			if !ok {
				goto stop
			}
			fmt.Println("fs event", event.Op, event.Name)
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Chmod == fsnotify.Chmod {
				c.refMap.Set(event.Name, "update")
				tmplChange = true
			}
		case err := <-c.tmplMon.Watcher.Errors:
			c.err = err
		case <-time.After(throttling):
			if tmplChange || cnfgChange {
				if cnfgChange {
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
					c.refMap.Assess()
					err := c.refMap.Sync(c.tmplMon)
					if err != nil {
						c.err = err
					}
				}

				for _, ref := range c.refMap.ChangedRefs() {
					for _, val := range ref.Files {
						val.Build(c)
					}
				}
				c.refMap.Finish()

				for _, act := range c.refMap.Execute() {
					fmt.Printf("%v\n%v", act.StdOut, act.StdErr)

					if act.Err != nil {
						c.err = act.Err
					}
				}

				if c.err != nil {
					c.err = errors.Wrap(c.err, "building")
				}

				cnfgChange = false
				tmplChange = false
			}
		}
	}

stop:
	c.cnfgMon.Stopped <- true
	c.tmplMon.Stopped <- true
}
