package block

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

func TestRandomBlock(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.NoError(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.TxIDs = TxIDs{}
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.Version = 2
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.StateHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.CommittersHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.LastReceiptsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.LastBlockHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.LastCommitHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.LastCommit.data.Round = b.data.LastCommit.data.Round + 1
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.LastCommit = nil
	assert.Error(t, b.SanityCheck())
}

func TestMarshaling(t *testing.T) {
	b1, _ := GenerateTestBlock(nil, nil)

	bz1, err := b1.MarshalCBOR()
	assert.NoError(t, err)
	var b2 Block
	err = b2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())

	bz2, _ := b1.MarshalCBOR()
	assert.Equal(t, bz1, bz2)
}

func TestJSONMarshaling(t *testing.T) {
	b1, _ := GenerateTestBlock(nil, nil)

	bz1, err := b1.MarshalJSON()
	assert.NoError(t, err)
	var b2 Block
	err = b2.UnmarshalJSON(bz1)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())
}

func TestDecode(t *testing.T) {
	var b1 Block
	d, _ := hex.DecodeString("a301a90101021a60002bd60358200853b05bf757daf780baf9c8cb527a6bbf7624ace66ad3e9d8412ac739deafb004582020bfc2f9c643f454fdd60295a5d35e550d20e26ec55cc4d80a9538172ef1a8b105582052c84ba3994eb2e15420214bcee1c0f242192a9c6818e3b9ffe642af563ac5dd065820b2bf3ef2b1dcda3dff9edb45f868b34046c52eb8a664f2eb192c38bc43dfbcfe0758205e9d55dadf1e68a8d97ab91edd57ee57abf8a36aff8fa1fb6d29f7f603972686085820ba28c8724585a4065e87cc8fb268c74a54ca1ed756399127630debd0e05cff960954ffa6444772e50aface1356e87f61201aa06d898e02a40158200853b05bf757daf780baf9c8cb527a6bbf7624ace66ad3e9d8412ac739deafb002020384a201000200a201010201a201020201a20103020104583004d65285102ea4f322fbc80e39512bca05201739afb309943ac7906fe8439be811bd33dee3af04d9d3365daa7c6be00c03a1018458201cb3f876dc020d745316aeda8612578b51e48b2acba5b5a76516ed8b5bf947d15820f9af3b0565449ee30e1427da8b7ca5ad264524cfe36f23eb4a8be2f67d3a4164582033c924108e42c47437a4a09ad0171cd282b5a47be02dcaddd0dbce581b558f76582081e02ff6fa7bd60b23edd3638055104e270cef262f2b0866770d0609503640f3")
	assert.NoError(t, b1.Decode(d))
	d2, _ := b1.Encode()
	assert.Equal(t, d, d2)
	h, _ := crypto.HashFromString("5a0809874dcbc8705d5eb9815f9d039665d5a01ed2332da20571b19ea86277c5")
	assert.True(t, b1.HashesTo(h))
}

func TestSanityCheck(t *testing.T) {
	tmp, _ := GenerateTestBlock(nil, nil)
	t.Run("Invalid block information, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		MakeBlock(util.Now(), tmp.TxIDs(),
			tmp.Header().LastBlockHash(),
			tmp.Header().CommittersHash(),
			tmp.Header().StateHash(),
			tmp.Header().LastReceiptsHash(), nil, tmp.Header().ProposerAddress())
	})
}
