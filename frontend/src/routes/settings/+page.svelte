<script lang="ts">
  let defaultMode = $state("copy");
</script>

<div class="p-6 max-w-2xl">
  <h2 class="text-2xl font-bold mb-6">设置</h2>

  <div class="space-y-6">
    <!-- Default install mode -->
    <div>
      <label class="text-sm font-medium" for="install-mode">默认安装模式</label>
      <p class="text-xs text-muted-foreground mb-2">
        Windows 推荐使用"复制"模式（symlink 需要管理员权限或开发者模式）
      </p>
      <select
        id="install-mode"
        bind:value={defaultMode}
        class="w-full px-3 py-2 border rounded-md bg-background"
      >
        <option value="copy">复制 (推荐)</option>
        <option value="symlink">符号链接</option>
      </select>
    </div>

    <!-- Info -->
    <div class="p-4 bg-muted rounded-lg text-sm">
      <h3 class="font-medium mb-2">关于符号链接</h3>
      <p class="text-muted-foreground">
        Windows 上创建符号链接需要以下条件之一：
      </p>
      <ul class="list-disc list-inside text-muted-foreground mt-1 space-y-1">
        <li>以管理员身份运行</li>
        <li>启用开发者模式（设置 → 更新和安全 → 开发者选项）</li>
      </ul>
      <p class="text-muted-foreground mt-2">
        如果不确定，请使用"复制"模式。
      </p>
    </div>
  </div>
</div>
