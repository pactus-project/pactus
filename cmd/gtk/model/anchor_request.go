//go:build gtk

package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/bech32m"
)

// MinAnchorDeposit is the deposit a new anchor needs (PIP-50).
const MinAnchorDeposit = amount.Amount(1_000_000_000)

const (
	// AnchorActionSet creates or updates the anchor of an account.
	AnchorActionSet = "set"
	// AnchorActionDelete removes the anchor of an account and refunds its deposit.
	AnchorActionDelete = "delete"

	// NetworkMainnet and NetworkTestnet are the network names of an anchor link.
	NetworkMainnet = "mainnet"
	NetworkTestnet = "testnet"

	anchorURIPrefix  = "pactus:"
	anchorURIPath    = "anchor"
	maxAnchorMemo    = 64
	maxAnchorURI     = 128
	maxAnchorTypeVal = 255
)

var (
	anchorParamName   = regexp.MustCompile(`^[a-z0-9_-]+$`)
	anchorPlainAmount = regexp.MustCompile(`^(0|[1-9]\d*)(?:\.(\d{1,9}))?$`)
	anchorTypeValue   = regexp.MustCompile(`^(0|[1-9]\d{0,2})$`)
	anchorRootHex     = regexp.MustCompile(`^(?:[0-9a-fA-F]{2}){32,64}$`)

	anchorSetOnlyParams = []string{"type", "root", "uri", "locked"}

	anchorNetworkHRP = map[string]string{
		NetworkMainnet: "pc",
		NetworkTestnet: "tpc",
	}
)

// ErrInvalidAnchorLink is wrapped by every error ParseAnchorRequest returns.
var ErrInvalidAnchorLink = errors.New("invalid anchor link")

// AnchorRequest holds the values of a PIP-57 anchor link (pactus:anchor?...).
// The wallet builds the transaction from them for the account the user picks.
type AnchorRequest struct {
	Network string
	Action  string
	Memo    string
	// Sender is a suggested account. It is empty when the link has none.
	Sender string

	// The fields below are set only for a "set" action.
	AnchorType  uint8
	RootHash    []byte
	ManifestURI string
	// Locked is the total deposit the anchor should hold after the transaction.
	Locked amount.Amount
}

// IsDelete reports whether the request removes the anchor.
func (r *AnchorRequest) IsDelete() bool {
	return r.Action == AnchorActionDelete
}

func rejectAnchorLink(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidAnchorLink, fmt.Sprintf(format, args...))
}

// ParseAnchorRequest parses a PIP-57 anchor link for a wallet on walletNetwork.
// It rejects a malformed link as a whole and never repairs a value.
func ParseAnchorRequest(link, walletNetwork string) (*AnchorRequest, error) {
	params, err := anchorLinkParams(link)
	if err != nil {
		return nil, err
	}

	network := params["network"]
	action := params["action"]
	if _, ok := anchorNetworkHRP[network]; !ok {
		return nil, rejectAnchorLink("network missing or unknown")
	}
	if action != AnchorActionSet && action != AnchorActionDelete {
		return nil, rejectAnchorLink("action missing or unknown")
	}
	if network != walletNetwork {
		return nil, rejectAnchorLink("the link is for %s, but this wallet is on %s", network, walletNetwork)
	}

	memo := params["memo"]
	if len(memo) > maxAnchorMemo {
		return nil, rejectAnchorLink("memo longer than %d bytes", maxAnchorMemo)
	}
	for _, name := range []string{"uri", "memo"} {
		if hasForbiddenRune(params[name]) {
			return nil, rejectAnchorLink("control or formatting character in %s", name)
		}
	}

	sender, hasSender := params["sender"]
	if hasSender && !isAnchorAccountFor(sender, network) {
		return nil, rejectAnchorLink("sender is not an account address for %s", network)
	}

	req := &AnchorRequest{
		Network: network,
		Action:  action,
		Memo:    memo,
		Sender:  sender,
	}

	if action == AnchorActionDelete {
		for _, name := range anchorSetOnlyParams {
			if _, ok := params[name]; ok {
				return nil, rejectAnchorLink("a delete cannot carry %s", name)
			}
		}

		return req, nil
	}

	if err := fillAnchorSet(req, params); err != nil {
		return nil, err
	}

	return req, nil
}

// anchorLinkParams splits the query of an anchor link into its decoded values.
// It does not use url.ParseQuery, which turns "+" into a space.
func anchorLinkParams(link string) (map[string]string, error) {
	rest, ok := strings.CutPrefix(link, anchorURIPrefix)
	if !ok {
		return nil, rejectAnchorLink("not a pactus: link")
	}
	path, query, _ := strings.Cut(rest, "?")
	if path != anchorURIPath {
		return nil, rejectAnchorLink("not an anchor link")
	}

	params := make(map[string]string)
	for part := range strings.SplitSeq(query, "&") {
		if part == "" {
			continue
		}
		name, value, _ := strings.Cut(part, "=")
		// Names are not percent-decoded, so "%61ction" cannot pass for "action".
		if !anchorParamName.MatchString(name) {
			return nil, rejectAnchorLink("bad parameter name %q", name)
		}
		if _, dup := params[name]; dup {
			return nil, rejectAnchorLink("parameter %s repeated", name)
		}
		decoded, err := percentDecode(value)
		if err != nil {
			return nil, rejectAnchorLink("%s is not UTF-8", name)
		}
		params[name] = decoded
	}

	for name := range params {
		if strings.HasPrefix(name, "req-") {
			return nil, rejectAnchorLink("unknown required parameter %s", name)
		}
	}

	return params, nil
}

