package builder

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/aleroux85/utils"
)

type Project struct {
	Description string   `json:"description"`
	Mode        string   `json:"-"`
	Secrets     []string `json:"-"`
	*Entity
	Error *error `json:"-"`
}

func NewProject(err ...*error) *Project {
	var newError error
	p := new(Project)

	if len(err) == 0 {
		p.Error = &newError
	} else {
		p.Error = err[0]
	}
	return p
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
			secrets = append(secrets, utils.RandString(16))
		}

		b, err := json.Marshal(secrets)
		if err != nil {
			*p.Error = err
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
			*p.Error = err
			return
		}
	}

	for i := 0; i < 10; i++ {
		p.Secrets = append(p.Secrets, secrets[i])
	}
}
