import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    hobbit: []
  },
  mutations: {
    setHobbits (state, hobbits) {
      state.hobbits = hobbits
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
    }
  }
})
