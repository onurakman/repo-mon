<template>
  <button
    class="w-10 h-10 rounded-lg flex items-center justify-center transition-colors cursor-pointer"
    :style="{
      backgroundColor: active ? 'var(--color-primary)' : 'transparent',
      color: active ? '#fff' : 'var(--color-text-secondary)',
    }"
    v-tooltip.right="label"
  >
    <i :class="icon" class="text-lg" />
  </button>
</template>

<script setup lang="ts">
defineProps<{
  icon: string
  label: string
  active: boolean
}>()
</script>
