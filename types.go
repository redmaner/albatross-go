package albatross

import (
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
