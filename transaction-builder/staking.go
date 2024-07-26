package transactionbuilder

import (
	"github.com/redmaner/albatross-go/types"
)

func (b *Builder) NewAddStakeTransaction(
	sender types.Address,
	recipient types.Address,
	value, fee types.Coin,
	network types.NetworkId,
	validityStartHeight types.Uint32,
) *Builder {

	b.txType = TxTypeExtended
	b.isValid = true
	b.Sender = sender
	b.Recipient = types.STAKING_CONTRACT_ADDRESS
	b.RecipientType = types.AccountTypeStaking
	b.RecipientData = append([]byte{6}, recipient[:]...)
	b.Value = value
	b.Fee = fee
	b.NetworkId = network
	b.ValidityStartHeight = validityStartHeight

	return b
}
