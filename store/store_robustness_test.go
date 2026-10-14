package store

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// putTxIndex writes a transaction index entry directly in the database.
func putTxIndex(t *testing.T, td *testData, txID []byte, height types.Height, offset, length uint32) {
	t.Helper()

	buf := new(bytes.Buffer)
	require.NoError(t, encoding.WriteElements(buf, &height, &offset, &length))
	require.NoError(t, td.store.db.Put(append(txPrefix, txID...), buf.Bytes(), nil))
}

func TestIterateAccountsCanBeStopped(t *testing.T) {
	td := setup(t, nil)

	for i := int32(0); i < 5; i++ {
		td.store.UpdateAccount(td.RandAccAddress(), account.NewAccount(i))
	}
	require.NoError(t, td.store.WriteBatch())

	calls := 0
	td.store.IterateAccounts(func(crypto.Address, *account.Account) bool {
		calls++

		return true
	})
	assert.Equal(t, 1, calls)

	// The iterator must be released, even if the consumer stopped the iteration.
	// A leaked iterator keeps the database files open for the whole life of the node.
	alive, err := td.store.db.GetProperty("leveldb.aliveiters")
	require.NoError(t, err)
	assert.Equal(t, "0", alive)
}

func TestIterateValidatorsCanBeStopped(t *testing.T) {
	td := setup(t, nil)

	for i := int32(0); i < 5; i++ {
		pub, _ := td.RandBLSKeyPair()
		td.store.UpdateValidator(validator.NewValidator(pub, i))
	}

	calls := 0
	td.store.IterateValidators(func(*validator.Validator) bool {
		calls++

		return true
	})
	assert.Equal(t, 1, calls)
}

func TestNewStoreFailures(t *testing.T) {
	t.Run("data path is a file", func(t *testing.T) {
		conf := testConfig()
		conf.Path = filepath.Join(t.TempDir(), "not-a-directory")
		require.NoError(t, os.WriteFile(conf.Path, []byte("data"), 0o600))

		_, err := NewStore(conf)
		require.Error(t, err)
	})

	t.Run("a block of the cache window is missing", func(t *testing.T) {
		td := setup(t, nil)
		require.NoError(t, td.store.db.Delete(blockKey(5), nil))
		td.store.Close()

		_, err := NewStore(td.store.config)
		require.Error(t, err)

		// Failing to open the store must not leave the database locked,
		// otherwise it is not possible to open it again after fixing the problem.
		db, err := leveldb.OpenFile(td.store.config.StorePath(), &opt.Options{ErrorIfMissing: true})
		require.NoError(t, err, "the database is still locked after the failure")
		require.NoError(t, db.Close())
	})

	t.Run("a block of the cache window is corrupted", func(t *testing.T) {
		td := setup(t, nil)
		// A valid block hash followed by data that is not a block.
		require.NoError(t, td.store.db.Put(blockKey(5), append(td.RandHash().Bytes(), 1, 2, 3), nil))
		td.store.Close()

		_, err := NewStore(td.store.config)
		require.Error(t, err)
	})
}

func TestTransactionIndexIsCorrupted(t *testing.T) {
	td := setup(t, nil)

	lastHeight := td.store.LastCertificate().Height()
	cBlk, err := td.store.Block(lastHeight)
	require.NoError(t, err)
	blk, err := cBlk.ToBlock()
	require.NoError(t, err)
	txID := blk.Transactions()[0].ID().Bytes()

	t.Run("offset is out of the block", func(t *testing.T) {
		putTxIndex(t, td, txID, lastHeight, 1<<20, 10)

		_, err := td.store.Transaction(blk.Transactions()[0].ID())
		require.ErrorIs(t, err, ErrBadOffset)
	})

	t.Run("offset and length overflow", func(t *testing.T) {
		putTxIndex(t, td, txID, lastHeight, 4, math.MaxUint32)

		require.NotPanics(t, func() {
			_, err := td.store.Transaction(blk.Transactions()[0].ID())
			require.ErrorIs(t, err, ErrBadOffset)
		})
	})

	t.Run("block of the transaction doesn't exist", func(t *testing.T) {
		putTxIndex(t, td, txID, lastHeight+100, 1, 1)

		_, err := td.store.Transaction(blk.Transactions()[0].ID())
		require.Error(t, err)
	})

	t.Run("index entry is truncated", func(t *testing.T) {
		require.NoError(t, td.store.db.Put(append(txPrefix, txID...), []byte{1, 2}, nil))

		_, err := td.store.Transaction(blk.Transactions()[0].ID())
		require.Error(t, err)
	})
}

