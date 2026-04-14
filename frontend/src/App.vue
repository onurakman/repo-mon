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
      <Sidebar @add-repo="showAddRepo = true" />
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
import { ref, onMounted } from 'vue'
import TitleBar from './components/TitleBar.vue'
import Sidebar from './components/Sidebar.vue'
import AddRepoModal from './components/AddRepoModal.vue'
import { useSettingsStore } from './stores/settingsStore'
import { useRepoStore } from './stores/repoStore'
import { useTagStore } from './stores/tagStore'
import { usePolling } from './composables/usePolling'

const settingsStore = useSettingsStore()
const repoStore = useRepoStore()
const tagStore = useTagStore()
const showAddRepo = ref(false)

onMounted(async () => {
  settingsStore.init()
  await Promise.all([
    repoStore.fetchRepositories(),
    tagStore.fetchTags(),
  ])
  // Fetch cached statuses immediately after repos load
  await repoStore.fetchStatuses()
  // Then start polling for live updates
  usePolling(5000)
})
</script>
