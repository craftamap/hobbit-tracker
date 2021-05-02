import { createStore } from 'vuex'
import { Hobbit, NumericRecord } from '../models/index'

export interface State {
  hobbits: {
    hobbits: {[key: number]: Hobbit};
    initialLoaded: boolean;
  };
  auth: {
    authenticated: boolean;
    username?: string;
    userId?: number;
  };
}

export const store = createStore<State>({
  state: {
    hobbits: {
      hobbits: {},
      initialLoaded: false,
    },
    auth: {
      authenticated: false,
      username: undefined,
      userId: undefined,
    },
  },
  getters: {
    getHobbits: (state) => (): Hobbit[] => {
      return Object.values(state.hobbits.hobbits).sort((a, b) => {
        return a.id - b.id
      })
    },
    getHobbitById: (state) => (id: number): Hobbit => {
      return state.hobbits.hobbits[id]
    },
    getHobbitsByUser: (state) => (userId: number): Hobbit[] => {
      return Object.values(state.hobbits.hobbits).filter((value) => {
        return value.user.id === userId
      })
    },
    getRecordById: (_, getters) => (id: number, recordId: number): NumericRecord => {
      console.log('getRecordById - hobbitById:', getters.getHobbitById(id))
      return getters.getHobbitById(id)?.records?.find((value: NumericRecord) => {
        return value.id === recordId
      })
    },
  },
  mutations: {
    setAuth(state, payload) {
      console.log(payload)
      state.auth = payload
      console.log(state)
    },
    setInitialLoaded(state, { load }) {
      state.hobbits.initialLoaded = load
    },
    setHobbits(state, hobbits: Hobbit[]) {
      state.hobbits.hobbits = Object.assign({}, state.hobbits.hobbits, ...hobbits.map((x: Hobbit) => ({ [x.id]: Object.assign({}, state.hobbits.hobbits[x.id], x) })))
      console.log(state)
    },
    setHobbit(state, hobbit: Hobbit) {
      state.hobbits.hobbits[hobbit.id] = Object.assign({}, state.hobbits.hobbits[hobbit.id], hobbit)
    },
    setRecordsForHobbit(state, { hobbitId, records }: {hobbitId: number; records: NumericRecord[]}) {
      const selectedHobbit = state.hobbits.hobbits[hobbitId]
      selectedHobbit.records = records
      console.log(state)
    },
    setRecordsForHeatmapForHobbit(state, { hobbitId, records }: {hobbitId: number; records: NumericRecord[]}) {
      const selectedHobbit = state.hobbits.hobbits[hobbitId]
      console.log('selectedHobbit', selectedHobbit.id)
      selectedHobbit.heatmap = records
      console.log(state)
    },
  },
  actions: {
    async fetchHobbitsByUser({ commit }, { userId }) {
      // TODO: Add endpoint for this
      console.log(userId)
      await fetch('/api/profile/me/hobbits/')
        .then(res => {
          return res.json()
        }).then(json => {
          commit('setHobbits', json)
        })
    },
    async fetchHobbits({ commit }) {
      await fetch('/api/hobbits/')
        .then(res => {
          return res.json()
        }).then(json => {
          commit('setHobbits', json)
          commit('setInitialLoaded', { load: true })
        })
    },
    async fetchHobbit({ commit }, { id }) {
      await fetch(`/api/hobbits/${id}`)
        .then(res => {
          return res.json()
        }).then(json => {
          commit('setHobbit', json)
        })
    },
    async fetchAuth({ commit }) {
      await fetch('/api/auth')
        .then(res => {
          return res.json()
        }).then(json => {
          commit('setAuth', json)
        })
    },
    async fetchHeatmapData({ commit }, payload) {
      await fetch(`/api/hobbits/${payload}/records/heatmap`)
        .then(res => {
          return res.json()
        }).then(json => {
          return commit('setRecordsForHeatmapForHobbit', { hobbitId: payload, records: json })
        })
    },
    async fetchRecords({ commit }, payload) {
      console.log('fetchRecords')
      return fetch(`/api/hobbits/${payload}/records/`)
        .then(res => {
          return res.json()
        }).then(json => {
          return commit('setRecordsForHobbit', { hobbitId: payload, records: json })
        })
    },
    async postRecord(_, { id, timestamp, value, comment }) {
      await fetch(`/api/hobbits/${id}/records/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          timestamp,
          value,
          comment,
        }),
      }).then(res => {
        return res.json()
      }).then(json => {
        console.log(json)
        // TODO: Store in store
      })
    },
    async putRecord(_, { id: hobbitId, recordId, timestamp, value, comment }) {
      await fetch(`/api/hobbits/${hobbitId}/records/${recordId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          timestamp,
          value,
          comment,
        }),
      }).then(res => {
        return res.json()
      }).then(json => {
        console.log(json)
        // TODO: Store in store
      })
    },
    async postHobbit({ commit }, { name, description, image }) {
      await fetch('/api/hobbits/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          description,
          image,
        }),
      }).then(res => {
        return res.json()
      }).then(json => {
        console.log(json)
        commit('setHobbit', json)
      })
    },
    async putHobbit({ commit }, { id, name, description, image }) {
      await fetch(`/api/hobbits/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          description,
          image,
        }),
      }).then(res => {
        return res.json()
      }).then(json => {
        console.log(json)
        commit('setHobbit', json)
      })
    },
  },
})
