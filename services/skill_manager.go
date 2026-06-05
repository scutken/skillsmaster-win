package services

import (
	"fmt"
	"os"
	"path/filepath"

	"skillsmaster-win/models"
)

type SkillManager struct {
	Resolver  models.PathResolver
	Scanner   *SkillScanner
	Installer *SkillInstaller
	LockFile  *LockFileManager
}

func NewSkillManager(resolver models.PathResolver) *SkillManager {
	lockFile := NewLockFileManager(resolver)
	scanner := NewSkillScanner(resolver)
	installer := NewSkillInstaller(resolver, lockFile)

	return &SkillManager{
		Resolver:  resolver,
		Scanner:   scanner,
		Installer: installer,
		LockFile:  lockFile,
	}
}

// GetAgents returns all agent configurations.
func (sm *SkillManager) GetAgents() []models.AgentConfig {
	return models.AllAgentConfigs()
}

// GetAllSkills returns skills for all agents as DTOs.
func (sm *SkillManager) GetAllSkills() map[string][]models.SkillDTO {
	allSkills := sm.Scanner.ScanAllAgents()
	lockEntries := sm.LockFile.GetEntries()

	result := make(map[string][]models.SkillDTO)
	for agentKey, skills := range allSkills {
		dtos := make([]models.SkillDTO, 0, len(skills))
		for _, skill := range skills {
			dto := models.SkillDTO{
				Name:        skill.Name,
				Description: skill.Manifest.Description,
				Version:     skill.Manifest.Version,
				Author:      skill.Manifest.Author,
				Tags:        skill.Manifest.Tags,
				Platforms:   skill.Manifest.Platforms,
				Path:        skill.Path,
				AgentType:   string(skill.AgentType),
				Body:        skill.Body,
			}
			// Check if installed via lockfile
			key := string(skill.AgentType) + "/" + skill.Name
			if _, ok := lockEntries[key]; ok {
				dto.IsInstalled = true
			}
			dtos = append(dtos, dto)
		}
		result[agentKey] = dtos
	}
	return result
}

// InstallSkill installs a skill to an agent's skills directory.
func (sm *SkillManager) InstallSkill(sourceDir string, agentType models.AgentType, mode InstallMode) error {
	return sm.Installer.Install(sourceDir, agentType, mode)
}

// UninstallSkill removes a skill from an agent.
func (sm *SkillManager) UninstallSkill(skillName string, agentType models.AgentType) error {
	skillsDir := sm.Resolver.AgentSkillsDir(agentType)
	targetDir := filepath.Join(skillsDir, skillName)

	if err := sm.Installer.Uninstall(targetDir); err != nil {
		return err
	}

	return sm.LockFile.RemoveEntry(string(agentType), skillName)
}

// ScanAgentSkills scans skills for a specific agent.
func (sm *SkillManager) ScanAgentSkills(agentType models.AgentType) []models.Skill {
	return sm.Scanner.ScanAgentSkills(agentType)
}

// BrowseAgentFiles returns a file tree of the agent's skills directory.
func (sm *SkillManager) BrowseAgentFiles(agentType models.AgentType) []models.FileNode {
	skillsDir := sm.Resolver.AgentSkillsDir(agentType)
	return buildFileTree(skillsDir, 3) // depth limit 3
}

// ReadFile reads a file and returns its contents as string.
func (sm *SkillManager) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read file %s: %w", path, err)
	}
	return string(data), nil
}

// SaveFile writes content to a file.
func (sm *SkillManager) SaveFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create directory %s: %w", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("cannot write file %s: %w", path, err)
	}
	return nil
}

func buildFileTree(dir string, maxDepth int) []models.FileNode {
	if maxDepth <= 0 {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var nodes []models.FileNode
	for _, entry := range entries {
		if entry.Name() == ".git" || entry.Name() == "node_modules" {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		node := models.FileNode{
			Name:  entry.Name(),
			Path:  fullPath,
			IsDir: entry.IsDir(),
		}

		if !entry.IsDir() {
			info, err := entry.Info()
			if err == nil {
				node.Size = info.Size()
			}
		}

		if entry.IsDir() {
			node.Children = buildFileTree(fullPath, maxDepth-1)
		}

		nodes = append(nodes, node)
	}

	return nodes
}
