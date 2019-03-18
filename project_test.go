package builder_test

import (
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
