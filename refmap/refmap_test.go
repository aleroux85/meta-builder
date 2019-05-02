package refmap

import (
	"testing"

	"github.com/aleroux85/utils"
)

func TestWriting(t *testing.T) {
	r := NewRefMap()
	r.Start()

	l := "s"
	r.Set("location", l)

	val := new(StubRefVal)
	val.hash = "test_hash"
	r.Write("a", "b", val)
	val = new(StubRefVal)
	val.hash = "another_test_hash"
	r.Write("a", "c", val)

	t.Run(`reading unavailable ref, test for non-existance`, func(t *testing.T) {
		link := r.Read("unavailable")
		if link != nil {
			t.Errorf(`expected "<nil>", got "%+v"`, link)
		}
	})

	t.Run(`reading available "s/a:b" ref, test for existance`, func(t *testing.T) {
		link := r.Read("s/a")
		if link == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link)
			t.FailNow()
		}
		if link.Change != DataAdded {
			t.Errorf(`expected "%s", got "%s"`, status[DataAdded], status[link.Change])
		}
		if link.Files["unavailable"] != nil {
			t.Errorf(`expected "<nil>", got "%+v"`, link.Files["unavailable"])
		}
		if link.Files["b"] == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link.Files["b"])
			t.FailNow()
		}
		if link.Files["b"].Hash() != "test_hash" {
			t.Errorf(`expected "test_hash", got "%+v"`, link.Files["b"].Hash())
		}
		if link.Files["b"].Change() != DataAdded {
			t.Errorf(`expected "%s", got "%s"`, status[DataAdded], status[link.Files["b"].Change()])
		}
		if link.Files["c"] == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link.Files["c"])
			t.FailNow()
		}
		if link.Files["c"].Hash() != "another_test_hash" {
			t.Errorf(`expected "another_test_hash", got "%+v"`, link.Files["c"].Hash())
		}
		if link.Files["c"].Change() != DataAdded {
			t.Errorf(`expected "%s", got "%s"`, status[DataAdded], status[link.Files["b"].Change()])
		}
	})

	val = new(StubRefVal)
	val.hash = "changed_test_hash"
	r.Write("a", "c", val)

	t.Run(`reading available "s/a:c" ref, test for hash change while status still DataAdded (not Finished)`, func(t *testing.T) {
		link := r.Read("s/a")
		if link == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link)
			t.FailNow()
		}
		if link.Files["c"] == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link.Files["c"])
			t.FailNow()
		}
		if link.Files["c"].Hash() != "changed_test_hash" {
			t.Errorf(`expected "changed_test_hash", got "%+v"`, link.Files["c"].Hash())
		}
		if link.Files["c"].Change() != DataAdded {
			t.Errorf(`expected "%s", got "%s"`, status[DataAdded], status[link.Files["c"].Change()])
		}
	})

	r.Finish()
	val = new(StubRefVal)
	val.hash = "changed_test_hash"
	r.Write("a", "c", val)

	t.Run(`reading available "s/a:c" ref, test for no hash change after Finish()`, func(t *testing.T) {
		link := r.Read("s/a")
		if link == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link)
			t.FailNow()
		}
		if link.Files["c"] == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link.Files["c"])
			t.FailNow()
		}
		if link.Files["c"].Hash() != "changed_test_hash" {
			t.Errorf(`expected "changed_test_hash", got "%+v"`, link.Files["c"].Hash())
		}
		if link.Files["c"].Change() != DataFlagged {
			t.Errorf(`expected "%s", got "%s"`, status[DataFlagged], status[link.Files["c"].Change()])
		}
	})

	r.Finish()
	val = new(StubRefVal)
	val.hash = "another_changed_test_hash"
	r.Write("a", "c", val)

	t.Run(`reading available "s/a:c" ref, test for hash change after Finish()`, func(t *testing.T) {
		link := r.Read("s/a")
		if link == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link)
			t.FailNow()
		}
		if link.Files["c"] == nil {
			t.Errorf(`expected "%+v", got "<nil>"`, link.Files["c"])
			t.FailNow()
		}
		if link.Files["c"].Hash() != "another_changed_test_hash" {
			t.Errorf(`expected "another_changed_test_hash", got "%+v"`, link.Files["c"].Hash())
		}
		if link.Files["c"].Change() != DataUpdated {
			t.Errorf(`expected "%s", got "%s"`, status[DataUpdated], status[link.Files["c"].Change()])
		}
	})
}

