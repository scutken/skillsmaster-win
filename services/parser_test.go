package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSkillFile_Valid(t *testing.T) {
	content := `---
name: test-skill
description: A test skill
version: "1.0.0"
author: testuser
tags:
  - testing
  - example
platforms:
  - linux
  - windows
---
# Test Skill

This is the body of the skill.
`
	skill, err := ParseSkillFile(content, "/some/path/test-skill/SKILL.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got %q", skill.Name)
	}
	if skill.Manifest.Description != "A test skill" {
		t.Errorf("expected description 'A test skill', got %q", skill.Manifest.Description)
	}
	if skill.Manifest.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", skill.Manifest.Version)
	}
	if len(skill.Manifest.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(skill.Manifest.Tags))
	}
	if !strings.Contains(skill.Body, "This is the body") {
		t.Error("body should contain 'This is the body'")
	}
}

func TestParseSkillFile_NoFrontmatter(t *testing.T) {
	content := `# Just a markdown file

No frontmatter here.
`
	skill, err := ParseSkillFile(content, "/some/path/my-skill/SKILL.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "my-skill" {
		t.Errorf("expected name 'my-skill', got %q", skill.Name)
	}
	if skill.Body != strings.TrimSpace(content) {
		t.Error("body should be entire content when no frontmatter")
	}
}

func TestParseSkillFile_MalformedYAML(t *testing.T) {
	content := `---
name: test
  invalid: [yaml: broken
---
Body here
`
	skill, err := ParseSkillFile(content, "/some/path/SKILL.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should fall back to treating entire content as body
	if !strings.Contains(skill.Body, "---") {
		t.Error("malformed YAML should result in full content as body")
	}
}

func TestParseSkillFile_EmptyContent(t *testing.T) {
	_, err := ParseSkillFile("", "/some/path/SKILL.md")
	if err != ErrEmptyContent {
		t.Errorf("expected ErrEmptyContent, got %v", err)
	}
}

func TestParseSkillFile_WhitespaceOnly(t *testing.T) {
	_, err := ParseSkillFile("   \n\t  ", "/some/path/SKILL.md")
	if err != ErrEmptyContent {
		t.Errorf("expected ErrEmptyContent, got %v", err)
	}
}

func TestParseSkillFile_FrontmatterOnly(t *testing.T) {
	content := `---
name: empty-body
description: No body
---
`
	skill, err := ParseSkillFile(content, "/some/path/SKILL.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "empty-body" {
		t.Errorf("expected name 'empty-body', got %q", skill.Name)
	}
}

func TestParseSkillFile_DeriveNameFromDir(t *testing.T) {
	// Create a real temp dir to test name derivation
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my-awesome-skill")
	os.MkdirAll(skillDir, 0755)
	path := filepath.Join(skillDir, "SKILL.md")

	skill, err := ParseSkillFile("# Simple skill", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "my-awesome-skill" {
		t.Errorf("expected name 'my-awesome-skill', got %q", skill.Name)
	}
}
