package state

import (
	"encoding/json"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/util"
)

type lastInfo struct {
	LastHeight int
	LastCommit *block.Commit
}

func (st *state) saveLastInfo(height int, commit *block.Commit) error {
	path := st.config.Store.Path + "/last_info.json"
	li := lastInfo{
		LastHeight: height,
		LastCommit: commit,
	}

	bs, err := json.Marshal(&li)
	if err != nil {
		return err
	}
	util.WriteFile(path, bs)
	return nil
}

func (st *state) loadLastInfo() (int, *block.Commit, error) {
	path := st.config.Store.Path + "/last_info.json"
	if !util.PathExists(path) {
		return 0, nil, nil
	}
	bs, err := util.ReadFile(path)
	if err != nil {
		return 0, nil, err
	}
	li := new(lastInfo)
	err = json.Unmarshal(bs, li)
	if err != nil {
		return 0, nil, err
	}
	return li.LastHeight, li.LastCommit, nil
}
