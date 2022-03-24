package tx

import (
	"encoding/hex"
	"strings"

	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/tx"
	grpcclient "github.com/zarbchain/zarb-go/www/grpc/client"
)

func signAndPublish(trx *tx.Tx, keyfile, auth string, rpcEndpoint *string) {
	//sign transaction
	k, err := key.DecryptKeyFile(keyfile, auth)
	if err != nil {
		cmd.PrintErrorMsg("Couldn't retrieve the key: %v", err)
		return
	}
	k.ToSigner().SignMsg(trx)

	//show
	cmd.PrintWarnMsg("Your transaction:")
	cmd.PrintJSONObject(trx)
	cmd.PrintLine()

	signedTrx, _ := trx.Bytes()

	if rpcEndpoint == nil || *rpcEndpoint == "" {
		//no endpoint specified just print the raw payload and leave
		cmd.PrintInfoMsg("raw signed transaction payload:\n%v", hex.EncodeToString(signedTrx))

	} else {

		confirm := cmd.PromptInput("This operation is \"not reversible\". Are you sure [yes/no]? ")
		if !strings.HasPrefix(strings.ToLower(confirm), "yes") {
			cmd.PrintWarnMsg("Opration aborted!")
			return
		}
		// publish
		if id, err := grpcclient.SendTx(*rpcEndpoint, signedTrx); err != nil {
			cmd.PrintErrorMsg("Couldn't publish transaction: %v", err)
		} else {
			cmd.PrintSuccessMsg("Transaction sent with ID: %v", id)
		}
	}
}

func promptRPCEndpoint(rpcEndpoint *string) string {
	if rpcEndpoint == nil || len(*rpcEndpoint) <= 0 {
		*rpcEndpoint = cmd.PromptInput("gRPC Endpoint: ")
	}
	return *rpcEndpoint
}
