package builder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/pkg/errors"
)

type FSTemplate struct {
	Name string `json:"name"`
	File string `json:"file"`
	Body string `json:"body"`
}

type FSFile struct {
	Name      string            `json:"name"`
	Copy      bool              `json:"copy"`
	Update    string            `json:"update"`
	Source    string            `json:"source"`
	Templates map[string]string `json:"templates"`
	Parent    BackRef           `json:"-"`
	changeDetector
}

func (file *FSFile) CalculateHash() error {
	var err error

	fileTemp := *file
	err = file.changeDetector.CalculateHash(fileTemp)
	if err != nil {
		return err
	}
	return nil
}

func (file *FSFile) Build(c refmap.Config) {
	if c.Error() != nil {
		return
	}

	dstFilename := strings.TrimSuffix(file.Name, ".tmpl")
	srcFilename := file.Name

	if file.Source != "" {
		srcFilename = filepath.Base(file.Source)
	}

	parentDS := file.Parent.(*FSDirectory)
	branch := parentDS.Branch
	sourcePath := parentDS.SourcePath
	destination := parentDS.DestinationPath

	dstFileLocation := filepath.Join(c.Destination(), destination, dstFilename)

	if _, err := os.Stat(dstFileLocation); err == nil {
		if !c.Force() && !c.Watching() {
			return
		}
	} else if os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(c.Destination(), destination), os.ModePerm)
	} else {
		c.Error(errors.Wrap(err, "building file, stating"))
		return
	}

	f, err := os.Create(dstFileLocation)
	if err != nil {
		c.Error(errors.Wrap(err, "building file, creating"))
		return
	}

	if file.Copy || parentDS.Copy {
		r, err := os.Open(filepath.Join(c.Source(), sourcePath, srcFilename))
		if err != nil {
			c.Error(err)
			return
		}
		defer r.Close()
		_, err = io.Copy(f, r)
		if err != nil {
			c.Error(err)
			return
		}
		err = f.Sync()
		if err != nil {
			c.Error(err)
			return
		}
		goto compile
	}

	fmt.Println("rebuilding", dstFileLocation)

	if parentDS.Template == nil || c.Watching() {
		parentDS.Template = new(Templax)
		err := parentDS.Template.Prepare(filepath.Join(c.Source(), sourcePath))
		if err != nil {
			c.Error(err)
			return
		}
	}

	for _, templates := range file.Templates {
		err := parentDS.Template.Prepare(filepath.Join(c.Source(), templates))
		if err != nil {
			c.Error(err)
			return
		}
	}

	if file.Source != "" {
		if filepath.Dir(file.Source) != "." {
			err := parentDS.Template.Prepare(filepath.Join(c.Source(), file.Source))
			if err != nil {
				c.Error(err)
				return
			}
		}
	}

	branch.SetFile(file)
	err = parentDS.Template.FExecute(f, srcFilename, branch)
	if err != nil {
		c.Error(err)
		return
	}
	f.Close()

compile:

	for walker := file.Parent; walker != nil; walker = walker.Up() {
		for name, e := range walker.CmdMatch() {
			r, _ := regexp.Compile(e.Pattern)
			if r.MatchString(dstFilename) {
				c.RegisterCmd(name, e.Cmd, e.Deps, e.Timeout)
			}
		}
	}
}
