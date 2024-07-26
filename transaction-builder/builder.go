package transactionbuilder

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/redmaner/albatross-go/types"
	"golang.org/x/crypto/blake2b"
)

type Builder struct {
	types.Transaction
	isValid   bool
	signature []byte
	publicKey ed25519.PublicKey
}

func (b *Builder) BasicTransaction(
	sender types.Address,
	recipient types.Address,
	value, fee types.Coin,
	network types.NetworkId,
	validityStartHeight types.ValidityStartHeight,
) *Builder {
	b.isValid = true
	b.Sender = sender
	b.Recipient = recipient
	b.Value = value
	b.Fee = fee
	b.NetworkId = network
	b.ValidityStartHeight = validityStartHeight

	return b
}

func (b *Builder) Sign(keypair ed25519.PrivateKey) error {

	if !b.isValid {
		return fmt.Errorf("invalid tx")
	}

	// TODO:
	// double check signature generation
	// highly doubt this will work
	hasher, err := blake2b.New256(nil)
	if err != nil {
		return err
	}

	if _, err = hasher.Write([]byte{uint8(len(b.RecipientData))}); err != nil {
		return err
	}

	if _, err = hasher.Write(b.RecipientData); err != nil {
		return err
	}

	if _, err := hasher.Write(b.Sender[:]); err != nil {
		return err
	}

	if _, err := hasher.Write(b.SenderType.AsBytes()); err != nil {
		return err
	}

	if _, err := hasher.Write(b.Recipient[:]); err != nil {
		return err
	}

	if _, err := hasher.Write(b.RecipientType.AsBytes()); err != nil {
		return err
	}

	if _, err := hasher.Write(b.Value.AsBytes()); err != nil {
		return err
	}

	if _, err := hasher.Write(b.Fee.AsBytes()); err != nil {
		return err
	}

	if _, err := hasher.Write(b.ValidityStartHeight.AsBytes()); err != nil {
		return err
	}

	if _, err := hasher.Write(b.NetworkId.AsBytes()); err != nil {
		return err
	}

	if _, err := hasher.Write(b.Flags.AsBytes()); err != nil {
		return err
	}

	if b.NetworkId.IsAlbatross() {
		if _, err := hasher.Write(b.SenderData); err != nil {
			return err
		}
	}

	payload := hasher.Sum(nil)

	b.signature = ed25519.Sign(keypair, payload)
	b.publicKey = keypair.Public().(ed25519.PublicKey)
	return nil
}

func (b *Builder) ToHex() (string, error) {
	if !b.isValid {
		return "", fmt.Errorf("invalid tx")
	}

	if len(b.signature) == 0 {
		return "", fmt.Errorf("transaction is not signed")
	}

	buf := bytes.NewBuffer(nil)

	// TODO: break out encoding based on Basic or Extended.
	// currently we only do basic

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
