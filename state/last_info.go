package state

import (
	"encoding/json"
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

// TODO:
// If we can replay the last block, then we don''t need this ugly structure
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
