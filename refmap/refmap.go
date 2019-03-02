package refmap

const (
	DataStable uint = iota
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
}

func NewRefMap() *RefMap {
	m := new(RefMap)
	m.Reads = make(chan *ReadOp)
	m.Writes = make(chan *WriteOp)
	m.Sets = make(chan *SetOp)
	m.Sync = make(chan *SyncOp)
	m.Changed = make(chan *ChangedOp)
	return m
}

func (m *RefMap) Start() {
	go func() {
		refs := make(map[string]*RefLink)
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
			}
		}
	}()
}

type Config interface {
	Error(...error) error
	Destination() string
	Source() string
	Force() bool
}

type RefVal interface {
	Build(Config)
	SetChange(...uint) uint
	Hash() string
}

type RefLink struct {
	Files  map[string]RefVal
	Change uint
}

func NewRefLink() *RefLink {
	r := new(RefLink)
	r.Files = make(map[string]RefVal)
	r.Change = DataAdded
	return r
}
