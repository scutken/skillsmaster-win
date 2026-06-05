<script lang="ts">
  import type { FileNode, AgentConfig } from "$lib/types";
  import { browseAgentFiles, readFile } from "$lib/wails-bindings";

  interface Props {
    agent: AgentConfig;
  }

  let { agent }: Props = $props();
  let files = $state<FileNode[]>([]);
  let selectedFile = $state<string | null>(null);
  let fileContent = $state<string>("");
  let loading = $state(false);

  $effect(() => {
    loadFiles(agent.type);
  });

  async function loadFiles(agentType: string) {
    loading = true;
    try {
      files = await browseAgentFiles(agentType);
    } catch (e) {
      console.error("Failed to load files:", e);
      files = [];
    } finally {
      loading = false;
    }
  }

  async function openFile(path: string) {
    selectedFile = path;
    try {
      fileContent = await readFile(path);
    } catch (e) {
      fileContent = `Error reading file: ${e}`;
    }
  }
</script>

<div class="flex h-full">
  <!-- File tree -->
  <div class="w-64 border-r overflow-auto p-2 text-sm">
    <div class="font-medium px-2 py-1 text-xs text-muted-foreground uppercase tracking-wider mb-1">
      {agent.displayName}
    </div>
    {#if loading}
      <div class="text-muted-foreground px-2">加载中...</div>
    {:else if files.length === 0}
      <div class="text-muted-foreground px-2">无文件</div>
    {:else}
      {#each files as node}
        {@render fileNode(node, 0)}
      {/each}
    {/if}
  </div>

  <!-- File preview -->
  <div class="flex-1 overflow-auto">
    {#if selectedFile}
      <div class="p-4">
        <div class="text-xs text-muted-foreground mb-2 font-mono">{selectedFile}</div>
        <pre class="text-sm bg-muted p-4 rounded-md overflow-auto max-h-[600px] whitespace-pre-wrap">{fileContent}</pre>
      </div>
    {:else}
      <div class="flex items-center justify-center h-full text-muted-foreground">
        选择文件查看内容
      </div>
    {/if}
  </div>
</div>

{#snippet fileNode(node: FileNode, depth: number)}
  <div style="padding-left: {depth * 16}px">
    {#if node.isDir}
      <div class="flex items-center gap-1 px-2 py-1 rounded hover:bg-accent/50 cursor-default">
        <span class="text-xs">📁</span>
        <span class="font-medium">{node.name}</span>
      </div>
      {#if node.children}
        {#each node.children as child}
          {@render fileNode(child, depth + 1)}
        {/each}
      {/if}
    {:else}
      <button
        class="flex items-center gap-1 px-2 py-1 rounded w-full text-left hover:bg-accent/50
          {selectedFile === node.path ? 'bg-accent' : ''}"
        onclick={() => openFile(node.path)}
      >
        <span class="text-xs">📄</span>
        <span class="truncate">{node.name}</span>
        <span class="text-xs text-muted-foreground ml-auto">{formatSize(node.size)}</span>
      </button>
    {/if}
  </div>
{/snippet}

{#snippet formatSize(bytes: number)}
  {#if bytes < 1024}
    {bytes}B
  {:else if bytes < 1024 * 1024}
    {(bytes / 1024).toFixed(1)}K
  {:else}
    {(bytes / 1024 / 1024).toFixed(1)}M
  {/if}
{/snippet}
