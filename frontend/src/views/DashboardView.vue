<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold">Dashboard</h1>
      <div class="flex items-center gap-2">
        <!-- Search -->
        <div class="relative">
          <Icon
            icon="codicon:search"
            width="14" height="14"
            class="absolute left-2 top-1/2 -translate-y-1/2 pointer-events-none"
            :style="{ color: 'var(--color-text-secondary)' }"
          />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search..."
            class="h-7 w-44 pl-7 pr-7 rounded-md text-xs border outline-none transition-colors"
            :style="{
              borderColor: 'var(--color-border)',
              backgroundColor: 'var(--color-surface)',
              color: 'var(--color-text)',
            }"
          />
          <button
            v-if="searchQuery"
            class="absolute right-1.5 top-1/2 -translate-y-1/2 cursor-pointer rounded-sm hover:opacity-80"
            :style="{ color: 'var(--color-text-secondary)' }"
            @click="searchQuery = ''"
            title="Clear search"
          >
            <Icon icon="codicon:close" width="14" height="14" />
          </button>
        </div>
        <button
          class="w-7 h-7 rounded-md text-xs cursor-pointer transition-colors border flex items-center justify-center"
          :style="{
            borderColor: selectMode ? 'var(--color-primary)' : 'var(--color-border)',
            color: selectMode ? 'var(--color-primary)' : 'var(--color-text-secondary)',
            backgroundColor: selectMode ? 'var(--color-primary)' + '20' : 'transparent',
          }"
          @click="toggleSelectMode"
          title="Select repos"
        >
          <Icon icon="codicon:checklist" width="14" height="14" />
        </button>
        <button
          class="w-7 h-7 rounded-md text-xs cursor-pointer transition-colors border flex items-center justify-center"
          :style="{
            borderColor: 'var(--color-border)',
            color: 'var(--color-text-secondary)',
          }"
          @click="refreshAll"
          title="Refresh All"
        >
          <Icon
            icon="codicon:refresh"
            width="14" height="14"
            class="transition-transform"
            :class="{ 'refresh-spinning': refreshingAll }"
          />
        </button>
        <ViewToggle v-model="viewMode" @update:model-value="onViewModeChange" />
      </div>
    </div>

    <!-- Bulk action bar -->
    <transition name="page">
      <div
        v-if="selectMode && repoStore.selectedIds.size > 0"
        class="flex items-center gap-3 mb-4 px-3 py-2 rounded-lg"
        :style="{ backgroundColor: 'var(--color-surface)', border: '1px solid var(--color-border)' }"
      >
        <span class="text-xs font-medium" :style="{ color: 'var(--color-text-secondary)' }">
          {{ repoStore.selectedIds.size }} selected
        </span>
        <button
          class="text-xs cursor-pointer underline"
          :style="{ color: 'var(--color-primary)' }"
          @click="repoStore.selectAll()"
        >
          Select all
        </button>
        <button
          class="text-xs cursor-pointer underline"
          :style="{ color: 'var(--color-text-secondary)' }"
          @click="repoStore.clearSelection()"
        >
          Clear
        </button>

        <div class="flex-1" />

        <!-- Bulk tag assign -->
        <div class="flex items-center gap-1">
          <span class="text-xs" :style="{ color: 'var(--color-text-secondary)' }">Tag:</span>
          <button
            v-for="tag in tagStore.tags"
            :key="tag.ID"
            class="px-2 py-0.5 rounded text-[10px] text-white cursor-pointer hover:opacity-80 transition-opacity"
            :style="{ backgroundColor: tag.color }"
            @click="bulkAssignTag(tag.ID)"
          >
            {{ tag.name }}
          </button>
        </div>
      </div>
    </transition>

    <!-- Tag filters -->
    <div v-if="tagStore.tags.length && !selectMode" class="flex flex-wrap gap-2 mb-4">
      <TagChip
        name="All"
        color="var(--color-primary)"
        :active="selectedTagIds.length === 0"
        @toggle="selectedTagIds = []"
      />
      <TagChip
        v-for="tag in tagStore.tags"
        :key="tag.ID"
        :name="tag.name"
        :color="tag.color"
        :active="selectedTagIds.includes(tag.ID)"
        @toggle="toggleTagFilter(tag.ID)"
      />
    </div>

    <!-- Empty state -->
    <div
      v-if="repoStore.repositories.length === 0"
      class="flex flex-col items-center justify-center py-20"
    >
      <Icon icon="codicon:repo" width="48" height="48" class="mb-4" :style="{ color: 'var(--color-text-secondary)' }" />
      <p class="text-lg mb-2" :style="{ color: 'var(--color-text-secondary)' }">No repositories added yet</p>
      <button
        class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white"
        :style="{ backgroundColor: 'var(--color-primary)' }"
        @click="$emit('add-repo')"
      >
        <i class="pi pi-plus mr-1" /> Add Repository
      </button>
    </div>

    <!-- Grid view: draggable when no filter, static when filtered -->
    <template v-else-if="viewMode === 'grid'">
      <draggable
        v-if="!isFiltered"
        v-model="repoStore.repositories"
        item-key="ID"
        :animation="200"
        ghost-class="opacity-30"
        drag-class="rotate-2"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
        @end="onDragEnd"
      >
        <template #item="{ element: repo }">
          <RepoCard
            :repo="repo"
            :status="repoStore.statuses[repo.ID]"
            :selectable="selectMode"
            :selected="repoStore.selectedIds.has(repo.ID)"
            @refresh="repoStore.refreshRepo(repo.ID)"
            @remove="repoStore.removeRepo(repo.ID)"
            @toggle-select="repoStore.toggleSelect(repo.ID)"
          />
        </template>
      </draggable>
      <TransitionGroup
        v-else
        name="card-list"
        tag="div"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
      >
        <RepoCard
          v-for="repo in filteredRepos"
          :key="repo.ID"
          :repo="repo"
          :status="repoStore.statuses[repo.ID]"
          :selectable="selectMode"
          :selected="repoStore.selectedIds.has(repo.ID)"
          @refresh="repoStore.refreshRepo(repo.ID)"
          @toggle-select="repoStore.toggleSelect(repo.ID)"
        />
      </TransitionGroup>
    </template>

    <!-- List view: draggable when no filter, static when filtered -->
    <template v-else>
      <draggable
        v-if="!isFiltered"
        v-model="repoStore.repositories"
        item-key="ID"
        :animation="200"
        ghost-class="opacity-30"
        class="flex flex-col gap-2"
        @end="onDragEnd"
      >
        <template #item="{ element: repo }">
          <RepoListItem
            :repo="repo"
            :status="repoStore.statuses[repo.ID]"
            @refresh="repoStore.refreshRepo(repo.ID)"
          />
        </template>
      </draggable>
      <TransitionGroup
        v-else
        name="card-list"
        tag="div"
        class="flex flex-col gap-2"
      >
        <RepoListItem
          v-for="repo in filteredRepos"
          :key="repo.ID"
          :repo="repo"
          :status="repoStore.statuses[repo.ID]"
          @refresh="repoStore.refreshRepo(repo.ID)"
        />
      </TransitionGroup>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { refDebounced } from '@vueuse/core'
