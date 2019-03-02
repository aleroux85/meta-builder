package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

import (
	"github.com/aleroux85/meta-builder/refmap"
	"github.com/aleroux85/utils"
)

type Project_Default struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Files       map[string]*FSDirectory `json:"files"`
	Blackboard  string                  `json:"-"`
	Mode        string                  `json:"-"`
	Secrets     []string                `json:"-"`
	changeDetector
	Error *error `json:"-"`
}

func (p *Project_Default) Load(fn string, m *refmap.RefMap) error {
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

// func (p *Project) Process(fn string, m *RefMap) error {
// 	if *p.Error != nil {
// 		return *p.Error
// 	}
// }

type BackRef interface {
	FileStructure() map[string]*FSDirectory
}

func (p Project_Default) FileStructure() map[string]*FSDirectory {
	return p.Files
}

type DataBranch interface {
	Files() map[string]*FSDirectory
	SetFile(*FSFile)
	File() *FSFile
	Project() *Project
}

// type Entity struct {
// 	Name  string                  `json:"name"`
// 	Files map[string]*FSDirectory `json:"files"`
// 	changeDetector
// 	Error *error `json:"-"`
// }

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

type FSDirectory struct {
	Name            string                  `json:"name"`
	Source          string                  `json:"from"`
	Destination     string                  `json:"dest"`
	Directories     map[string]*FSDirectory `json:"directories"`
	Files           map[string]*FSFile      `json:"files"`
	Copy            bool                    `json:"copyfiles"`
	Update          string                  `json:"update"`
	Parent          BackRef                 `json:"-"`
	Template        *utils.Templax          `json:"-"`
	Branch          DataBranch              `json:"-"`
	SourcePath      string                  `json:"-"`
	DestinationPath string                  `json:"-"`
	changeDetector
}

// const (
// 	DataStable uint = iota
// 	DataFlagged
// 	DataUpdated
// 	DataAdded
// 	DataRemove
// )

type changeDetector struct {
	Hash   string `json:"-"`
	Change uint   `json:"-"`
}
