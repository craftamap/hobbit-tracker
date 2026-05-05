import { Hobbit, NumericRecord } from '../models'
import { defineStore } from 'pinia'

export const useHobbitsStore = defineStore('hobbits', {
  state: () => ({
    hobbits: {} as { [key: number]: Hobbit },
    records: {} as { [key: number]: NumericRecord[] },
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
    getRecordsByHobbitId: (state) => (hobbitId: number): NumericRecord[] => {
      return state.records[hobbitId]
    },
    getRecordById: (state) => (hobbitId: number, recordId: number): NumericRecord | undefined => {
      return state.records[hobbitId]?.find((value: NumericRecord) => {
        return value.id === recordId
      })
    },
    getHeatmapByHobbitId: (state) => (hobbitId: number): { date: Temporal.PlainDate, count: number }[] => {
      const records = state.records[hobbitId] || [];
      const heatmap = new Map<Temporal.PlainDate, number>();
      const tz = Temporal.Now.timeZoneId();
      for (const record of records) {
        const date = Temporal.Instant.from(record.timestamp).toZonedDateTimeISO(tz).toPlainDate();
        heatmap.set(date, (heatmap.get(date) || 0) + record.value);
      }
      return Array.from(heatmap.entries()).map(([date, count]) => ({
        date, count,
      }));
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
      if (!this.records[hobbitId]) {
        return;
      }
      this.records[hobbitId] = this.records[hobbitId].filter((value: NumericRecord) => {
        return value.id !== recordId
      });
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
    async fetchRecords(hobbitId: number) {
      const res = await fetch(`/api/hobbits/${hobbitId}/records/`)
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      this.records[hobbitId] = resJson
    },
    async postRecord({ id, timestamp, value, comment }: {
      id: number;
      timestamp: Temporal.Instant;
      value: number;
      comment: string
    }) {
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
    async putRecord({ hobbitId, recordId, timestamp, value, comment }: {
      hobbitId: number;
      recordId: number;
      timestamp: Temporal.Instant;
      value: number;
      comment: string
    }) {
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
    async putHobbit({ id, name, description, image, archivedAt }: {
      id: number;
      name: string;
      description: string;
      image: string;
      archivedAt: string | null;
    }) {
      const res = await fetch(`/api/hobbits/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          description,
          image,
          archivedAt,
        }),
      })
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      const resJson = await res.json()
      this.setHobbit(resJson)
    },
    async deleteHobbit({ id }: { id: number }) {
      const res = await fetch(`/api/hobbits/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      if (!res.ok) {
        throw new Error(res.statusText)
      }
      delete this.hobbits[id]
    },
  },
})
