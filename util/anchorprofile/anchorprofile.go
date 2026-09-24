// Package anchorprofile implements PAC-ANCHOR-1, the off-chain convention of
// PIP-50 for the root hash an anchor commits to.
//
// The node never uses this package and never validates the profile. It is a
// reference for wallets and applications, so independent clients compute the
// same root for the same files.
//
// A single file is committed by its BLAKE2b-256 digest (BlobRoot).
// Several files are committed by a Merkle tree with the shape of RFC 9162, and
// the root also commits to the number of files:
//
//	leaf = BLAKE2b-256(0x00 || uint16be(len(name)) || name || BLAKE2b-256(content))
//	node = BLAKE2b-256(0x01 || left || right)
//	root = BLAKE2b-256(0x02 || uint64be(n) || tree)
//
// Items are ordered by name (bytewise). With n > 1 items, the left subtree holds
// the first k items, where k is the largest power of two smaller than n. No node
// is ever duplicated, and leaves and nodes cannot be confused, so two different
// item lists never share a root, no inner node can pass for a file, and a proof
// cannot claim another number of files.
package anchorprofile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/pactus-project/pactus/crypto/hash"
)

const (
	leafPrefix = 0x00
	nodePrefix = 0x01
	rootPrefix = 0x02

	// MaxNameLen is the longest item name, in bytes.
	MaxNameLen = math.MaxUint16

	// ManifestVersion and ManifestAlg are the values a PAC-ANCHOR-1 manifest carries.
	ManifestVersion = 1
	ManifestAlg     = "blake2b-256"
)

var (
	// ErrNoItem is returned when a tree has no item.
	ErrNoItem = errors.New("no item")

	// ErrDuplicateName is returned when two items share a name.
	ErrDuplicateName = errors.New("duplicate item name")

	// ErrInvalidName is returned when a name is empty, too long or not valid UTF-8.
	ErrInvalidName = errors.New("invalid item name")

	// ErrItemNotFound is returned when no item has the requested name.
	ErrItemNotFound = errors.New("item not found")

	// ErrInvalidManifest is returned when a manifest does not follow PAC-ANCHOR-1.
	ErrInvalidManifest = errors.New("invalid manifest")
)

// BlobRoot returns the root hash that commits to a single file.
func BlobRoot(content []byte) hash.Hash {
	return hash.CalcHash(content)
}

// Item is one file of a multi-file anchor.
type Item struct {
	Name string
	Hash hash.Hash // BLAKE2b-256 of the file content
}

// NewItem returns the item for a named file.
func NewItem(name string, content []byte) Item {
	return Item{Name: name, Hash: hash.CalcHash(content)}
}

// LeafHash returns the tree leaf of an item. It binds the name to the content.
func LeafHash(item Item) hash.Hash {
	data := make([]byte, 0, 1+2+len(item.Name)+hash.HashSize)
	data = append(data, leafPrefix)
	data = binary.BigEndian.AppendUint16(data, uint16(len(item.Name)))
	data = append(data, item.Name...)
	data = append(data, item.Hash.Bytes()...)

	return hash.CalcHash(data)
}

func nodeHash(left, right hash.Hash) hash.Hash {
	data := make([]byte, 0, 1+2*hash.HashSize)
	data = append(data, nodePrefix)
	data = append(data, left.Bytes()...)
	data = append(data, right.Bytes()...)

	return hash.CalcHash(data)
}

// sizedRoot binds the number of items to the tree root.
func sizedRoot(size int, tree hash.Hash) hash.Hash {
	data := make([]byte, 0, 1+8+hash.HashSize)
	data = append(data, rootPrefix)
	data = binary.BigEndian.AppendUint64(data, uint64(size))
	data = append(data, tree.Bytes()...)

	return hash.CalcHash(data)
}

// sortedLeaves checks the items and returns their leaves in name order.
func sortedLeaves(items []Item) ([]Item, []hash.Hash, error) {
	if len(items) == 0 {
		return nil, nil, ErrNoItem
	}

	sorted := slices.Clone(items)
	slices.SortFunc(sorted, func(left, right Item) int {
		return strings.Compare(left.Name, right.Name)
	})

	leaves := make([]hash.Hash, len(sorted))
	for i, item := range sorted {
		if !validName(item.Name) {
			return nil, nil, fmt.Errorf("%w: %q", ErrInvalidName, item.Name)
		}
		if i > 0 && sorted[i-1].Name == item.Name {
			return nil, nil, fmt.Errorf("%w: %q", ErrDuplicateName, item.Name)
		}
		leaves[i] = LeafHash(item)
	}

	return sorted, leaves, nil
}

func validName(name string) bool {
	return name != "" && len(name) <= MaxNameLen && utf8.ValidString(name)
}

