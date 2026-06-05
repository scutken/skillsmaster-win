<script lang="ts">
  import { page } from "$app/stores";

  const navItems = [
    { label: "已安装 Skills", path: "/skills", icon: "📦" },
    { label: "市场", path: "/marketplace", icon: "🏪" },
    { label: "Agents", path: "/agents", icon: "🤖" },
    { label: "设置", path: "/settings", icon: "⚙️" },
  ];
</script>

<nav class="flex-1 p-3 space-y-1">
  {#each navItems as item}
    <a
      href={item.path}
      class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors
        {$page.url.pathname.startsWith(item.path)
          ? 'bg-accent text-accent-foreground'
          : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'}"
    >
      <span class="text-lg">{item.icon}</span>
      <span>{item.label}</span>
    </a>
  {/each}
</nav>
