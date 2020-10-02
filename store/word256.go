package store

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"
	"sort"

	"github.com/tmthrgd/go-hex"
)

var (
	Zero256 = Word256{}
	One256  = LeftPadWord256([]byte{1})
)

const Word256Length = 32

var BigWord256Length = big.NewInt(Word256Length)

var trimCutSet = string([]byte{0})

type Word256 [Word256Length]byte

func (w *Word256) UnmarshalText(hexBytes []byte) error {
	bs, err := hex.DecodeString(string(hexBytes))
	if err != nil {
		return err
	}
	copy(w[:], bs)
	return nil
}

func (w Word256) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeUpperToString(w[:])), nil
}

func (w Word256) String() string {
	return string(w[:])
}

func (w Word256) Copy() Word256 {
	return w
}

func (w Word256) Bytes() []byte {
	return w[:]
}

// copied.
func (w Word256) Prefix(n int) []byte {
	return w[:n]
}

func (w Word256) Postfix(n int) []byte {
	return w[32-n:]
}

func (w Word256) IsZero() bool {
	accum := byte(0)
	for _, byt := range w {
		accum |= byt
	}
	return accum == 0
}

func (w Word256) Compare(other Word256) int {
	return bytes.Compare(w[:], other[:])
}

func (w Word256) UnpadLeft() []byte {
	return bytes.TrimLeft(w[:], trimCutSet)
}

func (w Word256) UnpadRight() []byte {
	return bytes.TrimRight(w[:], trimCutSet)
}

func Uint64ToWord256(i uint64) Word256 {
	buf := [8]byte{}
	PutUint64BE(buf[:], i)
	return LeftPadWord256(buf[:])
}

func Int64ToWord256(i int64) Word256 {
	buf := [8]byte{}
	PutInt64BE(buf[:], i)
	return LeftPadWord256(buf[:])
}

func RightPadWord256(bz []byte) (word Word256) {
	copy(word[:], bz)
	return
}

func LeftPadWord256(bz []byte) (word Word256) {
	copy(word[32-len(bz):], bz)
	return
}

func Uint64FromWord256(word Word256) uint64 {
	buf := word.Postfix(8)
	return GetUint64BE(buf)
}

func Int64FromWord256(word Word256) int64 {
	buf := word.Postfix(8)
	return GetInt64BE(buf)
}

func PutUint64LE(dest []byte, i uint64) {
	binary.LittleEndian.PutUint64(dest, i)
}

func GetUint64LE(src []byte) uint64 {
	return binary.LittleEndian.Uint64(src)
}

func PutUint64BE(dest []byte, i uint64) {
	binary.BigEndian.PutUint64(dest, i)
}

func GetUint64BE(src []byte) uint64 {
	return binary.BigEndian.Uint64(src)
}

func PutInt64LE(dest []byte, i int64) {
	binary.LittleEndian.PutUint64(dest, uint64(i))
}

func GetInt64LE(src []byte) int64 {
	return int64(binary.LittleEndian.Uint64(src))
}

func PutInt64BE(dest []byte, i int64) {
	binary.BigEndian.PutUint64(dest, uint64(i))
}

func GetInt64BE(src []byte) int64 {
	return int64(binary.BigEndian.Uint64(src))
}

// Returns whether a + b would be a uint64 overflow
func IsUint64SumOverflow(a, b uint64) bool {
	return math.MaxUint64-a < b
}

//-------------------------------------

type Words256 []Word256

func (ws Words256) Len() int {
	return len(ws)
}

func (ws Words256) Less(i, j int) bool {
	return ws[i].Compare(ws[j]) < 0
}

func (ws Words256) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}

type Tuple256 struct {
	First  Word256
	Second Word256
}

func (tuple Tuple256) Compare(other Tuple256) int {
	firstCompare := tuple.First.Compare(other.First)
	if firstCompare == 0 {
		return tuple.Second.Compare(other.Second)
	} else {
		return firstCompare
	}
}

func Tuple256Split(t Tuple256) (Word256, Word256) {
	return t.First, t.Second
}

type Tuple256Slice []Tuple256

func (p Tuple256Slice) Len() int { return len(p) }
func (p Tuple256Slice) Less(i, j int) bool {
	return p[i].Compare(p[j]) < 0
}
func (p Tuple256Slice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p Tuple256Slice) Sort()         { sort.Sort(p) }
