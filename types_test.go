package albatross

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLunaToNIM(t *testing.T) {
	assert := assert.New(t)

	// 0 luna is 0 NIM
	assert.Equal(NIM("0"), FormatNIM(0), "NIM not formatted correctly")

	// 1 luna is 0.00001 NIM
	assert.Equal(NIM("0.00001"), FormatNIM(1), "NIM not formatted correctly")

	// 100000 luna is 1 NIM
	assert.Equal(NIM("1"), FormatNIM(100000), "NIM not formatted correctly")

	// 1234567 is 12.34567 NIM
	assert.Equal(NIM("12.34567"), FormatNIM(1234567), "NIM not formatted correctly")

	// 1200000 is 12 NIM
	assert.Equal(NIM("12"), FormatNIM(1200000), "NIM not formatted correctly")

	// 123456789 is 1234.56789 NIM
	assert.Equal(NIM("1234.56789"), FormatNIM(123456789), "NIM not formatted correctly")
}

func TestNIMToLuna(t *testing.T) {
	// 0 NIM is 0 Luna
	if l, _ := FormatLuna("0"); l != 0 {
		t.Fail()
	}
	// 0.00001 NIM is 1 Luna
	if l, _ := FormatLuna("0.00001"); l != 1 {
		t.Fail()
	}
	// 1 NIM is 100000 Luna
	if l, _ := FormatLuna("1"); l != 100000 {
		t.Fail()
	}
	// 12.34567 NIM is 1234567 Luna
	if l, _ := FormatLuna("12.34567"); l != 1234567 {
		t.Fail()
	}
	// 12 NIM is 1200000 Luna
	if l, _ := FormatLuna("12"); l != 1200000 {
		t.Fail()
	}
	// 1234.56789 NIM is 123456789 Luna
	if l, _ := FormatLuna("1234.56789"); l != 123456789 {
		t.Fail()
	}
}

// The Nimiq Network has been designed for a total supply of 21 Billion NIM.
// The smallest unit of NIM is called Luna and 100’000 (1e5) Luna equal 1 NIM,
// which results in a total supply of 21e14 Luna
func TestMaxLuna(t *testing.T) {
	max, _ := FormatLuna("21000000000")
	if max != 2100000000000000 {
		t.Fail()
	}
}
