import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'
import { AppPassword } from '@/models'

export interface ProfileState {
  apppassword: {
    apppasswords: AppPassword[];
  };
}

export const mutations: MutationTree<ProfileState> = {
  setAppPasswords(state, { appPasswords }) {
    state.apppassword.apppasswords = appPasswords
  },
  removeAppPassword(state, { appPassword }) {
    const idx = state.apppassword.apppasswords.findIndex((ap) => {
      return ap.id === appPassword.id
    })
    state.apppassword.apppasswords.splice(idx, 1)
  },
  addAppPassword(state, { appPassword }) {
    state.apppassword.apppasswords.push(appPassword)
  },
}

export const getters: GetterTree<ProfileState, rootState> = {

}

export const actions: ActionTree<ProfileState, rootState> = {
  async fetchAppPasswords({ commit }) {
    const appPasswords: AppPassword = await fetch('/api/profile/me/apppassword').then((res) => {
      return res.json()
    })
    commit(mutations.setAppPasswords.name, { appPasswords: appPasswords })
  },
  async deleteAppPassword({ commit }, { id }) {
    const deletedAppPassword = await fetch(`/api/profile/me/apppassword/${id}`, {
      method: 'DELETE',
    }).then((res) => {
      return res.json()
    })

    commit(mutations.removeAppPassword.name, { appPassword: deletedAppPassword })
  },
  async postAppPassword({ commit }, { description }) {
    const newAppPassword: AppPassword = await fetch('/api/profile/me/apppassword/', {
      method: 'POST',
      body: JSON.stringify({ description: description }),
    }).then((res) => {
      return res.json()
    })

    commit(mutations.addAppPassword.name, { appPassword: newAppPassword })

    return newAppPassword.secret
  },
}

export const profileModule: Module<ProfileState, rootState> = {
  namespaced: true,
  state: {
    apppassword: {
      apppasswords: [],
    },
  },
  actions: actions,
  getters: getters,
  mutations: mutations,
}

export default profileModule
