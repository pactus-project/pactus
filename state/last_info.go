package state

import (
	"encoding/json"
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

//
// Some thoughts about this structure:
//
// This structure helps the node to restore the last state upon restarting the node.
// We could also replay the last block to recover the last state, but there is a tradeoff here.
//
// There are two ways to safely replay the last block:
// 1- Have a snapshot of the database for current state and previous state. (which is not beautiful)
// 2- Commit block after receiving the next block. Note that the next block has adjustment or proof for the previous block
// However postponing committing a block has its own complexity.
//
// For now, we keep this structure as it is.
//
type lastInfo struct {
	LastHeight      int
	LastCommit      block.Commit
	LastReceiptHash crypto.Hash
	Committee       []int
	NextProposer    crypto.Address
}

func (st *state) saveLastInfo(height int, commit block.Commit, lastReceiptHash crypto.Hash, committee []int, proposer crypto.Address) {
	path := st.config.Store.Path + "/last_info.json"
	li := lastInfo{
		LastHeight:      height,
		LastCommit:      commit,
		LastReceiptHash: lastReceiptHash,
		Committee:       committee,
		NextProposer:    proposer,
	}

	bs, _ := json.Marshal(&li)

	if err := util.WriteFile(path, bs); err != nil {
		st.logger.Error("Unable to write last sate info", "err", err)
	}
}

func (st *state) loadLastInfo() (*lastInfo, error) {
	path := st.config.Store.Path + "/last_info.json"
	if !util.PathExists(path) {
		return nil, fmt.Errorf("Unable to load %v", path)
	}
	bs, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}
	li := new(lastInfo)
	err = json.Unmarshal(bs, li)
	if err != nil {
		return nil, err
	}
	return li, nil
}
