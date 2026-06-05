<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { SkillDTO } from "$lib/types";

  interface Props {
    skills: Record<string, SkillDTO[]>;
    selectedName?: string;
  }

  let { skills, selectedName = $bindable("") }: Props = $props();
  const dispatch = createEventDispatcher<{ select: SkillDTO }>();

  let search = $state("");

  let allSkills = $derived(
    Object.entries(skills).flatMap(([agent, list]) =>
      list.map((s) => ({ ...s, agentType: agent }))
    )
  );

  let filtered = $derived(
    allSkills.filter(
      (s) =>
        s.name.toLowerCase().includes(search.toLowerCase()) ||
        s.description?.toLowerCase().includes(search.toLowerCase())
    )
  );

  function select(skill: SkillDTO) {
    selectedName = skill.name;
    dispatch("select", skill);
  }
</script>

<div class="flex flex-col h-full">
  <!-- Search -->
  <div class="p-3 border-b">
    <input
      type="text"
      bind:value={search}
      placeholder="搜索 Skills..."
      class="w-full px-3 py-2 text-sm border rounded-md bg-background focus:outline-none focus:ring-2 focus:ring-ring"
    />
  </div>

  <!-- Skill list -->
  <div class="flex-1 overflow-auto p-2 space-y-1">
    {#if filtered.length === 0}
      <div class="text-center text-muted-foreground py-8 text-sm">
        {search ? "没有匹配的 Skill" : "暂无已安装的 Skills"}
      </div>
    {/if}

    {#each filtered as skill (skill.name + skill.agentType)}
      <button
        class="w-full text-left p-3 rounded-lg border transition-colors
          {selectedName === skill.name
            ? 'border-primary bg-accent'
            : 'border-transparent hover:bg-accent/50'}"
        onclick={() => select(skill)}
      >
        <div class="font-medium text-sm">{skill.name}</div>
        {#if skill.description}
          <div class="text-xs text-muted-foreground line-clamp-2 mt-1">
            {skill.description}
          </div>
        {/if}
        <div class="flex gap-1.5 mt-2 flex-wrap">
          <span
            class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-primary/10 text-primary"
          >
            {skill.agentType}
          </span>
          {#if skill.version}
            <span
              class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-secondary text-secondary-foreground"
            >
              v{skill.version}
            </span>
          {/if}
        </div>
      </button>
    {/each}
  </div>

  <!-- Count -->
  <div class="p-2 border-t text-xs text-muted-foreground text-center">
    {filtered.length} / {allSkills.length} Skills
  </div>
</div>
