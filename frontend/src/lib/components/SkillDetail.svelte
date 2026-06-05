<script lang="ts">
  import type { SkillDTO } from "$lib/types";
  import { uninstallSkill } from "$lib/wails-bindings";

  interface Props {
    skill: SkillDTO | null;
    onUninstall?: () => void;
  }

  let { skill, onUninstall }: Props = $props();
  let uninstalling = $state(false);

  async function handleUninstall() {
    if (!skill) return;
    uninstalling = true;
    try {
      await uninstallSkill(skill.name, skill.agentType);
      onUninstall?.();
    } catch (e) {
      alert(`卸载失败: ${e}`);
    } finally {
      uninstalling = false;
    }
  }
</script>

{#if skill}
  <div class="p-6 space-y-4 max-w-3xl">
    <!-- Header -->
    <div class="flex items-start justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold">{skill.name}</h2>
        {#if skill.description}
          <p class="text-muted-foreground mt-1">{skill.description}</p>
        {/if}
      </div>
      <div class="flex gap-2 shrink-0">
        <button
          class="px-3 py-1.5 text-sm border rounded-md hover:bg-accent transition-colors"
        >
          编辑
        </button>
        <button
          class="px-3 py-1.5 text-sm bg-destructive text-destructive-foreground rounded-md hover:opacity-90 transition-opacity disabled:opacity-50"
          disabled={uninstalling}
          onclick={handleUninstall}
        >
          {uninstalling ? "卸载中..." : "卸载"}
        </button>
      </div>
    </div>

    <!-- Badges -->
    <div class="flex gap-2 flex-wrap">
      <span class="px-2.5 py-1 text-xs font-medium rounded-full bg-primary text-primary-foreground">
        {skill.agentType}
      </span>
      {#if skill.version}
        <span class="px-2.5 py-1 text-xs font-medium rounded-full bg-secondary text-secondary-foreground">
          v{skill.version}
        </span>
      {/if}
      {#if skill.author}
        <span class="px-2.5 py-1 text-xs font-medium rounded-full border">
          {skill.author}
        </span>
      {/if}
      {#if skill.tags?.length}
        {#each skill.tags as tag}
          <span class="px-2.5 py-1 text-xs rounded-full bg-muted text-muted-foreground">
            {tag}
          </span>
        {/each}
      {/if}
    </div>

    <hr class="border-border" />

    <!-- Metadata table -->
    <div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
      <span class="text-muted-foreground">路径</span>
      <code class="text-xs bg-muted px-2 py-1 rounded">{skill.path}</code>
      {#if skill.platforms?.length}
        <span class="text-muted-foreground">平台</span>
        <span>{skill.platforms.join(", ")}</span>
      {/if}
      <span class="text-muted-foreground">安装状态</span>
      <span>{skill.isInstalled ? "✅ 已安装" : "❌ 未安装"}</span>
    </div>

    <hr class="border-border" />

    <!-- Body preview -->
    {#if skill.body}
      <div>
        <h3 class="text-sm font-medium mb-2">SKILL.md 内容</h3>
        <pre class="text-xs bg-muted p-4 rounded-md overflow-auto max-h-96 whitespace-pre-wrap">{skill.body}</pre>
      </div>
    {/if}
  </div>
{:else}
  <div class="flex items-center justify-center h-full text-muted-foreground">
    选择一个 Skill 查看详情
  </div>
{/if}
