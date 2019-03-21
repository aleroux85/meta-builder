package builder

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

func buildBranch(m branchBuilder) {
	stepper := m.(BackRef)

	for {
		switch v := stepper.(type) {
		case *project:
			m.SetBranch(&PrjData{
				Prj: v,
			})
			return
		case *FSDirectory:
			stepper = v.Parent
		}
	}
}
