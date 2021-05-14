package tx

import (
	"strings"

	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/tx"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func signAndPublish(trx *tx.Tx, keyfile, auth string, grpcClient zarb.ZarbClient) {
	//sign transaction
	k, err := key.DecryptKeyFile(keyfile, auth)
	if err != nil {
		cmd.PrintErrorMsg("Couldn't retrieve Key: %v", err)
		return
	}
	k.ToSigner().SignMsg(trx)

	cmd.PrintWarnMsg("you are about to publish:")
	cmd.PrintJSONObject(trx)
	confirm := cmd.PromptInput("are you sure [y/n]:")
	if !strings.HasPrefix(strings.ToLower(confirm), "y") {
		cmd.PrintWarnMsg("Opration aborted!")
		return
	}

	// publish
	signedTrx, _ := trx.Encode()
	if id, err := cmd.SendTx(grpcClient, signedTrx); err != nil {
		cmd.PrintErrorMsg("Couldn't publish transaction: %v", err)
		return
	} else {
		cmd.PrintSuccessMsg("transaction sent with Id: %v", id)
	}
}
