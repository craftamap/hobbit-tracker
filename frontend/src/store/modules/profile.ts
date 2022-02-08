import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'
import { AppPassword, User } from '@/models'

export interface ProfileState {
  apppassword: {
    apppasswords: AppPassword[];
  };
  follows: User[];
}

export const mutations: MutationTree<ProfileState> = {
}

export const getters: GetterTree<ProfileState, rootState> = {
}

export const actions: ActionTree<ProfileState, rootState> = {
}

export const profileModule: Module<ProfileState, rootState> = {
  namespaced: true,
  state: {
    apppassword: {
      apppasswords: [],
    },
    follows: [],
  },
  actions: actions,
  getters: getters,
  mutations: mutations,
}

export default profileModule
