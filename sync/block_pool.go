package sync

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
)

type BlockPool struct {
	lk deadlock.RWMutex

	blocks  map[int]*block.Block
	commits map[crypto.Hash]*block.Commit
}

func NewBlockPool() *BlockPool {
	return &BlockPool{
		blocks:  make(map[int]*block.Block),
		commits: make(map[crypto.Hash]*block.Commit),
	}
}

func (pool *BlockPool) AppendCommit(blockHash crypto.Hash, commit *block.Commit) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	bc, has := pool.commits[blockHash]
	if has {
		if !bc.Hash().EqualsTo(commit.Hash()) {
			logger.Debug("Different commit for the same block", "hash", blockHash)
		}
	}
	pool.commits[blockHash] = commit
}

func (pool *BlockPool) AppendBlock(height int, block block.Block) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	bp, has := pool.blocks[height]
	if has {
		if !bp.Hash().EqualsTo(block.Hash()) {
			logger.Warn("Different block for the same height, overwrite the previous one", "height", height)
		}
	}

	pool.blocks[height] = &block
}

func (pool *BlockPool) Block(height int) *block.Block {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return pool.blocks[height]
}

func (pool *BlockPool) Commit(hash crypto.Hash) *block.Commit {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return pool.commits[hash]
}

func (pool *BlockPool) RemoveCommit(hash crypto.Hash) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	delete(pool.commits, hash)
}

func (pool *BlockPool) RemoveBlock(height int) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	delete(pool.blocks, height)
}

func (pool *BlockPool) BlockLen() int {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return len(pool.commits)
}

func (pool *BlockPool) CommitLen() int {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return len(pool.commits)
}
