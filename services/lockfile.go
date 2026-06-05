package services

import (
	"fmt"
	"os"
	"path/filepath"

	"skillsmaster-win/models"

	"gopkg.in/yaml.v3"
)

type LockFileManager struct {
	Resolver models.PathResolver
}

func NewLockFileManager(resolver models.PathResolver) *LockFileManager {
	return &LockFileManager{Resolver: resolver}
}

func (lf *LockFileManager) lockFilePath() string {
	return lf.Resolver.LockFilePath()
}

func (lf *LockFileManager) loadOrCreate() (*models.LockFile, error) {
	path := lf.lockFilePath()

	lockFile := &models.LockFile{
		Version: "1",
		Entries: make(map[string]models.LockEntry),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return lockFile, nil
		}
		return nil, fmt.Errorf("cannot read lockfile: %w", err)
	}

	if err := yaml.Unmarshal(data, lockFile); err != nil {
		return nil, fmt.Errorf("cannot parse lockfile: %w", err)
	}

	if lockFile.Entries == nil {
		lockFile.Entries = make(map[string]models.LockEntry)
	}

	return lockFile, nil
}

func (lf *LockFileManager) save(lockFile *models.LockFile) error {
	path := lf.lockFilePath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cannot create lockfile directory: %w", err)
	}

	data, err := yaml.Marshal(lockFile)
	if err != nil {
		return fmt.Errorf("cannot marshal lockfile: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("cannot write lockfile: %w", err)
	}

	return nil
}

// RecordInstall adds or updates an entry in the lockfile.
func (lf *LockFileManager) RecordInstall(entry models.LockEntry) error {
	lockFile, err := lf.loadOrCreate()
	if err != nil {
		return err
	}

	key := entry.AgentType + "/" + entry.SkillName
	lockFile.Entries[key] = entry

	return lf.save(lockFile)
}

// RemoveEntry removes an entry from the lockfile.
func (lf *LockFileManager) RemoveEntry(agentType, skillName string) error {
	lockFile, err := lf.loadOrCreate()
	if err != nil {
		return err
	}

	key := agentType + "/" + skillName
	delete(lockFile.Entries, key)

	return lf.save(lockFile)
}

// GetEntries returns all lockfile entries.
func (lf *LockFileManager) GetEntries() map[string]models.LockEntry {
	lockFile, err := lf.loadOrCreate()
	if err != nil {
		return make(map[string]models.LockEntry)
	}
	return lockFile.Entries
}
