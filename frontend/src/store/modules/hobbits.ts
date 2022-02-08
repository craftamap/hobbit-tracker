import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'
import { Hobbit, NumericRecord } from '@/models'

export interface HobbitsState {
  hobbits: { [key: number]: Hobbit };
  initialLoaded: boolean;
}

export const mutations: MutationTree<HobbitsState> = {
}

export const getters: GetterTree<HobbitsState, rootState> = {
}

export const actions: ActionTree<HobbitsState, rootState> = {
}

export const profileModule: Module<HobbitsState, rootState> = {
  namespaced: true,
  state: {
    hobbits: {},
    initialLoaded: false,
  },
  actions: actions,
  getters: getters,
  mutations: mutations,
}

export default profileModule
