package version_test

import (
	"testing"

	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestAgentString(t *testing.T) {
	agent := version.Agent{
		AppType:         "gui",
		Version:         version.Version{Major: 1, Minor: 2, Patch: 3, Meta: "beta"},
		ProtocolVersion: 2,
		OS:              "linux",
		Arch:            "amd64",
	}

	expected := "node=gui/node-version=1.2.3-beta/protocol-version=2/os=linux/arch=amd64"
	result := agent.String()

	assert.Equal(t, expected, result)
}

func TestParseAgent(t *testing.T) {
	tests := []struct {
		name        string
		agentStr    string
		expected    version.Agent
		expectedErr bool
	}{
		{
			name:     "Valid Agent String",
			agentStr: "node=gui/node-version=1.2.3-beta/protocol-version=2/os=linux/arch=amd64",
			expected: version.Agent{
				AppType: "gui", Version: version.Version{
					Major: 1, Minor: 2, Patch: 3, Meta: "beta",
				},
				ProtocolVersion: 2, OS: "linux", Arch: "amd64",
			},
			expectedErr: false,
		},
		{
			name:        "Invalid Agent String (Invalid Protocol Version)",
			agentStr:    "node=gui/node-version=1.2.3-beta/protocol-version=abc/os=linux/arch=amd64",
			expected:    version.Agent{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := version.ParseAgent(tt.agentStr)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, agent)
			}
		})
	}
}
