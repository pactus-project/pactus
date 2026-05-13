package mock

import (
	"testing"

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

type MockGRPCServer struct {
	MockState     *state.MockState
	MockSync      *sync.MockSync
	MockCons      *consmgr.MockReader
	MockConsMgr   *consmgr.MockManagerReader
	MockWalletMgr *wltmgr.MockIManager
	Server        *grpc.Server
}

func SetupServer(t *testing.T, ts *testsuite.TestSuite, conf *grpc.Config) *MockGRPCServer {
	t.Helper()

	// for saving test wallets in temp directory
	t.Chdir(util.TempDirPath())

	mockState := state.MockingState(ts)
	mockNet := network.MockingNetwork(ts, ts.RandPeerID())
	mockSync := sync.MockingSync(ts)
	mockConsMgr := consmgr.NewMockManagerReader(ts.Ctrl)
	mockCons := consmgr.NewMockReader(ts.Ctrl)

	pub, _ := ts.RandBLSKeyPair()
	mockCons.EXPECT().ConsensusKey().Return(pub).AnyTimes()

	mockConsMgr.EXPECT().Instances().Return([]consmgr.Reader{mockCons}).AnyTimes()

	mockState.CommitTestBlocks(10)
	mockWalletMgr := wltmgr.NewMockIManager(ts.MockingController())

	zmqPublishers := []zmq.Publisher{
		zmq.MockingPublisher("zmq_address", "zmq_topic", 100),
	}

	server := grpc.NewServer(t.Context(), conf,
		mockState, mockSync, mockNet, mockConsMgr,
		mockWalletMgr, zmqPublishers,
	)

	t.Cleanup(func() {
		server.StopServer()
	})

	return &MockGRPCServer{
		MockState:     mockState,
		MockSync:      mockSync,
		MockCons:      mockCons,
		MockConsMgr:   mockConsMgr,
		MockWalletMgr: mockWalletMgr,
		Server:        server,
	}
}
