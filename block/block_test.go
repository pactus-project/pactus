package block

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/util"
)

func TestSanityCheck(t *testing.T) {
	b := GenerateTestBlock(nil, nil)
	assert.NoError(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Txs = Txs{}
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.StateRoot = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxsRoot = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.PrevBlockHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.PrevCertHash = hash.UndefHash
	b.data.Header.data.PrevBlockHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.SortitionSeed = sortition.UndefVerifiableSeed
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.ProposerAddress = crypto.TreasuryAddress
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.PrevCertHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.PrevCert = GenerateTestCertificate(hash.UndefHash)
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxsRoot = hash.GenerateTestHash()
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.PrevCert = GenerateTestCertificate(hash.GenerateTestHash())
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxsRoot = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.PrevCert = nil
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	invalidSigner := bls.GenerateTestSigner()
	invalidSigner.SignMsg(b.data.Txs[0])
	assert.Error(t, b.SanityCheck())
}

func TestMarshaling(t *testing.T) {
	b1 := GenerateTestBlock(nil, nil)

	bz1, err := b1.Encode()
	assert.NoError(t, err)
	fmt.Printf("%x\n", bz1)
	var b2 Block
	err = b2.Decode(bz1)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())

	assert.Equal(t, b1.Hash(), b2.Hash())
	assert.Equal(t, b1.Header().Time(), b2.Header().Time())
	assert.Equal(t, b1.Header().Version(), b2.Header().Version())
	assert.Equal(t, b2.Header().Version(), 1)
	assert.Equal(t, b1.Hash(), b1.Header().Hash())
}

func TestJSONMarshaling(t *testing.T) {
	b1 := GenerateTestBlock(nil, nil)

	bz, err := b1.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, bz)
}