func TestPublicKeyIsCorrupted(t *testing.T) {
	td := setup(t, nil)

	blsPub, _ := td.RandBLSKeyPair()
	ed25519Pub, _ := td.RandEd25519KeyPair()
	secp256k1Pub, _ := td.RandSecp256k1KeyPair()

	var unknownTypeAddr crypto.Address
	unknownTypeAddr[0] = 99

	addrs := map[string]crypto.Address{
		"validator":         blsPub.ValidatorAddress(),
		"bls account":       blsPub.AccountAddress(),
		"ed25519 account":   ed25519Pub.AccountAddress(),
		"secp256k1 account": secp256k1Pub.AccountAddress(),
	}

	for name, addr := range addrs {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, td.store.db.Put(publicKeyKey(addr), []byte{1, 2, 3}, nil))

			assert.True(t, td.store.HasPublicKey(addr))
			require.NotPanics(t, func() {
				_, err := td.store.PublicKey(addr)
				require.Error(t, err)
			})
		})
	}

	t.Run("unknown address type", func(t *testing.T) {
		require.NoError(t, td.store.db.Put(publicKeyKey(unknownTypeAddr), []byte{1, 2, 3}, nil))

		_, err := td.store.PublicKey(unknownTypeAddr)
		require.ErrorIs(t, err, PublicKeyNotFoundError{Address: unknownTypeAddr})
	})
}

func TestCommittedDataIsCorrupted(t *testing.T) {
	td := setup(t, nil)

	cBlk := &CommittedBlock{store: td.store, Height: 1, Data: []byte{1, 2, 3}}
	_, err := cBlk.ToBlock()
	require.Error(t, err)

	cTrx := &CommittedTx{store: td.store, Height: 1, Data: []byte{1, 2, 3}}
	_, err = cTrx.ToTx()
	require.Error(t, err)
}

func TestPruneFailures(t *testing.T) {
	conf := testConfig()
	conf.RetentionDays = 1
	td := setup(t, conf)

	// A valid block hash followed by data that is not a block.
	require.NoError(t, td.store.db.Put(blockKey(3), append(td.RandHash().Bytes(), 1, 2, 3), nil))

	t.Run("pruning a corrupted block", func(t *testing.T) {
		pruned, err := td.store.pruneBlock(3)
		require.Error(t, err)
		assert.False(t, pruned)
	})

	t.Run("pruning stops with the first error", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(blockPerDay + 7)
		td.store.SaveBlock(blk, cert)
		require.NoError(t, td.store.WriteBatch())

		err := td.store.Prune(func(bool, types.Height) bool { return false })
		require.Error(t, err)
	})
}

func TestClosedStoreDoesNotPanic(t *testing.T) {
	td := setup(t, nil)

	blk, err := td.store.Block(1)
	require.NoError(t, err)
	trxID := func() hash.Hash {
		b, _ := blk.ToBlock()

		return b.Transactions()[0].ID()
	}()
	unknownAddr := td.RandAccAddress()

	td.store.Close()
	td.store.Close() // Closing twice is harmless

	require.NotPanics(t, func() {
		_, err := td.store.Block(1)
		require.Error(t, err)

		_, err = td.store.Transaction(trxID)
		require.Error(t, err)

		_, err = td.store.PublicKey(unknownAddr)
		require.Error(t, err)

		_, err = td.store.Account(unknownAddr)
		require.Error(t, err)

		assert.Equal(t, hash.UndefHash, td.store.BlockHash(1))
		assert.Zero(t, td.store.BlockHeight(blk.BlockHash))
		assert.False(t, td.store.HasPublicKey(unknownAddr))
		assert.False(t, td.store.HasAccount(unknownAddr))
	})
}

func TestAccountStoreWithoutCache(t *testing.T) {
	td := setup(t, nil)

	assert.Panics(t, func() { newAccountStore(td.store.db, 0) })
}

func TestErrorMessages(t *testing.T) {
	td := setup(t, nil)
	addr := td.RandAccAddress()

	assert.Equal(t, "bad path", ConfigError{Reason: "bad path"}.Error())
	assert.Contains(t, PublicKeyNotFoundError{Address: addr}.Error(), addr.String())
}

// If the last certificate is saved but can't be read, the store must not pretend to be at the genesis height:
// the node would write the genesis state over the state of the blockchain.
func TestCorruptedLastInfo(t *testing.T) {
	corruptions := map[string][]byte{
		"empty value":            {},
		"truncated version":      {2, 0},
		"truncated certificate":  {2, 0, 0, 0, 1, 0},
		"certificate is garbage": {2, 0, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	for name, data := range corruptions {
		t.Run(name, func(t *testing.T) {
			td := setup(t, nil)
			require.NoError(t, td.store.db.Put(lastInfoKey, data, nil))
			td.store.Close()

			_, err := NewStore(td.store.config)
			require.Error(t, err)

			// The database is released, and can be opened again once the problem is fixed.
			db, err := leveldb.OpenFile(td.store.config.StorePath(), &opt.Options{ErrorIfMissing: true})
			require.NoError(t, err)
			require.NoError(t, db.Close())
		})
	}

	t.Run("empty database is at the genesis height", func(t *testing.T) {
		store, err := NewStore(testConfig())
		require.NoError(t, err)
		assert.Nil(t, store.LastCertificate())
		store.Close()
	})

	t.Run("healthy database is reopened", func(t *testing.T) {
		td := setup(t, nil)
		td.store.Close()

		store, err := NewStore(td.store.config)
		require.NoError(t, err)
		assert.Equal(t, types.Height(10), store.LastCertificate().Height())
		store.Close()
	})
}
