import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    authenticated: false,
    userId: null as number | null,
    username: '',
  }),
  actions: {
    async fetchAuthDetails() {
      const response = await fetch('/api/auth')
      if (!response.ok) {
        throw new Error('TODO: message here')
      }
      const responseJson = await response.json()
      this.$patch(responseJson)
    },
  },
})
