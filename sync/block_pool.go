package sync

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
)

type BlockPool struct {
	lk deadlock.RWMutex

	blocks map[int]*block.Block
	logger *logger.Logger
}

func NewBlockPool(logger *logger.Logger) *BlockPool {
	return &BlockPool{
		blocks: make(map[int]*block.Block),
		logger: logger,
	}
}

func (pool *BlockPool) AppendBlock(height int, block block.Block) error {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	bp, has := pool.blocks[height]
	if has {
		if !bp.Hash().EqualsTo(block.Hash()) {
			pool.logger.Warn("Different blocks for same height", "height", height)
			delete(pool.blocks, height)
			return errors.Error(errors.ErrInvalidBlock)
		}
	} else {
		pool.blocks[height] = &block
	}
	return nil
}

func (pool *BlockPool) Block(height int) *block.Block {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return pool.blocks[height]
}

func (pool *BlockPool) RemoveBlock(height int) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	delete(pool.blocks, height)
}

func (pool *BlockPool) Size() int {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return len(pool.blocks)
}
