package anchorprofile

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/stretchr/testify/require"
)

func namedItems(count int) []Item {
	items := make([]Item, count)
	for i := range items {
		items[i] = NewItem(fmt.Sprintf("file-%03d.pdf", i), []byte(fmt.Sprintf("content %d", i)))
	}

	return items
}

func TestBlobRoot(t *testing.T) {
	content := []byte("contract.pdf bytes")
	require.Equal(t, hash.CalcHash(content), BlobRoot(content))

	// A one-file tree is not the blob root: the leaf binds the name.
	root, err := Root([]Item{NewItem("contract.pdf", content)})
	require.NoError(t, err)
	require.NotEqual(t, BlobRoot(content), root)
}

// The tree must have the RFC 9162 shape. The expected roots are built by hand
// from the formulas, independently of the recursive code.
func TestTreeShape(t *testing.T) {
	items := namedItems(5)
	leaf := make([]hash.Hash, len(items))
	for i, item := range items {
		leaf[i] = LeafHash(item)
	}
	node := nodeHash

	want := map[int]hash.Hash{
		1: leaf[0],
		2: node(leaf[0], leaf[1]),
		3: node(node(leaf[0], leaf[1]), leaf[2]),
		4: node(node(leaf[0], leaf[1]), node(leaf[2], leaf[3])),
		5: node(node(node(leaf[0], leaf[1]), node(leaf[2], leaf[3])), leaf[4]),
	}
	for count, tree := range want {
		root, err := Root(items[:count])
		require.NoError(t, err)
		require.Equal(t, sizedRoot(count, tree), root, "%d items", count)
	}
}

func TestLeafLayout(t *testing.T) {
	item := NewItem("a.txt", []byte("hello"))
	data := append([]byte{0x00, 0x00, 0x05}, []byte("a.txt")...)
	data = append(data, item.Hash.Bytes()...)
	require.Equal(t, hash.CalcHash(data), LeafHash(item))

	left, right := hash.CalcHash([]byte("l")), hash.CalcHash([]byte("r"))
	nodeData := append(append([]byte{0x01}, left.Bytes()...), right.Bytes()...)
	require.Equal(t, hash.CalcHash(nodeData), nodeHash(left, right))

	rootData := append([]byte{0x02, 0, 0, 0, 0, 0, 0, 0, 3}, left.Bytes()...)
	require.Equal(t, hash.CalcHash(rootData), sizedRoot(3, left))
}

func TestRootIgnoresInputOrder(t *testing.T) {
	items := namedItems(9)
	root, err := Root(items)
	require.NoError(t, err)

	reversed := slices.Clone(items)
	slices.Reverse(reversed)
	again, err := Root(reversed)
	require.NoError(t, err)
	require.Equal(t, root, again)
}

// The two weaknesses of a Bitcoin-style tree must be gone.
func TestNoAmbiguousRoots(t *testing.T) {
	fileA, fileB, fileC := NewItem("a", []byte("A")), NewItem("b", []byte("B")), NewItem("c", []byte("C"))

	t.Run("a repeated last file is a different list", func(t *testing.T) {
		three, err := Root([]Item{fileA, fileB, fileC})
		require.NoError(t, err)
		four, err := Root([]Item{fileA, fileB, fileC, {Name: "d", Hash: fileC.Hash}})
		require.NoError(t, err)
		require.NotEqual(t, three, four)

		_, err = Root([]Item{fileA, fileB, fileC, fileC})
		require.ErrorIs(t, err, ErrDuplicateName)
	})

	t.Run("an inner node cannot pass for a file", func(t *testing.T) {
		root, err := Root([]Item{fileA, fileB})
		require.NoError(t, err)

		joined := append(LeafHash(fileA).Bytes(), LeafHash(fileB).Bytes()...)
		forged := NewItem("forged", joined)
		require.False(t, Verify(root, forged, Proof{Index: 0, Size: 1}))
		require.False(t, Verify(root, Item{Name: "forged", Hash: hash.CalcHash(joined)}, Proof{Index: 0, Size: 1}))
	})

	t.Run("names are bound to contents", func(t *testing.T) {
		root, err := Root([]Item{fileA, fileB})
		require.NoError(t, err)
		swapped, err := Root([]Item{{Name: "a", Hash: fileB.Hash}, {Name: "b", Hash: fileA.Hash}})
		require.NoError(t, err)
		require.NotEqual(t, root, swapped)
	})
}

