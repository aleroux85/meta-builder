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

func (fs *FSDirectory) BuildBranch() {
	buildBranch(fs)
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

func ProcessFS(fs *FSDirectory, m *refmap.RefMap) error {
	// bb(fs)
	fs.BuildBranch()

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
		err := ProcessFS(dir, m)
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
