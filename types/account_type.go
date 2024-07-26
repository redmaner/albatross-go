package types

type AccountType uint8

func (a AccountType) AsBytes() []byte {
	return []byte{uint8(a)}
}

const (
	AccountTypeBasic AccountType = iota
	AccountTypeVesting
	AccountTypeHTLC
	AccountTypeStaking
)