func fillAnchorSet(req *AnchorRequest, params map[string]string) error {
	typ, hasType := params["type"]
	root, hasRoot := params["root"]
	if !hasType || !hasRoot {
		return rejectAnchorLink("type or root missing")
	}

	if !anchorTypeValue.MatchString(typ) {
		return rejectAnchorLink("type out of range")
	}
	typeVal, err := strconv.Atoi(typ)
	if err != nil || typeVal > maxAnchorTypeVal {
		return rejectAnchorLink("type out of range")
	}

	if !anchorRootHex.MatchString(root) {
		return rejectAnchorLink("root is not 32 to 64 bytes of hex")
	}
	rootHash, err := hex.DecodeString(root)
	if err != nil {
		return rejectAnchorLink("root is not hex")
	}

	manifestURI := params["uri"]
	if len(manifestURI) > maxAnchorURI {
		return rejectAnchorLink("uri longer than %d bytes", maxAnchorURI)
	}

	locked := MinAnchorDeposit
	if text, ok := params["locked"]; ok {
		locked, err = parsePlainPAC(text)
		if err != nil {
			return err
		}
	}

	req.AnchorType = uint8(typeVal)
	req.RootHash = rootHash
	req.ManifestURI = manifestURI
	req.Locked = locked

	return nil
}

// parsePlainPAC parses an amount in PAC written as plain decimal digits.
func parsePlainPAC(text string) (amount.Amount, error) {
	match := anchorPlainAmount.FindStringSubmatch(text)
	if match == nil {
		return 0, rejectAnchorLink("locked is not a plain PAC amount")
	}

	whole, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil || whole > int64(amount.MaxNanoPAC/amount.NanoPACPerPAC) {
		return 0, rejectAnchorLink("locked too large")
	}
	frac := int64(0)
	if match[2] != "" {
		frac, err = strconv.ParseInt(match[2]+strings.Repeat("0", 9-len(match[2])), 10, 64)
		if err != nil {
			return 0, rejectAnchorLink("locked is not a plain PAC amount")
		}
	}

	locked := amount.Amount(whole*int64(amount.NanoPACPerPAC) + frac)
	if locked > amount.MaxNanoPAC {
		return 0, rejectAnchorLink("locked too large")
	}

	return locked, nil
}

// percentDecode turns "%xx" into a byte and keeps anything else, a lone "%" or a
// "+" too, as it is. The result must be valid UTF-8.
func percentDecode(value string) (string, error) {
	var buf strings.Builder
	for i := 0; i < len(value); i++ {
		if value[i] == '%' && i+2 < len(value) {
			if b, err := hex.DecodeString(value[i+1 : i+3]); err == nil {
				buf.WriteByte(b[0])
				i += 2

				continue
			}
		}
		buf.WriteByte(value[i])
	}

	decoded := buf.String()
	if !utf8.ValidString(decoded) {
		return "", errors.New("not UTF-8")
	}

	return decoded, nil
}

// hasForbiddenRune reports whether the text has a control or formatting character
// (Unicode categories Cc, Cf, Zl and Zp), which could change how it reads.
func hasForbiddenRune(text string) bool {
	for _, r := range text {
		if unicode.In(r, unicode.Cc, unicode.Cf, unicode.Zl, unicode.Zp) {
			return true
		}
	}

	return false
}

// isAnchorAccountFor reports whether addr is an account address of the network
// that can own an anchor. Validators and the treasury cannot.
func isAnchorAccountFor(addr, network string) bool {
	if addr != strings.ToLower(addr) {
		return false
	}
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(addr)
	if err != nil || hrp != anchorNetworkHRP[network] || len(data) != 20 {
		return false
	}

	switch crypto.AddressType(typ) {
	case crypto.AddressTypeBLSAccount,
		crypto.AddressTypeEd25519Account,
		crypto.AddressTypeSecp256k1Account:
		return true
	case crypto.AddressTypeTreasury, crypto.AddressTypeValidator:
		return false
	default:
		return false
	}
}

// AnchorDepositToAdd returns the deposit a set adds so the anchor holds the
// target total. A new anchor gets at least MinAnchorDeposit; a target below what
// is already locked adds nothing, since a set never releases funds.
func AnchorDepositToAdd(target, currentLocked amount.Amount) amount.Amount {
	target = max(target, MinAnchorDeposit)
	if target > currentLocked {
		return target - currentLocked
	}

	return 0
}

// CanDeleteAnchor reports whether a delete is valid: the account has an anchor
// and its locked deposit covers the fee, which is taken from the refund.
func CanDeleteAnchor(currentLocked, fee amount.Amount) bool {
	return currentLocked > 0 && currentLocked >= fee
}

// KnownAnchorTypes lists the anchor types PIP-50 assigns, in display order.
var KnownAnchorTypes = []uint8{0x00, 0x01, 0x02, 0x03, 0xFF}

// AnchorTypeName returns the PIP-50 name of an anchor type.
func AnchorTypeName(anchorType uint8) string {
	switch anchorType {
	case 0x00:
		return "Generic commitment"
	case 0x01:
		return "Service manifest"
	case 0x02:
		return "DID document"
	case 0x03:
		return "Credential or revocation registry"
	case 0xFF:
		return "Private or encrypted"
	default:
		return "Reserved"
	}
}

// AnchorTypeLabel returns the anchor type as "<number> (<name>)".
func AnchorTypeLabel(anchorType uint8) string {
	return fmt.Sprintf("%d (%s)", anchorType, AnchorTypeName(anchorType))
}
