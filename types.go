package albatross

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// Luna is the smallest unit of NIM and 100’000 (1e5) Luna equals 1 NIM
type Luna uint64

const nimInLuna int64 = 100000

// FormatLuna is a function to format NIM to Luna
func FormatLuna(n NIM) (Luna, error) {

	nim, err := decimal.NewFromString(string(n))
	if err != nil {
		return 0, err
	}

	luna := nim.Mul(decimal.NewFromInt(nimInLuna)).IntPart()
	return Luna(luna), nil
}

// ToNIM converts Luna to NIM
func (l *Luna) ToNIM() NIM {
	return FormatNIM(*l)
}

// NIM is the token transacted within Nimiq as a store and transfer of value: it acts as digital cash
type NIM string

// FormatNIM is a function to format Luna to NIM
func FormatNIM(l Luna) NIM {
	nim := decimal.NewFromInt(int64(l))
	nim = nim.Div(decimal.NewFromInt(nimInLuna))
	return NIM(nim.String())
}

// ToLuna converts NIM to Luna
func (n *NIM) ToLuna() (Luna, error) {
	return FormatLuna(*n)
}

// Block represents a block on the Nimiq 2.0 blockchain
type Block struct {
	Number    int   `json:"number"`
	Epoch     int   `json:"epoch"`
	Batch     int   `json:"batch"`
	Timestamp int64 `json:"timestamp"`

	Type            string `json:"type"`
	IsElectionBlock bool   `json:"isElectionBlock"`

	ExtraData    []byte          `json:"extraData"`
	Transactions json.RawMessage `json:"transactions"`

	// Producer is only returned for Micro blocks
	Producer *Slot `json:"slot,omitempty"`

	// Slots is only returned in an election block and contains
	// the slot distribution for the next epoch
	Slots []Slots `json:"slots,omitempty"`
}

// Slot represents a slot used to produce a micro block
type Slot struct {
	SlotNumber int    `json:"slotNumber"`
	Validator  string `json:"validator"`
	PublicKey  string `json:"publicKey"`
}

type Slots struct {
	FirstSlotNumber int    `json:"FirstSlotNumber"`
	NumSlots        int    `json:"numSlots"`
	Validator       string `json:"validator"`
	PublicKey       string `json:"publicKey"`
}

// Account represents an account on the Nimiq 2.0 blockchain
// Currently only Basic account fields are unmarshalled
type Account struct {
	Address string `json:"address"`
	Balance Luna   `json:"balance"`
	Type    string `json:"type"`
}

// ReturnAccount holds information of an account that is returned when
// a new account is created through the RPC interface
type ReturnAccount struct {
	Address    string `json:"address"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"PrivateKey"`
}
