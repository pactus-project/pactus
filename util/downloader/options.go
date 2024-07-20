package downloader

import "net/http"

type options struct {
	client *http.Client
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		client: http.DefaultClient,
	}
}

func WithCustomClient(client *http.Client) Option {
	return func(o *options) {
		o.client = client
	}
}
