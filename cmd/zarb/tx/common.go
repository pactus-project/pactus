package tx

import (
	"encoding/hex"
	"strings"

	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func signAndPublish(trx *tx.Tx, keyfile, auth, rpcEndpoint string) {
	//sign transaction
	k, err := key.DecryptKeyFile(keyfile, auth)
	if err != nil {
		cmd.PrintErrorMsg("Couldn't retrieve the key: %v", err)
		return
	}
	k.ToSigner().SignMsg(trx)

	//show
	cmd.PrintWarnMsg("you are about to publish:")
	cmd.PrintJSONObject(trx)

	cmd.PrintLine()
	signedTrx, _ := trx.Encode()
	cmd.PrintInfoMsg("raw signed transaction payload:\n%v", hex.EncodeToString(signedTrx))

	cmd.PrintLine()
	confirm := cmd.PromptInput("This operation is \"not reversible\". Are you sure [yes/no]? ")
	if !strings.HasPrefix(strings.ToLower(confirm), "yes") {
		cmd.PrintWarnMsg("Opration aborted!")
		return
	}

	// publish
	if id, err := util.SendTx(rpcEndpoint, signedTrx); err != nil {
		cmd.PrintErrorMsg("Couldn't publish transaction: %v", err)
	} else {
		cmd.PrintSuccessMsg("Transaction sent with ID: %v", id)
	}
}

func promptRPCEndpoint(rpcEndpoint string) string {
	if len(rpcEndpoint) < 0 {
		return cmd.PromptInput("gRPC Endpoint: ")
	}
	return rpcEndpoint
}
