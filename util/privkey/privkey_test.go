package privkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyTypeFromString(t *testing.T) {
	ed25519Priv := "SECRET1R74ET5MG73Z2AR82EUQS9RWA4G56Q35CM2QEYRNSCLXC3EDAWNEKQ0QC792"
	blsPriv := "SECRET1PYZ76D4XTGP7SQUUU808JNAD6NX6D3Y4ZYRAUFKCY4PNN4U4NTMGSXXHPXA"
	unknownPriv := "SECRET1D74ET5MG73Z2AR82EUQS9RWA4G56Q35CM2QEYRNSCLXC3EDAWNEKQ0QC792"

	tests := []struct {
		name string
		priv string
		typ  PrivateKeyType
		err  bool
	}{
		{"ed25519 private key", ed25519Priv, PrivateKeyTypeEd25519, false},
		{"bls private key", blsPriv, PrivateKeyTypeBLS, false},
		{"unknown private key", unknownPriv, PrivateKeyTypeUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ, err := PrivateKeyTypeFromString(tt.priv)
			assert.Equal(t, tt.typ, typ)
			assert.Equal(t, tt.err, err != nil)
		})
	}
}
