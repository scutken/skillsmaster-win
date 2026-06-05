/**
 * Wails Go backend bindings.
 * In production, Wails auto-generates these on window.go.services.SkillManager.*
 * This wrapper provides typed async functions for Svelte components.
 */

import type { AgentConfig, SkillDTO, FileNode, InstallMode } from './types';

// Access the Wails-bound Go backend
function getBackend() {
  return (window as any).go?.services?.SkillManager;
}

export async function getAgents(): Promise<AgentConfig[]> {
  const backend = getBackend();
  if (!backend) return mockAgents();
  return backend.GetAgents();
}

export async function getAllSkills(): Promise<Record<string, SkillDTO[]>> {
  const backend = getBackend();
  if (!backend) return {};
  return backend.GetAllSkills();
}

export async function installSkill(
  sourceDir: string,
  agentType: string,
  mode: InstallMode
): Promise<void> {
  const backend = getBackend();
  if (!backend) throw new Error("Backend not available");
  return backend.InstallSkill(sourceDir, agentType, mode);
}

export async function uninstallSkill(
  skillName: string,
  agentType: string
): Promise<void> {
  const backend = getBackend();
  if (!backend) throw new Error("Backend not available");
  return backend.UninstallSkill(skillName, agentType);
}

export async function browseAgentFiles(agentType: string): Promise<FileNode[]> {
  const backend = getBackend();
  if (!backend) return [];
  return backend.BrowseAgentFiles(agentType);
}

export async function readFile(path: string): Promise<string> {
  const backend = getBackend();
  if (!backend) throw new Error("Backend not available");
  return backend.ReadFile(path);
}

export async function saveFile(path: string, content: string): Promise<void> {
  const backend = getBackend();
  if (!backend) throw new Error("Backend not available");
  return backend.SaveFile(path, content);
}

// Mock data for development without Wails
function mockAgents(): AgentConfig[] {
  return [
    { type: "claude-code", displayName: "Claude Code", configDir: "", skillsDir: "", brandColor: "#D97706" },
    { type: "codex", displayName: "Codex", configDir: "", skillsDir: "", brandColor: "#10A37F" },
    { type: "hermes", displayName: "Hermes", configDir: "", skillsDir: "", brandColor: "#E11D48" },
    { type: "cursor", displayName: "Cursor", configDir: "", skillsDir: "", brandColor: "#00D4AA" },
  ];
}
