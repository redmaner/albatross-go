package types

import "encoding/binary"

type Varint int64

func VarintFromInt(i int) Varint {
	return Varint(i)
}

func (v Varint) AsBytes() []byte {
	return binary.AppendVarint(nil, int64(v))
}
