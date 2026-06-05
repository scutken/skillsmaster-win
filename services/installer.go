package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"skillsmaster-win/models"
)

type InstallMode string

const (
	InstallModeSymlink InstallMode = "symlink"
	InstallModeCopy    InstallMode = "copy"
)

type SkillInstaller struct {
	Resolver   models.PathResolver
	LockFile   *LockFileManager
}

func NewSkillInstaller(resolver models.PathResolver, lockFile *LockFileManager) *SkillInstaller {
	return &SkillInstaller{
		Resolver: resolver,
		LockFile: lockFile,
	}
}

// Install copies or symlinks a skill directory to the agent's skills directory.
func (si *SkillInstaller) Install(sourceDir string, agentType models.AgentType, mode InstallMode) error {
	if mode == "" {
		mode = InstallModeCopy
	}

	// Read skill manifest from source
	skillPath := filepath.Join(sourceDir, "SKILL.md")
	data, err := os.ReadFile(skillPath)
	if err != nil {
		return fmt.Errorf("cannot read SKILL.md from %s: %w", sourceDir, err)
	}

	skill, err := ParseSkillFile(string(data), skillPath)
	if err != nil {
		return fmt.Errorf("cannot parse SKILL.md: %w", err)
	}

	// Determine target directory
	skillsDir := si.Resolver.AgentSkillsDir(agentType)
	targetDir := filepath.Join(skillsDir, skill.Name)

	// Create skills directory if needed
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("cannot create skills dir: %w", err)
	}

	// Remove existing installation if present
	if _, err := os.Stat(targetDir); err == nil {
		os.RemoveAll(targetDir)
	}

	switch mode {
	case InstallModeSymlink:
		absSource, err := filepath.Abs(sourceDir)
		if err != nil {
			return fmt.Errorf("cannot resolve source path: %w", err)
		}
		if err := os.Symlink(absSource, targetDir); err != nil {
			return fmt.Errorf("symlink failed: %w", err)
		}
	case InstallModeCopy:
		if err := copyDir(sourceDir, targetDir); err != nil {
			return fmt.Errorf("copy failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown install mode: %s", mode)
	}

	// Record in lockfile
	if si.LockFile != nil {
		entry := models.LockEntry{
			AgentType:   string(agentType),
			SkillName:   skill.Name,
			SourcePath:  sourceDir,
			TargetPath:  targetDir,
			Mode:        string(mode),
			InstalledAt: time.Now().UTC().Format(time.RFC3339),
		}
		if err := si.LockFile.RecordInstall(entry); err != nil {
			return fmt.Errorf("failed to record install: %w", err)
		}
	}

	return nil
}

// Uninstall removes an installed skill.
func (si *SkillInstaller) Uninstall(targetPath string) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return nil // Already gone
	}
	if err := os.RemoveAll(targetPath); err != nil {
		return fmt.Errorf("failed to remove %s: %w", targetPath, err)
	}
	return nil
}

// CheckInstallStatus checks if a skill is installed for an agent.
// Returns: "installed", "not_installed", or "conflict"
func (si *SkillInstaller) CheckInstallStatus(sourceDir string, agentType models.AgentType) string {
	skillPath := filepath.Join(sourceDir, "SKILL.md")
	data, err := os.ReadFile(skillPath)
	if err != nil {
		return "not_installed"
	}

	skill, err := ParseSkillFile(string(data), skillPath)
	if err != nil {
		return "not_installed"
	}

	skillsDir := si.Resolver.AgentSkillsDir(agentType)
	targetDir := filepath.Join(skillsDir, skill.Name)

	if _, err := os.Stat(targetDir); err != nil {
		return "not_installed"
	}
	return "installed"
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	srcInfo, _ := os.Stat(src)
	if srcInfo != nil {
		os.Chmod(dst, srcInfo.Mode())
	}
	return nil
}
