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

	header := block.NewHeader(
		uint(ch.Version()),
		time.Unix(ch.Time(), 0),
		bytesToHash(ch.TxsHash()),
		bytesToHash(ch.LastBlockHash()),
		bytesToHash(ch.NextValidatorsHash()),
		bytesToHash(ch.StateHash()),
		bytesToHash(ch.LastReceiptsHash()),
		bytesToHash(ch.LastCommitHash()),
		bytesToAddress(ch.ProposerAddress()),
	)

	txs := block.NewTxHashes()
	hashesList, _ := ctxs.Hashes()
	for i := 0; i < hashesList.Len(); i += 1 {
		txs.Append(bytesToHash(hashesList.At(i)))
	}

	commitersList, _ := clc.Commiters()
	commiters := make([]crypto.Address, commitersList.Len())
	for i := 0; i < commitersList.Len(); i += 1 {
		commiters[i] = bytesToAddress(commitersList.At(i))
	}
	signaturesList, _ := clc.Signatures()
	signatures := make([]crypto.Signature, signaturesList.Len())
	for i := 0; i < signaturesList.Len(); i += 1 {
		signatures[i] = bytesToSignature(signaturesList.At(i))
	}
	lastCommit := block.NewCommit(int(clc.Round()), commiters, signatures)
	if bi.Height == 1 {
		lastCommit = nil
	}
	block, err := block.NewBlock(header, txs, lastCommit)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	bi.Block = *block
	bi.Time = block.Header().Time()

	j, _ := json.MarshalIndent(bi, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(j))
}

func (s *Server) BlockByHeightHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Recovered in capnp call", "r", r)
		}
	}()
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
		io.WriteString(w, err.Error())
		return
	}

	s.WriteBlock(cbi, w)
}

func (s *Server) BlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Recovered in capnp call", "r", r)
		}
	}()
	b := s.server.Block(s.ctx, func(p capnp.ZarbServer_block_Params) error {
		vars := mux.Vars(r)
		hash, err := crypto.HashFromRawBytes([]byte(vars["hash"]))
		if err != nil {
			return err
		}
		p.SetHash(hash.RawBytes())
		return nil
	}).BlockInfo()

	cbi, err := b.Struct()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	s.WriteBlock(cbi, w)
}
