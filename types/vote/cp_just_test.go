package vote_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestBasicCheckCPJust(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("JustInitNo: no quorum certificate", func(t *testing.T) {
		just := vote.JustInitNo{
			QCert: nil,
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, vote.BasicCheckError{Reason: "quorum certificate is not set"})
	})

	t.Run("JustInitNo: invalid quorum certificate", func(t *testing.T) {
		just := vote.JustInitNo{
			QCert: ts.GenerateTestCertificate(0),
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, certificate.BasicCheckError{Reason: "height is not positive: 0"})
	})

	t.Run("JustPreVoteSoft: no quorum certificate", func(t *testing.T) {
		just := vote.JustPreVoteSoft{
			QCert: nil,
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, vote.BasicCheckError{Reason: "quorum certificate is not set"})
	})

	t.Run("JustPreVoteHard: no quorum certificate", func(t *testing.T) {
		just := vote.JustPreVoteHard{
			QCert: nil,
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, vote.BasicCheckError{Reason: "quorum certificate is not set"})
	})

	t.Run("JustMainVoteConflict: no quorum certificate", func(t *testing.T) {
		just := vote.JustMainVoteConflict{
			JustNo:  nil,
			JustYes: nil,
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, vote.BasicCheckError{Reason: "quorum certificate is not set"})
	})

	t.Run("JustMainVoteConflict: invalid quorum certificate", func(t *testing.T) {
		just := vote.JustMainVoteConflict{
			JustYes: &vote.JustInitYes{},
			JustNo: &vote.JustInitNo{
				QCert: nil,
			},
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, vote.BasicCheckError{Reason: "quorum certificate is not set"})
	})

	t.Run("JustMainVoteNoConflict: no quorum certificate", func(t *testing.T) {
		just := vote.JustMainVoteNoConflict{
			QCert: nil,
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, vote.BasicCheckError{Reason: "quorum certificate is not set"})
	})

	t.Run("JustDecided: no quorum certificate", func(t *testing.T) {
		just := vote.JustDecided{
			QCert: nil,
		}

		err := just.BasicCheck()
		require.ErrorIs(t, err, vote.BasicCheckError{Reason: "quorum certificate is not set"})
	})
}
