package wallet

import (
	_ "embed"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

var (
	/// ErrWalletExits describes an error in which there is a wallet
	/// exists in the given path
	ErrWalletExits = errors.New("wallet exists")

	/// ErrWalletExits describes an error in which the wallet CRC is
	/// invalid
	ErrInvalidCRC = errors.New("invalid CRC")

	/// ErrWalletExits describes an error in which the network is not
	/// valid
	ErrInvalidNetwork = errors.New("invalid network")

	/// ErrWalletExits describes an error in which the address doesn't
	/// exist in wallet
	ErrAddressNotFound = errors.New("address not found")

	/// ErrWalletExits describes an error in which the address already
	/// exist in wallet
	ErrAddressExists = errors.New("address already exists")
)

type Wallet struct {
	path   string
	store  *Store
	client *grpcClient
}

type serverInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}
type servers = map[string][]serverInfo

//go:embed servers.json
var serversJSON []byte

/// OpenWallet generates an empty wallet and save the seed string
func OpenWallet(path string) (*Wallet, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := new(Store)
	err = json.Unmarshal(data, s)
	exitOnErr(err)

	if s.VaultCRC != s.calcVaultCRC() {
		exitOnErr(ErrInvalidCRC)
	}

	return newWallet(path, s, true)
}

/// Recover recovers a wallet from mnemonic (seed phrase)
func RecoverWallet(path, mnemonic string, net int) (*Wallet, error) {
	path = util.MakeAbs(path)
	if util.PathExists(path) {
		return nil, ErrWalletExits
	}
	s, err := RecoverStore(mnemonic, net)
	if err != nil {
		return nil, err
	}
	w, err := newWallet(path, s, false)
	if err != nil {
		return nil, err
	}

	err = w.saveToFile()
	if err != nil {
		return nil, err
	}

	return w, nil
}

/// CreateWallet generates an empty wallet and save the seed string
func CreateWallet(path, passphrase string, net int) (*Wallet, error) {
	path = util.MakeAbs(path)
	if util.PathExists(path) {
		return nil, ErrWalletExits
	}
	s, err := NewStore(passphrase, net)
	if err != nil {
		return nil, err
	}
	w, err := newWallet(path, s, false)
	if err != nil {
		return nil, err
	}

	err = w.saveToFile()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newWallet(path string, store *Store, online bool) (*Wallet, error) {
	w := &Wallet{
		store: store,
		path:  path,
	}

	err := w.connectToRandomServer()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Wallet) connectToRandomServer() error {
	serversInfo := servers{}
	err := json.Unmarshal(serversJSON, &serversInfo)
	exitOnErr(err)

	var netServers []serverInfo
	switch w.store.Network {
	case 0:
		{ // mainnet
			netServers = serversInfo["mainnet"]
		}
	case 1:
		{ // testnet
			netServers = serversInfo["testnet"]
		}

	default:
		{
			return ErrInvalidNetwork
		}
	}

	for i := 0; i < 3; i++ {
		n := util.RandInt32(int32(len(netServers)))
		serverInfo := netServers[n]
		client, err := gewGRPCClient(serverInfo.IP)
		if err == nil {
			w.client = client
			return nil
		}
	}

	return errors.New("unable to connect to the servers")
}

func (w *Wallet) Path() string {
	return w.path
}

func (w *Wallet) IsEncrypted() bool {
	return w.store.Encrypted
}

func (w *Wallet) saveToFile() error {
	w.store.VaultCRC = w.store.calcVaultCRC()

	bs, err := json.MarshalIndent(w.store, "  ", "  ")
	exitOnErr(err)

	return util.WriteFile(w.path, bs)
}

func (w *Wallet) ImportPrivateKey(passphrase string, prvStr string) error {
	prv, err := bls.PrivateKeyFromString(prvStr)
	if err != nil {
		return err
	}
	err = w.store.ImportPrivateKey(passphrase, prv)
	if err != nil {
		return err
	}
	return w.saveToFile()
}

func (w *Wallet) NewAddress(passphrase, label string) (string, error) {
	addr, err := w.store.NewAddress(passphrase, label)
	if err != nil {
		return "", err
	}
	err = w.saveToFile()
	if err != nil {
		return "", err
	}

	return addr, nil
}

func (w *Wallet) GetBalance(addrStr string) (int64, int64, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return 0, 0, err
	}

	balance, _ := w.client.getAccountBalance(addr)
	//exitOnErr(err)
	stake, _ := w.client.getValidatorStake(addr)
	//exitOnErr(err)

	return balance, stake, nil
}

func (w *Wallet) PrivateKey(passphrase, addr string) (string, error) {
	prv, err := w.store.PrivateKey(passphrase, addr)
	if err != nil {
		return "", err
	}
	return prv.String(), nil
}

func (w *Wallet) PublicKey(passphrase, addr string) (string, error) {
	prv, err := w.store.PrivateKey(passphrase, addr)
	if err != nil {
		return "", err
	}
	return prv.PublicKey().String(), nil
}