func TestEveryProofVerifies(t *testing.T) {
	for count := 1; count <= 70; count++ {
		items := namedItems(count)
		root, err := Root(items)
		require.NoError(t, err)

		for _, item := range items {
			proof, err := Prove(items, item.Name)
			require.NoError(t, err)
			require.Equal(t, count, proof.Size)
			require.LessOrEqual(t, len(proof.Path), bitsLen(count-1), "path of %d items", count)
			require.True(t, Verify(root, item, proof), "%s in %d items", item.Name, count)
		}
	}
}

func bitsLen(value int) int {
	bits := 0
	for value > 0 {
		bits++
		value >>= 1
	}

	return bits
}

func TestTamperedProofsFail(t *testing.T) {
	for _, count := range []int{2, 3, 5, 8, 13, 33} {
		items := namedItems(count)
		root, err := Root(items)
		require.NoError(t, err)

		for _, item := range items {
			proof, err := Prove(items, item.Name)
			require.NoError(t, err)

			tampered := make([]Proof, 0, 6+len(proof.Path))
			tampered = append(tampered,
				Proof{Index: proof.Index + 1, Size: proof.Size, Path: proof.Path},
				Proof{Index: proof.Index - 1, Size: proof.Size, Path: proof.Path},
				Proof{Index: proof.Index, Size: proof.Size + 1, Path: proof.Path},
				Proof{Index: proof.Index, Size: proof.Size - 1, Path: proof.Path},
				Proof{Index: proof.Index, Size: proof.Size, Path: proof.Path[:len(proof.Path)-1]},
				Proof{Index: proof.Index, Size: proof.Size, Path: append(slices.Clone(proof.Path), root)},
			)
			for step := range proof.Path {
				path := slices.Clone(proof.Path)
				path[step][0] ^= 0x01
				tampered = append(tampered, Proof{Index: proof.Index, Size: proof.Size, Path: path})
			}
			for i, bad := range tampered {
				require.False(t, Verify(root, item, bad), "%s in %d items, change %d", item.Name, count, i)
			}

			otherContent := Item{Name: item.Name, Hash: hash.CalcHash([]byte("other"))}
			require.False(t, Verify(root, otherContent, proof))
			renamed := Item{Name: item.Name + "x", Hash: item.Hash}
			require.False(t, Verify(root, renamed, proof))
		}
	}
}

func TestInvalidInput(t *testing.T) {
	_, err := Root(nil)
	require.ErrorIs(t, err, ErrNoItem)

	for _, name := range []string{"", string([]byte{0xFF}), strings.Repeat("n", MaxNameLen+1)} {
		_, err := Root([]Item{{Name: name}})
		require.ErrorIs(t, err, ErrInvalidName)
		require.False(t, Verify(LeafHash(Item{Name: name}), Item{Name: name}, Proof{Index: 0, Size: 1}))
	}
	longest := Item{Name: strings.Repeat("n", MaxNameLen)}
	root, err := Root([]Item{longest})
	require.NoError(t, err)
	require.True(t, Verify(root, longest, Proof{Index: 0, Size: 1}))

	_, err = Prove(namedItems(3), "missing.pdf")
	require.ErrorIs(t, err, ErrItemNotFound)

	item := namedItems(1)[0]
	require.False(t, Verify(LeafHash(item), item, Proof{Index: 0, Size: 0}))
	require.False(t, Verify(LeafHash(item), item, Proof{Index: -1, Size: 1}))
}

