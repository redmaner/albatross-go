package types

import "encoding/binary"

type Varint uint64

func VarintFromInt(i int) Varint {
	return Varint(i)
}

func (v Varint) AsBytes() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(v))
	return buf[:n]
}
