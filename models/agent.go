package models

import (
	"os"
	"path/filepath"
)

type AgentType string

const (
	AgentClaudeCode      AgentType = "claude-code"
	AgentCodex           AgentType = "codex"
	AgentGeminiCLI       AgentType = "gemini-cli"
	AgentGitHubCopilot   AgentType = "github-copilot"
	AgentOpencode        AgentType = "opencode"
	AgentCursor          AgentType = "cursor"
	AgentHermes          AgentType = "hermes"
	AgentOpenClaw        AgentType = "openclaw"
	AgentTrae            AgentType = "trae"
	AgentAntiGravity     AgentType = "antigravity"
	AgentKiroCLI         AgentType = "kiro-cli"
	AgentCodeBuddy       AgentType = "codebuddy"
)

var AllAgentTypes = []AgentType{
	AgentClaudeCode,
	AgentCodex,
	AgentGeminiCLI,
	AgentGitHubCopilot,
	AgentOpencode,
	AgentCursor,
	AgentHermes,
	AgentOpenClaw,
	AgentTrae,
	AgentAntiGravity,
	AgentKiroCLI,
	AgentCodeBuddy,
}

type AgentConfig struct {
	Type        AgentType
	DisplayName string
	ConfigDir   string
	SkillsDir   string
	BrandColor  string
}

func GetAgentConfig(agentType AgentType) AgentConfig {
	home, _ := os.UserHomeDir()
	configs := agentConfigDefs(home)
	if cfg, ok := configs[agentType]; ok {
		return cfg
	}
	return AgentConfig{
		Type:        agentType,
		DisplayName: string(agentType),
		ConfigDir:   filepath.Join(home, ".config", string(agentType)),
		SkillsDir:   filepath.Join(home, ".config", string(agentType), "skills"),
		BrandColor:  "#666666",
	}
}

func AllAgentConfigs() []AgentConfig {
	home, _ := os.UserHomeDir()
	configs := agentConfigDefs(home)
	result := make([]AgentConfig, 0, len(AllAgentTypes))
	for _, t := range AllAgentTypes {
		if cfg, ok := configs[t]; ok {
			result = append(result, cfg)
		}
	}
	return result
}

func agentConfigDefs(home string) map[AgentType]AgentConfig {
	return map[AgentType]AgentConfig{
		AgentClaudeCode: {
			Type:        AgentClaudeCode,
			DisplayName: "Claude Code",
			ConfigDir:   filepath.Join(home, ".claude"),
			SkillsDir:   filepath.Join(home, ".claude", "skills"),
			BrandColor:  "#D97706",
		},
		AgentCodex: {
			Type:        AgentCodex,
			DisplayName: "Codex",
			ConfigDir:   filepath.Join(home, ".codex"),
			SkillsDir:   filepath.Join(home, ".codex", "skills"),
			BrandColor:  "#10A37F",
		},
		AgentGeminiCLI: {
			Type:        AgentGeminiCLI,
			DisplayName: "Gemini CLI",
			ConfigDir:   filepath.Join(home, ".gemini"),
			SkillsDir:   filepath.Join(home, ".gemini", "skills"),
			BrandColor:  "#4285F4",
		},
		AgentGitHubCopilot: {
			Type:        AgentGitHubCopilot,
			DisplayName: "GitHub Copilot",
			ConfigDir:   filepath.Join(home, ".github-copilot"),
			SkillsDir:   filepath.Join(home, ".github-copilot", "skills"),
			BrandColor:  "#238636",
		},
		AgentOpencode: {
			Type:        AgentOpencode,
			DisplayName: "OpenCode",
			ConfigDir:   filepath.Join(home, ".opencode"),
			SkillsDir:   filepath.Join(home, ".opencode", "skills"),
			BrandColor:  "#8B5CF6",
		},
		AgentCursor: {
			Type:        AgentCursor,
			DisplayName: "Cursor",
			ConfigDir:   filepath.Join(home, ".cursor"),
			SkillsDir:   filepath.Join(home, ".cursor", "skills"),
			BrandColor:  "#00D4AA",
		},
		AgentHermes: {
			Type:        AgentHermes,
			DisplayName: "Hermes",
			ConfigDir:   filepath.Join(home, ".hermes"),
			SkillsDir:   filepath.Join(home, ".hermes", "skills"),
			BrandColor:  "#E11D48",
		},
		AgentOpenClaw: {
			Type:        AgentOpenClaw,
			DisplayName: "OpenClaw",
			ConfigDir:   filepath.Join(home, ".openclaw"),
			SkillsDir:   filepath.Join(home, ".openclaw", "skills"),
			BrandColor:  "#F59E0B",
		},
		AgentTrae: {
			Type:        AgentTrae,
			DisplayName: "Trae",
			ConfigDir:   filepath.Join(home, ".trae"),
			SkillsDir:   filepath.Join(home, ".trae", "skills"),
			BrandColor:  "#3B82F6",
		},
		AgentAntiGravity: {
			Type:        AgentAntiGravity,
			DisplayName: "AntiGravity",
			ConfigDir:   filepath.Join(home, ".antigravity"),
			SkillsDir:   filepath.Join(home, ".antigravity", "skills"),
			BrandColor:  "#6366F1",
		},
		AgentKiroCLI: {
			Type:        AgentKiroCLI,
			DisplayName: "Kiro CLI",
			ConfigDir:   filepath.Join(home, ".kiro"),
			SkillsDir:   filepath.Join(home, ".kiro", "skills"),
			BrandColor:  "#EC4899",
		},
		AgentCodeBuddy: {
			Type:        AgentCodeBuddy,
			DisplayName: "CodeBuddy",
			ConfigDir:   filepath.Join(home, ".codebuddy"),
			SkillsDir:   filepath.Join(home, ".codebuddy", "skills"),
			BrandColor:  "#14B8A6",
		},
	}
}
