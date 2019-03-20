package builder

import (
	"path/filepath"
	"strings"

	"github.com/aleroux85/meta-builder/refmap"
	"github.com/aleroux85/utils"
)

type FSDirectory struct {
	Source          string             `json:"from"`
	Destination     string             `json:"dest"`
	Files           map[string]*FSFile `json:"files"`
	Copy            bool               `json:"copyfiles"`
	Update          string             `json:"update"`
	Template        *utils.Templax     `json:"-"`
	SourcePath      string             `json:"-"`
	DestinationPath string             `json:"-"`
	Entity
}

func (fs *FSDirectory) SetBranch(branch ...DataBranch) DataBranch {
	if len(branch) > 0 {
		fs.Branch = branch[0]
	}
	return fs.Branch
}

func (dir *FSDirectory) CalculateHash() error {
	var err error

	dirTemp := *dir
	dirTemp.Directories = nil
	dirTemp.Files = nil
	err = dir.changeDetector.CalculateHash(dirTemp)
	if err != nil {
		return err
	}
	return nil
}

func ProcessFSs(location string, b BackRef, m *refmap.RefMap, opts ...string) error {
	for name, fs := range b.FileStructure() {
		fs.Parent = b

		fs.SourcePath = filepath.Join(location, name)
		fs.DestinationPath = filepath.Join(location, name)

		for _, opt := range opts {
			if opt == "skip top directory name" {
				fs.SourcePath = location
				fs.DestinationPath = location
				break
			}
		}

		fs.Name = name
		err := processFS(fs, m)
		if err != nil {
			return err
		}
	}
	return nil
}

type PrjData struct {
	Prj *project
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

func (d PrjData) Project() *project {
	return d.Prj
}

type branchBuilder interface {
	SetBranch(...DataBranch) DataBranch
}

func buildBranch(m branchBuilder) DataBranch {
	stepper := m.(BackRef)

	for {
		switch v := stepper.(type) {
		case *project:
			return m.SetBranch(&PrjData{
				Prj: v,
			})
		case *FSDirectory:
			stepper = v.Parent
		}
	}

	return nil
}

func processFS(fs *FSDirectory, m *refmap.RefMap) error {
	buildBranch(fs)

	err := fs.CalculateHash()
	if err != nil {
		return err
	}

	fs.SourcePath = path(fs.SourcePath, fs.Source)
	fs.DestinationPath = path(fs.DestinationPath, fs.Destination)

	for name, dir := range fs.Directories {
		dir.Parent = fs
		dir.SourcePath = filepath.Join(fs.SourcePath, name)
		dir.DestinationPath = filepath.Join(fs.DestinationPath, name)
		dir.Name = name
		err := processFS(dir, m)
		if err != nil {
			return err
		}
	}

	for name, file := range fs.Files {
		file.Name = name
		file.Parent = fs

		err := file.CalculateHash()
		if err != nil {
			return err
		}

		filename := name
		if file.Source != "" {
			filename = file.Source
		}

		if m != nil {
			source := filepath.Join(fs.SourcePath, filename)
			destination := filepath.Join(fs.DestinationPath, name)
			m.Write(source, destination, file)
		}
	}

	return nil
}

func path(path, modify string) string {
	if strings.HasPrefix(modify, ".") {
		return filepath.Join(filepath.Dir(path), modify)
	}
	if strings.HasPrefix(modify, "/") {
		return strings.TrimPrefix(modify, "/")
	}
	return filepath.Join(path, modify)
}
