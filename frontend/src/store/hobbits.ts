import { Hobbit, NumericRecord } from '../models'
import { defineStore } from 'pinia'

export const useHobbitsStore = defineStore('hobbits', {
  state: () => ({
    hobbits: {} as { [key: number]: Hobbit },
    initialLoaded: false,
  }),
  getters: {
    getHobbits: (state) => {
      return Object.values(state.hobbits).sort((a, b) => {
        return a.id - b.id
      })
    },
    getHobbitById: (state) => (id: number): Hobbit => {
      return state.hobbits[id]
    },
    getHobbitsByUser: (state) => (userId: number): Hobbit[] => {
      return Object.values(state.hobbits).filter((value) => {
        return value.user.id === userId
      })
    },
    getRecordById: (state) => (hobbitId: number, recordId: number): NumericRecord | undefined => {
      return state.hobbits[hobbitId]?.records?.find((value: NumericRecord) => {
        return value.id === recordId
      })
    },
  },
  actions: {
    setHobbit(hobbit: Hobbit) {
      console.log('setHobbit', hobbit)
      this.hobbits[hobbit.id] = Object.assign({}, this.hobbits[hobbit.id], hobbit)
    },
    setHobbits(hobbits: Hobbit[]) {
      this.hobbits = Object.assign({}, this.hobbits, ...hobbits.map((x: Hobbit) => ({ [x.id]: Object.assign({}, this.hobbits[x.id], x) })))
    },
    deleteRecordForHobbit({ hobbitId, recordId }: { hobbitId: number; recordId: number }) {
      const selectedHobbit = this.hobbits[hobbitId]
      selectedHobbit.records = selectedHobbit.records.filter((record) => {
        return record.id !== recordId
      })
    },
    async fetchHobbitsByUser(userId: number | string = 'me') {
      const res = await fetch(`/api/profile/${userId}/hobbits`)
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      this.setHobbits(resJson)
    },
    async fetchHobbits() {
      const res = await fetch('/api/hobbits/')
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      console.log('fetchHobbits', resJson)
      this.setHobbits(resJson)
      this.initialLoaded = true
    },
    async fetchHobbit(hobbitId: number) {
      const res = await fetch(`/api/hobbits/${hobbitId}/`)
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      this.setHobbit(resJson)
    },
    async fetchHeatmapData(hobbitId: number) {
      const res = await fetch(`/api/hobbits/${hobbitId}/records/heatmap`)
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      if (this.hobbits[hobbitId]) {
        this.hobbits[hobbitId].heatmap = resJson
      }
    },
    async fetchRecords(hobbitId: number) {
      const res = await fetch(`/api/hobbits/${hobbitId}/records/`)
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      if (this.hobbits[hobbitId]) {
        this.hobbits[hobbitId].records = resJson
      }
    },
    async postRecord({ id, timestamp, value, comment }: { id: number; timestamp: Date; value: number; comment: string }) {
      const res = await fetch(`/api/hobbits/${id}/records/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          timestamp,
          value,
          comment,
        }),
      })
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      // TODO: Put in store
    },
    async putRecord({ hobbitId, recordId, timestamp, value, comment }:
      { hobbitId: number; recordId: number; timestamp: Date; value: number; comment: string }) {
      const res = await fetch(`/api/hobbits/${hobbitId}/records/${recordId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          timestamp,
          value,
          comment,
        }),
      })
      if (!res.ok) {
        throw new Error(res.statusText)
      }
    },
    async deleteRecord({ hobbitId, recordId }: { hobbitId: number; recordId: number }) {
      const res = await fetch(`/api/hobbits/${hobbitId}/records/${recordId}`, {
        method: 'DELETE',
      })
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      this.deleteRecordForHobbit({ hobbitId, recordId })
    },
    async postHobbit({ name, description, image }: { name: string; description: string; image: string }) {
      const res = await fetch('/api/hobbits/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          description,
          image,
        }),
      })
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      this.setHobbit(resJson)
    },
    async putHobbit({ id, name, description, image }: { id: number; name: string; description: string; image: string }) {
      const res = await fetch(`/api/hobbits/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          description,
          image,
        }),
      })
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      this.setHobbit(resJson)
    },
  },
})
