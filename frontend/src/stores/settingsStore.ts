import { defineStore } from 'pinia'
import { computed } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { SetGlobalPollInterval, SetPollingEnabled } from '../../wailsjs/go/main/App'

export const useSettingsStore = defineStore('settings', () => {
  const theme = useLocalStorage('repo-mon-theme', 'neutral-carbon')
  const darkMode = useLocalStorage('repo-mon-dark-mode', true)
  const viewMode = useLocalStorage('repo-mon-view-mode', 'grid')
  const globalPollInterval = useLocalStorage('repo-mon-poll-interval', 30)
  const pollingEnabled = useLocalStorage('repo-mon-polling-enabled', true)

  const settings = computed(() => ({
    theme: theme.value,
    darkMode: darkMode.value,
    viewMode: viewMode.value,
    globalPollInterval: globalPollInterval.value,
    pollingEnabled: pollingEnabled.value,
  }))

  function applyTheme() {
    const root = document.documentElement
    root.classList.remove(
      'theme-neutral-carbon', 'theme-slate-blue', 'theme-deep-purple',
      'dark', 'light',
    )
    root.classList.add(`theme-${theme.value}`)
    root.classList.add(darkMode.value ? 'dark' : 'light')
  }

  async function updateSettings(partial: Partial<{ theme: string; darkMode: boolean; viewMode: string; globalPollInterval: number; pollingEnabled: boolean }>) {
    if (partial.theme !== undefined) theme.value = partial.theme
    if (partial.darkMode !== undefined) darkMode.value = partial.darkMode
    if (partial.viewMode !== undefined) viewMode.value = partial.viewMode
    if (partial.pollingEnabled !== undefined) {
      pollingEnabled.value = partial.pollingEnabled
      await SetPollingEnabled(partial.pollingEnabled)
    }
    if (partial.globalPollInterval !== undefined) {
      globalPollInterval.value = partial.globalPollInterval
      await SetGlobalPollInterval(partial.globalPollInterval)
    }
    applyTheme()
  }

  function init() {
    applyTheme()
    // Sync polling state to backend on startup
    SetPollingEnabled(pollingEnabled.value)
  }

  return { settings, updateSettings, init }
})
