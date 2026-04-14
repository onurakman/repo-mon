<template>
  <nav
    class="flex flex-col items-center py-4 gap-2 w-14 shrink-0 border-r"
    :style="{
      backgroundColor: 'var(--color-sidebar)',
      borderColor: 'var(--color-border)',
    }"
  >
    <SidebarIcon
      v-for="item in navItems"
      :key="item.route"
      :icon="item.icon"
      :label="item.label"
      :active="route.name === item.route"
      @click="router.push({ name: item.route })"
    />

    <!-- Separator -->
    <div class="w-6 border-t my-1" :style="{ borderColor: 'var(--color-border)' }" />

    <!-- Add repo button -->
    <SidebarIcon
      icon="pi pi-plus"
      label="Add Repository"
      :active="false"
      @click="$emit('add-repo')"
    />
  </nav>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import SidebarIcon from './SidebarIcon.vue'

const router = useRouter()
const route = useRoute()

defineEmits<{ 'add-repo': [] }>()

const navItems = [
  { route: 'dashboard', icon: 'pi pi-objects-column', label: 'Dashboard' },
  { route: 'tags', icon: 'pi pi-tags', label: 'Tags' },
  { route: 'settings', icon: 'pi pi-cog', label: 'Settings' },
]
</script>
