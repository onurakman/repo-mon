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

    <div class="flex-1" />

    <SidebarIcon
      icon="pi pi-cog"
      label="Settings"
      :active="route.name === 'settings'"
      @click="router.push({ name: 'settings' })"
    />
  </nav>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import SidebarIcon from './SidebarIcon.vue'

const router = useRouter()
const route = useRoute()

const navItems = [
  { route: 'dashboard', icon: 'pi pi-objects-column', label: 'Dashboard' },
  { route: 'tags', icon: 'pi pi-tags', label: 'Tags' },
]
</script>
