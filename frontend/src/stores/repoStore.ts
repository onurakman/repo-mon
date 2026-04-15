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
  UpdateSortOrder,
  AssignTag,
  UnassignTag,
  AssignTagToRepos,
} from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
import { monitor } from '../../wailsjs/go/models'

export type Repository = models.Repository
export type RepoStatus = monitor.RepoStatus

export const useRepoStore = defineStore('repo', () => {
  const repositories = ref<Repository[]>([])
  const statuses = ref<Record<number, RepoStatus>>({})
  const loading = ref(false)
  const selectedIds = ref<Set<number>>(new Set())
  const checkingIds = ref<Set<number>>(new Set())

  function setChecking(id: number) {
    checkingIds.value = new Set([...checkingIds.value, id])
  }

  function clearChecking(id: number) {
    const next = new Set(checkingIds.value)
    next.delete(id)
    checkingIds.value = next
  }

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
    selectedIds.value.delete(id)
    await fetchRepositories()
  }

  async function refreshRepo(id: number) {
    await RefreshRepository(id)
  }

  async function refreshAll() {
    await RefreshAll()
  }

  async function updateInterval(id: number, seconds: number) {
    await UpdatePollInterval(id, seconds)
    await fetchRepositories()
  }

  async function updateSortOrder() {
    const ids = repositories.value.map(r => r.ID)
    await UpdateSortOrder(ids)
  }

  async function assignTag(repoID: number, tagID: number) {
    await AssignTag(repoID, tagID)
    await fetchRepositories()
  }

  async function unassignTag(repoID: number, tagID: number) {
    await UnassignTag(repoID, tagID)
    await fetchRepositories()
  }

  async function assignTagToRepos(repoIDs: number[], tagID: number) {
    await AssignTagToRepos(repoIDs, tagID)
    await fetchRepositories()
  }

  function toggleSelect(id: number) {
    const next = new Set(selectedIds.value)
    if (next.has(id)) {
      next.delete(id)
    } else {
      next.add(id)
    }
    selectedIds.value = next
  }

  function clearSelection() {
    selectedIds.value = new Set()
  }

  function selectAll() {
    selectedIds.value = new Set(repositories.value.map(r => r.ID))
  }

  return {
    repositories,
    statuses,
    loading,
    selectedIds,
    checkingIds,
    setChecking,
    clearChecking,
    fetchRepositories,
    fetchStatuses,
    addRepo,
    removeRepo,
    refreshRepo,
    refreshAll,
    updateInterval,
    updateSortOrder,
    assignTag,
    unassignTag,
    assignTagToRepos,
    toggleSelect,
    clearSelection,
    selectAll,
  }
})
