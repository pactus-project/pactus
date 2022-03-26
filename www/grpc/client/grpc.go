package grpcclient

import (
	"context"
	"encoding/hex"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc"
)

var (
	zclient zarb.ZarbClient
)

func GetStamp(rpcEndpoint string) (hash.Stamp, error) {
	client, err := GetRPCClient(rpcEndpoint)
	if err != nil {
		return hash.Stamp{}, err
	}

	info, err := client.GetBlockchainInfo(context.Background(), &zarb.BlockchainInfoRequest{})
	if err != nil {
		return hash.Stamp{}, err
	}
	return hash.StampFromString(string(info.LastBlockHash))
}

func GetSequence(rpcEndpoint string, addr crypto.Address) (int32, error) {
	client, err := GetRPCClient(rpcEndpoint)
	if err != nil {
		return 0, err
	}

	acc, err := client.GetAccount(context.Background(), &zarb.AccountRequest{Address: addr.Bytes()})
	if err != nil {
		return 0, err
	}

	return acc.Account.Sequence + 1, nil
}

func SendTx(rpcEndpoint string, payload []byte) (string, error) {
	client, err := GetRPCClient(rpcEndpoint)
	if err != nil {
		return "", err
	}

	res, err := client.SendRawTransaction(context.Background(), &zarb.SendRawTransactionRequest{
		Data: hex.EncodeToString(payload),
	})

	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func GetRPCClient(rpcEndpoint string) (zarb.ZarbClient, error) {
	if zclient != nil {
		return zclient, nil
	}

	conn, err := grpc.Dial(rpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	zclient = zarb.NewZarbClient(conn)
	return zclient, nil
}
