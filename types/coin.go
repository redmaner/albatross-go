package types

import "encoding/binary"

type Coin uint64

func (c Coin) AsBytes() []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(c))
	return buf[:]
}
