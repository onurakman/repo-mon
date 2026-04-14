import { watch } from 'vue'
import { useIntervalFn } from '@vueuse/core'
import { useRepoStore } from '../stores/repoStore'
import { useSettingsStore } from '../stores/settingsStore'

export function usePolling(intervalMs: number = 5000) {
  const repoStore = useRepoStore()
  const settingsStore = useSettingsStore()

  const { pause, resume, isActive } = useIntervalFn(
    async () => {
      await repoStore.fetchStatuses()
    },
    intervalMs,
    { immediate: true },
  )

  // React to polling toggle
  watch(() => settingsStore.settings.pollingEnabled, (enabled) => {
    if (enabled) {
      resume()
    } else {
      pause()
    }
  }, { immediate: true })

  return { pause, resume, isActive }
}
