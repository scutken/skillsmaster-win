<script lang="ts">
  import { onMount } from "svelte";
  import SkillList from "$lib/components/SkillList.svelte";
  import SkillDetail from "$lib/components/SkillDetail.svelte";
  import { getAllSkills } from "$lib/wails-bindings";
  import type { SkillDTO } from "$lib/types";

  let skills: Record<string, SkillDTO[]> = $state({});
  let selectedSkill: SkillDTO | null = $state(null);
  let loading = $state(true);

  onMount(async () => {
    await loadSkills();
  });

  async function loadSkills() {
    loading = true;
    try {
      skills = await getAllSkills();
    } catch (e) {
      console.error("Failed to load skills:", e);
    } finally {
      loading = false;
    }
  }

  function handleSelect(e: CustomEvent<SkillDTO>) {
    selectedSkill = e.detail;
  }
</script>

<div class="flex h-full">
  <!-- Skill list -->
  <div class="w-80 border-r shrink-0">
    {#if loading}
      <div class="flex items-center justify-center h-32 text-muted-foreground">
        加载中...
      </div>
    {:else}
      <SkillList
        {skills}
        selectedName={selectedSkill?.name ?? ""}
        on:select={handleSelect}
      />
    {/if}
  </div>

  <!-- Skill detail -->
  <div class="flex-1 overflow-auto">
    <SkillDetail
      skill={selectedSkill}
      onUninstall={() => {
        selectedSkill = null;
        loadSkills();
      }}
    />
  </div>
</div>