// splitPoint returns the largest power of two smaller than n (n > 1).
func splitPoint(n int) int {
	k := 1
	for k<<1 < n {
		k <<= 1
	}

	return k
}

func subtreeRoot(leaves []hash.Hash) hash.Hash {
	if len(leaves) == 1 {
		return leaves[0]
	}
	k := splitPoint(len(leaves))

	return nodeHash(subtreeRoot(leaves[:k]), subtreeRoot(leaves[k:]))
}

// Root returns the root hash that commits to the items. The input order does
// not matter; the items are ordered by name.
func Root(items []Item) (hash.Hash, error) {
	_, leaves, err := sortedLeaves(items)
	if err != nil {
		return hash.Hash{}, err
	}

	return sizedRoot(len(leaves), subtreeRoot(leaves)), nil
}

// Proof shows that one item belongs to a root.
type Proof struct {
	Index int         // position of the item in name order
	Size  int         // number of items in the tree
	Path  []hash.Hash // sibling subtree roots, from the leaf up
}

func inclusionPath(index int, leaves []hash.Hash) []hash.Hash {
	if len(leaves) == 1 {
		return nil
	}
	k := splitPoint(len(leaves))
	if index < k {
		return append(inclusionPath(index, leaves[:k]), subtreeRoot(leaves[k:]))
	}

	return append(inclusionPath(index-k, leaves[k:]), subtreeRoot(leaves[:k]))
}

// Prove returns the inclusion proof of the named item.
func Prove(items []Item, name string) (Proof, error) {
	sorted, leaves, err := sortedLeaves(items)
	if err != nil {
		return Proof{}, err
	}

	index := slices.IndexFunc(sorted, func(item Item) bool {
		return item.Name == name
	})
	if index < 0 {
		return Proof{}, fmt.Errorf("%w: %q", ErrItemNotFound, name)
	}

	return Proof{
		Index: index,
		Size:  len(leaves),
		Path:  inclusionPath(index, leaves),
	}, nil
}

// Verify reports whether the proof shows that the item belongs to the root, in
// a tree of exactly proof.Size items. The path is checked with the algorithm of
// RFC 9162, section 2.1.3.2.
func Verify(root hash.Hash, item Item, proof Proof) bool {
	if !validName(item.Name) || proof.Index < 0 || proof.Index >= proof.Size {
		return false
	}

	// position and lastPosition are fn and sn in RFC 9162.
	position, lastPosition := uint64(proof.Index), uint64(proof.Size-1)
	current := LeafHash(item)
	for _, sibling := range proof.Path {
		if lastPosition == 0 {
			return false
		}
		if position&1 == 1 || position == lastPosition {
			current = nodeHash(sibling, current)
			for position&1 == 0 && position != 0 {
				position >>= 1
				lastPosition >>= 1
			}
		} else {
			current = nodeHash(current, sibling)
		}
		position >>= 1
		lastPosition >>= 1
	}

	return lastPosition == 0 && sizedRoot(proof.Size, current) == root
}

// Manifest is the JSON document a multi-file anchor's ManifestURI points to.
type Manifest struct {
	Version int            `json:"v"`
	Alg     string         `json:"alg"`
	Items   []ManifestItem `json:"items"`
}

// ManifestItem is one file of a manifest; Hash is the hex BLAKE2b-256 digest.
type ManifestItem struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

// NewManifest returns the manifest of the items, in name order.
func NewManifest(items []Item) (*Manifest, error) {
	sorted, _, err := sortedLeaves(items)
	if err != nil {
		return nil, err
	}

	manifest := &Manifest{
		Version: ManifestVersion,
		Alg:     ManifestAlg,
		Items:   make([]ManifestItem, 0, len(sorted)),
	}
	for _, item := range sorted {
		manifest.Items = append(manifest.Items, ManifestItem{Name: item.Name, Hash: item.Hash.String()})
	}

	return manifest, nil
}

// Entries returns the manifest items, checking the version, the algorithm and
// every hash.
func (m *Manifest) Entries() ([]Item, error) {
	if m.Version != ManifestVersion || m.Alg != ManifestAlg {
		return nil, fmt.Errorf("%w: version %d, algorithm %q", ErrInvalidManifest, m.Version, m.Alg)
	}

	items := make([]Item, 0, len(m.Items))
	for _, entry := range m.Items {
		digest, err := hash.FromString(entry.Hash)
		if err != nil {
			return nil, fmt.Errorf("%w: item %q: %w", ErrInvalidManifest, entry.Name, err)
		}
		items = append(items, Item{Name: entry.Name, Hash: digest})
	}

	return items, nil
}

// Root returns the root hash the manifest commits to.
func (m *Manifest) Root() (hash.Hash, error) {
	items, err := m.Entries()
	if err != nil {
		return hash.Hash{}, err
	}

	return Root(items)
}
