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

// Sign signs the transaction using the provided privatekey
// Currently only supports signing using ed25519
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

// Encode transaction to a raw hex encoded transaction
// Will return an error if the transaction has not been build yet,
// or has not been signed yet
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
	case TxTypeExtended:
		return b.encodeExtended()
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

func (b *Builder) encodeExtended() (string, error) {
	buf := bytes.NewBuffer(nil)

	if err := buf.WriteByte(byte(b.txType)); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.Sender[:]); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.SenderType.AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(types.Varint(len(b.SenderData)).AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.SenderData); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.Recipient[:]); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.RecipientType.AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(types.Varint(len(b.RecipientData)).AsBytes()); err != nil {
		return "", err
	}

	if _, err := buf.Write(b.RecipientData); err != nil {
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

	if _, err := buf.Write(b.Flags.AsBytes()); err != nil {
		return "", err
	}

	proof, err := b.encodeProof()
	if err != nil {
		return "", err
	}

	if _, err := buf.Write(proof); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func (b *Builder) encodeProof() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// TODO Webauth is not supported as of yet
	// signatureType is hardcoded to ed25519
	var signatureType uint8 = 0
	var flag uint8 = 0
	signatureType |= flag << 4

	if err := buf.WriteByte(signatureType); err != nil {
		return nil, err
	}

	if _, err := buf.Write(b.publicKey); err != nil {
		return nil, err
	}

	// merkle path, we don't have that
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}

	if _, err := buf.Write(b.signature); err != nil {
		return nil, err
	}

	// should this be hex encoded first?
	return buf.Bytes(), nil
}
