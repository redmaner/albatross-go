package types

import "encoding/binary"

type Uint32 uint32

func (v Uint32) AsBytes() []byte {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v))
	return buf[:]
}
