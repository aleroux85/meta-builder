package builder_test

import (
	"fmt"
	"testing"

	builder "github.com/aleroux85/meta-builder"
)

func TestNewProject(t *testing.T) {
	t.Run("new project without initialized error", func(t *testing.T) {
		project := builder.NewProject()
		if project.Error == nil {
			t.Errorf("project error not initialized")
		}
	})

	t.Run("new project with initialized error", func(t *testing.T) {
		var initError error
		project := builder.NewProject(&initError)
		if project.Error == nil {
			t.Errorf("project error not initialized")
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("load project with error", func(t *testing.T) {
		project := builder.NewProject()
		*project.Error = fmt.Errorf("pre-existing error")
		project.Load("testdata/meta.json")
		if *project.Error == nil {
			t.Errorf("expected project error")
		}
		if (*project.Error).Error() != "pre-existing error" {
			t.Errorf(`expected "pre-existing error" error, got "%s"`, (*project.Error).Error())
		}
	})

	t.Run("load project with non-existing file", func(t *testing.T) {
		project := builder.NewProject()
		project.Load("testdata/non-existing.json")
		if *project.Error == nil {
			t.Errorf("expected project error")
		}
		errorString := "open testdata/non-existing.json: no such file or directory"
		if (*project.Error).Error() != errorString {
			t.Errorf(`expected "%s" error, got "%s"`, errorString, (*project.Error).Error())
		}
	})

	t.Run("load project with incorrect file", func(t *testing.T) {
		project := builder.NewProject()
		project.Load("testdata/incorrect.json")
		if *project.Error == nil {
			t.Errorf("expected project error")
		}
		errorString := "unexpected end of JSON input"
		if (*project.Error).Error() != errorString {
			t.Errorf(`expected "%s" error, got "%s"`, errorString, (*project.Error).Error())
		}
	})

	t.Run("load project", func(t *testing.T) {
		project := builder.NewProject()
		project.Load("testdata/meta.json")
		if *project.Error != nil {
			t.Errorf("got project error")
		}
		name := "Abc"
		if project.Name != name {
			t.Errorf(`expected "%s", got "%s"`, name, project.Name)
		}
	})
}
