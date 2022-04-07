package wallet

import (
	"encoding/base64"
	"strconv"
)

type params map[string]string

func newParams() params {
	return make(map[string]string)
}
func (p params) SetUint8(key string, val uint8) {
	p.SetUint32(key, uint32(val))
}

func (p params) SetUint32(key string, val uint32) {
	p[key] = strconv.FormatInt(int64(val), 10)
}

func (p params) SetBytes(key string, val []byte) {
	p[key] = base64.StdEncoding.EncodeToString(val)
}

func (p params) SetString(key string, val string) {
	p[key] = val
}

func (p params) GetUint8(key string) uint8 {
	return uint8(p.GetUint32(key))
}
func (p params) GetUint32(key string) uint32 {
	val, err := strconv.ParseUint(p[key], 10, 32)
	exitOnErr(err)
	return uint32(val)
}

func (p params) GetBytes(key string) []byte {
	val, err := base64.StdEncoding.DecodeString(p[key])
	exitOnErr(err)
	return val
}

func (p params) GetString(key string) string {
	return p[key]
}
