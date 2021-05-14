package cmd

import (
	"context"
	"encoding/hex"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc"
)

var (
	zclient zarb.ZarbClient
)

func GetStamp(client zarb.ZarbClient) (crypto.Hash, error) {
	info, err := client.GetBlockchainInfo(context.Background(), &zarb.BlockchainInfoRequest{})
	if err != nil {
		return crypto.Hash{}, err
	}
	return crypto.HashFromString(info.Stamp)
}

func GetSequence(client zarb.ZarbClient, addr crypto.Address) (int, error) {
	a, err := client.GetAccount(context.Background(), &zarb.AccountRequest{Address: addr.String()})
	if err != nil {
		return 0, err
	}
	acc := account.NewAccount(addr, 0)

	err = acc.Decode(a.Data)
	if err != nil {
		return 0, err
	}
	return acc.Sequence() + 1, nil
}

func SendTx(client zarb.ZarbClient, payload []byte) (string, error) {
	res, err := client.SendRawTransaction(context.Background(), &zarb.SendRawTransactionRequest{
		Data: hex.EncodeToString(payload),
	})

	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func GetRPCClient(rpc string) (zarb.ZarbClient, error) {
	if zclient != nil {
		return zclient, nil
	}

	conn, err := grpc.Dial(rpc, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	zclient = zarb.NewZarbClient(conn)
	return zclient, nil
}
