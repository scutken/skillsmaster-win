export interface SkillDTO {
  name: string;
  description: string;
  version: string;
  author: string;
  tags: string[];
  platforms: string[];
  path: string;
  agentType: string;
  isInstalled: boolean;
  body?: string;
}

export interface AgentConfig {
  type: string;
  displayName: string;
  configDir: string;
  skillsDir: string;
  brandColor: string;
}

export interface FileNode {
  name: string;
  path: string;
  isDir: boolean;
  size: number;
  children?: FileNode[];
}

export type InstallMode = "copy" | "symlink";
