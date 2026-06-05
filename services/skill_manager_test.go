package services

import (
	"os"
	"path/filepath"
	"testing"

	"skillsmaster-win/models"
)

func setupManagerTest(t *testing.T) (*SkillManager, string) {
	t.Helper()
	tmpDir := t.TempDir()
	resolver := models.NewTestPathResolver(tmpDir)
	manager := NewSkillManager(resolver)
	return manager, tmpDir
}

func createSkillInDir(t *testing.T, dir, name, description string) string {
	t.Helper()
	skillDir := filepath.Join(dir, name)
	os.MkdirAll(skillDir, 0755)
	content := `---
name: ` + name + `
description: ` + description + `
version: "1.0.0"
tags:
  - test
---
# ` + name + `

Content of ` + name + `
`
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644)
	return skillDir
}

func TestGetAgents(t *testing.T) {
	manager, _ := setupManagerTest(t)
	agents := manager.GetAgents()
	if len(agents) < 12 {
		t.Errorf("expected at least 12 agents, got %d", len(agents))
	}

	found := make(map[string]bool)
	for _, a := range agents {
		found[string(a.Type)] = true
	}
	for _, expected := range []string{"claude-code", "codex", "gemini-cli", "cursor", "hermes"} {
		if !found[expected] {
			t.Errorf("expected agent %q not found", expected)
		}
	}
}

func TestInstallAndUninstall(t *testing.T) {
	manager, tmpDir := setupManagerTest(t)
	sourceDir := createSkillInDir(t, tmpDir, "my-skill", "A great skill")

	// Install
	err := manager.InstallSkill(sourceDir, models.AgentClaudeCode, InstallModeCopy)
	if err != nil {
		t.Fatalf("InstallSkill failed: %v", err)
	}

	// Verify installed
	skills := manager.ScanAgentSkills(models.AgentClaudeCode)
	found := false
	for _, s := range skills {
		if s.Name == "my-skill" {
			found = true
		}
	}
	if !found {
		t.Error("skill should be found after install")
	}

	// Uninstall
	err = manager.UninstallSkill("my-skill", models.AgentClaudeCode)
	if err != nil {
		t.Fatalf("UninstallSkill failed: %v", err)
	}

	// Verify removed
	skills = manager.ScanAgentSkills(models.AgentClaudeCode)
	for _, s := range skills {
		if s.Name == "my-skill" {
			t.Error("skill should not be found after uninstall")
		}
	}
}

func TestGetAllSkills(t *testing.T) {
	manager, _ := setupManagerTest(t)

	// Create skills for multiple agents
	for _, agent := range []models.AgentType{models.AgentClaudeCode, models.AgentCodex} {
		skillsDir := manager.Resolver.AgentSkillsDir(agent)
		createSkillInDir(t, skillsDir, "skill-"+string(agent), "Skill for "+string(agent))
	}

	allSkills := manager.GetAllSkills()
	if len(allSkills) < 2 {
		t.Errorf("expected skills for at least 2 agents, got %d", len(allSkills))
	}
}

func TestReadFile(t *testing.T) {
	manager, tmpDir := setupManagerTest(t)
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("hello world"), 0644)

	content, err := manager.ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("expected 'hello world', got %q", content)
	}
}

func TestReadFile_NotFound(t *testing.T) {
	manager, _ := setupManagerTest(t)
	_, err := manager.ReadFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestSaveFile(t *testing.T) {
	manager, tmpDir := setupManagerTest(t)
	testFile := filepath.Join(tmpDir, "subdir", "saved.txt")

	err := manager.SaveFile(testFile, "saved content")
	if err != nil {
		t.Fatalf("SaveFile failed: %v", err)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("cannot read saved file: %v", err)
	}
	if string(data) != "saved content" {
		t.Errorf("expected 'saved content', got %q", string(data))
	}
}

func TestBrowseAgentFiles(t *testing.T) {
	manager, _ := setupManagerTest(t)
	skillsDir := manager.Resolver.AgentSkillsDir(models.AgentClaudeCode)
	createSkillInDir(t, skillsDir, "browsable-skill", "Browse me")

	nodes := manager.BrowseAgentFiles(models.AgentClaudeCode)
	if len(nodes) == 0 {
		t.Error("expected at least one file node")
	}

	found := false
	for _, n := range nodes {
		if n.Name == "browsable-skill" && n.IsDir {
			found = true
			if len(n.Children) == 0 {
				t.Error("expected children in skill directory")
			}
		}
	}
	if !found {
		t.Error("expected browsable-skill directory node")
	}
}

func TestIntegration_Lifecycle(t *testing.T) {
	manager, tmpDir := setupManagerTest(t)

	// 1. Create source skill
	sourceDir := createSkillInDir(t, tmpDir, "lifecycle-skill", "Full lifecycle test")

	// 2. Scan - should not be in any agent yet
	agentSkills := manager.ScanAgentSkills(models.AgentClaudeCode)
	for _, s := range agentSkills {
		if s.Name == "lifecycle-skill" {
			t.Error("skill should not exist in agent before install")
		}
	}

	// 3. Install
	err := manager.InstallSkill(sourceDir, models.AgentClaudeCode, InstallModeCopy)
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	// 4. Scan again - should find it
	agentSkills = manager.ScanAgentSkills(models.AgentClaudeCode)
	found := false
	for _, s := range agentSkills {
		if s.Name == "lifecycle-skill" {
			found = true
			if s.Manifest.Description != "Full lifecycle test" {
				t.Errorf("wrong description: %q", s.Manifest.Description)
			}
		}
	}
	if !found {
		t.Fatal("skill not found after install")
	}

	// 5. Check lockfile
	entries := manager.LockFile.GetEntries()
	key := string(models.AgentClaudeCode) + "/lifecycle-skill"
	if _, ok := entries[key]; !ok {
		t.Error("expected lockfile entry")
	}

	// 6. GetAllSkills should show it
	allSkills := manager.GetAllSkills()
	claudeSkills, ok := allSkills[string(models.AgentClaudeCode)]
	if !ok {
		t.Fatal("expected claude-code in all skills")
	}
	foundDTO := false
	for _, dto := range claudeSkills {
		if dto.Name == "lifecycle-skill" {
			foundDTO = true
		}
	}
	if !foundDTO {
		t.Error("expected lifecycle-skill in GetAllSkills")
	}

	// 7. Uninstall
	err = manager.UninstallSkill("lifecycle-skill", models.AgentClaudeCode)
	if err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}

	// 8. Verify gone
	agentSkills = manager.ScanAgentSkills(models.AgentClaudeCode)
	for _, s := range agentSkills {
		if s.Name == "lifecycle-skill" {
			t.Error("skill should be gone after uninstall")
		}
	}

	// 9. Lockfile entry removed
	entries = manager.LockFile.GetEntries()
	if _, ok := entries[key]; ok {
		t.Error("lockfile entry should be removed after uninstall")
	}
}
