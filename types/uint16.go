package types

import "encoding/binary"

type Uint16 uint16

func Uint16FromInt(i int) Uint16 {
	return Uint16(i)
}

func (u Uint16) AsBytes() []byte {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(u))
	return buf[:]
}
