//go111:build gtk

package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
)

// ConfigModel manages local node configuration file state for the config editor.
type ConfigModel struct {
	isLocal      bool
	configPath   string
	savedContent string
	defConfig    *config.Config
}

// NewConfigModel loads genesis and config.toml from the node working directory.
func NewConfigModel(workingDir string, isLocal bool) (*ConfigModel, error) {
	if !isLocal {
		return &ConfigModel{
			isLocal: false,
		}, nil
	}

	gen, err := genesis.LoadFromFile(cmd.PactusGenesisPath(workingDir))
	if err != nil {
		return nil, fmt.Errorf("failed to load genesis: %w", err)
	}

	chainType := gen.ChainType()
	configPath := cmd.PactusConfigPath(workingDir)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	savedContent := normalizeConfigContent(string(data))
	defConfig := config.DefaultConfigForChain(chainType)

	return &ConfigModel{
		isLocal:      isLocal,
		configPath:   configPath,
		savedContent: savedContent,
		defConfig:    defConfig,
	}, nil
}

// SavedContent returns the configuration currently stored on disk.
func (m *ConfigModel) SavedContent() string {
	return m.savedContent
}

// DefaultTOML returns the default configuration document for this node's network.
func (m *ConfigModel) DefaultTOML() string {
	defToml, _ := m.defConfig.ToTOML()

	return defToml
}

// IsDirty reports whether editor content differs from the saved file.
func (m *ConfigModel) IsDirty(editorContent string) bool {
	return normalizeConfigContent(editorContent) != m.savedContent
}

// Validate checks TOML syntax and configuration semantics for this node's network.
func (m *ConfigModel) Validate(editorContent string) error {
	conf, err := config.LoadFromToml(editorContent, true, m.defConfig)
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return conf.BasicCheck()
}

// Save validates and writes the configuration to disk, then updates the saved baseline.
func (m *ConfigModel) Save(editorContent string) error {
	if err := m.Validate(editorContent); err != nil {
		return err
	}

	if err := util.WriteFile(m.configPath, []byte(editorContent)); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	m.savedContent = normalizeConfigContent(editorContent)

	return nil
}

func normalizeConfigContent(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}
