package objbytes

import "unsafe"

// Header memory storage structure
type Header struct {
    Magic     uint32
    PairCount uint32
}

// PairOffset relative offset
type PairOffset struct {
    From uint32
    To   uint32
}

// OverHead message head
type OverHead struct {
    Head  Header
    Pairs []PairOffset
    Data  uintptr
}

var emptyHeader = Header{}

const (
    HeaderSize = uint64(unsafe.Sizeof(emptyHeader))
    Magic      = 2112
    Align      = 8
)
