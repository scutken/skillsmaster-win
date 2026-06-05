package models

type Skill struct {
	Name       string        `json:"name"`
	Path       string        `json:"path"`
	AgentType  AgentType     `json:"agentType"`
	Manifest   SkillManifest `json:"manifest"`
	Body       string        `json:"body"`
	IsLinked   bool          `json:"isLinked"`
	IsInherited bool         `json:"isInherited"`
	Source     string        `json:"source"`
}

type SkillManifest struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Version     string   `yaml:"version" json:"version"`
	Author      string   `yaml:"author" json:"author"`
	Tags        []string `yaml:"tags" json:"tags"`
	Platforms   []string `yaml:"platforms" json:"platforms"`
}

type LockEntry struct {
	AgentType   string `json:"agentType"`
	SkillName   string `json:"skillName"`
	SourcePath  string `json:"sourcePath"`
	TargetPath  string `json:"targetPath"`
	Mode        string `json:"mode"`
	InstalledAt string `json:"installedAt"`
}

type LockFile struct {
	Version string               `yaml:"version" json:"version"`
	Entries map[string]LockEntry `yaml:"entries" json:"entries"`
}

type SkillDTO struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Platforms   []string `json:"platforms"`
	Path        string   `json:"path"`
	AgentType   string   `json:"agentType"`
	IsInstalled bool     `json:"isInstalled"`
	Body        string   `json:"body,omitempty"`
}

type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Size     int64      `json:"size"`
	Children []FileNode `json:"children,omitempty"`
}
