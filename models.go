package main

import (
	"context"
	"time"

	"github.com/gobwas/glob"
)

// BackupConfig stores configuration for how to handle backups
type BackupConfig struct {
	KeepSaves    int `yaml:"keep_saves" json:"keepSaves"`
	MaxMBPerGame int `yaml:"max_mb_per_game" json:"maxMBPerGame"`
}

// CompactionConfig stores configuration for how to handle backups when they enter a state for compaction
type CompactionConfig struct {
	KeepSaves        int `yaml:"keep_saves" json:"keepSaves"`
	CompactAfterDays int `yaml:"compact_after_days" json:"compactAfterDays"`
}

// Config stores application configuration
type Config struct {
	DisabledRules    []string          `yaml:"disabled_rules" json:"disabledRules"`
	PathSeparator    string            `json:"pathSeparator"`
	Backups          *BackupConfig     `yaml:"backups" json:"backups"`
	Compaction       *CompactionConfig `yaml:"compaction" json:"compaction"`
	RulesLastUpdated time.Time         `yaml:"rules_last_updated" json:"rulesLastUpdated"`
	RulesAutoUpdate  bool              `yaml:"rules_auto_update" json:"rulesAutoUpdate"`
}

// RulePlatform is the "platform" section of a rule
type RulePlatform struct {
	Executable string   `yaml:"executable" json:"executable"`
	Savegames  []string `yaml:"savegames" json:"savegames"`
}

// Rule for how to manage a game
type Rule struct {
	Name      string                  `yaml:"name" json:"name"`
	Issues    string                  `yaml:"issues" json:"issues"`
	Platforms map[string]RulePlatform `yaml:"platforms" json:"platforms"`
}

// ActiveRule is a game rule for a game that has been detected as running
type ActiveRule struct {
	RuleFilename   string       `json:"ruleFilename"`
	Name           string       `json:"name"`
	Issues         string       `yaml:"issues" json:"issues"`
	Platform       RulePlatform `json:"platform"`
	executableGlob glob.Glob
}

// BackupMetadata stores the metadata for a backed up savefile
type BackupMetadata struct {
	Filename     string    `json:"filename" yaml:"-"`
	Source       string    `yaml:"source" json:"source"`
	BackupTime   time.Time `yaml:"backup_time" json:"backupTime"`
	LastModified time.Time `yaml:"last_modified" json:"lastModified"`
}

// Monitor information for what paths we're monitoring
type Monitor struct {
	Path         string `json:"path"`
	RuleFilename string `json:"ruleFilename"`
}

// App is the root application
type App struct {
	ctx            context.Context
	Config         *Config                     `json:"config"`
	Rules          map[string]ActiveRule       `json:"rules"`
	Backups        map[string][]BackupMetadata `json:"activeBackups"`
	ActiveMonitors []Monitor                   `json:"activeMonitors"`
	ActiveRules    map[string]ActiveRule       `json:"activeRules"`
	ProcessList    []string                    `json:"processList"`
	BasePath       string                      `json:"basePath"`
	Errors         []string                    `json:"errors"`
	Events         []string                    `json:"events"`
	exit           chan bool
}
