package builder

import (
	"github.com/aleroux85/utils"
)

type BackRef interface {
	FileStructure() map[string]*FSDirectory
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

type changeDetector struct {
	Hash   string `json:"-"`
	Change uint   `json:"-"`
}
