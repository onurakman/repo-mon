<template>
  <div
    class="flex items-center justify-between p-3 rounded-lg transition-all duration-300"
    :style="{
      backgroundColor: 'var(--color-surface)',
      borderLeft: '3px solid ' + borderColor,
    }"
  >
    <!-- Left: info -->
    <div class="flex items-center gap-4 min-w-0">
      <div class="min-w-0">
        <h3 class="font-bold text-sm truncate" :style="{ color: 'var(--color-text)' }">{{ repo.name }}</h3>
        <div class="flex items-center gap-3 mt-0.5">
          <span class="flex items-center gap-1">
            <Icon icon="codicon:git-branch" class="text-[10px]" :style="{ color: 'var(--color-primary)' }" />
            <span class="text-xs" :style="{ color: 'var(--color-text-secondary)' }">{{ status?.currentBranch ?? '...' }}</span>
          </span>
          <span class="flex items-center gap-1">
            <Icon icon="codicon:history" class="text-[10px]" :style="{ color: 'var(--color-text-secondary)' }" />
            <span class="text-[10px]" :style="{ color: 'var(--color-text-secondary)' }">{{ lastCheckedText }}</span>
          </span>
        </div>
      </div>
      <div v-if="repo.tags?.length" class="flex gap-1 shrink-0">
        <span
          v-for="tag in repo.tags"
          :key="tag.ID"
          class="px-1.5 py-0.5 rounded text-[10px] text-white"
          :style="{ backgroundColor: tag.color }"
        >
          {{ tag.name }}
        </span>
      </div>
    </div>

    <!-- Right: badges + refresh -->
    <div class="flex items-center gap-2 shrink-0">
      <div class="flex flex-wrap gap-1">
        <template v-if="status">
          <StatusBadge v-if="isClean && !status.checkingRemote" type="clean" />
          <StatusBadge v-if="status.hasConflicts" type="conflict" />
          <StatusBadge v-if="status.modifiedFiles > 0" type="modified" :count="status.modifiedFiles" />
          <StatusBadge v-if="status.stagedFiles > 0" type="staged" :count="status.stagedFiles" />
          <StatusBadge v-if="status.unpushedCommits > 0" type="ahead" :count="status.unpushedCommits" />
          <StatusBadge v-if="status.unpulledCommits > 0" type="behind" :count="status.unpulledCommits" />
          <StatusBadge v-if="status.checkingRemote" type="checking" />
          <StatusBadge v-else-if="!status.remoteAccessible && status.remotes?.length > 0" type="unreachable" />
        </template>
      </div>
      <button
        class="p-1.5 rounded cursor-pointer hover:opacity-80"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="refresh"
        title="Refresh"
      >
        <Icon
          icon="codicon:refresh"
          class="text-sm transition-transform"
          :class="{ 'refresh-spinning': refreshing }"
        />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Icon } from '@iconify/vue'
import StatusBadge from './StatusBadge.vue'
import type { Repository, RepoStatus } from '../stores/repoStore'
import { useTimeAgo } from '@vueuse/core'

const props = defineProps<{
  repo: Repository
  status: RepoStatus | undefined
}>()

const emit = defineEmits<{ refresh: [] }>()

const refreshing = ref(false)

function refresh() {
  refreshing.value = true
  emit('refresh')
  setTimeout(() => { refreshing.value = false }, 600)
}

const isClean = computed(() => {
  if (!props.status) return false
  return (
    props.status.uncommittedChanges === 0 &&
    !props.status.hasConflicts &&
    props.status.unpushedCommits === 0 &&
    props.status.unpulledCommits === 0
  )
})

const borderColor = computed(() => {
  if (!props.status) return 'var(--color-border)'
  if (props.status.hasConflicts) return 'var(--color-danger)'
  if (!props.status.remoteAccessible && props.status.remotes?.length > 0) return 'var(--color-muted)'
  if (props.status.uncommittedChanges > 0 || props.status.unpushedCommits > 0) return 'var(--color-warning)'
  return 'var(--color-success)'
})

const lastCheckedText = computed(() => {
  if (!props.status?.lastChecked) return 'Not checked yet'
  return useTimeAgo(new Date(props.status.lastChecked)).value
})
</script>
