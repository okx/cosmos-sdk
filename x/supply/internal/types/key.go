package types

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "supply"

	// StoreKey is the store key string for supply
	StoreKey = ModuleName

	// RouterKey is the message route for supply
	RouterKey = ModuleName

	// QuerierRoute is the querier route for supply
	QuerierRoute = ModuleName
)

var (
	PrefixTokenSupplyKey = []byte{0x00}
)

// GetTokenSupplyKey gets the store key of a supply for a token
func GetTokenSupplyKey(tokenSymbol string) []byte {
	return append(PrefixTokenSupplyKey, []byte(tokenSymbol)...)
}
