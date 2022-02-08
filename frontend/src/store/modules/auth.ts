import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'

export interface AuthenticationState {
    authenticated: boolean;
    username?: string;
    userId?: number;
}

export const mutations: MutationTree<AuthenticationState> = {
}

export const getters: GetterTree<AuthenticationState, rootState> = {
}

export const actions: ActionTree<AuthenticationState, rootState> = {
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
