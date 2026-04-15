<template>
  <div
    class="flex items-center justify-between h-14 select-none shrink-0"
    style="--wails-draggable: drag"
    :style="{
      backgroundColor: 'var(--color-sidebar)',
      borderBottom: '1px solid var(--color-border)',
    }"
  >
    <!-- App title -->
    <div class="flex items-center gap-2.5 pl-3">
      <img src="../assets/images/logo.png" alt="Repo Monitor" class="w-9 h-9" />
      <span class="text-sm font-semibold tracking-wide" :style="{ color: 'var(--color-text-secondary)' }">
        Repo Monitor
      </span>
    </div>

    <!-- Quick settings + Window controls -->
    <div class="flex h-full" style="--wails-draggable: no-drag">
      <!-- Theme cycle -->
      <button
        class="w-10 h-full flex items-center justify-center transition-colors hover:bg-white/10 cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="cycleTheme"
        v-tooltip.bottom="currentThemeLabel"
      >
        <div
          class="w-3.5 h-3.5 rounded-full border border-white/30"
          :style="{ backgroundColor: currentThemeColor }"
        />
      </button>
      <!-- Dark/Light toggle -->
      <button
        class="w-10 h-full flex items-center justify-center transition-colors hover:bg-white/10 cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="toggleDarkMode"
        v-tooltip.bottom="settings.darkMode ? 'Light mode' : 'Dark mode'"
      >
        <i :class="settings.darkMode ? 'pi pi-sun' : 'pi pi-moon'" class="text-sm" />
      </button>
      <!-- Minimize -->
      <button
        class="w-12 h-full flex items-center justify-center transition-colors hover:bg-white/10 cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="minimise"
        v-tooltip.bottom="'Minimize'"
      >
        <svg width="10" height="1" viewBox="0 0 10 1">
          <rect fill="currentColor" width="10" height="1" />
        </svg>
      </button>
      <button
        class="w-12 h-full flex items-center justify-center transition-colors hover:bg-white/10 cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="toggleMaximise"
        v-tooltip.bottom="'Maximize'"
      >
        <svg width="10" height="10" viewBox="0 0 10 10">
          <rect fill="none" stroke="currentColor" stroke-width="1" x="0.5" y="0.5" width="9" height="9" rx="1" />
        </svg>
      </button>
      <button
        class="w-12 h-full flex items-center justify-center transition-colors hover:bg-red-500 hover:text-white cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="showCloseConfirm = true"
        v-tooltip.bottom="'Close'"
      >
        <svg width="10" height="10" viewBox="0 0 10 10">
          <line stroke="currentColor" stroke-width="1.2" x1="1" y1="1" x2="9" y2="9" />
          <line stroke="currentColor" stroke-width="1.2" x1="9" y1="1" x2="1" y2="9" />
        </svg>
      </button>
    </div>

    <ConfirmDialog
      :visible="showCloseConfirm"
      title="Quit Application"
      message="Are you sure you want to quit?"
      confirmText="Quit"
      @confirm="close"
      @cancel="showCloseConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { WindowMinimise, WindowToggleMaximise, WindowClose } from '../../wailsjs/go/main/App'
import ConfirmDialog from './ConfirmDialog.vue'
import { useSettingsStore } from '../stores/settingsStore'

const settingsStore = useSettingsStore()
const settings = computed(() => settingsStore.settings)

const themes = [
  { value: 'neutral-carbon', label: 'Carbon', color: '#10b981' },
  { value: 'slate-blue', label: 'Slate Blue', color: '#3b82f6' },
  { value: 'deep-purple', label: 'Purple', color: '#8b5cf6' },
]

const currentThemeLabel = computed(() => {
  const t = themes.find(t => t.value === settings.value.theme)
  return t ? t.label : 'Theme'
})

const currentThemeColor = computed(() => {
  const t = themes.find(t => t.value === settings.value.theme)
  return t ? t.color : '#10b981'
})

function cycleTheme() {
  const idx = themes.findIndex(t => t.value === settings.value.theme)
  const next = themes[(idx + 1) % themes.length]
  settingsStore.updateSettings({ theme: next.value })
}

function toggleDarkMode() {
  settingsStore.updateSettings({ darkMode: !settings.value.darkMode })
}

const showCloseConfirm = ref(false)

function minimise() {
  WindowMinimise()
}

function toggleMaximise() {
  WindowToggleMaximise()
}

function close() {
  showCloseConfirm.value = false
  WindowClose()
}
</script>
