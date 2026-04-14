<template>
  <Teleport to="body">
    <transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @click.self="$emit('cancel')"
      >
        <div class="absolute inset-0 bg-black/50" />

        <transition name="modal-panel" appear>
          <div
            class="relative w-full max-w-sm rounded-xl p-5 z-10"
            :style="{
              backgroundColor: 'var(--color-surface)',
              border: '1px solid var(--color-border)',
            }"
          >
            <div class="flex items-start gap-3 mb-4">
              <div
                class="w-9 h-9 rounded-full flex items-center justify-center shrink-0"
                :style="{ backgroundColor: 'var(--color-danger)' + '20' }"
              >
                <Icon icon="codicon:warning" width="18" height="18" :style="{ color: 'var(--color-danger)' }" />
              </div>
              <div>
                <h3 class="font-bold text-sm mb-1">{{ title }}</h3>
                <p class="text-xs" :style="{ color: 'var(--color-text-secondary)' }">{{ message }}</p>
              </div>
            </div>

            <div class="flex justify-end gap-2">
              <button
                class="px-3 py-1.5 rounded-lg text-xs cursor-pointer"
                :style="{
                  backgroundColor: 'var(--color-bg)',
                  color: 'var(--color-text-secondary)',
                  border: '1px solid var(--color-border)',
                }"
                @click="$emit('cancel')"
              >
                Cancel
              </button>
              <button
                class="px-3 py-1.5 rounded-lg text-xs cursor-pointer text-white"
                :style="{ backgroundColor: 'var(--color-danger)' }"
                @click="$emit('confirm')"
              >
                {{ confirmText }}
              </button>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{
  visible: boolean
  title: string
  message: string
  confirmText?: string
}>()

defineEmits<{
  confirm: []
  cancel: []
}>()
</script>
