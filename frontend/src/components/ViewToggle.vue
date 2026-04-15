<template>
  <div class="flex items-center rounded-md overflow-hidden border h-7" :style="{ borderColor: 'var(--color-border)' }">
    <button
      class="w-7 h-full flex items-center justify-center cursor-pointer transition-colors"
      :style="{
        backgroundColor: modelValue === 'grid' ? 'var(--color-primary)' : 'transparent',
        color: modelValue === 'grid' ? '#fff' : 'var(--color-text-secondary)',
      }"
      @click="$emit('update:modelValue', 'grid')"
      v-tooltip.bottom="'Grid view'"
    >
      <i class="pi pi-th-large text-xs" />
    </button>
    <button
      class="w-7 h-full flex items-center justify-center cursor-pointer transition-colors"
      :style="{
        backgroundColor: modelValue === 'list' ? 'var(--color-primary)' : 'transparent',
        color: modelValue === 'list' ? '#fff' : 'var(--color-text-secondary)',
      }"
      @click="$emit('update:modelValue', 'list')"
      v-tooltip.bottom="'List view'"
    >
      <i class="pi pi-list text-xs" />
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: 'grid' | 'list'
}>()

defineEmits<{
  'update:modelValue': [value: 'grid' | 'list']
}>()
</script>