func TestDecode(t *testing.T) {
	var b Block
	/*
		{
			1: {
				1: 1,
				2: 1646972445,
				3: h'307CDF86A913BAF3FBF00165DCEA2428B7BA91507178CF1DA3D7FC5F4498C65D',
				4: h'BC0D63C3292FD7C8BF5B5A95E5C1F7D9F70F57E23130960AC6DD0D44E0BD4A50',
				5: h'CB0CCC31A192186E5E01645AAB5CE886BC9D1560A01E8100617C3FB7D6370CB8',
				6: h'9A664DBC4F48F92BCC4AE364F80A81EC4124BA78BC1484B0782FCEED1E505AAE',
				7: h'B4B43DD593FFB8938E710DC6DACC9982182B5DEC5F6A5D1C4A4D516F0F85130FD32D38D34E1AA3BAED1F51A9E06DEC14',
				8: h'019913B0689EEFC8349183EC71FCF7FDDFD2F035A5'
			},
			2: {
				1: h'6FD848BDA84B325277EED41DA0C36B84AC8B4EFE96606E87C1A8634B2AC63D78',
				2: 7,
				3: [10, 18, 12, 16],
				4: [18],
				5: h'B0F50FA16F29747513C29EE2880812E119A3300B226F8BF3BF2F25A04157F6B2EC7A7093E882F176219A61D4B58C12FD'
			}, 3: [
				{1: 1, 2: h'9BBAE371', 3: 110, 4: 1000, 5: 1, 6: h'A30155017C6646C1C19839192E8C76B9DA6F8EAD5A31E7A40255014F6C19966599CE7796CD8C6D4C7DFEB24EAC1941031903E8', 7: "test send-tx", 8: h'83F16A0F72BE00F653F364FE3756FCD4D569B276256C9E0719C12C6DEF2A82E92649475B00A699BD3046523B9C993E9D07F7B4BDCCA90121B609D4F7CE3174C50AE4D88E39132A45D826AB7ACA66A63083DCE33EC7382C29384085AD62DD9AAE', 9: h'8B4A3AEF75C1BF419E79A33763D49BA3770138278627829539C42ABEBFC580F357F75A7ED825414CFD35C80634695690'},
				{1: 1, 2: h'2FAB6D29', 3: 110, 4: 1000, 5: 1, 6: h'A301550198E01D56F0C0C346E5A45B6CDDDE89C98AE7CE4B02550192C557E2850F632E9901B52E1F3502D994776686031903E8', 7: "test send-tx", 8: h'B522CE35CE82962C20B5618990490A1F60BAE007A582443303DE0942625A16A5DEB2BF816A43A867B499CF2C5DD3A33E057D6972872F3F741BDB72C8AD2CBB148F91055046324BFDEDAA51C671A6D8CB61A3070BBA07FF25B32D15326655313A', 9: h'B3EE7FA1DD553C0822DA3E7CC3A6FFDF79DC7CBC7F10548801C691E0A2E2B0BEA473C4E154E4749C3E559C6BF6D03065'},
				{1: 1, 2: h'98981DAC', 3: 110, 4: 1000, 5: 1, 6: h'A3015501E302B94CBEEDDCDFAE85853ED1FF3A03F311769E02550174B82826F6E3988BBD9D825295CDA46D3BEA68A2031903E8', 7: "test send-tx", 8: h'B2BDF71B26B4EAE9F41BF6BAFC429E72EE26AF86DAEF1F2FA2C2BEFB9D7CA959686B0BDA8D05F9E260EE8B81CE93540704DFE47855FE6787D3705FCB60403EB52B54CAC6E9719484C33822064FB1BAB829BD759F64C82A48ECBF79C845B16927', 9: h'B9439E12F760F17EE4BC63AEC971302920FB58232717057C69E796AD3F4C5A50FA01462C657D14B73C759F6C96C267DA'},
				{1: 1, 2: h'2946F5D8', 3: 110, 4: 1000, 5: 1, 6: h'A3015501636D35BC986E7A77EC1780D0F230975243B683A70255018878552180638549221000395DF8254D9355A8C9031903E8', 7: "test send-tx", 8: h'A19DC11452D6104FE8B9AD89BCC9C3F9F896C0DFB2AE908438FABF84BF1AB678A534A8881B13D0558FDBB60965ABC09D148614C148DE96AF20EF3B006C1B8524BC03CC00E20A195F42E2CDEF2325113EEA190E97B9BFC48EBE94DF0CB5AF4C84', 9: h'B994008A76E30AEBC14330C5307B61D0424513CA771253B8D4849FEA3EEAED37D3EB5E56A50D9005500859C84F09AAB7'}
			]
		}
	*/
	d, _ := hex.DecodeString("a301a80101021a622ace1d035820307cdf86a913baf3fbf00165dcea2428b7ba91507178cf1da3d7fc5f4498c65d045820bc0d63c3292fd7c8bf5b5a95e5c1f7d9f70f57e23130960ac6dd0d44e0bd4a50055820cb0ccc31a192186e5e01645aab5ce886bc9d1560a01e8100617c3fb7d6370cb80658209a664dbc4f48f92bcc4ae364f80a81ec4124ba78bc1484b0782fceed1e505aae075830b4b43dd593ffb8938e710dc6dacc9982182b5dec5f6a5d1c4a4d516f0f85130fd32d38d34e1aa3baed1f51a9e06dec140855019913b0689eefc8349183ec71fcf7fddfd2f035a502a50158206fd848bda84b325277eed41da0c36b84ac8b4efe96606e87c1a8634b2ac63d78020703840a120c10048112055830b0f50fa16f29747513c29ee2880812e119a3300b226f8bf3bf2f25a04157f6b2ec7a7093e882f176219a61d4b58c12fd0384a9010102449bbae37103186e041903e80501065833a30155017c6646c1c19839192e8c76b9da6f8ead5a31e7a40255014f6c19966599ce7796cd8c6d4c7dfeb24eac1941031903e8076c746573742073656e642d747808586083f16a0f72be00f653f364fe3756fcd4d569b276256c9e0719c12c6def2a82e92649475b00a699bd3046523b9c993e9d07f7b4bdcca90121b609d4f7ce3174c50ae4d88e39132a45d826ab7aca66a63083dce33ec7382c29384085ad62dd9aae0958308b4a3aef75c1bf419e79a33763d49ba3770138278627829539c42abebfc580f357f75a7ed825414cfd35c80634695690a9010102442fab6d2903186e041903e80501065833a301550198e01d56f0c0c346e5a45b6cddde89c98ae7ce4b02550192c557e2850f632e9901b52e1f3502d994776686031903e8076c746573742073656e642d7478085860b522ce35ce82962c20b5618990490a1f60bae007a582443303de0942625a16a5deb2bf816a43a867b499cf2c5dd3a33e057d6972872f3f741bdb72c8ad2cbb148f91055046324bfdedaa51c671a6d8cb61a3070bba07ff25b32d15326655313a095830b3ee7fa1dd553c0822da3e7cc3a6ffdf79dc7cbc7f10548801c691e0a2e2b0bea473c4e154e4749c3e559c6bf6d03065a90101024498981dac03186e041903e80501065833a3015501e302b94cbeeddcdfae85853ed1ff3a03f311769e02550174b82826f6e3988bbd9d825295cda46d3bea68a2031903e8076c746573742073656e642d7478085860b2bdf71b26b4eae9f41bf6bafc429e72ee26af86daef1f2fa2c2befb9d7ca959686b0bda8d05f9e260ee8b81ce93540704dfe47855fe6787d3705fcb60403eb52b54cac6e9719484c33822064fb1bab829bd759f64c82a48ecbf79c845b16927095830b9439e12f760f17ee4bc63aec971302920fb58232717057c69e796ad3f4c5a50fa01462c657d14b73c759f6c96c267daa9010102442946f5d803186e041903e80501065833a3015501636d35bc986e7a77ec1780d0f230975243b683a70255018878552180638549221000395df8254d9355a8c9031903e8076c746573742073656e642d7478085860a19dc11452d6104fe8b9ad89bcc9c3f9f896c0dfb2ae908438fabf84bf1ab678a534a8881b13d0558fdbb60965abc09d148614c148de96af20ef3b006c1b8524bc03cc00e20a195f42e2cdef2325113eea190e97b9bfc48ebe94df0cb5af4c84095830b994008a76e30aebc14330c5307b61d0424513ca771253b8d4849fea3eeaed37d3eb5e56a50d9005500859c84f09aab7")
	assert.NoError(t, b.Decode(d))
	d2, _ := b.Encode()
	assert.Equal(t, d, d2)

	// block header: a80101021a622ace1d035820307cdf86a913baf3fbf00165dcea2428b7ba91507178cf1da3d7fc5f4498c65d045820bc0d63c3292fd7c8bf5b5a95e5c1f7d9f70f57e23130960ac6dd0d44e0bd4a50055820cb0ccc31a192186e5e01645aab5ce886bc9d1560a01e8100617c3fb7d6370cb80658209a664dbc4f48f92bcc4ae364f80a81ec4124ba78bc1484b0782fceed1e505aae075830b4b43dd593ffb8938e710dc6dacc9982182b5dec5f6a5d1c4a4d516f0f85130fd32d38d34e1aa3baed1f51a9e06dec140855019913b0689eefc8349183ec71fcf7fddfd2f035a5
	expected1 := hash.CalcHash(d[2:225])
	assert.True(t, b.Hash().EqualsTo(expected1))

	expected2, _ := hash.FromString("bf630c1d3b5b3fb31be95250c56fc56b6072bf871dd50eb4c75b4d32f4bbfb2c")
	assert.Equal(t, b.Hash(), expected1)
	assert.Equal(t, b.Hash(), expected2)
	assert.Equal(t, b.Stamp(), hash.Stamp{0xbf, 0x63, 0x0c, 0x1d})
	assert.Equal(t, b.Header().TxsRoot(), b.Transactions().Root())
}

func TestMakeBlock(t *testing.T) {
	tmp := GenerateTestBlock(nil, nil)
	t.Run("Invalid block information, should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// Certificate is missed,
		MakeBlock(1, util.Now(), tmp.Transactions(),
			tmp.Header().PrevBlockHash(),
			tmp.Header().StateRoot(),
			nil,
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})

	t.Run("Invalid block information, should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// Invalid state root
		MakeBlock(1, util.Now(), tmp.Transactions(),
			tmp.Header().PrevBlockHash(),
			hash.UndefHash,
			tmp.PrevCertificate(),
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})

	t.Run("Invalid block information, should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// Invalid previous block hash
		MakeBlock(1, util.Now(), tmp.Transactions(),
			hash.UndefHash,
			tmp.Header().PrevBlockHash(),
			tmp.PrevCertificate(),
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})
}
