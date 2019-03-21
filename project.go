package builder

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/aleroux85/utils"
	"github.com/pkg/errors"
)

type Project interface {
	CalculateHash()
	Load(string)
	LoadSecrets(string)
	Process(*refmap.RefMap)
}

type project struct {
	Description string   `json:"description"`
	Mode        string   `json:"-"`
	Secrets     []string `json:"-"`
	Entity
	Error *error `json:"-"`
}

func NewProject(err ...*error) *project {
	var newError error
	p := &project{
		Entity: Entity{
			changeDetector: changeDetector{},
		},
	}

	if len(err) == 0 {
		p.Error = &newError
	} else {
		p.Error = err[0]
	}
	return p
}

func (p *project) CalculateHash() {
	if *p.Error != nil {
		return
	}

	pTemp := *p
	pTemp.Directories = nil
	err := p.Entity.changeDetector.CalculateHash(pTemp)
	if err != nil {
		*p.Error = err
	}
}

func (p *project) Load(fn string) {
	if *p.Error != nil {
		return
	}

	f, err := ioutil.ReadFile(fn)
	if err != nil {
		*p.Error = err
		return
	}

	err = json.Unmarshal(f, p)
	if err != nil {
		*p.Error = err
		return
	}
}

func (p *project) Process(m *refmap.RefMap) {

	p.CalculateHash()

	for name, fs := range p.Directories {
		fs.Parent = p

		fs.SourcePath = filepath.Join("", name)
		fs.DestinationPath = filepath.Join("", name)

		fs.Name = name
		err := ProcessFS(buildBranch, fs, m)
		if err != nil {
			*p.Error = err
			return
		}
	}
}

func (p *project) LoadSecrets(fn string) {
	if *p.Error != nil {
		return
	}

	if len(p.Secrets) != 0 {
		return
	}

	_, err := os.Stat(fn)
	if err != nil {
		if !os.IsNotExist(err) {
			*p.Error = err
			return
		}
	}

	var secrets []string

	if os.IsNotExist(err) {
		f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			*p.Error = err
			return
		}

		for i := 0; i < 10; i++ {
			secrets = append(secrets, utils.RandString(16))
		}

		b, err := json.Marshal(secrets)
		if err != nil {
			*p.Error = errors.Wrap(err, "marshalling")
			return
		}

		_, err = f.Write(b)
		if err != nil {
			*p.Error = err
			return
		}

		err = f.Close()
		if err != nil {
			*p.Error = err
			return
		}
	} else {
		read, err := ioutil.ReadFile(fn)
		if err != nil {
			*p.Error = err
			return
		}

		err = json.Unmarshal(read, &secrets)
		if err != nil {
			*p.Error = errors.Wrap(err, "unmarshalling")
			return
		}
	}

	for i := 0; i < 10; i++ {
		p.Secrets = append(p.Secrets, secrets[i])
	}
}
