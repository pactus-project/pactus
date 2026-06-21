package fake

import (
	"testing"

	"github.com/pactus-project/pactus/consensus"
	consmgr "github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	wltmgr "github.com/pactus-project/pactus/wallet/manager"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/zmq"
)

type FakeGRPCServer struct {
	FakeState     *state.FakeState
	FakeSync      *sync.FakeSync
	FakeCons      *consensus.FakeConsensus
	FakeConsMgr   *consmgr.FakeConsensusManager
	FakeWalletMgr *wltmgr.FakeWalletManager
	Server        *grpc.Server
}

func NewFakeGRPCServer(t *testing.T, ts *testsuite.TestSuite, conf *grpc.Config) *FakeGRPCServer {
	t.Helper()

	// for saving test wallets in temp directory
	t.Chdir(util.TempDirPath())

	cmt, _ := ts.GenerateTestCommittee(51)
	fakeState := state.NewFakeState(ts, cmt)
	fakeNetwork := network.MockingNetwork(ts, ts.RandPeerID())
	fakeSync := sync.NewFakeSync(ts)
	fakeConsMgr := consmgr.NewFakeConsensusManager(ts)
	fakeCons := consensus.NewFakeConsensus(ts)
	fakeWalletMgr := wltmgr.NewFakeWalletManager(ts)

	pub, _ := ts.RandBLSKeyPair()
	fakeCons.EXPECT().ConsensusKey().Return(pub).AnyTimes()

	fakeConsMgr.EXPECT().Instances().Return([]consensus.ConsensusReader{fakeCons}).AnyTimes()

	zmqPublishers := []zmq.Publisher{
		zmq.MockingPublisher("zmq_address", "zmq_topic", 100),
	}

	server := grpc.NewServer(
		t.Context(), conf,
		fakeState, fakeSync, fakeNetwork, fakeConsMgr,
		fakeWalletMgr, zmqPublishers,
	)

	t.Cleanup(func() {
		server.StopServer()
	})

	return &FakeGRPCServer{
		FakeState:     fakeState,
		FakeSync:      fakeSync,
		FakeCons:      fakeCons,
		FakeConsMgr:   fakeConsMgr,
		FakeWalletMgr: fakeWalletMgr,
		Server:        server,
	}
}
