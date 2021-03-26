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
      initialLoaded: false
    },
    auth: {
      authenticated: false,
      username: undefined,
      userId: undefined
    }
  },
  getters: {
    getHobbitById: (state) => (id: number): Hobbit => {
      return state.hobbits.hobbits[id]
    }
  },
  mutations: {
    setAuth (state, payload) {
      console.log(payload)
      state.auth = payload
      console.log(state)
    },
    setInitialLoaded (state, { load }) {
      state.hobbits.initialLoaded = load
    },
    setHobbits (state, hobbits: Hobbit[]) {
      state.hobbits.hobbits = Object.assign({}, ...hobbits.map((x: Hobbit) => ({ [x.id]: Object.assign({}, state.hobbits.hobbits[x.id], x) })))
      console.log(state)
    },
    setHobbit (state, hobbit: Hobbit) {
      state.hobbits.hobbits[hobbit.id] = Object.assign({}, state.hobbits.hobbits[hobbit.id], hobbit)
    },
    setRecordsForHobbit (state, { hobbitId, records }: {hobbitId: number; records: NumericRecord[]}) {
      const selectedHobbit = state.hobbits.hobbits[hobbitId]
      selectedHobbit.records = records
      console.log(state)
    },
    setRecordsForHeatmapForHobbit (state, { hobbitId, records }: {hobbitId: number; records: NumericRecord[]}) {
      const selectedHobbit = state.hobbits.hobbits[hobbitId]
      console.log('selectedHobbit', selectedHobbit.id)
      selectedHobbit.heatmap = records
      console.log(state)
    }
  },
  actions: {
    async fetchHobbits ({ commit }) {
      fetch('/api/hobbits/')
        .then(res => {
          return res.json()
        }).then(json => {
          commit('setHobbits', json)
          commit('setInitialLoaded', { load: true })
        })
    },
    async fetchHobbit ({ commit }, { id }) {
      fetch(`/api/hobbits/${id}`)
        .then(res => {
          return res.json()
        }).then(json => {
          commit('setHobbit', json)
        })
    },
    async fetchAuth ({ commit }) {
      fetch('/api/auth')
        .then(res => {
          return res.json()
        }).then(json => {
          commit('setAuth', json)
        })
    },
    async fetchHeatmapData ({ commit }, payload) {
      return fetch(`/api/hobbits/${payload}/records/heatmap`)
        .then(res => {
          return res.json()
        }).then(json => {
          return commit('setRecordsForHeatmapForHobbit', { hobbitId: payload, records: json })
        })
    },
    async fetchRecords ({ commit }, payload) {
      console.log('fetchRecords')
      return fetch(`/api/hobbits/${payload}/records/`)
        .then(res => {
          return res.json()
        }).then(json => {
          return commit('setRecordsForHobbit', { hobbitId: payload, records: json })
        })
    },
    async postRecord ({ commit }, { id, timestamp, value, comment }) {
      return fetch(`/api/hobbits/${id}/records/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          timestamp,
          value,
          comment
        })
      }).then(res => {
        return res.json()
      }).then(json => {
        console.log(json)
      })
    }
  }
})
