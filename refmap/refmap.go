package refmap

const (
	DataStable uint8 = iota
	DataFlagged
	DataUpdated
	DataAdded
	DataRemove
)

type RefMap struct {
	Reads   chan *ReadOp
	Writes  chan *WriteOp
	Sets    chan *SetOp
	Sync    chan *SyncOp
	Changed chan *ChangedOp
	Exec    chan *ExecOp
}

func NewRefMap() *RefMap {
	m := new(RefMap)
	m.Reads = make(chan *ReadOp)
	m.Writes = make(chan *WriteOp)
	m.Sets = make(chan *SetOp)
	m.Sync = make(chan *SyncOp)
	m.Changed = make(chan *ChangedOp)
	m.Exec = make(chan *ExecOp)
	return m
}

type command struct {
	Cmd     []string
	TimeOut int
}

func (m *RefMap) Start() {
	go func() {
		refs := make(map[string]*RefLink)
		execs := make(map[string]command)
		startLocation := ""

		for {
			select {
			case read := <-m.Reads:
				read.Rsp <- refs[read.Src]
			case write := <-m.Writes:
				write.handle(startLocation, refs)
			case setter := <-m.Sets:
				setter.handle(&startLocation, refs)
			case sync := <-m.Sync:
				sync.handle(refs)
			case changed := <-m.Changed:
				changed.handle(refs)
			case exec := <-m.Exec:
				exec.handle(execs)
			}
		}
	}()
}

type Config interface {
	Error(...error) error
	Destination(...string) string
	Source(...string) string
	Force() bool
	RegisterCmd(string, []string, ...int)
}

type RefVal interface {
	Build(Config)
	SetChange(...uint8) uint8
	GetHash() string
}

type RefLink struct {
	Files  map[string]RefVal
	Change uint8
}

func NewRefLink() *RefLink {
	r := new(RefLink)
	r.Files = make(map[string]RefVal)
	r.Change = DataAdded
	return r
}
