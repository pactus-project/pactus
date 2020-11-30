package sortition

import (
	"github.com/zarbchain/zarb-go/crypto"
)

type Sortition struct {
	vrf VRF
}

func NewSortition(signer crypto.Signer) *Sortition {
	return &Sortition{
		//vrf: NewVRF(signer),
	}
}

// // Evaluate return the vrf for self choosing to be a validator
// func (s *Sortition) Evaluate(blockHash crypto.Hash) {
// 	addr := s.signer.Address()
// 	totalStake, valStake := s.getTotalStake(addr)
// 	s.vrf.SetMax(totalStake)
// 	index, proof := s.vrf.Evaluate(blockHash)

// 	if index < valStake {
// 		logger.Info("This validator is chosen to be in set", "height", blockHeight, "address", addr, "stake", valStake)

// 		/// TODO: better way????
// 		val, err := s.state.GetValidator(s.Address())
// 		if err != nil {
// 			return
// 		}

// 		tx, _ := tx.NewSortitionTx(
// 			s.signer.Address(),
// 			blockHeight,
// 			val.Sequence()+1,
// 			s.sortitionFee,
// 			index,
// 			proof)

// 		txEnv := txs.Enclose(s.chainID, tx)
// 		err = txEnv.Sign(s.signer)
// 		if err != nil {
// 			return
// 		}

// 		// TODO:: better way?????
// 		codec := txs.NewAminoCodec()
// 		bs, err := codec.MarshalBinaryLengthPrefixed(txEnv)
// 		if err != nil {
// 			return
// 		}

// 		res, err := tmRPC.BroadcastTxAsync(bs)
// 		if err != nil {
// 			return
// 		}

// 		if res != nil {
// 			/// TODO: log result
// 		}
// 	}
// }

// func (s *Sortition) Verify(blockHash []byte, pb crypto.PublicKey, index uint64, proof []byte) bool {
// 	addr := pb.ValidatorAddress()
// 	totalStake, valStake := s.getTotalStake(addr)

// 	// Note: totalStake can be changed by time on verifying
// 	// So we calculate the index again
// 	s.vrf.SetMax(totalStake)

// 	index2, result := s.vrf.Verify(blockHash, pb, proof)
// 	if !result {
// 		logger.Warn("Unable to verify a sortition tx", "blockhash", blockHash, "Address", addr)
// 		return false
// 	}

// 	return index2 < valStake
// }

// func (s *Sortition) Address() crypto.Address {
// 	return s.signer.Address()
// }

// func (s *Sortition) getTotalStake(addr crypto.Address) (totalStake uint64, validatorStake uint64) {
// 	totalStake = 0
// 	validatorStake = 0

// 	s.state.IterateValidators(func(validator *validator.Validator) (stop bool) {
// 		totalStake += validator.Stake()

// 		if addr == validator.Address() {
// 			validatorStake = validator.Stake()
// 		}

// 		return false
// 	})

// 	return
// }
