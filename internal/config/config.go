package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultPanelRatio = 0.4
	minPanelRatio     = 0.2
	maxPanelRatio     = 0.8
)

// FocusConfig holds pomodoro/focus timer settings.
type FocusConfig struct {
	WorkDuration       int  `json:"work_duration_min"`
	ShortBreakDuration int  `json:"short_break_duration_min"`
	LongBreakDuration  int  `json:"long_break_duration_min"`
	LongBreakInterval  int  `json:"long_break_interval"`
	DailyGoal          int  `json:"daily_goal"`
	AutoStartBreak     bool `json:"auto_start_break"`
	Sound              bool `json:"sound"`
}

// Config holds user-configurable settings for the application.
type Config struct {
	PanelRatio float64     `json:"panel_ratio"`
	Focus      FocusConfig `json:"focus"`
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() Config {
	return Config{
		PanelRatio: defaultPanelRatio,
		Focus: FocusConfig{
			WorkDuration:       25,
			ShortBreakDuration: 5,
			LongBreakDuration:  15,
			LongBreakInterval:  4,
			DailyGoal:          8,
			Sound:              true,
		},
	}
}

// validate clamps PanelRatio to the allowed range and applies defaults for
// zero values.
func (c *Config) validate() {
	if c.PanelRatio == 0 {
		c.PanelRatio = defaultPanelRatio
	}
	if c.PanelRatio < minPanelRatio {
		c.PanelRatio = minPanelRatio
	}
	if c.PanelRatio > maxPanelRatio {
		c.PanelRatio = maxPanelRatio
	}

	if c.Focus.WorkDuration == 0 {
		c.Focus.WorkDuration = 25
	}
	if c.Focus.ShortBreakDuration == 0 {
		c.Focus.ShortBreakDuration = 5
	}
	if c.Focus.LongBreakDuration == 0 {
		c.Focus.LongBreakDuration = 15
	}
	if c.Focus.LongBreakInterval == 0 {
		c.Focus.LongBreakInterval = 4
	}
	if c.Focus.DailyGoal == 0 {
		c.Focus.DailyGoal = 8
	}

	// Clamp durations.
	if c.Focus.WorkDuration < 1 {
		c.Focus.WorkDuration = 1
	}
	if c.Focus.WorkDuration > 120 {
		c.Focus.WorkDuration = 120
	}
	if c.Focus.ShortBreakDuration < 1 {
		c.Focus.ShortBreakDuration = 1
	}
	if c.Focus.LongBreakDuration < 1 {
		c.Focus.LongBreakDuration = 1
	}
	if c.Focus.LongBreakInterval < 1 {
		c.Focus.LongBreakInterval = 1
	}
	if c.Focus.LongBreakInterval > 10 {
		c.Focus.LongBreakInterval = 10
	}
}

// Path returns the absolute path to the config file (~/.todo-app/config.json).
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("config path: %w", err)
	}
	return filepath.Join(home, ".todo-app", "config.json"), nil
}

// Load reads the configuration from the config file. If the file does not
// exist, it returns DefaultConfig without error.
func Load() (Config, error) {
	p, err := Path()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), fmt.Errorf("parse config: %w", err)
	}
	cfg.validate()
	return cfg, nil
}

// Save writes the configuration to the config file as formatted JSON. It
// creates the parent directory if necessary.
func Save(cfg Config) error {
	cfg.validate()

	p, err := Path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(p, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
