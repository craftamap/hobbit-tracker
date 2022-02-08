import { User } from '@/models'
import { defineStore } from 'pinia'

export const useUsersStore = defineStore('users', {
  state: () => ({
    users: {} as { [userId: string]: User },
  }),
  getters: {
    getUserById: (state) => (id: number): User | undefined => {
      return state.users[id]
    },
  },
  actions: {
    async fetchUser({ id }: { id: number }) {
      const response = await fetch(`/api/profile/${id}/`)
      if (!response.ok) {
        throw new Error(response.statusText)
      }

      const responseJson = await response.json() as User

      if (responseJson?.id) {
        this.users[responseJson.id] = responseJson
      }
    },
  },
})
