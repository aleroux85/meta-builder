package builder_test

import (
	"testing"

	builder "github.com/aleroux85/meta-builder"
)

func TestConfigCreateModify(t *testing.T) {
	config := builder.NewConfig("a", "b")

	if config.Source() != "a" {
		t.Errorf(`expected "a", got "%s"`, config.Source())
	}
	if config.Destination() != "b" {
		t.Errorf(`expected "b", got "%s"`, config.Destination())
	}
	if config.Source("c") != "c" {
		t.Errorf(`expected "c", got "%s"`, config.Source())
	}
	if config.Destination("d") != "d" {
		t.Errorf(`expected "d", got "%s"`, config.Destination())
	}

	if config.Force() {
		t.Errorf(`expected false, got true`)
	}
	if !config.Force(true) {
		t.Errorf(`expected true, got false`)
	}
	if !config.Force() {
		t.Errorf(`expected true, got false`)
	}
	if config.Force() {
		t.Errorf(`expected false, got true`)
	}
	if config.Force(false) {
		t.Errorf(`expected false, got true`)
	}
	if config.Force() {
		t.Errorf(`expected false, got true`)
	}
}
