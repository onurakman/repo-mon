<template>
  <div class="max-w-lg mx-auto">
    <h1 class="text-xl font-bold mb-6">Settings</h1>

    <div class="space-y-6">
      <div>
        <label class="block text-sm font-medium mb-2" :style="{ color: 'var(--color-text-secondary)' }">
          Theme
        </label>
        <div class="grid grid-cols-3 gap-2">
          <button
            v-for="t in themes"
            :key="t.value"
            class="p-3 rounded-lg border-2 cursor-pointer transition-colors text-center"
            :style="{
              backgroundColor: 'var(--color-surface)',
              borderColor: settings.theme === t.value ? 'var(--color-primary)' : 'var(--color-border)',
            }"
            @click="updateTheme(t.value)"
          >
            <div class="flex justify-center gap-1 mb-2">
              <div
                v-for="color in t.preview"
                :key="color"
                class="w-4 h-4 rounded-full"
                :style="{ backgroundColor: color }"
              />
            </div>
            <span class="text-xs">{{ t.label }}</span>
          </button>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <label class="text-sm font-medium" :style="{ color: 'var(--color-text-secondary)' }">
          Dark Mode
        </label>
        <button
          class="w-12 h-6 rounded-full cursor-pointer transition-colors relative"
          :style="{
            backgroundColor: settings.darkMode ? 'var(--color-primary)' : 'var(--color-border)',
          }"
          @click="settingsStore.updateSettings({ darkMode: !settings.darkMode })"
        >
          <div
            class="w-5 h-5 rounded-full bg-white absolute top-0.5 transition-all"
            :style="{ left: settings.darkMode ? '26px' : '2px' }"
          />
        </button>
      </div>

      <div>
        <label class="block text-sm font-medium mb-2" :style="{ color: 'var(--color-text-secondary)' }">
          Default Poll Interval: {{ settings.globalPollInterval }}s
        </label>
        <input
          type="range"
          :min="5"
          :max="300"
          :step="5"
          :value="settings.globalPollInterval"
          class="w-full"
          @input="onIntervalChange"
        />
        <div class="flex justify-between text-xs" :style="{ color: 'var(--color-text-secondary)' }">
          <span>5s</span>
          <span>300s</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSettingsStore } from '../stores/settingsStore'

const settingsStore = useSettingsStore()
const settings = computed(() => settingsStore.settings)

const themes = [
  { value: 'neutral-carbon', label: 'Carbon', preview: ['#171717', '#10b981', '#1c1c1c'] },
  { value: 'slate-blue', label: 'Slate Blue', preview: ['#0f172a', '#3b82f6', '#1e293b'] },
  { value: 'deep-purple', label: 'Purple', preview: ['#13111c', '#8b5cf6', '#1a1726'] },
]

function updateTheme(theme: string) {
  settingsStore.updateSettings({ theme })
}

function onIntervalChange(e: Event) {
  const value = parseInt((e.target as HTMLInputElement).value)
  settingsStore.updateSettings({ globalPollInterval: value })
}
</script>
