package builder

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/pkg/errors"
)

type ProjectLoader interface {
	CalculateHash()
	Load(string)
	LoadSecrets(string)
	Process(*refmap.RefMap)
}

type Project struct {
	Repo    string   `json:"repository"`
	Mode    string   `json:"-"`
	Secrets []string `json:"-"`
	Entity
}

func NewProject(err ...*error) *Project {
	var newError error
	p := &Project{
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

func (p *Project) CalculateHash() {
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

func (p *Project) Load(fn string) {
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

func (p *Project) Process(m *refmap.RefMap) {
	p.CalculateHash()

	for name, fs := range p.Directories {
		fs.Name = name
		fs.Parent = p
		fs.SourcePath = name
		fs.DestinationPath = name
		fs.Error = p.Error
		ProcessFS(buildBranch, fs, m)
	}
}

func (p *Project) LoadSecrets(fn string) {
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
			secrets = append(secrets, RandString(16))
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

type PrjData struct {
	TemplateMethods
	Prj *Project
	FSF *FSFile
}

func (d PrjData) Files() map[string]*FSDirectory {
	return d.Prj.Directories
}

func (d *PrjData) SetFile(file *FSFile) {
	d.FSF = file
}

func (d PrjData) File() *FSFile {
	return d.FSF
}

func (d PrjData) Project() *Project {
	return d.Prj
}

type BranchSetter interface {
	SetBranch(...DataBranch) DataBranch
}

const (
	numalphaLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits       = 6
	letterIdxMask       = 1<<letterIdxBits - 1
	letterIdxMax        = 63 / letterIdxBits
)

var rs = rand.NewSource(time.Now().UnixNano())

// RandString generates a random string
func RandString(n int) string {
	// solution from http://stackoverflow.com/a/31832326
	b := make([]byte, n)
	for i, cache, remain := n-1, rs.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rs.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(numalphaLetterBytes) {
			b[i] = numalphaLetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
