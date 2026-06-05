package services

import (
	"errors"
	"path/filepath"
	"strings"

	"skillsmaster-win/models"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyContent = errors.New("empty content")
)

// ParseSkillFile parses a SKILL.md file with optional YAML frontmatter.
func ParseSkillFile(content string, path string) (*models.Skill, error) {
	if strings.TrimSpace(content) == "" {
		return nil, ErrEmptyContent
	}

	skill := &models.Skill{
		Path: path,
		Name: deriveName(path),
	}

	body := content
	manifest := models.SkillManifest{}

	// Check for YAML frontmatter between --- delimiters
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		trimmed := strings.TrimSpace(content)
		rest := trimmed[3:] // skip first ---
		if idx := strings.Index(rest, "---"); idx >= 0 {
			yamlPart := rest[:idx]
			body = rest[idx+3:]
			if err := yaml.Unmarshal([]byte(yamlPart), &manifest); err != nil {
				// If YAML is malformed, treat entire content as body
				body = content
				manifest = models.SkillManifest{}
			}
		}
	}

	if manifest.Name != "" {
		skill.Name = manifest.Name
	}
	skill.Manifest = manifest
	skill.Body = strings.TrimSpace(body)
	return skill, nil
}

func deriveName(path string) string {
	dir := filepath.Dir(path)
	name := filepath.Base(dir)
	if name == "." || name == "/" {
		name = filepath.Base(path)
	}
	return name
}
