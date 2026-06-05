package services

import (
	"os"
	"path/filepath"
	"testing"

	"skillsmaster-win/models"
)

func setupInstallerTest(t *testing.T) (*SkillInstaller, string) {
	t.Helper()
	tmpDir := t.TempDir()
	resolver := models.NewTestPathResolver(tmpDir)
	lockFile := NewLockFileManager(resolver)
	installer := NewSkillInstaller(resolver, lockFile)
	return installer, tmpDir
}

func createTestSkillDir(t *testing.T, parentDir, name string) string {
	t.Helper()
	skillDir := filepath.Join(parentDir, name)
	os.MkdirAll(skillDir, 0755)
	content := `---
name: ` + name + `
description: Test skill ` + name + `
version: "1.0.0"
---
# ` + name + `

Body of ` + name + `
`
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644)
	// Add an extra file
	os.WriteFile(filepath.Join(skillDir, "config.json"), []byte(`{"key": "value"}`), 0644)
	return skillDir
}

func TestInstall_Copy(t *testing.T) {
	installer, tmpDir := setupInstallerTest(t)
	sourceDir := createTestSkillDir(t, tmpDir, "source-skill")

	err := installer.Install(sourceDir, models.AgentClaudeCode, InstallModeCopy)
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	// Verify files exist
	skillsDir := installer.Resolver.AgentSkillsDir(models.AgentClaudeCode)
	targetDir := filepath.Join(skillsDir, "source-skill")

	if _, err := os.Stat(filepath.Join(targetDir, "SKILL.md")); err != nil {
		t.Error("SKILL.md should exist in target")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "config.json")); err != nil {
		t.Error("config.json should exist in target")
	}
}

func TestInstall_Symlink(t *testing.T) {
	installer, tmpDir := setupInstallerTest(t)
	sourceDir := createTestSkillDir(t, tmpDir, "linked-skill")

	err := installer.Install(sourceDir, models.AgentClaudeCode, InstallModeSymlink)
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	skillsDir := installer.Resolver.AgentSkillsDir(models.AgentClaudeCode)
	targetDir := filepath.Join(skillsDir, "linked-skill")

	info, err := os.Lstat(targetDir)
	if err != nil {
		t.Fatalf("target should exist: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("target should be a symlink")
	}
}

func TestInstall_DefaultMode(t *testing.T) {
	installer, tmpDir := setupInstallerTest(t)
	sourceDir := createTestSkillDir(t, tmpDir, "default-skill")

	// Empty mode should default to copy
	err := installer.Install(sourceDir, models.AgentClaudeCode, "")
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	skillsDir := installer.Resolver.AgentSkillsDir(models.AgentClaudeCode)
	targetDir := filepath.Join(skillsDir, "default-skill")
	if _, err := os.Stat(filepath.Join(targetDir, "SKILL.md")); err != nil {
		t.Error("SKILL.md should exist in target")
	}
}

func TestUninstall(t *testing.T) {
	installer, tmpDir := setupInstallerTest(t)
	sourceDir := createTestSkillDir(t, tmpDir, "remove-me")

	err := installer.Install(sourceDir, models.AgentClaudeCode, InstallModeCopy)
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	skillsDir := installer.Resolver.AgentSkillsDir(models.AgentClaudeCode)
	targetDir := filepath.Join(skillsDir, "remove-me")

	// Verify it exists
	if _, err := os.Stat(targetDir); err != nil {
		t.Fatal("skill should be installed before uninstall")
	}

	// Uninstall
	if err := installer.Uninstall(targetDir); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}

	// Verify it's gone
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Error("skill should be gone after uninstall")
	}
}

func TestUninstall_AlreadyGone(t *testing.T) {
	installer, _ := setupInstallerTest(t)
	err := installer.Uninstall("/nonexistent/path")
	if err != nil {
		t.Errorf("uninstalling nonexistent path should not error: %v", err)
	}
}

func TestCheckInstallStatus(t *testing.T) {
	installer, tmpDir := setupInstallerTest(t)
	sourceDir := createTestSkillDir(t, tmpDir, "status-skill")

	// Should be not installed
	status := installer.CheckInstallStatus(sourceDir, models.AgentClaudeCode)
	if status != "not_installed" {
		t.Errorf("expected 'not_installed', got %q", status)
	}

	// Install it
	installer.Install(sourceDir, models.AgentClaudeCode, InstallModeCopy)

	// Should be installed
	status = installer.CheckInstallStatus(sourceDir, models.AgentClaudeCode)
	if status != "installed" {
		t.Errorf("expected 'installed', got %q", status)
	}
}

func TestInstall_RecordsInLockfile(t *testing.T) {
	installer, tmpDir := setupInstallerTest(t)
	sourceDir := createTestSkillDir(t, tmpDir, "locked-skill")

	err := installer.Install(sourceDir, models.AgentClaudeCode, InstallModeCopy)
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	entries := installer.LockFile.GetEntries()
	found := false
	for _, e := range entries {
		if e.SkillName == "locked-skill" && e.AgentType == string(models.AgentClaudeCode) {
			found = true
			if e.Mode != "copy" {
				t.Errorf("expected mode 'copy', got %q", e.Mode)
			}
		}
	}
	if !found {
		t.Error("expected lockfile entry for installed skill")
	}
}
