//go:build gtk

package controller

import "time"

const (
	refreshNodeProgressInterval = 1 * time.Second
	refreshNodeInfoInterval     = 10 * time.Second
	refreshCommitteeInterval    = 15 * time.Second
	refreshNetworkInterval      = 15 * time.Second
	refreshValidatorsInterval   = 15 * time.Second
	refreshWalletInterval       = 15 * time.Second
)
