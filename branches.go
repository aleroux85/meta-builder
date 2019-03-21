package builder

func buildBranch(m BranchSetter) {
	stepper := m.(BackRef)

	for {
		switch v := stepper.(type) {
		case *Project:
			m.SetBranch(&PrjData{
				Prj: v,
			})
			return
		case *FSDirectory:
			stepper = v.Parent
		}
	}
}
