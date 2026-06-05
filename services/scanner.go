package services

import (
	"os"
	"path/filepath"
	"strings"

	"skillsmaster-win/models"
)

var skipDirs = map[string]bool{
	".git":          true,
	"node_modules":  true,
	".DS_Store":     true,
	"__pycache__":   true,
}

type SkillScanner struct {
	Resolver models.PathResolver
}

func NewSkillScanner(resolver models.PathResolver) *SkillScanner {
	return &SkillScanner{Resolver: resolver}
}

// ScanDirectory scans a directory for SKILL.md files (non-recursive into subdirs beyond one level).
func (s *SkillScanner) ScanDirectory(dir string) []models.Skill {
	var skills []models.Skill

	entries, err := os.ReadDir(dir)
	if err != nil {
		return skills
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if skipDirs[entry.Name()] {
			continue
		}

		skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
		data, err := os.ReadFile(skillPath)
		if err != nil {
			// Also check for lowercase
			skillPath = filepath.Join(dir, entry.Name(), "skill.md")
			data, err = os.ReadFile(skillPath)
			if err != nil {
				continue
			}
		}

		skill, err := ParseSkillFile(string(data), skillPath)
		if err != nil {
			continue
		}
		skills = append(skills, *skill)
	}

	// Also check if dir itself contains SKILL.md (flat layout)
	skillPath := filepath.Join(dir, "SKILL.md")
	if data, err := os.ReadFile(skillPath); err == nil {
		if skill, err := ParseSkillFile(string(data), skillPath); err == nil {
			skills = append(skills, *skill)
		}
	}

	return skills
}

// ScanAgentSkills scans the skills directory for a specific agent type.
func (s *SkillScanner) ScanAgentSkills(agentType models.AgentType) []models.Skill {
	skillsDir := s.Resolver.AgentSkillsDir(agentType)
	skills := s.ScanDirectory(skillsDir)
	for i := range skills {
		skills[i].AgentType = agentType
	}
	return skills
}

// ScanAllAgents scans skills for all known agent types.
func (s *SkillScanner) ScanAllAgents() map[string][]models.Skill {
	result := make(map[string][]models.Skill)
	for _, agentType := range models.AllAgentTypes {
		skills := s.ScanAgentSkills(agentType)
		if len(skills) > 0 {
			result[string(agentType)] = skills
		}
	}
	// Also scan the canonical skills directory
	canonicalDir := s.Resolver.CanonicalSkillsDir()
	canonicalSkills := s.ScanDirectory(canonicalDir)
	if len(canonicalSkills) > 0 {
		for i := range canonicalSkills {
			canonicalSkills[i].Source = "canonical"
		}
		result["canonical"] = canonicalSkills
	}
	return result
}

// ScanSourceSkills scans a source directory for skills to install.
func (s *SkillScanner) ScanSourceSkills(sourceDir string) []models.Skill {
	return s.ScanDirectory(sourceDir)
}

// shouldSkip returns true if the path contains directories to skip.
func shouldSkip(path string) bool {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		if skipDirs[part] {
			return true
		}
	}
	return false
}
