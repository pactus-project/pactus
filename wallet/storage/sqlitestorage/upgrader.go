package sqlitestorage

import (
	"fmt"

	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/pactus/crypto"
)

func (s *Storage) upgrade() error {
	if !s.info.Network.IsMainnet() {
		crypto.ToTestnetHRP()
	}

	switch s.info.Version {
	case Version1:
		vlt := s.Vault()
		vlt.Purposes.PurposeBIP44.NextSexp256k1Index = 0
		err := s.saveVault(vlt)
		if err != nil {
			return err
		}

		s.info.Version = Version2
		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d", Version1, Version2))

		return s.updateWalletEntry(keyVersion, fmt.Sprintf("%d", s.info.Version))

	case Version2:
		// Latest version, no need to upgrade.
		return nil

	default:
		return UnsupportedVersionError{
			WalletVersion:    s.info.Version,
			SupportedVersion: VersionLatest,
		}
	}
}
