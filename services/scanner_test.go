package services

import (
	"os"
	"path/filepath"
	"testing"

	"skillsmaster-win/models"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Create a skill directory with SKILL.md
	skillDir := filepath.Join(tmpDir, "test-skill")
	os.MkdirAll(skillDir, 0755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(`---
name: test-skill
description: A test skill
version: "1.0.0"
---
# Test Skill
Body content here.
`), 0644)

	// Create another skill
	skillDir2 := filepath.Join(tmpDir, "another-skill")
	os.MkdirAll(skillDir2, 0755)
	os.WriteFile(filepath.Join(skillDir2, "SKILL.md"), []byte(`# Another Skill
No frontmatter, just content.
`), 0644)

	// Create a .git directory (should be skipped)
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0755)
	os.WriteFile(filepath.Join(tmpDir, ".git", "SKILL.md"), []byte(`---
name: git-skill
---
Should not be found.
`), 0644)

	// Create node_modules (should be skipped)
	os.MkdirAll(filepath.Join(tmpDir, "node_modules"), 0755)

	return tmpDir
}

func TestScanDirectory(t *testing.T) {
	tmpDir := setupTestDir(t)
	resolver := models.NewTestPathResolver(tmpDir)
	scanner := NewSkillScanner(resolver)

	skills := scanner.ScanDirectory(tmpDir)
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	names := make(map[string]bool)
	for _, s := range skills {
		names[s.Name] = true
	}
	if !names["test-skill"] {
		t.Error("expected to find 'test-skill'")
	}
	if !names["another-skill"] {
		t.Error("expected to find 'another-skill'")
	}
}

func TestScanDirectory_SkipsGitAndNodeModules(t *testing.T) {
	tmpDir := setupTestDir(t)
	resolver := models.NewTestPathResolver(tmpDir)
	scanner := NewSkillScanner(resolver)

	skills := scanner.ScanDirectory(tmpDir)
	for _, s := range skills {
		if s.Name == "git-skill" {
			t.Error("should not find skills in .git directory")
		}
	}
}

func TestScanAgentSkills(t *testing.T) {
	tmpDir := t.TempDir()
	resolver := models.NewTestPathResolver(tmpDir)
	scanner := NewSkillScanner(resolver)

	// Create agent skills directory
	skillsDir := resolver.AgentSkillsDir(models.AgentClaudeCode)
	skillDir := filepath.Join(skillsDir, "my-skill")
	os.MkdirAll(skillDir, 0755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(`---
name: my-skill
description: For Claude
---
# My Skill
`), 0644)

	skills := scanner.ScanAgentSkills(models.AgentClaudeCode)
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].AgentType != models.AgentClaudeCode {
		t.Errorf("expected agent type 'claude-code', got %q", skills[0].AgentType)
	}
}

func TestScanAllAgents(t *testing.T) {
	tmpDir := t.TempDir()
	resolver := models.NewTestPathResolver(tmpDir)
	scanner := NewSkillScanner(resolver)

	// Create skills for two different agents
	for _, agent := range []models.AgentType{models.AgentClaudeCode, models.AgentCodex} {
		skillsDir := resolver.AgentSkillsDir(agent)
		skillDir := filepath.Join(skillsDir, "skill-"+string(agent))
		os.MkdirAll(skillDir, 0755)
		os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(`---
name: skill
---
# Skill
`), 0644)
	}

	result := scanner.ScanAllAgents()
	if len(result) < 2 {
		t.Errorf("expected at least 2 agents with skills, got %d", len(result))
	}
}

func TestScanDirectory_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	resolver := models.NewTestPathResolver(tmpDir)
	scanner := NewSkillScanner(resolver)

	skills := scanner.ScanDirectory(tmpDir)
	if len(skills) != 0 {
		t.Errorf("expected 0 skills in empty dir, got %d", len(skills))
	}
}

func TestScanDirectory_NonexistentDir(t *testing.T) {
	resolver := models.NewTestPathResolver("/nonexistent/path")
	scanner := NewSkillScanner(resolver)

	skills := scanner.ScanDirectory("/nonexistent/path/nowhere")
	if len(skills) != 0 {
		t.Errorf("expected 0 skills for nonexistent dir, got %d", len(skills))
	}
}
