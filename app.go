package main

import (
	"skillsmaster-win/models"
	"skillsmaster-win/services"
)

// App is the main application struct bound to the Wails frontend.
type App struct {
	Manager *services.SkillManager
}

func NewApp(resolver models.PathResolver) *App {
	return &App{
		Manager: services.NewSkillManager(resolver),
	}
}

// GetAgents returns all configured AI agent definitions.
func (a *App) GetAgents() []models.AgentConfig {
	return a.Manager.GetAgents()
}

// GetAllSkills returns skills across all agents as DTOs.
func (a *App) GetAllSkills() map[string][]models.SkillDTO {
	return a.Manager.GetAllSkills()
}

// ScanAgentSkills scans skills for a given agent type.
func (a *App) ScanAgentSkills(agentType string) []models.Skill {
	return a.Manager.ScanAgentSkills(models.AgentType(agentType))
}

// InstallSkill installs a skill to an agent.
func (a *App) InstallSkill(sourceDir, agentType, mode string) error {
	return a.Manager.InstallSkill(sourceDir, models.AgentType(agentType), services.InstallMode(mode))
}

// UninstallSkill removes a skill from an agent.
func (a *App) UninstallSkill(skillName, agentType string) error {
	return a.Manager.UninstallSkill(skillName, models.AgentType(agentType))
}

// BrowseAgentFiles returns a file tree for an agent's skills directory.
func (a *App) BrowseAgentFiles(agentType string) []models.FileNode {
	return a.Manager.BrowseAgentFiles(models.AgentType(agentType))
}

// ReadFile reads a file and returns its contents.
func (a *App) ReadFile(path string) (string, error) {
	return a.Manager.ReadFile(path)
}

// SaveFile writes content to a file.
func (a *App) SaveFile(path, content string) error {
	return a.Manager.SaveFile(path, content)
}
