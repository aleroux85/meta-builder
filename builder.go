package builder

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

type BackRef interface {
	FileStructure() map[string]*FSDirectory
	CmdMatch() map[string]*Exec
	Up() BackRef
}

type DataBranch interface {
	Files() map[string]*FSDirectory
	SetFile(*FSFile)
	File() *FSFile
	Project() *Project
}

type PrjData struct {
	Prj *Project
	FSF *FSFile
}

type Entity struct {
	Name   string                  `json:"name"`
	Files  map[string]*FSDirectory `json:"files"`
	Execs  map[string]*Exec        `json:"execs"`
	Branch DataBranch              `json:"-"`
	Parent BackRef                 `json:"-"`
	changeDetector
}

func (m Entity) FileStructure() map[string]*FSDirectory {
	return m.Files
}

func (m Entity) CmdMatch() map[string]*Exec {
	return m.Execs
}

func (m Entity) Up() BackRef {
	return m.Parent
}

type Exec struct {
	File string   `json:"file"`
	Exec []string `json:"exec"`
}

type changeDetector struct {
	hash   string `json:"-"`
	change uint8  `json:"-"`
}

func hash(d interface{}) (string, error) {
	json, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	_, err = h.Write(json)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
