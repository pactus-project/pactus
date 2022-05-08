package capnp

func (zs *zarbServer) GetConsensusInfo(args ZarbServer_getConsensusInfo) error {
	height, round := zs.consensus.HeightRound()
	votes := zs.consensus.AllVotes()
	res, _ := args.Results.NewResult()
	res.SetHeight(height)
	res.SetRound(round)

	capVotes, _ := res.NewVotes(int32(len(votes)))
	for i, vote := range votes {
		capVote := capVotes.At(i)
		capVote.SetType(int8(vote.Type()))
		capVote.SetRound(vote.Round())
		err := capVote.SetVoter(vote.Signer().String())
		if err != nil {
			return err
		}
		err = capVote.SetBlockHash(vote.BlockHash().Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}
