package transactionbuilder

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/redmaner/albatross-go/types"
)

type TxType byte

const (
	TxTypeBasic TxType = iota
	TxTypeExtended
)

type Builder struct {
	types.Transaction
	isValid   bool
	txType    TxType
	signature []byte
	publicKey ed25519.PublicKey
}

func (b *Builder) BasicTransaction(
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

func (b *Builder) Sign(keypair ed25519.PrivateKey) (err error) {

	if !b.isValid {
		return fmt.Errorf("invalid tx")
	}

	buf := bytes.NewBuffer(nil)
	if _, err = buf.Write(types.Uint16FromInt(len(b.RecipientData)).AsBytes()); err != nil {
		return err
	}

	if _, err = buf.Write(b.RecipientData); err != nil {
		return err
	}

	if _, err := buf.Write(b.Sender[:]); err != nil {
		return err
	}

	if _, err := buf.Write(b.SenderType.AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(b.Recipient[:]); err != nil {
		return err
	}

	if _, err := buf.Write(b.RecipientType.AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(b.Value.AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(b.Fee.AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(b.ValidityStartHeight.AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(b.NetworkId.AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(b.Flags.AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(types.VarintFromInt(len(b.SenderData)).AsBytes()); err != nil {
		return err
	}

	if _, err := buf.Write(b.SenderData); err != nil {
		return err
	}

	b.signature = ed25519.Sign(keypair, buf.Bytes())
	b.publicKey = keypair.Public().(ed25519.PublicKey)
	return nil
}

func (b *Builder) Encode() (string, error) {
	if !b.isValid {
		return "", fmt.Errorf("invalid tx")
	}

	if len(b.signature) == 0 {
		return "", fmt.Errorf("transaction is not signed")
	}

	switch b.txType {
	case TxTypeBasic:
		return b.encodeBasic()
	}

	return "", fmt.Errorf("unsupported tx type")
}

func (b *Builder) encodeBasic() (string, error) {
	buf := bytes.NewBuffer(nil)

	if err := buf.WriteByte(byte(b.txType)); err != nil {
		return "", err
	}

	// TODO Webauth is not supported as of yet
	// signatureType is hardcoded to ed25519
	var signatureType uint8 = 0
	var flag uint8 = 0
	signatureType |= flag << 4

	if err := buf.WriteByte(signatureType); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.publicKey); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.Recipient[:]); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.Value.AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.Fee.AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.ValidityStartHeight.AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.NetworkId.AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.signature); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf.Bytes()), nil
}
