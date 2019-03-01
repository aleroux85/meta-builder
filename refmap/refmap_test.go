package refmap

import "testing"

func TestWrite(t *testing.T) {
	r := NewRefMap()
	r.Start()

	l := "s"
	r.Set("location", l)

	val := new(DummyRefVal)
	val.hash = "i"
	r.Write("a", "b", val)

	link := r.Read("x")
	if link != nil {
		t.Error("expected \"nil\"")
	}

	link = r.Read("s/a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Change != DataAdded {
		t.Error("expected DataAdded, got", link.Change)
	}
	if link.Files["x"] != nil {
		t.Error("expected \"nil\"")
	}
	if link.Files["b"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["b"].Hash() != "i" {
		t.Error("expected \"i\", got", link.Files["b"].Hash())
	}
	if link.Files["b"].SetChange() != DataAdded {
		t.Error("expected DataAdded, got", link.Files["b"].SetChange())
	}

	val = new(DummyRefVal)
	val.hash = "j"
	r.Write("a", "c", val)

	link = r.Read("s/a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"].Hash() != "j" {
		t.Error("expected \"j\", got", link.Files["c"].Hash())
	}
	if link.Files["c"].SetChange() != DataAdded {
		t.Error("expected DataAdded, got", link.Files["c"].SetChange())
	}

	val = new(DummyRefVal)
	val.hash = "k"
	r.Write("a", "c", val)

	link = r.Read("s/a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"].Hash() != "j" {
		t.Error("expected \"j\", got", link.Files["c"].Hash())
	}
	if link.Files["c"].SetChange() != DataAdded {
		t.Error("expected DataAdded, got", link.Files["c"].SetChange())
	}

	r.Finish()
	val = new(DummyRefVal)
	val.hash = "j"
	r.Write("a", "c", val)

	link = r.Read("s/a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"].Hash() != "j" {
		t.Error("expected \"j\", got", link.Files["c"].Hash())
	}
	if link.Files["c"].SetChange() != DataFlagged {
		t.Error("expected DataFlagged, got", link.Files["c"].SetChange())
	}

	r.Finish()
	val = new(DummyRefVal)
	val.hash = "k"
	r.Write("a", "c", val)

	link = r.Read("s/a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["c"].Hash() != "k" {
		t.Error("expected \"k\", got", link.Files["c"].Hash())
	}
	if link.Files["c"].SetChange() != DataUpdated {
		t.Error("expected DataUpdated, got", link.Files["c"].SetChange())
	}
}

func TestRemove(t *testing.T) {
	r := NewRefMap()
	r.Start()

	val := new(DummyRefVal)
	val.hash = "i"
	r.Write("a", "b", val)

	r.Finish()
	link := r.Read("a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Change != DataStable {
		t.Error("expected DataStable, got", link.Change)
	}
	if link.Files["b"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["b"].SetChange() != DataStable {
		t.Error("expected DataStable, got", link.Files["b"].SetChange())
	}

	val = new(DummyRefVal)
	val.hash = "i"
	r.Write("a", "c", val)

	r.Assess()
	link = r.Read("a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Change != DataFlagged {
		t.Error("expected DataFlagged, got", link.Change)
	}
	if link.Files["b"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["b"].SetChange() != DataRemove {
		t.Error("expected DataRemove, got", link.Files["b"].SetChange())
	}

	r.Finish()
	link = r.Read("a")
	if link.Files["b"] != nil {
		t.Error("expected \"nil\"")
	}

	r.Assess()
	link = r.Read("a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Change != DataRemove {
		t.Error("expected DataRemove, got", link.Change)
	}

	r.Finish()
	link = r.Read("a")
	if link != nil {
		t.Error("expected \"nil\"")
	}
}

func TestSetUpdate(t *testing.T) {
	r := NewRefMap()
	r.Start()

	val := new(DummyRefVal)
	val.hash = "i"
	r.Write("a", "b", val)

	r.Finish()
	err := r.Set("a", "xyz")
	if err == nil {
		t.Error("expected error")
	}
	r.Set("a", "update")

	link := r.Read("a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Change != DataUpdated {
		t.Error("expected DataUpdated, got", link.Change)
	}
	if link.Files["b"] == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["b"].SetChange() != DataStable {
		t.Error("expected DataStable, got", link.Files["b"].SetChange())
	}

	r.Assess()
	r.Finish()
	link = r.Read("a")
	if link == nil {
		t.Error("expected not \"nil\"")
	}
	if link.Files["b"] != nil {
		t.Error("expected \"nil\"")
	}
}

type DummyRefVal struct {
	built  bool
	change uint
	hash   string
}

func (r *DummyRefVal) Build(c config) {
	r.built = true
}

func (r *DummyRefVal) SetChange(v ...uint) uint {
	if len(v) > 0 {
		r.change = v[0]
	}
	return r.change
}

func (r *DummyRefVal) Hash() string {
	return r.hash
}
