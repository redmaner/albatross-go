package types

type NetworkId uint8

func (a NetworkId) AsBytes() []byte {
	return []byte{uint8(a)}
}

const (
	// Nimiq 1.0
	NetworkIdTest   NetworkId = 1
	NetworkIdDev    NetworkId = 2
	NetworkIdBounty NetworkId = 3
	NetworkIdDummy  NetworkId = 4
	NetworkIdMain   NetworkId = 42

	// Nimiq 2.0 Albatross
	NetworkIdTestAlbatross NetworkId = 5
	NetworkIdDevAlbatross  NetworkId = 6
	NetworkIdUnitAlbatross NetworkId = 7
	NetworkIdMainAlbatross NetworkId = 24
)

// Returns whether the network id is albatross (Nimiq v2)
func (n NetworkId) IsAlbatross() bool {
	return n == NetworkIdTestAlbatross || n == NetworkIdDevAlbatross || n == NetworkIdUnitAlbatross || n == NetworkIdMainAlbatross
}
