<script lang="ts">
  import { onMount } from "svelte";
  import AgentFiles from "$lib/components/AgentFiles.svelte";
  import { getAgents } from "$lib/wails-bindings";
  import type { AgentConfig } from "$lib/types";

  let agents: AgentConfig[] = $state([]);
  let selectedAgent: AgentConfig | null = $state(null);
  let loading = $state(true);

  onMount(async () => {
    try {
      agents = await getAgents();
      if (agents.length > 0) selectedAgent = agents[0];
    } finally {
      loading = false;
    }
  });
</script>

<div class="flex h-full">
  <!-- Agent list -->
  <div class="w-56 border-r overflow-auto p-2">
    <div class="text-xs font-medium text-muted-foreground uppercase tracking-wider px-3 py-2">
      已安装的 Agent
    </div>
    {#if loading}
      <div class="text-muted-foreground px-3 text-sm">加载中...</div>
    {:else}
      {#each agents as agent}
        <button
          class="flex items-center gap-2 w-full text-left px-3 py-2 rounded-md text-sm transition-colors
            {selectedAgent?.type === agent.type
              ? 'bg-accent text-accent-foreground'
              : 'hover:bg-accent/50'}"
          onclick={() => (selectedAgent = agent)}
        >
          <span
            class="w-3 h-3 rounded-full shrink-0"
            style="background-color: {agent.brandColor}"
          ></span>
          <span class="truncate">{agent.displayName}</span>
        </button>
      {/each}
    {/if}
  </div>

  <!-- Agent files -->
  <div class="flex-1 overflow-auto">
    {#if selectedAgent}
      {#key selectedAgent.type}
        <AgentFiles agent={selectedAgent} />
      {/key}
    {:else}
      <div class="flex items-center justify-center h-full text-muted-foreground">
        选择一个 Agent 浏览文件
      </div>
    {/if}
  </div>
</div>
