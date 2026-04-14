<template>
  <Teleport to="body">
    <transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @click.self="$emit('close')"
      >
        <div class="absolute inset-0 bg-black/50" />

        <transition name="modal-panel" appear>
          <div
            class="relative w-full max-w-lg rounded-xl p-6 z-10"
            :style="{
              backgroundColor: 'var(--color-surface)',
              border: '1px solid var(--color-border)',
            }"
          >
            <!-- Header -->
            <div class="flex items-center justify-between mb-5">
              <h2 class="text-lg font-bold">Add Repository</h2>
              <button
                class="p-1 rounded cursor-pointer hover:opacity-80"
                :style="{ color: 'var(--color-text-secondary)' }"
                @click="$emit('close')"
              >
                <i class="pi pi-times" />
              </button>
            </div>

            <!-- Mode tabs -->
            <div class="flex gap-1 mb-4 p-1 rounded-lg" :style="{ backgroundColor: 'var(--color-bg)' }">
              <button
                class="flex-1 px-3 py-1.5 rounded-md text-xs font-medium cursor-pointer transition-colors"
                :style="{
                  backgroundColor: mode === 'single' ? 'var(--color-primary)' : 'transparent',
                  color: mode === 'single' ? '#fff' : 'var(--color-text-secondary)',
                }"
                @click="mode = 'single'; scannedRepos = []"
              >
                Single Repo
              </button>
              <button
                class="flex-1 px-3 py-1.5 rounded-md text-xs font-medium cursor-pointer transition-colors"
                :style="{
                  backgroundColor: mode === 'scan' ? 'var(--color-primary)' : 'transparent',
                  color: mode === 'scan' ? '#fff' : 'var(--color-text-secondary)',
                }"
                @click="mode = 'scan'; path = ''"
              >
                Scan Folder
              </button>
            </div>

            <!-- Single mode -->
            <template v-if="mode === 'single'">
              <div class="mb-4">
                <label class="block text-sm font-medium mb-1.5" :style="{ color: 'var(--color-text-secondary)' }">
                  Repository Path
                </label>
                <div class="flex gap-2">
                  <input
                    v-model="path"
                    type="text"
                    readonly
                    placeholder="Select a git repository folder..."
                    class="flex-1 px-3 py-2 rounded-lg text-sm border outline-none"
                    :style="{
                      backgroundColor: 'var(--color-bg)',
                      borderColor: 'var(--color-border)',
                      color: 'var(--color-text)',
                    }"
                  />
                  <button
                    class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white shrink-0"
                    :style="{ backgroundColor: 'var(--color-primary)' }"
                    @click="selectDir"
                  >
                    Browse
                  </button>
                </div>
              </div>
            </template>

            <!-- Scan mode -->
            <template v-else>
              <div class="mb-4">
                <label class="block text-sm font-medium mb-1.5" :style="{ color: 'var(--color-text-secondary)' }">
                  Parent Folder
                </label>
                <div class="flex gap-2">
                  <input
                    v-model="scanPath"
                    type="text"
                    readonly
                    placeholder="Select a folder to scan..."
                    class="flex-1 px-3 py-2 rounded-lg text-sm border outline-none"
                    :style="{
                      backgroundColor: 'var(--color-bg)',
                      borderColor: 'var(--color-border)',
                      color: 'var(--color-text)',
                    }"
                  />
                  <button
                    class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white shrink-0"
                    :style="{ backgroundColor: 'var(--color-primary)' }"
                    @click="scanDir"
                  >
                    {{ scanning ? 'Scanning...' : 'Scan' }}
                  </button>
                </div>
              </div>

              <!-- Scan results -->
              <div v-if="scannedRepos.length > 0" class="mb-4 max-h-52 overflow-auto rounded-lg border" :style="{ borderColor: 'var(--color-border)' }">
                <div class="flex items-center justify-between px-3 py-2" :style="{ backgroundColor: 'var(--color-bg)' }">
                  <span class="text-xs font-medium" :style="{ color: 'var(--color-text-secondary)' }">
                    {{ newRepos.length }} new / {{ existingPaths.size }} already added
                  </span>
                  <button
                    v-if="newRepos.length > 0"
                    class="text-xs cursor-pointer underline"
                    :style="{ color: 'var(--color-primary)' }"
                    @click="toggleAllScanned"
                  >
                    {{ selectedScanned.size === newRepos.length ? 'Deselect all' : 'Select all' }}
                  </button>
                </div>
                <div
                  v-for="repo in scannedRepos"
                  :key="repo"
                  class="flex items-center gap-2 px-3 py-2 transition-colors border-t"
                  :class="existingPaths.has(repo) ? 'opacity-40' : 'cursor-pointer hover:opacity-80'"
                  :style="{ borderColor: 'var(--color-border)' }"
                  @click="!existingPaths.has(repo) && toggleScanned(repo)"
                >
                  <div
                    v-if="existingPaths.has(repo)"
                    class="w-4 h-4 rounded flex items-center justify-center shrink-0"
                    :style="{ backgroundColor: 'var(--color-muted)' }"
                  >
                    <Icon icon="codicon:check" width="10" height="10" style="color: #fff" />
                  </div>
                  <div
                    v-else
                    class="w-4 h-4 rounded border flex items-center justify-center shrink-0 transition-colors"
                    :style="{
                      backgroundColor: selectedScanned.has(repo) ? 'var(--color-primary)' : 'transparent',
                      borderColor: selectedScanned.has(repo) ? 'var(--color-primary)' : 'var(--color-border)',
                    }"
                  >
                    <Icon v-if="selectedScanned.has(repo)" icon="codicon:check" width="10" height="10" style="color: #fff" />
                  </div>
                  <span class="text-xs truncate flex-1" :style="{ color: 'var(--color-text)' }">{{ repo }}</span>
                  <span v-if="existingPaths.has(repo)" class="text-[10px] shrink-0" :style="{ color: 'var(--color-text-secondary)' }">already added</span>
                </div>
              </div>
              <div v-else-if="scanPath && !scanning" class="mb-4 text-xs text-center py-4" :style="{ color: 'var(--color-text-secondary)' }">
                No git repositories found
              </div>
            </template>

            <!-- Error / Success -->
            <transition name="page">
              <p v-if="error" class="text-sm mb-4" :style="{ color: 'var(--color-danger)' }">{{ error }}</p>
            </transition>
            <transition name="page">
              <p v-if="success" class="text-sm mb-4 flex items-center gap-1.5" :style="{ color: 'var(--color-success)' }">
                <i class="pi pi-check-circle" /> {{ successMessage }}
              </p>
            </transition>

            <!-- Actions -->
            <div class="flex justify-end gap-2">
              <button
                class="px-4 py-2 rounded-lg text-sm cursor-pointer"
                :style="{
                  backgroundColor: 'var(--color-bg)',
                  color: 'var(--color-text-secondary)',
                  border: '1px solid var(--color-border)',
                }"
                @click="$emit('close')"
              >
                Cancel
              </button>
              <button
                v-if="mode === 'single'"
                :disabled="!path || repoStore.loading"
                class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white disabled:opacity-50"
                :style="{ backgroundColor: 'var(--color-primary)' }"
                @click="addSingleRepo"
              >
                {{ repoStore.loading ? 'Adding...' : 'Add Repository' }}
              </button>
              <button
                v-else
                :disabled="selectedScanned.size === 0 || repoStore.loading"
                class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white disabled:opacity-50"
                :style="{ backgroundColor: 'var(--color-primary)' }"
                @click="addScannedRepos"
              >
                {{ repoStore.loading ? 'Adding...' : `Add ${selectedScanned.size} Repos` }}
              </button>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useRepoStore } from '../stores/repoStore'
