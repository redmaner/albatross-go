package types

type TransactionFlag byte

func (a TransactionFlag) AsBytes() []byte {
	return []byte{byte(a)}
}

const (
	TransactionFlagSignaling        TransactionFlag = 0b10
	TransactionFlagContractCreation TransactionFlag = 0b1
)

type Transaction struct {
	Sender              Address
	SenderType          AccountType
	SenderData          []byte
	Recipient           Address
	RecipientType       AccountType
	RecipientData       []byte
	Value               Coin
	Fee                 Coin
	ValidityStartHeight Uint32
	NetworkId           NetworkId
	Flags               TransactionFlag
	Proof               []byte
}
