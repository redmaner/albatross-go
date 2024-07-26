package types

import "encoding/binary"

type Coin uint64

func (c Coin) AsBytes() []byte {
	var valueBuf []byte
	binary.AppendUvarint(valueBuf, uint64(c))
	return valueBuf
}
