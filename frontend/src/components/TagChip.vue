<template>
  <button
    class="px-3 py-1 rounded-full text-xs font-medium transition-opacity cursor-pointer"
    :class="active ? 'opacity-100' : 'opacity-50 hover:opacity-75'"
    :style="{ backgroundColor: color, color: '#fff' }"
    @click="$emit('toggle')"
  >
    {{ name }}
  </button>
</template>

<script setup lang="ts">
defineProps<{
  name: string
  color: string
  active: boolean
}>()

defineEmits<{
  toggle: []
}>()
</script>
