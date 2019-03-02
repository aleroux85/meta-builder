package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/aleroux85/utils"
)

type ProjectDefault struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Files       map[string]*FSDirectory `json:"files"`
	// Blackboard  string                  `json:"-"`
	Mode    string   `json:"-"`
	Secrets []string `json:"-"`
	changeDetector
	Error *error `json:"-"`
}

func (p ProjectDefault) FileStructure() map[string]*FSDirectory {
	return p.Files
}

func (p *ProjectDefault) Load(fn string, m *refmap.RefMap) error {
	if *p.Error != nil {
		return *p.Error
	}

	read, err := ioutil.ReadFile(fn)
	if err != nil {
		*p.Error = err
		fmt.Println("error", p.Error)
		return err
	}

	err = json.Unmarshal(read, p)
	if err != nil {
		*p.Error = err
		return err
	}

	return nil
}

func (p *ProjectDefault) LoadSecrets(fn string) {
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
