package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerPattern(t *testing.T) {
	tests := []struct {
		basePath        string
		expectedRoot    string
		expectedAPI     string
		expectedSwagger string
	}{
		{
			basePath:        "",
			expectedRoot:    "/",
			expectedAPI:     "/api/",
			expectedSwagger: "/ui/",
		},
		{
			basePath:        "/",
			expectedRoot:    "/",
			expectedAPI:     "/api/",
			expectedSwagger: "/ui/",
		},
		{
			basePath:        "rest",
			expectedRoot:    "/rest/",
			expectedAPI:     "/rest/api/",
			expectedSwagger: "/rest/ui/",
		},
		{
			basePath:        "/rest",
			expectedRoot:    "/rest/",
			expectedAPI:     "/rest/api/",
			expectedSwagger: "/rest/ui/",
		},
		{
			basePath:        "rest/",
			expectedRoot:    "/rest/",
			expectedAPI:     "/rest/api/",
			expectedSwagger: "/rest/ui/",
		},
	}

	for _, tt := range tests {
		cfg := &Config{
			BasePath: tt.basePath,
		}

		assert.Equal(t, tt.expectedRoot, cfg.rootPattern())
		assert.Equal(t, tt.expectedAPI, cfg.apiPattern())
		assert.Equal(t, tt.expectedSwagger, cfg.swaggerPattern())
	}
}
