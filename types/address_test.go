package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	address, err := NewAddressFromHex("93ef3e945f99cf643f26a258a28832f898eda496")
	assert.NoError(t, err)

	friendlyAddress := address.String()
	assert.Equal(t, "NQ61 JFPK V52Y K77N 8FR6 L9CA 521J Y2CE T94N", friendlyAddress)

	decodedFromFriendly, err := NewAddressFromFriendly(friendlyAddress)
	assert.NoError(t, err)
	assert.Equal(t, decodedFromFriendly, address)
}
