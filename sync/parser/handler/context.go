package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

const FlagInitialBlockDownload = 0x1

// TODO: Better way ?
type publishFn = func(payload.Payload)
type syncedFn = func()

type HandlerContext struct {
	state                state.StateFacade
	consensus            consensus.Consensus
	cache                *cache.Cache
	peerSet              *peerset.PeerSet
	moniker              string
	publicKey            crypto.PublicKey
	selfID               peer.ID
	publishFn            publishFn
	syncedFn             syncedFn
	blockPerMessage      int
	initialBlockDownload bool
	requestBlockInterval int
	logger               *logger.Logger
}

func NewHandlerContext(state state.StateFacade,
	consensus consensus.Consensus,
	cache *cache.Cache,
	peerSet *peerset.PeerSet,
	moniker string,
	publicKey crypto.PublicKey,
	selfID peer.ID,
	publishFn publishFn,
	syncedFn syncedFn,
	blockPerMessage int,
	initialBlockDownload bool,
	requestBlockInterval int,
	logger *logger.Logger) *HandlerContext {
	return &HandlerContext{
		state:                state,
		consensus:            consensus,
		cache:                cache,
		peerSet:              peerSet,
		moniker:              moniker,
		publicKey:            publicKey,
		selfID:               selfID,
		publishFn:            publishFn,
		syncedFn:             syncedFn,
		blockPerMessage:      blockPerMessage,
		initialBlockDownload: initialBlockDownload,
		requestBlockInterval: requestBlockInterval,
		logger:               logger,
	}
}

// peerIsInTheCommittee checks if the peer is a member of committee
func (ctx *HandlerContext) peerIsInTheCommittee(id peer.ID) bool {
	p := ctx.peerSet.GetPeer(id)
	if p == nil {
		return false
	}

	addr := p.PublicKey().Address()
	return ctx.state.IsInCommittee(addr)
}

// weAreInTheCommittee checks if we are a member of committee
func (ctx *HandlerContext) weAreInTheCommittee() bool {
	return ctx.state.IsInCommittee(ctx.publicKey.Address())
}

func (ctx *HandlerContext) addBlocksToCache(from int, blocks []block.Block) {
	for _, block := range blocks {
		ctx.cache.AddCertificate(block.LastCertificate())
		ctx.cache.AddBlock(from, &block)

		from++
	}
}

func (ctx *HandlerContext) addTransactionsToCache(trxs []tx.Tx) {
	ctx.cache.AddTransactions(trxs)
}

func (ctx *HandlerContext) tryCommitBlocks() {
	for {
		ourHeight := ctx.state.LastBlockHeight()
		b := ctx.cache.GetBlock(ourHeight + 1)
		if b == nil {
			break
		}
		c := ctx.cache.GetCertificate(b.Hash())
		if c == nil {
			break
		}
		for _, id := range b.TxIDs().IDs() {
			if tx := ctx.cache.GetTransaction(id); tx != nil {
				ctx.state.AddPendingTx(tx)
			}
		}
		ctx.logger.Trace("Committing block", "height", ourHeight+1, "block", b)
		if err := ctx.state.CommitBlock(ourHeight+1, *b, *c); err != nil {
			ctx.logger.Warn("Committing block failed", "block", b, "err", err, "height", ourHeight+1)
			// We will ask peers to send this block later ...
			break
		}
	}
}

func (ctx *HandlerContext) prepareBlocksAndTransactions(from, count int) ([]block.Block, []tx.Tx) {
	ourHeight := ctx.state.LastBlockHeight()

	if from > ourHeight {
		ctx.logger.Debug("We don't have block at this height", "height", from)
		return nil, nil
	}

	if from+count > ourHeight {
		count = ourHeight - from + 1
	}

	blocks := make([]block.Block, 0, count)
	trxs := make([]tx.Tx, 0)

	for h := from; h < from+count; h++ {
		b := ctx.cache.GetBlock(h)
		if b == nil {
			ctx.logger.Warn("Unable to find a block", "height", h)
			return nil, nil
		}
		for _, id := range b.TxIDs().IDs() {
			trx := ctx.cache.GetTransaction(id)
			if trx != nil {
				trxs = append(trxs, *trx)
			} else {
				ctx.logger.Debug("Unable to find a transaction", "id", id.Fingerprint())
				return nil, nil
			}
		}

		blocks = append(blocks, *b)
	}

	return blocks, trxs
}

func (ctx *HandlerContext) prepareTransactions(ids []tx.ID) []tx.Tx {
	trxs := make([]tx.Tx, 0, len(ids))

	for _, id := range ids {
		trx := ctx.cache.GetTransaction(id)
		if trx == nil {
			ctx.logger.Debug("Unable to find a transaction", "id", id.Fingerprint())
			continue
		}
		trxs = append(trxs, *trx)
	}
	return trxs
}

func (ctx *HandlerContext) updateSession(code payload.ResponseCode, sessionID int, initiator peer.ID, target peer.ID) {
	if target != ctx.selfID {
		return
	}

	switch code {
	case payload.ResponseCodeRejected:
		ctx.logger.Debug("session rejected, close session", "session-id", sessionID)
		ctx.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeBusy:
		// TODO: handle this situation
		ctx.logger.Debug("Peer is busy. close session", "session-id", sessionID)
		ctx.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeMoreBlocks:
		ctx.logger.Debug("Peer responding us. keep session open", "session-id", sessionID)

	case payload.ResponseCodeNoMoreBlocks:
		ctx.logger.Debug("Peer has no more block. close session", "session-id", sessionID)
		ctx.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeSynced:
		ctx.logger.Debug("Peer infomed us we are synced. close session", "session-id", sessionID)
		ctx.peerSet.CloseSession(sessionID)
		ctx.syncedFn()
	}

	s := ctx.peerSet.FindSession(sessionID)
	if s == nil {
		ctx.logger.Debug("Session not found or closed", "session-id", sessionID)
	} else {
		s.LastResponseCode = code
		s.LastActivityAt = util.Now()
	}
}
