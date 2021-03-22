import { InjectionKey } from 'vue'
import { createStore, Store, useStore as vuexUseStore } from 'vuex'
import { Hobbit, NumericRecord } from '../models/index'

export interface State {
  hobbits: Hobbit[];
  auth: {
    authenticated: boolean;
  };
}

// eslint-disable-next-line symbol-description
export const key: InjectionKey<Store<State>> = Symbol()

export const store = createStore<State>({
  state: {
    hobbits: [],
    auth: {
      authenticated: false
    }
  },
  getters: {
    getHobbitById: (state) => (id: number): Hobbit => {
      return state.hobbits.filter((h) => h.id === id)[0]
    }
  },
  mutations: {
    setAuth (state, payload) {
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
    }
  }
})

export function useStore (): Store<State> {
  return vuexUseStore(key)
}
