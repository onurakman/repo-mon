<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold">Dashboard</h1>
      <div class="flex items-center gap-2">
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
            class="text-sm transition-transform"
            :class="{ 'refresh-spinning': refreshingAll }"
          />
        </button>
        <ViewToggle v-model="viewMode" @update:model-value="onViewModeChange" />
      </div>
    </div>

    <div v-if="tagStore.tags.length" class="flex flex-wrap gap-2 mb-4">
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

    <div
      v-if="repoStore.repositories.length === 0"
      class="flex flex-col items-center justify-center py-20"
    >
      <i class="pi pi-folder-open text-4xl mb-4" :style="{ color: 'var(--color-text-secondary)' }" />
      <p class="text-lg mb-2" :style="{ color: 'var(--color-text-secondary)' }">No repositories added yet</p>
      <button
        class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white"
        :style="{ backgroundColor: 'var(--color-primary)' }"
        @click="$emit('add-repo')"
      >
        <i class="pi pi-plus mr-1" /> Add Repository
      </button>
    </div>

    <TransitionGroup
      v-else-if="viewMode === 'grid'"
      name="card-list"
      tag="div"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
    >
      <RepoCard
        v-for="repo in filteredRepos"
        :key="repo.ID"
        :repo="repo"
        :status="repoStore.statuses[repo.ID]"
        @refresh="repoStore.refreshRepo(repo.ID)"
      />
    </TransitionGroup>

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
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Icon } from '@iconify/vue'
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
const selectedTagIds = ref<number[]>([])
const refreshingAll = ref(false)

function refreshAll() {
  refreshingAll.value = true
  repoStore.refreshAll()
  setTimeout(() => { refreshingAll.value = false }, 600)
}

const filteredRepos = computed(() => {
  if (selectedTagIds.value.length === 0) return repoStore.repositories
  return repoStore.repositories.filter((repo) =>
    repo.tags?.some((tag) => selectedTagIds.value.includes(tag.ID))
  )
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
