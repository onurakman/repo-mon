<template>
  <Teleport to="body">
    <transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @click.self="$emit('close')"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/50" />

        <!-- Modal panel -->
        <transition name="modal-panel" appear>
          <div
            class="relative w-full max-w-md rounded-xl p-6 z-10"
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

            <!-- Path selector -->
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

            <!-- Error -->
            <transition name="page">
              <p v-if="error" class="text-sm mb-4" :style="{ color: 'var(--color-danger)' }">
                {{ error }}
              </p>
            </transition>

            <!-- Success -->
            <transition name="page">
              <p v-if="success" class="text-sm mb-4 flex items-center gap-1.5" :style="{ color: 'var(--color-success)' }">
                <i class="pi pi-check-circle" /> Repository added successfully!
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
        </transition>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRepoStore } from '../stores/repoStore'
import { SelectDirectory } from '../../wailsjs/go/main/App'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ close: [] }>()

const repoStore = useRepoStore()
const path = ref('')
const error = ref('')
const success = ref(false)

watch(() => props.visible, (val) => {
  if (val) {
    path.value = ''
    error.value = ''
    success.value = false
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

async function addRepo() {
  if (!path.value) return
  error.value = ''
  success.value = false
  try {
    await repoStore.addRepo(path.value)
    success.value = true
    setTimeout(() => emit('close'), 800)
  } catch (e: any) {
    error.value = e?.message || 'Failed to add repository'
  }
}
</script>