func (w *Wallet) Mnemonic(passphrase string) (string, error) {
	return w.store.Mnemonic(passphrase)
}

func (w *Wallet) Addresses() map[string]string {
	return w.store.Addresses()
}

//
// MakeBondTx creates a new bond transaction based on the given parameters
func (w *Wallet) MakeBondTx(stampStr, seqStr, senderStr, valPubStr, stakeStr, feeStr, memo string) (*tx.Tx, error) {
	sender, err := crypto.AddressFromString(senderStr)
	if err != nil {
		return nil, err
	}
	valPub, err := bls.PublicKeyFromString(valPubStr)
	if err != nil {
		return nil, err
	}
	stake, err := strconv.ParseInt(stakeStr, 10, 64)
	if err != nil {
		return nil, err
	}
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsAccSeq(sender, seqStr)
	if err != nil {
		return nil, err
	}
	fee, err := w.parsFee(stake, feeStr)
	if err != nil {
		return nil, err
	}

	tx := tx.NewBondTx(stamp, seq, sender, valPub, stake, fee, memo)
	return tx, nil
}

// MakeUnbondTx creates a new unbond transaction based on the given parameters
func (w *Wallet) MakeUnbondTx(stampStr, seqStr, addrStr, memo string) (*tx.Tx, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return nil, err
	}
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsValSeq(addr, seqStr)
	if err != nil {
		return nil, err
	}

	tx := tx.NewUnbondTx(stamp, seq, addr, memo)
	return tx, nil
}

// TODO: write tests for me by mocking grpc server
// MakeWithdrawTx creates a new unbond transaction based on the given parameters
func (w *Wallet) MakeWithdrawTx(stampStr, seqStr, valAddrStr, accAddrStr, amountStr, feeStr, memo string) (*tx.Tx, error) {
	valAddr, err := crypto.AddressFromString(valAddrStr)
	if err != nil {
		return nil, err
	}
	accAddr, err := crypto.AddressFromString(accAddrStr)
	if err != nil {
		return nil, err
	}
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsValSeq(valAddr, seqStr)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return nil, err
	}
	fee, err := w.parsFee(amount, feeStr)
	if err != nil {
		return nil, err
	}

	tx := tx.NewWithdrawTx(stamp, seq, valAddr, accAddr, amount, fee, memo)
	return tx, nil
}

// MakeSendTx creates a new send transaction based on the given parameters
func (w *Wallet) MakeSendTx(stampStr, seqStr, senderStr, receiverStr, amountStr, feeStr, memo string) (*tx.Tx, error) {
	sender, err := crypto.AddressFromString(senderStr)
	if err != nil {
		return nil, err
	}
	receiver, err := crypto.AddressFromString(receiverStr)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return nil, err
	}
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsAccSeq(sender, seqStr)
	if err != nil {
		return nil, err
	}

	fee, err := w.parsFee(amount, feeStr)
	if err != nil {
		return nil, err
	}

	tx := tx.NewSendTx(stamp, seq, sender, receiver, amount, fee, memo)
	return tx, nil
}

func (w *Wallet) parsAccSeq(signer crypto.Address, seqStr string) (int32, error) {
	if seqStr != "" {
		seq, err := strconv.ParseInt(seqStr, 10, 32)
		if err != nil {
			return -1, err
		}
		return int32(seq), nil
	}

	return w.client.getAccountSequence(signer)
}

func (w *Wallet) parsFee(amount int64, feeStr string) (int64, error) {
	if feeStr != "" {
		fee, err := strconv.ParseInt(feeStr, 10, 64)
		if err != nil {
			return -1, err
		}
		return fee, nil
	}

	fee := amount / 10000
	if fee < 10000 {
		fee = 10000
	}
	return fee, nil
}

func (w *Wallet) parsValSeq(signer crypto.Address, seqStr string) (int32, error) {
	if seqStr != "" {
		seq, err := strconv.ParseInt(seqStr, 10, 32)
		if err != nil {
			return -1, err
		}
		return int32(seq), nil
	}

	return w.client.GetValidatorSequence(signer)
}

func (w *Wallet) parsStamp(stampStr string) (hash.Stamp, error) {
	if stampStr != "" {
		stamp, err := hash.StampFromString(stampStr)
		if err != nil {
			return hash.UndefHash.Stamp(), err
		}
		return stamp, nil
	}
	return w.client.getStamp()
}

func (w *Wallet) SignAndBroadcast(passphrase string, tx *tx.Tx) (string, error) {
	prv, err := w.store.PrivateKey(passphrase, tx.Payload().Signer().String())
	if err != nil {
		return "", err
	}

	signer := crypto.NewSigner(prv)
	signer.SignMsg(tx)
	b, err := tx.Bytes()
	if err != nil {
		return "", err
	}

	return w.client.sendTx(b)
}
