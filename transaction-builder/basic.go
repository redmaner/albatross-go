package transactionbuilder

import "github.com/redmaner/albatross-go/types"

// Create a Basic transaction
func (b *Builder) NewBasicTransaction(
	sender types.Address,
	recipient types.Address,
	value, fee types.Coin,
	network types.NetworkId,
	validityStartHeight types.Uint32,
) *Builder {
	b.txType = TxTypeBasic
	b.isValid = true
	b.Sender = sender
	b.Recipient = recipient
	b.Value = value
	b.Fee = fee
	b.NetworkId = network
	b.ValidityStartHeight = validityStartHeight

	return b
}