func TestManifestRoundTrip(t *testing.T) {
	items := namedItems(6)
	slices.Reverse(items)
	manifest, err := NewManifest(items)
	require.NoError(t, err)
	require.Equal(t, "file-000.pdf", manifest.Items[0].Name, "manifest items are in name order")

	raw, err := json.Marshal(manifest)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"v":1,"alg":"blake2b-256"`)

	var decoded Manifest
	require.NoError(t, json.Unmarshal(raw, &decoded))
	fromManifest, err := decoded.Root()
	require.NoError(t, err)
	want, err := Root(items)
	require.NoError(t, err)
	require.Equal(t, want, fromManifest)

	bad := []Manifest{
		{Version: 2, Alg: ManifestAlg, Items: decoded.Items},
		{Version: 1, Alg: "sha256", Items: decoded.Items},
		{Version: 1, Alg: ManifestAlg, Items: []ManifestItem{{Name: "a", Hash: "zz"}}},
	}
	for _, manifest := range bad {
		_, err := manifest.Root()
		require.ErrorIs(t, err, ErrInvalidManifest)
	}
	empty := Manifest{Version: 1, Alg: ManifestAlg}
	_, err = empty.Root()
	require.ErrorIs(t, err, ErrNoItem)
}

// Test vectors published in PIP-50 for other implementations. They were
// checked against an independent implementation written from the formulas.
func TestVectors(t *testing.T) {
	require.Equal(t, "324dcf027dd4a30a932c441f365a25e86b173defa4b8e58948253471b81b72cf",
		BlobRoot([]byte("hello")).String())

	three, err := Root([]Item{
		NewItem("c.txt", []byte("!")),
		NewItem("a.txt", []byte("hello")),
		NewItem("b.txt", []byte("world")),
	})
	require.NoError(t, err)
	require.Equal(t, "39125630b77b5d7035a1fb07e1b2854b0c911f4bfa7605a9b59476ad1bf163b7", three.String())

	five := make([]Item, 5)
	for i := range five {
		five[i] = NewItem(fmt.Sprintf("file-%d.txt", i), []byte(fmt.Sprintf("content %d", i)))
	}
	root, err := Root(five)
	require.NoError(t, err)
	require.Equal(t, "6fc9d50176e7a1abe6c0875d9e655142370e93487a69badf586e60f9afe35e64", root.String())
}

// FuzzVerify alters real proofs: only the exact index, size and path verify.
func FuzzVerify(f *testing.F) {
	f.Add(uint8(3), uint8(0), 0, 0, uint8(0))
	f.Add(uint8(3), uint8(0), 0, 1, uint8(0))
	f.Add(uint8(8), uint8(5), 1, 0, uint8(0))
	f.Add(uint8(13), uint8(12), 0, -1, uint8(0))
	f.Add(uint8(33), uint8(7), 0, 0, uint8(9))

	f.Fuzz(func(t *testing.T, count, pick uint8, indexShift, sizeShift int, flip uint8) {
		items := namedItems(int(count%40) + 1)
		item := items[int(pick)%len(items)]
		root, err := Root(items)
		require.NoError(t, err)
		proof, err := Prove(items, item.Name)
		require.NoError(t, err)

		changed := Proof{
			Index: proof.Index + indexShift,
			Size:  proof.Size + sizeShift,
			Path:  slices.Clone(proof.Path),
		}
		flipped := false
		if flip != 0 && len(changed.Path) > 0 {
			changed.Path[int(flip)%len(changed.Path)][int(flip)%hash.HashSize] ^= flip
			flipped = true
		}

		honest := indexShift == 0 && sizeShift == 0 && !flipped
		require.Equal(t, honest, Verify(root, item, changed))
	})
}

// The proof of b.txt in the three-file test vector. The JSON was produced by an
// independent implementation written from the formulas.
func TestProofDocumentVector(t *testing.T) {
	items := []Item{
		NewItem("a.txt", []byte("hello")),
		NewItem("b.txt", []byte("world")),
		NewItem("c.txt", []byte("!")),
	}
	root, err := Root(items)
	require.NoError(t, err)
	proof, err := Prove(items, "b.txt")
	require.NoError(t, err)

	raw, err := json.Marshal(NewProofDocument("", items[1], proof))
	require.NoError(t, err)
	require.JSONEq(t, `{
		"v": 1,
		"alg": "blake2b-256",
		"item": {"name": "b.txt", "hash": "9a3440c9d1529b122faceef33739b6e814616658d53faaf6e4f129fb20edfb13"},
		"index": 1,
		"size": 3,
		"path": [
			"26bb722ed5fac59aaa2de435d19877b74be151b15bffbe62296ba2168b7e6805",
			"427a194a554ee2fa4dd1c53edc68c8ad515ad8780c0b1d6941ab367f6da476b2"
		]
	}`, string(raw))
	require.NotContains(t, string(raw), "address", "an empty address is left out")

	var decoded ProofDocument
	require.NoError(t, json.Unmarshal(raw, &decoded))
	require.True(t, decoded.Verify(root))
}

func TestProofDocumentRoundTrip(t *testing.T) {
	for count := 1; count <= 20; count++ {
		items := namedItems(count)
		root, err := Root(items)
		require.NoError(t, err)

		for _, item := range items {
			proof, err := Prove(items, item.Name)
			require.NoError(t, err)

			raw, err := json.Marshal(NewProofDocument("pc1zexample", item, proof))
			require.NoError(t, err)
			require.Contains(t, string(raw), `"address":"pc1zexample"`)
			if count == 1 {
				require.Contains(t, string(raw), `"path":[]`, "a one-item proof has an empty path, not null")
			}

			var decoded ProofDocument
			require.NoError(t, json.Unmarshal(raw, &decoded))
			require.True(t, decoded.Verify(root), "%s in %d items", item.Name, count)

			gotItem, gotProof, err := decoded.Decode()
			require.NoError(t, err)
			require.Equal(t, item, gotItem)
			require.Equal(t, proof.Index, gotProof.Index)
			require.Equal(t, proof.Size, gotProof.Size)
			require.True(t, slices.Equal(proof.Path, gotProof.Path), "path of %s in %d items", item.Name, count)
		}
	}
}

func TestProofDocumentRejects(t *testing.T) {
	items := namedItems(5)
	root, err := Root(items)
	require.NoError(t, err)
	proof, err := Prove(items, items[2].Name)
	require.NoError(t, err)
	valid := NewProofDocument("", items[2], proof)
	require.True(t, valid.Verify(root))

	change := func(edit func(*ProofDocument)) ProofDocument {
		doc := valid
		doc.Path = slices.Clone(valid.Path)
		edit(&doc)

		return doc
	}
	bad := map[string]ProofDocument{
		"version":         change(func(d *ProofDocument) { d.Version = 2 }),
		"algorithm":       change(func(d *ProofDocument) { d.Alg = "sha256" }),
		"empty name":      change(func(d *ProofDocument) { d.Item.Name = "" }),
		"item hash":       change(func(d *ProofDocument) { d.Item.Hash = "zz" }),
		"short item hash": change(func(d *ProofDocument) { d.Item.Hash = d.Item.Hash[:62] }),
		"path hash":       change(func(d *ProofDocument) { d.Path[0] = "not-hex" }),
		"negative index":  change(func(d *ProofDocument) { d.Index = -1 }),
		"index too large": change(func(d *ProofDocument) { d.Index = d.Size }),
		"no item":         change(func(d *ProofDocument) { d.Size = 0 }),
		"path too long": change(func(d *ProofDocument) {
			d.Path = slices.Repeat([]string{d.Path[0]}, maxPathLen+1)
		}),
	}
	for name, doc := range bad {
		_, _, err := doc.Decode()
		require.ErrorIs(t, err, ErrInvalidProof, name)
		require.False(t, doc.Verify(root), name)
	}
}

// FuzzProofDocument feeds arbitrary JSON to a verifier. It must never panic, and
// it may only accept a document for a real item at its real position.
func FuzzProofDocument(f *testing.F) {
	items := namedItems(7)
	root, err := Root(items)
	require.NoError(f, err)
	for _, item := range items[:3] {
		proof, err := Prove(items, item.Name)
		require.NoError(f, err)
		raw, err := json.Marshal(NewProofDocument("pc1zexample", item, proof))
		require.NoError(f, err)
		f.Add(raw)
	}
	f.Add([]byte(`{"v":1,"alg":"blake2b-256","item":{"name":"x","hash":"00"},"index":0,"size":1,"path":[]}`))
	f.Add([]byte(`{"v":1,"index":-1,"size":-1,"path":null}`))
	f.Add([]byte(`not json`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var doc ProofDocument
		if json.Unmarshal(data, &doc) != nil {
			return
		}
		if !doc.Verify(root) {
			return
		}

		item, proof, err := doc.Decode()
		require.NoError(t, err)
		genuine, err := Prove(items, item.Name)
		require.NoError(t, err, "accepted an item that is not in the tree")
		require.Contains(t, items, item)
		require.Equal(t, genuine.Index, proof.Index)
		require.Equal(t, genuine.Size, proof.Size)
	})
}
