<script lang="ts">
  import "../app.css";
  import Sidebar from "$lib/components/Sidebar.svelte";

  let { children } = $props();
</script>

<svelte:head>
  <title>SkillsMaster</title>
</svelte:head>

<div class="flex h-screen bg-background text-foreground">
  <!-- Sidebar -->
  <aside class="w-56 border-r bg-muted/30 flex flex-col shrink-0">
    <div class="p-4 border-b">
      <h1 class="font-bold text-lg tracking-tight">⚡ SkillsMaster</h1>
      <p class="text-xs text-muted-foreground mt-0.5">AI Agent Skills 管理</p>
    </div>
    <Sidebar />
  </aside>

  <!-- Main content -->
  <main class="flex-1 overflow-auto">
    {@render children()}
  </main>
</div>
