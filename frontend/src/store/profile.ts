import { AppPassword, User } from '@/models'
import { defineStore } from 'pinia'

export const useProfileStore = defineStore('profile', {
  state: () => ({
    follows: [] as User[],
  }),
  getters: {
    followsUser: (state) => {
      return (id: number) => {
        const idx = state.follows.findIndex((u) => {
          return u.id === id
        })
        return idx !== -1
      }
    },
  },
  actions: {
    addFollow({ user }: { user: User }) {
      this.follows.push(user)
    },
    removeFollow({ user }: { user: User }) {
      const idx = this.follows.findIndex((u) => {
        return u.id === user.id
      })

      if (idx !== -1) {
        this.follows.splice(idx, 1)
      }
    },
    async fetchFollow({ id }: { id: number }) {
      const response = await fetch(`/api/profile/${id}/follow`)
      if (!response.ok) {
        throw new Error(response.statusText)
      }
      const content: {
        follows: boolean;
        user: User;
      } = await response.json()
      if (content.follows) {
        this.addFollow({ user: content.user })
      }
    },
    async followUser({ id }: { id: number }) {
      const response = await fetch(`/api/profile/${id}/follow`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: '{}',
      })
      if (!response.ok) {
        throw new Error(response.statusText)
      }
      this.addFollow({
        user: await response.json(),
      })
    },
    async unfollowUser({ id }: { id: number }) {
      const response = await fetch(`/api/profile/${id}/follow`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: '{}',
      })
      if (!response.ok) {
        throw new Error(response.statusText)
      }
      this.removeFollow({
        user: await response.json(),
      })
    },
  },
})

export const useAppPasswordStore = defineStore('apppassword', {
  state: () => ({
    appPasswords: [] as AppPassword[],
  }),
  actions: {
    removeAppPassword({ appPassword }: {appPassword: AppPassword}) {
      const idx = this.appPasswords.findIndex((ap) => {
        return ap.id === appPassword.id
      })
      this.appPasswords.splice(idx, 1)
    },
    addAppPassword({ appPassword }: {appPassword: AppPassword}) {
      this.appPasswords.push(appPassword)
    },
    async fetchAppPasswords() {
      const appPasswords: AppPassword[] = await fetch('/api/profile/me/apppassword').then((res) => {
        return res.json()
      })
      this.appPasswords = appPasswords
    },
    async deleteAppPassword({ id }: { id: string }) {
      const deletedAppPassword = await fetch(`/api/profile/me/apppassword/${id}`, {
        method: 'DELETE',
      }).then((res) => {
        return res.json()
      })

      this.removeAppPassword({ appPassword: deletedAppPassword })
    },
    async postAppPassword({ description }: { description: string }) {
      const newAppPassword: AppPassword = await fetch('/api/profile/me/apppassword/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ description: description }),
      }).then((res) => {
        return res.json()
      })

      this.addAppPassword({ appPassword: newAppPassword })

      return newAppPassword.secret
    },
  },
})
