//go:generate go run github.com/abice/go-enum --file=enigma_ciper_keys.go --names --nocase=false
package v2

/*
	ENUM(

synapse.default.crypto
)
*/
type CiperKey int
