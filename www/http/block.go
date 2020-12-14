package http

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) WriteBlock(cbi capnp.BlockInfo, w http.ResponseWriter) {
	cb, _ := cbi.Block()
	ch, _ := cb.Header()
	ctxs, _ := cb.Txs()
	clc, _ := cb.LastCommit()

	var bi BlockInfo

	d, _ := cbi.Data()
	bi.Hash = bytesToHash(cbi.Hash())
	bi.Height = int(cbi.Height())
	bi.Data = hex.EncodeToString(d)

	var lastCommit *block.Commit
	if bi.Height == 1 {
		lastCommit = nil
	} else {
		list, _ := clc.Committers()
		committers := make([]block.Committer, list.Len())
		for i := 0; i < list.Len(); i++ {
			c := list.At(i)
			committers[i].Address = bytesToAddress(c.Address())
			committers[i].Status = int(c.Status())
		}
		sig := bytesToSignature(clc.Signature())
		lastCommit = block.NewCommit(int(clc.Round()), committers, sig)
	}

	header := block.NewHeader(
		uint(ch.Version()),
		time.Unix(ch.Time(), 0),
		bytesToHash(ch.TxsHash()),
		bytesToHash(ch.LastBlockHash()),
		bytesToHash(ch.CommittersHash()),
		bytesToHash(ch.StateHash()),
		bytesToHash(ch.LastReceiptsHash()),
		bytesToHash(ch.LastCommitHash()),
		bytesToAddress(ch.ProposerAddress()))

	txs := block.NewTxHashes()
	hashesList, _ := ctxs.Hashes()
	for i := 0; i < hashesList.Len(); i++ {
		txs.Append(bytesToHash(hashesList.At(i)))
	}

	block, err := block.NewBlock(header, lastCommit, txs)
	if err != nil {
		if _, err := io.WriteString(w, err.Error()); err != nil {
			s.logger.Error("Unable to write string", "err", err)
		}
		return
	}
	bi.Block = *block
	bi.Time = block.Header().Time()

	j, _ := json.MarshalIndent(bi, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, string(j)); err != nil {
		s.logger.Error("Unable to write string", "err", err)
	}
}

func (s *Server) BlockByHeightHandler(w http.ResponseWriter, r *http.Request) {
	b := s.server.BlockAt(s.ctx, func(p capnp.ZarbServer_blockAt_Params) error {
		vars := mux.Vars(r)
		height, err := strconv.Atoi(vars["height"])
		if err != nil {
			return err
		}
		p.SetHeight(uint32(height))
		return nil
	}).BlockInfo()

	cbi, err := b.Struct()
	if err != nil {
		if _, err = io.WriteString(w, err.Error()); err != nil {
			s.logger.Error("Unable to write string", "err", err)
		}
		return
	}

	s.WriteBlock(cbi, w)
}

func (s *Server) BlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	b := s.server.Block(s.ctx, func(p capnp.ZarbServer_block_Params) error {
		vars := mux.Vars(r)
		hash, err := crypto.HashFromRawBytes([]byte(vars["hash"]))
		if err != nil {
			return err
		}
		if err := p.SetHash(hash.RawBytes()); err != nil {
			return err
		}
		return nil
	}).BlockInfo()

	cbi, err := b.Struct()
	if err != nil {
		if _, err = io.WriteString(w, err.Error()); err != nil {
			s.logger.Error("Unable to write string", "err", err)
		}
		return
	}

	s.WriteBlock(cbi, w)
}
