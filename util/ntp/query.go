package ntp

import "github.com/beevik/ntp"

type Querier interface {
	Query(address string) (*ntp.Response, error)
}

type RemoteQuerier struct{}

func (RemoteQuerier) Query(address string) (*ntp.Response, error) {
	return ntp.Query(address)
}
