package types

import "encoding/binary"

type TransactionFlag uint8

func (a TransactionFlag) AsBytes() []byte {
	return []byte{uint8(a)}
}

const (
	TransactionFlagSignaling        TransactionFlag = 0b10
	TransactionFlagContractCreation TransactionFlag = 0b1
)

type ValidityStartHeight uint64

func (v ValidityStartHeight) AsBytes() []byte {
	var buf []byte
	binary.AppendUvarint(buf[:], uint64(v))
	return buf[:]
}

type Transaction struct {
	Sender              Address
	SenderType          AccountType
	SenderData          []byte
	Recipient           Address
	RecipientType       AccountType
	RecipientData       []byte
	Value               Coin
	Fee                 Coin
	ValidityStartHeight ValidityStartHeight
	NetworkId           NetworkId
	Flags               TransactionFlag
	Proof               []byte
}
