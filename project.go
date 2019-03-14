package builder

import (
	"encoding/json"
	"io/ioutil"
)

type Project struct {
	Description string   `json:"description"`
	Mode        string   `json:"-"`
	Secrets     []string `json:"-"`
	Entity
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