import { Icon } from '@iconify/vue'
import draggable from 'vuedraggable'
import { useRepoStore } from '../stores/repoStore'
import { useTagStore } from '../stores/tagStore'
import { useSettingsStore } from '../stores/settingsStore'
import RepoCard from '../components/RepoCard.vue'
import RepoListItem from '../components/RepoListItem.vue'
import TagChip from '../components/TagChip.vue'
import ViewToggle from '../components/ViewToggle.vue'

defineEmits<{ 'add-repo': [] }>()

const repoStore = useRepoStore()
const tagStore = useTagStore()
const settingsStore = useSettingsStore()

const viewMode = ref<'grid' | 'list'>(settingsStore.settings.viewMode as 'grid' | 'list')
const searchQuery = ref('')
const debouncedSearch = refDebounced(searchQuery, 250)
const selectedTagIds = ref<number[]>([])
const refreshingAll = ref(false)
const selectMode = ref(false)

function toggleSelectMode() {
  selectMode.value = !selectMode.value
  if (!selectMode.value) {
    repoStore.clearSelection()
  }
}

function refreshAll() {
  refreshingAll.value = true
  repoStore.refreshAll()
  setTimeout(() => { refreshingAll.value = false }, 600)
}

function onDragEnd() {
  repoStore.updateSortOrder()
}

const isFiltered = computed(() => selectedTagIds.value.length > 0 || debouncedSearch.value.length > 0)

async function bulkAssignTag(tagId: number) {
  const ids = Array.from(repoStore.selectedIds)
  await repoStore.assignTagToRepos(ids, tagId)
  repoStore.clearSelection()
  selectMode.value = false
}

const filteredRepos = computed(() => {
  let repos = repoStore.repositories
  if (selectedTagIds.value.length > 0) {
    repos = repos.filter((repo) =>
      repo.tags?.some((tag) => selectedTagIds.value.includes(tag.ID))
    )
  }
  const q = debouncedSearch.value.toLowerCase()
  if (q) {
    repos = repos.filter((repo) =>
      repo.name.toLowerCase().includes(q) || repo.path.toLowerCase().includes(q)
    )
  }
  return repos
})

function toggleTagFilter(tagId: number) {
  const idx = selectedTagIds.value.indexOf(tagId)
  if (idx >= 0) {
    selectedTagIds.value.splice(idx, 1)
  } else {
    selectedTagIds.value.push(tagId)
  }
}

function onViewModeChange(mode: 'grid' | 'list') {
  settingsStore.updateSettings({ viewMode: mode })
}
</script>
