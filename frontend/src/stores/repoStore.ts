import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  GetRepositories,
  AddRepository,
  RemoveRepository,
  GetAllStatuses,
  RefreshRepository,
  RefreshAll,
  UpdatePollInterval,
  AssignTag,
  UnassignTag,
} from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
import { monitor } from '../../wailsjs/go/models'

export type Repository = models.Repository
export type RepoStatus = monitor.RepoStatus

export const useRepoStore = defineStore('repo', () => {
  const repositories = ref<Repository[]>([])
  const statuses = ref<Record<number, RepoStatus>>({})
  const loading = ref(false)

  async function fetchRepositories() {
    try {
      repositories.value = await GetRepositories() ?? []
    } catch (e) {
      console.error('Failed to fetch repos:', e)
    }
  }

  async function fetchStatuses() {
    try {
      const raw = await GetAllStatuses()
      if (raw) {
        statuses.value = raw as unknown as Record<number, RepoStatus>
      }
    } catch (e) {
      console.error('Failed to fetch statuses:', e)
    }
  }

  async function addRepo(path: string) {
    loading.value = true
    try {
      await AddRepository(path)
      await fetchRepositories()
    } finally {
      loading.value = false
    }
  }

  async function removeRepo(id: number) {
    await RemoveRepository(id)
    await fetchRepositories()
  }

  async function refreshRepo(id: number) {
    await RefreshRepository(id)
    await fetchStatuses()
  }

  async function refreshAll() {
    await RefreshAll()
    await fetchStatuses()
  }

  async function updateInterval(id: number, seconds: number) {
    await UpdatePollInterval(id, seconds)
    await fetchRepositories()
  }

  async function assignTag(repoID: number, tagID: number) {
    await AssignTag(repoID, tagID)
    await fetchRepositories()
  }

  async function unassignTag(repoID: number, tagID: number) {
    await UnassignTag(repoID, tagID)
    await fetchRepositories()
  }

  return {
    repositories,
    statuses,
    loading,
    fetchRepositories,
    fetchStatuses,
    addRepo,
    removeRepo,
    refreshRepo,
    refreshAll,
    updateInterval,
    assignTag,
    unassignTag,
  }
})
