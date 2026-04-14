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

    <!-- Window controls -->
    <div class="flex h-full" style="--wails-draggable: no-drag">
      <button
        class="w-12 h-full flex items-center justify-center transition-colors hover:bg-white/10 cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="minimise"
        title="Minimize"
      >
        <svg width="10" height="1" viewBox="0 0 10 1">
          <rect fill="currentColor" width="10" height="1" />
        </svg>
      </button>
      <button
        class="w-12 h-full flex items-center justify-center transition-colors hover:bg-white/10 cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="toggleMaximise"
        title="Maximize"
      >
        <svg width="10" height="10" viewBox="0 0 10 10">
          <rect fill="none" stroke="currentColor" stroke-width="1" x="0.5" y="0.5" width="9" height="9" rx="1" />
        </svg>
      </button>
      <button
        class="w-12 h-full flex items-center justify-center transition-colors hover:bg-red-500 hover:text-white cursor-pointer"
        :style="{ color: 'var(--color-text-secondary)' }"
        @click="close"
        title="Close"
      >
        <svg width="10" height="10" viewBox="0 0 10 10">
          <line stroke="currentColor" stroke-width="1.2" x1="1" y1="1" x2="9" y2="9" />
          <line stroke="currentColor" stroke-width="1.2" x1="9" y1="1" x2="1" y2="9" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { WindowMinimise, WindowToggleMaximise, WindowClose } from '../../wailsjs/go/main/App'

function minimise() {
  WindowMinimise()
}

function toggleMaximise() {
  WindowToggleMaximise()
}

function close() {
  WindowClose()
}
</script>
