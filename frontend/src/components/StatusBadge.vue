<template>
  <span
    class="inline-flex items-center gap-1 px-2 py-1 rounded text-xs font-medium transition-all duration-200"
    :class="{ 'badge-pulse': type === 'checking' }"
    :style="{ backgroundColor: bgColor, color: '#fff' }"
    v-tooltip.top="tooltip"
  >
    <Icon :icon="iconName" width="14" height="14" />
    <span v-if="displayText">{{ displayText }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  type: 'clean' | 'modified' | 'staged' | 'untracked' | 'conflict' | 'ahead' | 'behind' | 'stash' | 'unreachable' | 'checking'
  count?: number
}>()

const config = computed(() => {
  const map: Record<string, { bg: string; icon: string; text: string; tip: string }> = {
    clean:       { bg: 'var(--color-success)', icon: 'codicon:check',                text: '',                    tip: 'Working tree clean' },
    modified:    { bg: 'var(--color-warning)', icon: 'codicon:diff-modified',         text: `${props.count ?? 0}`, tip: `${props.count ?? 0} modified files` },
    staged:      { bg: 'var(--color-info)',    icon: 'codicon:diff-added',            text: `${props.count ?? 0}`, tip: `${props.count ?? 0} staged files` },
    untracked:   { bg: 'var(--color-muted)',   icon: 'codicon:diff-ignored',          text: `${props.count ?? 0}`, tip: `${props.count ?? 0} untracked files` },
    conflict:    { bg: 'var(--color-danger)',   icon: 'codicon:git-merge',             text: '',                    tip: 'Merge conflict' },
    ahead:       { bg: 'var(--color-info)',    icon: 'codicon:cloud-upload',          text: `${props.count ?? 0}`, tip: `${props.count ?? 0} commits ahead (unpushed)` },
    behind:      { bg: 'var(--color-warning)', icon: 'codicon:cloud-download',        text: `${props.count ?? 0}`, tip: `${props.count ?? 0} commits behind (unpulled)` },
    stash:       { bg: 'var(--color-muted)',   icon: 'codicon:archive',               text: `${props.count ?? 0}`, tip: `${props.count ?? 0} stash entries` },
    unreachable: { bg: 'var(--color-danger)',  icon: 'codicon:debug-disconnect',      text: '',                    tip: 'Remote unreachable' },
    checking:    { bg: 'var(--color-muted)',   icon: 'codicon:loading~spin',          text: '',                    tip: 'Checking remote...' },
  }
  return map[props.type] ?? map.clean
})

const bgColor = computed(() => config.value.bg)
const iconName = computed(() => config.value.icon)
const displayText = computed(() => config.value.text)
const tooltip = computed(() => config.value.tip)
</script>

<style scoped>
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
:deep([data-icon="codicon:loading~spin"]) {
  animation: spin 1s linear infinite;
}
</style>
