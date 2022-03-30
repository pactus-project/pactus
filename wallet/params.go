package wallet

import (
	"encoding/base64"
	"strconv"
)

type params map[string]string

func (p params) setUint8(key string, val uint8) {
	p.setUint32(key, uint32(val))
}

func (p params) setUint32(key string, val uint32) {
	p[key] = strconv.FormatInt(int64(val), 10)
}

func (p params) setBytes(key string, val []byte) {
	p[key] = base64.StdEncoding.EncodeToString([]byte(val))
}

func (p params) getUint8(key string) uint8 {
	return uint8(p.getUint32(key))
}
func (p params) getUint32(key string) uint32 {
	val, err := strconv.ParseUint(p[key], 10, 32)
	exitOnErr(err)
	return uint32(val)
}

func (p params) getBytes(key string) []byte {
	val, err := base64.StdEncoding.DecodeString(p[key])
	exitOnErr(err)
	return val
}
