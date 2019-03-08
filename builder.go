package builder

import (
	"github.com/aleroux85/utils"
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

type Project struct {
	Description string   `json:"description"`
	Mode        string   `json:"-"`
	Secrets     []string `json:"-"`
	Entity
	Error *error `json:"-"`
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

type FSDirectory struct {
	Source          string                  `json:"from"`
	Destination     string                  `json:"dest"`
	Directories     map[string]*FSDirectory `json:"directories"`
	Copy            bool                    `json:"copyfiles"`
	Update          string                  `json:"update"`
	Template        *utils.Templax          `json:"-"`
	SourcePath      string                  `json:"-"`
	DestinationPath string                  `json:"-"`
	Entity
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

type FSTemplate struct {
	Name string `json:"name"`
	File string `json:"file"`
	Body string `json:"body"`
}

type Exec struct {
	File string   `json:"file"`
	Exec []string `json:"exec"`
}

type changeDetector struct {
	Hash   string `json:"-"`
	Change uint   `json:"-"`
}
