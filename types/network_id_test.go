package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkIdIsAlbatross(t *testing.T) {
	assert.Equal(t, false, NetworkIdTest.IsAlbatross())
	assert.Equal(t, true, NetworkIdTestAlbatross.IsAlbatross())
	assert.Equal(t, false, NetworkIdDev.IsAlbatross())
	assert.Equal(t, true, NetworkIdDevAlbatross.IsAlbatross())
	assert.Equal(t, false, NetworkIdMain.IsAlbatross())
	assert.Equal(t, true, NetworkIdMainAlbatross.IsAlbatross())
}
