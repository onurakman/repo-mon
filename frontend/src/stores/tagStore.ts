import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GetTags, AddTag, RemoveTag } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
import { useRepoStore } from './repoStore'

export type Tag = models.Tag

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([])

  async function fetchTags() {
    try {
      tags.value = await GetTags() ?? []
    } catch (e) {
      console.error('Failed to fetch tags:', e)
    }
  }

  async function addTag(name: string, color: string) {
    await AddTag(name, color)
    await fetchTags()
  }

  async function removeTag(id: number) {
    await RemoveTag(id)
    await fetchTags()
    // Refresh repos so removed tag disappears from cards
    const repoStore = useRepoStore()
    await repoStore.fetchRepositories()
  }

  return { tags, fetchTags, addTag, removeTag }
})
