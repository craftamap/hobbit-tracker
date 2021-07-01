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
  addFollow(state, { user }: {user: User}) {
    state.follows.push(user)
  },
  removeFollow(state, { user }: {user: User}) {
    const idx = state.follows.findIndex((u) => {
      return u.id === user.id
    })

    if (idx !== -1) {
      state.follows.splice(idx, 1)
    }
  },
}

export const getters: GetterTree<ProfileState, rootState> = {
  followsUser(state) {
    return (id: number): boolean => {
      const idx = state.follows.findIndex((u) => {
        return u.id === id
      })
      return idx !== -1
    }
  },
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
  async fetchFollow({ commit }, { id }: {id: number}) {
    const response = await fetch(`/api/profile/${id}/follow`)
    if (!response.ok) {
      throw new Error(response.statusText)
    }
    const content: {
      follows: boolean;
      user: User;
    } = await response.json()
    if (content.follows) {
      commit('addFollow', { user: content.user })
    }
  },
  async followUser({ commit }, { id }: {id: number}) {
    const response = await fetch(`/api/profile/${id}/follow`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: '{}',
    })
    if (!response.ok) {
      throw new Error(response.statusText)
    }
    commit('addFollow', {
      user: await response.json(),
    })
  },
  async unfollowUser({ commit }, { id }: {id: number}) {
    const response = await fetch(`/api/profile/${id}/follow`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: '{}',
    })
    if (!response.ok) {
      throw new Error(response.statusText)
    }
    commit('removeFollow', {
      user: await response.json(),
    })
  },
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
