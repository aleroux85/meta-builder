package builder

import "testing"

func TestDSPathFunc(t *testing.T) {
	if p := path("a/b", ""); p != "a/b" {
		t.Errorf(`expected "a/b", got "%s"`, p)
	}
	if p := path("a/b", "c"); p != "a/b/c" {
		t.Errorf(`expected "a/b/c", got "%s"`, p)
	}
	if p := path("a/b", "."); p != "a" {
		t.Errorf(`expected "a", got "%s"`, p)
	}
	if p := path("a/b", "./c"); p != "a/c" {
		t.Errorf(`expected "a/c", got "%s"`, p)
	}
	if p := path("a/b", "/"); p != "" {
		t.Errorf(`expected "", got "%s"`, p)
	}
	if p := path("a/b", "/c"); p != "c" {
		t.Errorf(`expected "c", got "%s"`, p)
	}
}
