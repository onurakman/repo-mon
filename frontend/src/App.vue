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

usePolling(5000)

onMounted(async () => {
  await settingsStore.fetchSettings()
  await repoStore.fetchRepositories()
  await tagStore.fetchTags()
})
</script>
