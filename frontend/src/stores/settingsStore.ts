import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GetSettings, UpdateSettings } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'

export type UserSettings = models.UserSettings

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<UserSettings>(new models.UserSettings({
    ID: 1,
    theme: 'neutral-carbon',
    darkMode: true,
    viewMode: 'grid',
    globalPollInterval: 30,
  }))

  async function fetchSettings() {
    try {
      const s = await GetSettings()
      if (s) {
        settings.value = s
      }
    } catch (e) {
      console.error('Failed to fetch settings:', e)
    }
    applyTheme()
  }

  async function updateSettings(partial: Partial<{ theme: string; darkMode: boolean; viewMode: string; globalPollInterval: number }>) {
    Object.assign(settings.value, partial)
    applyTheme()
    try {
      await UpdateSettings(settings.value)
    } catch (e) {
      console.error('Failed to save settings:', e)
    }
  }

  function applyTheme() {
    const root = document.documentElement
    root.classList.remove(
      'theme-neutral-carbon', 'theme-slate-blue', 'theme-deep-purple',
      'dark', 'light',
    )
    root.classList.add(`theme-${settings.value.theme}`)
    root.classList.add(settings.value.darkMode ? 'dark' : 'light')
  }

  return { settings, fetchSettings, updateSettings }
})
