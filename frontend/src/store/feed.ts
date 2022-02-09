import { FeedEvent } from '@/models'
import { defineStore } from 'pinia'

export const useFeedStore = defineStore('feed', {
  state: () => ({
    feedEvents: [] as FeedEvent[],
  }),
  actions: {
    async fetchFeed() {
      await fetch('/api/profile/me/feed')
        .then(res => {
          return res.json()
        }).then(json => {
          this.feedEvents = json
        })
    },
  },
})
