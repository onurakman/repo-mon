<template>
  <div class="max-w-lg">
    <h1 class="text-xl font-bold mb-6">Add Repository</h1>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium mb-1" :style="{ color: 'var(--color-text-secondary)' }">
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
              backgroundColor: 'var(--color-surface)',
              borderColor: 'var(--color-border)',
              color: 'var(--color-text)',
            }"
          />
          <button
            class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white"
            :style="{ backgroundColor: 'var(--color-primary)' }"
            @click="selectDir"
          >
            Browse
          </button>
        </div>
      </div>

      <p v-if="error" class="text-sm" :style="{ color: 'var(--color-danger)' }">
        {{ error }}
      </p>

      <p v-if="success" class="text-sm" :style="{ color: 'var(--color-success)' }">
        Repository added successfully!
      </p>

      <button
        :disabled="!path || repoStore.loading"
        class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white disabled:opacity-50"
        :style="{ backgroundColor: 'var(--color-primary)' }"
        @click="addRepo"
      >
        <i class="pi pi-plus mr-1" />
        {{ repoStore.loading ? 'Adding...' : 'Add Repository' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRepoStore } from '../stores/repoStore'
import { SelectDirectory } from '../../wailsjs/go/main/App'

const repoStore = useRepoStore()
const path = ref('')
const error = ref('')
const success = ref(false)

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

async function addRepo() {
  if (!path.value) return
  error.value = ''
  success.value = false
  try {
    await repoStore.addRepo(path.value)
    success.value = true
    path.value = ''
  } catch (e: any) {
    error.value = e?.message || 'Failed to add repository'
  }
}
</script>
