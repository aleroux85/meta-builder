package builder

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
	Hash   string `json:"-"`
	Change uint   `json:"-"`
}
