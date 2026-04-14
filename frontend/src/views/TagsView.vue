<template>
  <div class="max-w-lg mx-auto">
    <h1 class="text-xl font-bold mb-6">Tags</h1>

    <div class="flex gap-2 mb-6">
      <input
        v-model="newTagName"
        type="text"
        placeholder="Tag name..."
        class="flex-1 px-3 py-2 rounded-lg text-sm border outline-none"
        :style="{
          backgroundColor: 'var(--color-surface)',
          borderColor: 'var(--color-border)',
          color: 'var(--color-text)',
        }"
        @keyup.enter="createTag"
      />
      <input
        v-model="newTagColor"
        type="color"
        class="w-10 h-10 rounded cursor-pointer border-0 p-0"
      />
      <button
        :disabled="!newTagName.trim()"
        class="px-4 py-2 rounded-lg text-sm cursor-pointer text-white disabled:opacity-50"
        :style="{ backgroundColor: 'var(--color-primary)' }"
        @click="createTag"
      >
        Add
      </button>
    </div>

    <div class="space-y-2">
      <div
        v-for="tag in tagStore.tags"
        :key="tag.ID"
        class="flex items-center justify-between p-3 rounded-lg"
        :style="{ backgroundColor: 'var(--color-surface)' }"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-4 h-4 rounded-full"
            :style="{ backgroundColor: tag.color }"
          />
          <span class="text-sm font-medium">{{ tag.name }}</span>
        </div>
        <button
          class="p-1 rounded cursor-pointer hover:opacity-80"
          :style="{ color: 'var(--color-danger)' }"
          @click="tagToDelete = tag; showDeleteConfirm = true"
          title="Delete tag"
        >
          <i class="pi pi-trash text-sm" />
        </button>
      </div>

      <p
        v-if="tagStore.tags.length === 0"
        class="text-sm py-8 text-center"
        :style="{ color: 'var(--color-text-secondary)' }"
      >
        No tags yet. Create one above.
      </p>
    </div>

    <ConfirmDialog
      :visible="showDeleteConfirm"
      title="Delete Tag"
      :message="'Delete tag \'' + (tagToDelete?.name ?? '') + '\'? It will be removed from all repositories.'"
      confirm-text="Delete"
      @confirm="deleteTag"
      @cancel="showDeleteConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useTagStore, type Tag } from '../stores/tagStore'
import ConfirmDialog from '../components/ConfirmDialog.vue'

const tagStore = useTagStore()
const newTagName = ref('')
const newTagColor = ref('#10b981')
const showDeleteConfirm = ref(false)
const tagToDelete = ref<Tag | null>(null)

async function createTag() {
  if (!newTagName.value.trim()) return
  await tagStore.addTag(newTagName.value.trim(), newTagColor.value)
  newTagName.value = ''
}

async function deleteTag() {
  if (tagToDelete.value) {
    await tagStore.removeTag(tagToDelete.value.ID)
  }
  showDeleteConfirm.value = false
  tagToDelete.value = null
}
</script>
