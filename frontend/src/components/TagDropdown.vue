<template>
  <div class="relative" ref="dropdownRef">
    <button
      class="p-1.5 rounded cursor-pointer hover:opacity-80"
      :style="{ color: 'var(--color-text-secondary)' }"
      @click.stop="open = !open"
      title="Manage tags"
    >
      <Icon icon="codicon:tag" width="14" height="14" />
    </button>

    <!-- Dropdown -->
    <transition name="page">
      <div
        v-if="open"
        class="absolute right-0 top-8 z-30 w-44 rounded-lg py-1 border"
        :style="{
          backgroundColor: 'var(--color-surface)',
          borderColor: 'var(--color-border)',
          boxShadow: '0 4px 12px rgba(0,0,0,0.3)',
        }"
      >
        <div v-if="tagStore.tags.length === 0" class="px-3 py-2 text-xs" :style="{ color: 'var(--color-text-secondary)' }">
          No tags yet
        </div>
        <button
          v-for="tag in tagStore.tags"
          :key="tag.ID"
          class="flex items-center gap-2 w-full px-3 py-1.5 text-xs cursor-pointer hover:opacity-80 transition-colors"
          :style="{ color: 'var(--color-text)' }"
          @click.stop="toggleTag(tag.ID)"
        >
          <div
            class="w-3 h-3 rounded-sm border flex items-center justify-center transition-colors"
            :style="{
              backgroundColor: hasTag(tag.ID) ? tag.color : 'transparent',
              borderColor: tag.color,
            }"
          >
            <Icon v-if="hasTag(tag.ID)" icon="codicon:check" width="10" height="10" style="color: #fff" />
          </div>
          <span>{{ tag.name }}</span>
        </button>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useRepoStore } from '../stores/repoStore'
import { useTagStore } from '../stores/tagStore'

const props = defineProps<{ repoId: number; currentTagIds: number[] }>()

const repoStore = useRepoStore()
const tagStore = useTagStore()
const open = ref(false)
const dropdownRef = ref<HTMLElement>()

function hasTag(tagId: number): boolean {
  return props.currentTagIds.includes(tagId)
}

async function toggleTag(tagId: number) {
  if (hasTag(tagId)) {
    await repoStore.unassignTag(props.repoId, tagId)
  } else {
    await repoStore.assignTag(props.repoId, tagId)
  }
}

function onClickOutside(e: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', onClickOutside))
onUnmounted(() => document.removeEventListener('click', onClickOutside))
</script>
