package models

import (
	"os"
	"path/filepath"
)

// PathResolver abstracts filesystem paths so tests can use /tmp/ directories.
type PathResolver interface {
	HomeDir() string
	SkillsMasterDir() string
	CanonicalSkillsDir() string
	LockFilePath() string
	AgentConfigDir(agentType AgentType) string
	AgentSkillsDir(agentType AgentType) string
}

// DefaultPathResolver is the production resolver using os.UserHomeDir().
type DefaultPathResolver struct {
	home string
}

func NewDefaultPathResolver() *DefaultPathResolver {
	home, _ := os.UserHomeDir()
	return &DefaultPathResolver{home: home}
}

func (d *DefaultPathResolver) HomeDir() string {
	return d.home
}

func (d *DefaultPathResolver) SkillsMasterDir() string {
	return filepath.Join(d.home, ".skillsmaster")
}

func (d *DefaultPathResolver) CanonicalSkillsDir() string {
	return filepath.Join(d.home, ".skillsmaster", "canonical")
}

func (d *DefaultPathResolver) LockFilePath() string {
	return filepath.Join(d.home, ".skillsmaster", "lockfile.yaml")
}

func (d *DefaultPathResolver) AgentConfigDir(agentType AgentType) string {
	cfg := GetAgentConfig(agentType)
	return cfg.ConfigDir
}

func (d *DefaultPathResolver) AgentSkillsDir(agentType AgentType) string {
	cfg := GetAgentConfig(agentType)
	return cfg.SkillsDir
}

// TestPathResolver rewrites all paths under a configurable root directory.
type TestPathResolver struct {
	Root string
}

func NewTestPathResolver(root string) *TestPathResolver {
	return &TestPathResolver{Root: root}
}

func (t *TestPathResolver) HomeDir() string {
	return t.Root
}

func (t *TestPathResolver) SkillsMasterDir() string {
	return filepath.Join(t.Root, ".skillsmaster")
}

func (t *TestPathResolver) CanonicalSkillsDir() string {
	return filepath.Join(t.Root, ".skillsmaster", "canonical")
}

func (t *TestPathResolver) LockFilePath() string {
	return filepath.Join(t.Root, ".skillsmaster", "lockfile.yaml")
}

func (t *TestPathResolver) AgentConfigDir(agentType AgentType) string {
	return filepath.Join(t.Root, ".config", string(agentType))
}

func (t *TestPathResolver) AgentSkillsDir(agentType AgentType) string {
	return filepath.Join(t.Root, ".config", string(agentType), "skills")
}
