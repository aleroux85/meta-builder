package refmap

type ReadOp struct {
	Src string
	Rsp chan *RefLink
}

func (r RefMap) Read(src string) *RefLink {
	read := &ReadOp{
		Src: src,
		Rsp: make(chan *RefLink),
	}
	r.Reads <- read
	return <-read.Rsp
}
