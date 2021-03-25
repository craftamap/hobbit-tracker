import { InjectionKey } from 'vue'
import { createStore, Store, useStore as vuexUseStore } from 'vuex'
import { Hobbit, NumericRecord } from '../models/index'

export interface State {
  hobbits: Hobbit[];
  auth: {
    authenticated: boolean;
    username?: string;
    userId?: number;
  };
}

export const store = createStore<State>({
  state: {
    hobbits: [],
    auth: {
      authenticated: false,
      username: undefined,
      userId: undefined
    }
  },
  getters: {
    getHobbitById: (state) => (id: number): Hobbit => {
      return state.hobbits.filter((h) => h.id === id)[0]
    }
  },
  mutations: {
    setAuth (state, payload) {
      console.log(payload)
      state.auth = payload
      console.log(state)
    },
    setHobbits (state, hobbits) {
      state.hobbits = hobbits
      console.log(state)
    },
    setRecordsForHobbit (state, { hobbitId, records }: {hobbitId: number; records: NumericRecord[]}) {
      const selectedHobbit = state.hobbits.filter((h) => h.id === hobbitId)[0]
      selectedHobbit.records = records
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
    async fetchRecords ({ commit }, payload) {
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