import { SelectDirectory, ScanForRepos, AddRepositories } from '../../wailsjs/go/main/App'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ close: [] }>()

const repoStore = useRepoStore()
const mode = ref<'single' | 'scan'>('single')
const path = ref('')
const scanPath = ref('')
const scanning = ref(false)
const scannedRepos = ref<string[]>([])
const selectedScanned = ref<Set<string>>(new Set())
const error = ref('')
const success = ref(false)
const successMessage = ref('')

const existingPaths = computed(() => new Set(repoStore.repositories.map(r => r.path)))
const newRepos = computed(() => scannedRepos.value.filter(r => !existingPaths.value.has(r)))

watch(() => props.visible, (val) => {
  if (val) {
    path.value = ''
    scanPath.value = ''
    scannedRepos.value = []
    selectedScanned.value = new Set()
    error.value = ''
    success.value = false
    mode.value = 'single'
  }
})

async function selectDir() {
  try {
    const dir = await SelectDirectory()
    if (dir) {
      path.value = dir
      error.value = ''
      success.value = false
    }
  } catch (e) {
    error.value = 'Failed to open directory picker'
  }
}

async function scanDir() {
  try {
    const dir = await SelectDirectory()
    if (!dir) return
    scanPath.value = dir
    scanning.value = true
    error.value = ''
    scannedRepos.value = []
    selectedScanned.value = new Set()
    const repos = await ScanForRepos(dir)
    scannedRepos.value = repos ?? []
    // Select only new repos by default
    const existing = new Set(repoStore.repositories.map(r => r.path))
    selectedScanned.value = new Set(scannedRepos.value.filter(r => !existing.has(r)))
  } catch (e: any) {
    error.value = typeof e === 'string' ? e : 'Scan failed'
  } finally {
    scanning.value = false
  }
}

function toggleScanned(repo: string) {
  if (selectedScanned.value.has(repo)) {
    selectedScanned.value.delete(repo)
  } else {
    selectedScanned.value.add(repo)
  }
  selectedScanned.value = new Set(selectedScanned.value) // trigger reactivity
}

function toggleAllScanned() {
  if (selectedScanned.value.size === newRepos.value.length) {
    selectedScanned.value = new Set()
  } else {
    selectedScanned.value = new Set(newRepos.value)
  }
}

async function addSingleRepo() {
  if (!path.value) return
  error.value = ''
  success.value = false
  try {
    await repoStore.addRepo(path.value)
    success.value = true
    successMessage.value = 'Repository added successfully!'
    setTimeout(() => emit('close'), 800)
  } catch (e: any) {
    error.value = typeof e === 'string' ? e : (e?.message || 'Failed to add repository')
  }
}

async function addScannedRepos() {
  if (selectedScanned.value.size === 0) return
  error.value = ''
  success.value = false
  repoStore.loading = true
  try {
    const paths = Array.from(selectedScanned.value)
    const count = await AddRepositories(paths)
    await repoStore.fetchRepositories()
    success.value = true
    successMessage.value = `${count} repositories added!`
    setTimeout(() => emit('close'), 800)
  } catch (e: any) {
    error.value = typeof e === 'string' ? e : 'Failed to add repositories'
  } finally {
    repoStore.loading = false
  }
}
</script>
