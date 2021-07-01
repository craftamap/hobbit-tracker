import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'
import { User } from '@/models'

export interface UsersState {
  users: User[];
}

export const mutations: MutationTree<UsersState> = {
  addUser(state, { user }: {user: User}): number {
    const index = state.users.findIndex((userToCompare) => {
      return user.id === userToCompare.id
    })
    if (index === -1) {
      return state.users.push(user)
    }
    state.users[index] = user
    return index
  },
}

export const getters: GetterTree<UsersState, rootState> = {
  getUserById: (state) => (id: number): User | undefined => {
    return state.users.find((value) => {
      console.log('value.id', value.id)
      return value.id === id
    })
  },
}

export const actions: ActionTree<UsersState, rootState> = {
  async fetchUser({ commit }, { id }: {id: number}) {
    const response = await fetch(`/api/profile/${id}`)
    if (!response.ok) {
      throw new Error(response.statusText)
    }

    const responseJson = await response.json()

    commit('addUser', { user: responseJson })
  },
}

export const usersModule: Module<UsersState, rootState> = {
  namespaced: true,
  state: {
    users: [],
  },
  actions: actions,
  getters: getters,
  mutations: mutations,
}

export default usersModule
