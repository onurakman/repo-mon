<template>
  <div
    class="flex flex-col h-full w-full overflow-hidden"
    :style="{
      backgroundColor: 'var(--color-bg)',
      color: 'var(--color-text)',
    }"
  >
    <TitleBar />
    <div class="flex flex-1 overflow-hidden">
      <Sidebar />
      <main class="flex-1 overflow-auto p-6">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" @add-repo="showAddRepo = true" />
          </transition>
        </router-view>
      </main>
    </div>
    <AddRepoModal :visible="showAddRepo" @close="showAddRepo = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import TitleBar from './components/TitleBar.vue'
import Sidebar from './components/Sidebar.vue'
import AddRepoModal from './components/AddRepoModal.vue'
import { useSettingsStore } from './stores/settingsStore'
import { useRepoStore } from './stores/repoStore'
import { useTagStore } from './stores/tagStore'
import { usePolling } from './composables/usePolling'
import { EventsOn } from '../wailsjs/runtime/runtime'

const settingsStore = useSettingsStore()
const repoStore = useRepoStore()
const tagStore = useTagStore()
const showAddRepo = ref(false)
let cancelChecking: (() => void) | null = null
let cancelChecked: (() => void) | null = null

// Register composable at setup scope (required by VueUse), start paused
const { resume: startPolling } = usePolling(5000, false)

// Debounced fetch for batching rapid repo:checked events
const debouncedFetch = useDebounceFn(() => repoStore.fetchStatuses(), 300)

onMounted(async () => {
  settingsStore.init()
  await Promise.all([
    repoStore.fetchRepositories(),
    tagStore.fetchTags(),
  ])
  await repoStore.fetchStatuses()

  // Track checking state instantly via events (no RPC latency)
  cancelChecking = EventsOn('repo:checking', (id: number) => {
    repoStore.setChecking(id)
  })
  cancelChecked = EventsOn('repo:checked', (id: number) => {
    repoStore.clearChecking(id)
    debouncedFetch()
  })

  startPolling()
})

onUnmounted(() => {
  cancelChecking?.()
  cancelChecked?.()
})
</script>
