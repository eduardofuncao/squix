package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eduardofuncao/squix/internal/styles"
	"gopkg.in/yaml.v2"
)

var CfgPath = os.ExpandEnv("$HOME/.config/squix/")
var CfgFile = filepath.Join(CfgPath, "config.yaml")

type KeybindingsConfig map[string][]string

type stringOrSlice []string

func (s *stringOrSlice) UnmarshalYAML(unmarshal func(any) error) error {
	var single string
	if err := unmarshal(&single); err == nil {
		*s = []string{single}
		return nil
	}
	var multi []string
	if err := unmarshal(&multi); err != nil {
		return err
	}
	*s = multi
	return nil
}

func (kc *KeybindingsConfig) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[string]stringOrSlice
	if err := unmarshal(&raw); err != nil {
		return err
	}
	result := make(map[string][]string, len(raw))
	for k, v := range raw {
		result[k] = []string(v)
	}
	*kc = result
	return nil
}

type Config struct {
	CurrentConnection  string                     `yaml:"current_connection"`
	Connections        map[string]*ConnectionYAML `yaml:"connections"`
	ColorScheme        string                     `yaml:"color_scheme"`
	CustomColorScheme  *styles.ColorScheme        `yaml:"custom_colors,omitempty"`
	History            History                    `yaml:"history"`
	DefaultRowLimit    int                        `yaml:"default_row_limit"`
	DefaultColumnWidth int                        `yaml:"default_column_width"`
	UIVisibility       UIVisibility               `yaml:"ui_visibility"`
	Keybindings        KeybindingsConfig          `yaml:"keybindings,omitempty"`
	KeyMap             *KeyMap                    `yaml:"-"`
}

type History struct {
	Size int `yaml:"size"`
}

type UIVisibility struct {
	QueryName         bool `yaml:"query_name"`
	QuerySQL          bool `yaml:"query_sql"`
	TypeDisplay       bool `yaml:"type_display"`
	KeyIcons          bool `yaml:"key_icons"`
	FooterCellContent bool `yaml:"footer_cell_content"`
	FooterStats       bool `yaml:"footer_stats"`
	FooterKeymaps     bool `yaml:"footer_keymaps"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Creating blank config file at", CfgFile)
			cfg := &Config{
				CurrentConnection:  "",
				Connections:        make(map[string]*ConnectionYAML),
				ColorScheme:        "default",
				History:            History{},
				DefaultRowLimit:    1000,
				DefaultColumnWidth: 15,
				UIVisibility: UIVisibility{
					QueryName:         true,
					QuerySQL:          true,
					TypeDisplay:       true,
					KeyIcons:          true,
					FooterCellContent: true,
					FooterStats:       true,
					FooterKeymaps:     true,
				},
			}
			err := cfg.Save()
			if err != nil {
				return nil, err
			}
			return cfg, nil
		}
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.DefaultColumnWidth == 0 {
		cfg.DefaultColumnWidth = 15
	}
	if cfg.DefaultRowLimit == 0 {
		cfg.DefaultRowLimit = 1000
	}

	// Set UI visibility defaults (all true by default)
	if !cfg.UIVisibility.QueryName && !cfg.UIVisibility.QuerySQL &&
		!cfg.UIVisibility.TypeDisplay && !cfg.UIVisibility.KeyIcons &&
		!cfg.UIVisibility.FooterCellContent && !cfg.UIVisibility.FooterStats &&
		!cfg.UIVisibility.FooterKeymaps {
		// All false means config is unset, use defaults
		cfg.UIVisibility.QueryName = true
		cfg.UIVisibility.QuerySQL = true
		cfg.UIVisibility.TypeDisplay = true
		cfg.UIVisibility.KeyIcons = true
		cfg.UIVisibility.FooterCellContent = true
		cfg.UIVisibility.FooterStats = true
		cfg.UIVisibility.FooterKeymaps = true
	}

	cfg.KeyMap = BuildKeyMap(cfg.Keybindings)

	return &cfg, nil
}

func (c *Config) Save() error {
	err := os.MkdirAll(CfgPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(CfgFile, data, 0644)
}
