# SkillsMaster-Windows

基于 [zhls-ayl/SkillsMaster](https://github.com/zhls-ayl/SkillsMaster) 的 Windows 版本，用 Wails v2 (Go + Svelte) 构建。

管理多个 AI 编程 Agent 的 Skills：扫描、安装、编辑、更新。

## 支持的 Agents

Claude Code · Codex · Gemini CLI · GitHub Copilot · OpenCode · Cursor · Hermes · OpenClaw · Trae · Antigravity · Kiro CLI · CodeBuddy

## 技术栈

- **后端:** Go + Wails v2
- **前端:** Svelte 5 + TypeScript + Tailwind CSS
- **构建:** 单个 exe，~10-15MB

## 开发

```bash
# 前置：Go 1.21+, Node.js 20+
go install github.com/wailsapp/wails/v2/cmd/wails@latest
cd frontend && npm install && cd ..
wails dev
```

## 构建

```bash
wails build -platform windows/amd64 -o SkillsMaster.exe
```

## 测试

```bash
go test ./... -v
```

## License

MIT
