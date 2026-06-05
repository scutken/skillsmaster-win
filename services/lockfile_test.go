package services

import (
	"os"
	"path/filepath"
	"testing"

	"skillsmaster-win/models"
)

func setupLockFileTest(t *testing.T) (*LockFileManager, string) {
	t.Helper()
	tmpDir := t.TempDir()
	resolver := models.NewTestPathResolver(tmpDir)
	lfm := NewLockFileManager(resolver)
	return lfm, tmpDir
}

func TestRecordInstall(t *testing.T) {
	lfm, _ := setupLockFileTest(t)

	entry := models.LockEntry{
		AgentType:   "claude-code",
		SkillName:   "test-skill",
		SourcePath:  "/source/test-skill",
		TargetPath:  "/target/test-skill",
		Mode:        "copy",
		InstalledAt: "2024-01-01T00:00:00Z",
	}

	if err := lfm.RecordInstall(entry); err != nil {
		t.Fatalf("RecordInstall failed: %v", err)
	}

	entries := lfm.GetEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	key := "claude-code/test-skill"
	e, ok := entries[key]
	if !ok {
		t.Fatalf("expected entry with key %q", key)
	}
	if e.SkillName != "test-skill" {
		t.Errorf("expected skill name 'test-skill', got %q", e.SkillName)
	}
	if e.Mode != "copy" {
		t.Errorf("expected mode 'copy', got %q", e.Mode)
	}
}

func TestRemoveEntry(t *testing.T) {
	lfm, _ := setupLockFileTest(t)

	entry := models.LockEntry{
		AgentType: "codex",
		SkillName: "remove-me",
		Mode:      "copy",
	}
	lfm.RecordInstall(entry)

	entries := lfm.GetEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if err := lfm.RemoveEntry("codex", "remove-me"); err != nil {
		t.Fatalf("RemoveEntry failed: %v", err)
	}

	entries = lfm.GetEntries()
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after removal, got %d", len(entries))
	}
}

func TestRemoveEntry_Nonexistent(t *testing.T) {
	lfm, _ := setupLockFileTest(t)

	// Should not error when removing non-existent entry
	if err := lfm.RemoveEntry("claude-code", "nonexistent"); err != nil {
		t.Errorf("removing nonexistent entry should not error: %v", err)
	}
}

func TestPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	resolver := models.NewTestPathResolver(tmpDir)
	lfm1 := NewLockFileManager(resolver)

	entry := models.LockEntry{
		AgentType:   "hermes",
		SkillName:   "persistent-skill",
		SourcePath:  "/src",
		TargetPath:  "/dst",
		Mode:        "symlink",
		InstalledAt: "2024-06-01T12:00:00Z",
	}
	lfm1.RecordInstall(entry)

	// Create a new manager pointing to the same path
	lfm2 := NewLockFileManager(resolver)
	entries := lfm2.GetEntries()

	key := "hermes/persistent-skill"
	e, ok := entries[key]
	if !ok {
		t.Fatal("entry should persist across manager instances")
	}
	if e.Mode != "symlink" {
		t.Errorf("expected mode 'symlink', got %q", e.Mode)
	}
}

func TestMultipleEntries(t *testing.T) {
	lfm, _ := setupLockFileTest(t)

	agents := []models.AgentType{models.AgentClaudeCode, models.AgentCodex, models.AgentHermes}
	for i, agent := range agents {
		entry := models.LockEntry{
			AgentType: string(agent),
			SkillName: "multi-skill",
			Mode:      "copy",
		}
		if err := lfm.RecordInstall(entry); err != nil {
			t.Fatalf("RecordInstall %d failed: %v", i, err)
		}
	}

	entries := lfm.GetEntries()
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestLockFileCreatedOnDisk(t *testing.T) {
	lfm, tmpDir := setupLockFileTest(t)

	entry := models.LockEntry{
		AgentType: "cursor",
		SkillName: "disk-skill",
		Mode:      "copy",
	}
	lfm.RecordInstall(entry)

	lockPath := filepath.Join(tmpDir, ".skillsmaster", "lockfile.yaml")
	if _, err := os.Stat(lockPath); err != nil {
		t.Errorf("lockfile should exist on disk at %s: %v", lockPath, err)
	}
}
