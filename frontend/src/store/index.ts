import { InjectionKey } from 'vue'
import { createStore, Store } from 'vuex'
import { Hobbit } from '../models/index'

export interface State {
  hobbits: Hobbit[];
  auth: object;
}

// eslint-disable-next-line symbol-description
export const key: InjectionKey<Store<State>> = Symbol()

export default createStore<State>({
  state: {
    hobbits: [],
    auth: {
      authenticated: false
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
    }
  }
})
