import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'

export interface AuthenticationState {
    authenticated: boolean;
    username?: string;
    userId?: number;
}

export const mutations: MutationTree<AuthenticationState> = {
  setAuthenticationDetails(state, payload: AuthenticationState) {
    state.authenticated = payload.authenticated
    state.username = payload.username
    state.userId = payload.userId
  },
}

export const getters: GetterTree<AuthenticationState, rootState> = {
  isAuthenticated: (state) => (): boolean => {
    return state.authenticated
  },
}

export const actions: ActionTree<AuthenticationState, rootState> = {
  async extractAuthenticationDetails({ commit }) {
    const text = document?.querySelector('#data')?.textContent
    if (!text) {
      return
    }
    const parsed = JSON.parse(text)

    commit('setAuthenticationDetails', parsed.auth)
    console.log('parsed', parsed)
  },
  async fetchAuthenticationDetails({ commit }) {
    const response = await fetch('/api/auth')
    if (!response.ok) {
      throw new Error('TODO: message here')
    }
    const responseJson = await response.json()
    commit('setAuthenticationDetails', responseJson)
  },
}

export const feedModule: Module<AuthenticationState, rootState> = {
  namespaced: true,
  state: {
    authenticated: false,
  },
  actions: actions,
  getters: getters,
  mutations: mutations,
}

export default feedModule
