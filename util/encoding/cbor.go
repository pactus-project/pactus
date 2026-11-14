package encoding

import (
	"io"

	"github.com/fxamacker/cbor/v2"
)

// strictCBOR decodes the CBOR data that comes from the network.
//
// One value has many encodings in CBOR. Each encoding of a signed message is a different
// message for the network layer, which identifies a message by its bytes, so the same vote
// or transaction could be sent again and again under new identifiers.
// This mode rejects the encodings that the encoder of a node never produces:
// items of indefinite length, maps with a duplicate key and tags.
var strictCBOR = func() cbor.DecMode {
	mode, err := cbor.DecOptions{
		DupMapKey:        cbor.DupMapKeyEnforcedAPF,
		IndefLength:      cbor.IndefLengthForbidden,
		TagsMd:           cbor.TagsForbidden,
		MaxArrayElements: 65_536, // default is 131072
		MaxMapPairs:      65_536, // default is 131072
	}.DecMode()
	if err != nil {
		panic(err) // unreachable: the options are valid
	}

	return mode
}()

// UnmarshalCBOR decodes CBOR data that comes from the network into v.
func UnmarshalCBOR(data []byte, v any) error {
	return strictCBOR.Unmarshal(data, v)
}

// NewCBORDecoder returns a decoder for CBOR data that comes from the network.
func NewCBORDecoder(r io.Reader) *cbor.Decoder {
	return strictCBOR.NewDecoder(r)
}
