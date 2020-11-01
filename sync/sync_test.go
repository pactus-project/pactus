package sync

import (
	"context"
	"reflect"
	"testing"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var (
	privVal          *validator.PrivValidator
	genDoc           *genesis.Genesis
	expectedPayloads []message.Payload
	ctx              context.Context
	generalTopic     *pubsub.Topic
	txTopic          *pubsub.Topic
	blockTopic       *pubsub.Topic
	consensusTopic   *pubsub.Topic
)

func init() {
	val, key := validator.GenerateTestValidator()
	acc := account.NewAccount(crypto.MintbaseAddress)
	acc.SetBalance(21000000000000)
	privVal = validator.NewPrivValidator(key)
	genDoc = genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val})
	expectedPayloads = make([]message.Payload, 0)
	ctx = context.Background()

}

func newTestSynchronizer(t *testing.T, privVal *validator.PrivValidator) (*Synchronizer, state.State) {
	syncConf := DefaultConfig()
	consConf := consensus.DefaultConfig()
	consConf.TimeoutPrevote = 1 * time.Millisecond
	consConf.TimeoutPrecommit = 1 * time.Millisecond
	consConf.TimeoutPropose = 1 * time.Millisecond
	consConf.NewRoundDeltaDuration = 0
	stateConf := state.DefaultConfig()
	stateConf.Store.Path = util.TempDirName()
	txPoolConf := txpool.DefaultConfig()

	loggerConfig := logger.DefaultConfig()
	loggerConfig.Levels["default"] = "error"
	logger.InitLogger(loggerConfig)

	ch := make(chan message.Message, 10)
	go func() {
		for {
			select {
			case <-ch:
			default:
			}
		}
	}()

	txPool, _ := txpool.NewTxPool(txPoolConf, ch)
	st, _ := state.LoadOrNewState(stateConf, genDoc, privVal.Address(), txPool)

	netConfig := network.DefaultConfig()
	netConfig.NodeKey = util.TempFilename()
	net, _ := network.NewNetwork(netConfig)
	cons, _ := consensus.NewConsensus(consConf, st, privVal, ch)
	sync, err := NewSynchronizer(syncConf, privVal.Address(), st, cons, txPool, net, ch)

	require.NoError(t, err)

	return sync, st
}

func newTestNetwork(t *testing.T) *network.Network {
	var err error
	netConfig := network.DefaultConfig()
	netConfig.NodeKey = util.TempFilename()
	net, _ := network.NewNetwork(netConfig)
	generalTopic, err = net.JoinTopic("general")
	assert.NoError(t, err)
	generalSub, err := generalTopic.Subscribe()
	assert.NoError(t, err)
	txTopic, err = net.JoinTopic("tx")
	assert.NoError(t, err)
	txSub, err := txTopic.Subscribe()
	assert.NoError(t, err)
	blockTopic, err = net.JoinTopic("block")
	assert.NoError(t, err)
	blockSub, err := blockTopic.Subscribe()
	assert.NoError(t, err)
	consensusTopic, err = net.JoinTopic("consensus")
	assert.NoError(t, err)
	consensusSub, err := consensusTopic.Subscribe()
	assert.NoError(t, err)

	go func() {
		for {
			m, err := txSub.Next(ctx)
			assert.NoError(t, err)
			parsMessage(t, m)
		}
	}()

	go func() {
		for {
			m, err := blockSub.Next(ctx)
			assert.NoError(t, err)
			parsMessage(t, m)
		}
	}()

	go func() {
		for {
			m, err := generalSub.Next(ctx)
			assert.NoError(t, err)
			parsMessage(t, m)
		}
	}()

	go func() {
		for {
			m, err := consensusSub.Next(ctx)
			assert.NoError(t, err)
			parsMessage(t, m)
		}
	}()

	return net
}

func remove(slice []message.Payload, s int) []message.Payload {
	return append(slice[:s], slice[s+1:]...)
}

func parsMessage(t *testing.T, m *pubsub.Message) {
	msg := new(message.Message)
	err := msg.UnmarshalCBOR(m.Data)
	assert.NoError(t, err)

	for i, p := range expectedPayloads {
		if reflect.DeepEqual(p, msg.Payload) {
			expectedPayloads = remove(expectedPayloads, i)
		}
	}

}

func TestReceiveInvalidBlock(t *testing.T) {
	_, _, key := crypto.GenerateTestKeyPair()
	privVal2 := validator.NewPrivValidator(key)

	sync1, st1 := newTestSynchronizer(t, privVal)
	sync2, _ := newTestSynchronizer(t, privVal2)
	net := newTestNetwork(t)
	net.Start()

	blocks := make([]block.Block, 11)
	for i := 0; i < 11; i++ {
		blocks[i] = st1.ProposeBlock()
		v := vote.NewPrecommit(i+1, 0, blocks[i].Hash(), privVal.Address())
		privVal.SignMsg(v)
		sig := v.Signature()
		commit := block.NewCommit(0, []crypto.Address{privVal.Address()}, []crypto.Signature{*sig})

		st1.ApplyBlock(blocks[i], *commit)
	}

	invalidBlock, _ := block.GenerateTestBlock(nil)
	sync2.blockPool.AppendBlock(1, invalidBlock)

	expectedPayloads = append(expectedPayloads, message.NewSalamMessage(11).Payload)
	expectedPayloads = append(expectedPayloads, message.NewSalamMessage(0).Payload)
	expectedPayloads = append(expectedPayloads, message.NewBlocksReqMessage(1, 2, crypto.UndefHash).Payload)

	sync1.Start()
	sync2.Start()

	time.Sleep(5000 * time.Millisecond)

}
