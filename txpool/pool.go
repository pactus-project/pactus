package txpool

import (
	"container/list"
	"fmt"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"

	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

type TxPoolReader interface {
	PendingTx(hash crypto.Hash) (*tx.Tx, bool)
}

type TxPool struct {
	lk deadlock.RWMutex

	config       *config.Config
	pendingsList *list.List
	pendingsMap  map[crypto.Hash]*list.Element
	broadcastCh  chan message.Message
	logger       *logger.Logger
}

func NewTxPool(
	conf *config.Config,
	broadcastCh chan message.Message) (*TxPool, error) {
	pool := &TxPool{
		config:       conf,
		pendingsList: list.New(),
		pendingsMap:  make(map[crypto.Hash]*list.Element),
		broadcastCh:  broadcastCh,
	}

	pool.logger = logger.NewLogger("_pool", pool)
	return pool, nil
}

func (pool *TxPool) AppendTxs(txs []tx.Tx) error {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	for _, tx := range txs {
		pool.appendTx(&tx)
	}
	return nil
}

func (pool *TxPool) AppendTx(tx *tx.Tx) error {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return pool.appendTx(tx)
}

func (pool *TxPool) appendTx(tx *tx.Tx) error {
	if pool.pendingsList.Len() >= pool.config.TxPool.MaxSize {
		return errors.Errorf(errors.ErrGeneric, "Tx pool is full (%d txs)", pool.pendingsList.Len())
	}

	_, found := pool.pendingsMap[tx.Hash()]
	if found {
		return errors.Errorf(errors.ErrGeneric, "We already have this transaction in our pool")
	}
	// TODO:
	// validate transaction

	el := pool.pendingsList.PushFront(tx)
	pool.pendingsMap[tx.Hash()] = el

	// TODO:
	//go pool.syncer.BroadcastTx(tx)

	return nil
}

func (pool *TxPool) RemoveTx(hash crypto.Hash) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	el, found := pool.pendingsMap[hash]
	if !found {
		return
	}

	pool.pendingsList.Remove(el)
	delete(pool.pendingsMap, hash)
}

func (pool *TxPool) PendingTx(hash crypto.Hash) (*tx.Tx, bool) {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	el, found := pool.pendingsMap[hash]
	if !found {
		return nil, false
	}

	return el.Value.(*tx.Tx), true
}

func (pool *TxPool) HasTx(hash crypto.Hash) bool {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	_, found := pool.pendingsMap[hash]
	return found
}

func (pool *TxPool) Fingerprint() string {
	return fmt.Sprintf("{%v}", pool.pendingsList.Len())
}
