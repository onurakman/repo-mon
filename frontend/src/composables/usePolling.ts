import { useIntervalFn } from '@vueuse/core'
import { useRepoStore } from '../stores/repoStore'

export function usePolling(intervalMs: number = 5000) {
  const repoStore = useRepoStore()

  const { pause, resume, isActive } = useIntervalFn(
    async () => {
      await repoStore.fetchStatuses()
    },
    intervalMs,
    { immediate: true },
  )

  return { pause, resume, isActive }
}
