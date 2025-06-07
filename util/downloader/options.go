package downloader

import "net/http"

type options struct {
	client     *http.Client
	maxRetries int
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		client:     http.DefaultClient,
		maxRetries: _defaultMaxRetries,
	}
}

func WithCustomClient(client *http.Client) Option {
	return func(o *options) {
		o.client = client
	}
}

func WithMaxRetries(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.maxRetries = n
		}
	}
}