func TestRemove(t *testing.T) {
	r := NewRefMap()
	r.Start()

	val := new(StubRefVal)
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
	if link.Files["b"].Change() != DataStable {
		t.Error("expected DataStable, got", link.Files["b"].Change())
	}

	val = new(StubRefVal)
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
	if link.Files["b"].Change() != DataRemove {
		t.Error("expected DataRemove, got", link.Files["b"].Change())
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

	val := new(StubRefVal)
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
	if link.Files["b"].Change() != DataStable {
		t.Error("expected DataStable, got", link.Files["b"].Change())
	}

	for _, link = range r.ChangedRefs() {
		if link.Change != DataUpdated {
			t.Error("expected DataUpdated, got", link.Change)
		}
		if link.Files["b"] == nil {
			t.Error("expected not \"nil\"")
		}
		if link.Files["b"].Change() != DataStable {
			t.Error("expected DataStable, got", link.Files["b"].Change())
		}
	}

	for link = range r.ChangedRefsChan() {
		if link.Change != DataUpdated {
			t.Error("expected DataUpdated, got", link.Change)
		}
		if link.Files["b"] == nil {
			t.Error("expected not \"nil\"")
		}
		if link.Files["b"].Change() != DataStable {
			t.Error("expected DataStable, got", link.Files["b"].Change())
		}
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

func TestExec_Registring_Executing(t *testing.T) {
	t.Run("send unknown operation, test return", func(t *testing.T) {
		r := NewRefMap()
		r.Start()
		register := &ExecOp{
			Op:  "unknown",
			Rsp: make(chan ExecRsp),
		}
		r.Exec <- register

		rsp := <-register.Rsp
		if rsp.Key != "" {
			t.Errorf(`expected empty key`)
		}
	})

	t.Run("register and execute successful command, test timeout", func(t *testing.T) {
		r := NewRefMap()
		r.Start()
		r.Register("a", []string{"sleep", `2`}, []string{}, 1)

		rsp := r.Execute()
		if rsp[0].Err.Error() != "signal: killed" {
			t.Errorf(`expected "signal: killed", got "%+v"`, rsp[0].Err.Error())
		}
	})

	t.Run("register and execute successful command, test std outs", func(t *testing.T) {
		r := NewRefMap()
		r.Start()
		r.Register("a", []string{"sh", "-c", `printf "hello here"; printf "error here" >&2`}, []string{"b"})

		rsp := r.Execute()
		if rsp[0].Err != nil {
			t.Errorf(`expected "<nil>", got "%+v"`, rsp[0].Err)
		}
		if rsp[0].Key != "a" {
			t.Errorf(`expected "a", got "%+v"`, rsp[0].Key)
		}
		if rsp[0].StdOut != "hello here" {
			t.Errorf(`expected "hello here", got "%+v"`, rsp[0].StdOut)
		}
		if rsp[0].StdErr != "error here" {
			t.Errorf(`expected "error here", got "%+v"`, rsp[0].StdErr)
		}
	})

	t.Run("register and execute multiple dependant commands with one error, test order and error", func(t *testing.T) {
		r := NewRefMap()
		r.Start()
		r.Register("a", []string{"printf", `print a`}, []string{"b"})
		r.Register("b", []string{"printf", `print b`}, []string{})
		r.Register("c", []string{"bash", "-c", `exit 2`}, []string{"a"})

		rsp := r.Execute()
		if rsp[0].Err != nil {
			t.Errorf(`expected "<nil>", got "%+v"`, rsp[0].Err)
		}
		if rsp[0].Key != "b" {
			t.Errorf(`expected "b", got "%+v"`, rsp[0].Key)
		}
		if rsp[0].StdOut != "print b" {
			t.Errorf(`expected "print b", got "%s"`, rsp[0].StdOut)
		}

		if rsp[1].Err != nil {
			t.Errorf(`expected "<nil>", got "%+v"`, rsp[1].Err)
		}
		if rsp[1].Key != "a" {
			t.Errorf(`expected "a", got "%+v"`, rsp[1].Key)
		}
		if rsp[1].StdOut != "print a" {
			t.Errorf(`expected "print a", got "%s"`, rsp[1].StdOut)
		}

		if rsp[2].Err.Error() != "exit status 2" {
			t.Errorf(`expected "exit status 2", got "%+v"`, rsp[2].Err)
		}
		if rsp[2].Key != "c" {
			t.Errorf(`expected "c", got "%+v"`, rsp[2].Key)
		}
		if rsp[2].StdOut != "" {
			t.Errorf(`expected "", got "%s"`, rsp[2].StdOut)
		}
	})

	t.Run("execute no commands, test empty list return", func(t *testing.T) {
		r := NewRefMap()
		r.Start()

		rsp := r.Execute()
		if len(rsp) != 0 {
			t.Errorf(`expected no results, got "%+v"`, rsp)
		}
	})
}

func TestSync(t *testing.T) {
	r := NewRefMap()
	r.Start()

	val := new(StubRefVal)
	val.hash = "test_hash"
	r.Write("a", "b", val)
	r.Assess()

	var e error
	w := &utils.Monitor{}
	w.Error = &e

	err := r.Sync(w)
	if err == nil {
		t.Errorf(`expected error "%v"`, err)
	}

	err = w.SetWatcher()
	if err != nil {
		t.Error(err)
	}

	err = r.Sync(w)
	if err != nil {
		t.Errorf(`unexpected error "%v"`, err)
	}
	r.Finish()

	r.Assess()
	err = r.Sync(w)
	if err != nil {
		t.Errorf(`unexpected error "%v"`, err)
	}
}

var status map[uint8]string = map[uint8]string{
	DataStable:  "DataStable",
	DataFlagged: "DataFlagged",
	DataUpdated: "DataUpdated",
	DataAdded:   "DataAdded",
	DataRemove:  "DataRemove",
}

type StubRefVal struct {
	built  bool
	change uint8
	hash   string
}

func (r *StubRefVal) Build(c Config) {
	r.built = true
}

func (r *StubRefVal) Change(v ...uint8) uint8 {
	if len(v) > 0 {
		r.change = v[0]
	}
	return r.change
}

func (r *StubRefVal) Hash() string {
	return r.hash
}
