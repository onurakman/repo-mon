<template>
  <div
    class="rounded-lg p-4 relative overflow-hidden transition-all duration-300"
    :style="{
      backgroundColor: 'var(--color-surface)',
      borderTop: '3px solid ' + borderColor,
      outline: selected ? '2px solid var(--color-primary)' : 'none',
      outlineOffset: '-2px',
    }"
  >
    <!-- Shine overlay while checking -->
    <div
      v-if="checking"
      class="absolute inset-0 pointer-events-none card-shine rounded-lg"
    />

    <!-- Select checkbox -->
    <div
      v-if="selectable"
      class="absolute top-2 left-2 cursor-pointer"
      @click.stop="$emit('toggle-select')"
    >
      <div
        class="w-4 h-4 rounded border flex items-center justify-center transition-colors"
        :style="{
          backgroundColor: selected ? 'var(--color-primary)' : 'transparent',
          borderColor: selected ? 'var(--color-primary)' : 'var(--color-border)',
        }"
      >
        <Icon v-if="selected" icon="codicon:check" width="10" height="10" style="color: #fff" />
      </div>
    </div>

    <!-- Header: name + actions -->
    <div class="flex items-start justify-between mb-3">
      <div class="min-w-0 flex-1" :class="{ 'pl-5': selectable }">
        <h3 class="font-bold text-sm truncate" :style="{ color: 'var(--color-text)' }">{{ repo.name }}</h3>
        <div class="flex items-center gap-1.5 mt-1">
          <Icon icon="codicon:git-branch" width="14" height="14" :style="{ color: 'var(--color-primary)' }" />
          <span class="text-xs" :style="{ color: 'var(--color-text-secondary)' }">
            {{ status?.currentBranch ?? repo.name }}
          </span>
        </div>
      </div>
      <div class="flex items-center gap-0.5 shrink-0">
        <TagDropdown :repo-id="repo.ID" :current-tag-ids="tagIds" />
        <button
          class="p-1.5 rounded cursor-pointer hover:opacity-80"
          :style="{ color: 'var(--color-text-secondary)' }"
          @click="showRemoveConfirm = true"
          title="Remove repository"
        >
          <Icon icon="codicon:trash" width="14" height="14" />
        </button>
        <button
          class="p-1.5 rounded cursor-pointer hover:opacity-80"
          :style="{ color: 'var(--color-text-secondary)' }"
          @click="$emit('refresh')"
          title="Refresh"
        >
          <Icon
            icon="codicon:refresh"
            width="16" height="16"
          />
        </button>
      </div>
    </div>

    <!-- Tags -->
    <div v-if="repo.tags?.length" class="flex flex-wrap gap-1 mb-3">
      <span
        v-for="tag in repo.tags"
        :key="tag.ID"
        class="px-1.5 py-0.5 rounded text-[10px] text-white"
        :style="{ backgroundColor: tag.color }"
      >
        {{ tag.name }}
      </span>
    </div>

    <!-- Status badges -->
    <div class="flex flex-wrap gap-1">
      <template v-if="status">
        <StatusBadge v-if="isClean && !status.checkingRemote" type="clean" />
        <StatusBadge v-if="status.hasConflicts" type="conflict" />
        <StatusBadge v-if="status.modifiedFiles > 0" type="modified" :count="status.modifiedFiles" />
        <StatusBadge v-if="status.stagedFiles > 0" type="staged" :count="status.stagedFiles" />
        <StatusBadge v-if="status.untrackedFiles > 0" type="untracked" :count="status.untrackedFiles" />
        <StatusBadge v-if="status.unpushedCommits > 0" type="ahead" :count="status.unpushedCommits" />
        <StatusBadge v-if="status.unpulledCommits > 0" type="behind" :count="status.unpulledCommits" />
        <StatusBadge v-if="status.stashCount > 0" type="stash" :count="status.stashCount" />
        <StatusBadge v-if="status.checkingRemote" type="checking" />
        <StatusBadge v-else-if="!status.remoteAccessible && status.remotes?.length > 0" type="unreachable" />
      </template>
    </div>

    <!-- Footer: last checked -->
    <div class="flex items-center gap-1 mt-3">
      <Icon icon="codicon:history" width="12" height="12" :style="{ color: 'var(--color-text-secondary)' }" />
      <span class="text-[10px]" :style="{ color: 'var(--color-text-secondary)' }">
        {{ lastCheckedText }}
      </span>
    </div>

    <ConfirmDialog
      :visible="showRemoveConfirm"
      title="Remove Repository"
      :message="'Remove \'' + repo.name + '\' from monitoring? This won\'t delete the actual repository.'"
      confirm-text="Remove"
      @confirm="emit('remove'); showRemoveConfirm = false"
      @cancel="showRemoveConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import StatusBadge from './StatusBadge.vue'
import TagDropdown from './TagDropdown.vue'
import ConfirmDialog from './ConfirmDialog.vue'
import type { Repository, RepoStatus } from '../stores/repoStore'
import { useTimeAgo } from '@vueuse/core'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

const props = defineProps<{
  repo: Repository
  status: RepoStatus | undefined
  selectable?: boolean
  selected?: boolean
}>()

const emit = defineEmits<{ refresh: []; remove: []; 'toggle-select': [] }>()

const showRemoveConfirm = ref(false)
const checking = ref(false)

onMounted(() => {
  EventsOn('repo:checking', (id: number) => {
    if (id === props.repo.ID) checking.value = true
  })
  EventsOn('repo:checked', (id: number) => {
    if (id === props.repo.ID) checking.value = false
  })
})

onUnmounted(() => {
  EventsOff('repo:checking')
  EventsOff('repo:checked')
})

const tagIds = computed(() => props.repo.tags?.map(t => t.ID) ?? [])

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
  if (!props.status?.lastChecked) return 'Waiting for check...'
  return useTimeAgo(new Date(props.status.lastChecked)).value
})
</script>
