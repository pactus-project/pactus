package version

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var NodeAgent = Agent{
	Version:         NodeVersion,
	ProtocolVersion: 1,
	OS:              runtime.GOOS,
	Arch:            runtime.GOARCH,
}

type Agent struct {
	AppType         string
	Version         Version
	ProtocolVersion uint
	OS              string
	Arch            string
}

// ParseAgent parses a string into an Agent struct.
func ParseAgent(agentStr string) (Agent, error) {
	var agent Agent

	parts := strings.Split(agentStr, "/")
	for _, part := range parts {
		fields := strings.Split(part, "=")
		if len(fields) != 2 {
			return agent, fmt.Errorf("invalid field format in agent string")
		}
		key := fields[0]
		value := fields[1]

		switch key {
		case "node":
			agent.AppType = value
		case "node-version":
			v, err := ParseVersion(value)
			if err != nil {
				return agent, fmt.Errorf("failed to parse version: %w", err)
			}
			agent.Version = v
		case "protocol-version":
			protocolVersion, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return agent, fmt.Errorf("failed to parse protocol version: %w", err)
			}
			agent.ProtocolVersion = uint(protocolVersion)
		case "os":
			agent.OS = value
		case "arch":
			agent.Arch = value
		default:
			return agent, fmt.Errorf("unknown key in agent string: %s", key)
		}
	}

	return agent, nil
}

func (a *Agent) String() string {
	return fmt.Sprintf("node=%s/node-version=%s/protocol-version=%d/os=%s/arch=%s",
		a.AppType, a.Version.String(), a.ProtocolVersion, a.OS, a.Arch)
}
