import { watch } from 'vue'
import { useIntervalFn } from '@vueuse/core'
import { useRepoStore } from '../stores/repoStore'
import { useSettingsStore } from '../stores/settingsStore'

export function usePolling(intervalMs: number = 5000, startImmediately = true) {
  const repoStore = useRepoStore()
  const settingsStore = useSettingsStore()

  const { pause, resume, isActive } = useIntervalFn(
    async () => {
      await repoStore.fetchStatuses()
    },
    intervalMs,
    { immediate: startImmediately },
  )

  // React to polling toggle (not immediate — App.vue controls initial start)
  watch(() => settingsStore.settings.pollingEnabled, (enabled) => {
    if (enabled) {
      resume()
    } else {
      pause()
    }
  })

  return { pause, resume, isActive }
}
