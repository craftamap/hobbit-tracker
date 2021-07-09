import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'
import { Hobbit, NumericRecord } from '@/models'

export interface HobbitsState {
  hobbits: { [key: number]: Hobbit };
  initialLoaded: boolean;
}

export const mutations: MutationTree<HobbitsState> = {
  setInitialLoaded(state, { load }) {
    state.initialLoaded = load
  },
  setHobbits(state, hobbits: Hobbit[]) {
    state.hobbits = Object.assign({}, state.hobbits, ...hobbits.map((x: Hobbit) => ({ [x.id]: Object.assign({}, state.hobbits[x.id], x) })))
  },
  setHobbit(state, hobbit: Hobbit) {
    state.hobbits[hobbit.id] = Object.assign({}, state.hobbits[hobbit.id], hobbit)
  },
  setRecordsForHobbit(state, { hobbitId, records }: { hobbitId: number; records: NumericRecord[] }) {
    const selectedHobbit = state.hobbits[hobbitId]
    selectedHobbit.records = records
  },
  setRecordsForHeatmapForHobbit(state, { hobbitId, records }: { hobbitId: number; records: NumericRecord[] }) {
    const selectedHobbit = state.hobbits[hobbitId]
    selectedHobbit.heatmap = records
  },
  deleteRecordForHobbit(state, { hobbitId, recordId }: { hobbitId: number; recordId: number }) {
    const selectedHobbit = state.hobbits[hobbitId]
    selectedHobbit.records = selectedHobbit.records.filter((record) => {
      return record.id !== recordId
    })
  },
}

export const getters: GetterTree<HobbitsState, rootState> = {
  getHobbits: (state) => {
    return Object.values(state.hobbits).sort((a, b) => {
      return a.id - b.id
    })
  },
  getHobbitById: (state) => (id: number): Hobbit => {
    return state.hobbits[id]
  },
  getHobbitsByUser: (state) => (userId: number): Hobbit[] => {
    return Object.values(state.hobbits).filter((value) => {
      return value.user.id === userId
    })
  },
  getRecordById: (_, getters) => (id: number, recordId: number): NumericRecord => {
    console.log('getRecordById - hobbitById:', getters.getHobbitById(id))
    return getters.getHobbitById(id)?.records?.find((value: NumericRecord) => {
      return value.id === recordId
    })
  },
}

export const actions: ActionTree<HobbitsState, rootState> = {
  async fetchHobbitsByUser({ commit }, { userId }) {
    if (!userId) {
      userId = 'me'
    }
    const res = await fetch(`/api/profile/${userId}/hobbits/`)
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    const resJson = await res.json()
    commit('setHobbits', resJson)
  },
  async fetchHobbits({ commit }) {
    const res = await fetch('/api/hobbits/')
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    const resJson = await res.json()
    commit('setHobbits', resJson)
    commit('setInitialLoaded', { load: true })
  },
  async fetchHobbit({ commit }, { id }) {
    const res = await fetch(`/api/hobbits/${id}`)
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    const resJson = await res.json()
    commit('setHobbit', resJson)
  },
  async fetchHeatmapData({ commit }, payload) {
    const res = await fetch(`/api/hobbits/${payload}/records/heatmap`)
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    const resJson = await res.json()
    return commit('setRecordsForHeatmapForHobbit', { hobbitId: payload, records: resJson })
  },
  async fetchRecords({ commit }, payload) {
    const res = await fetch(`/api/hobbits/${payload}/records/`)
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    const resJson = await res.json()
    return commit('setRecordsForHobbit', { hobbitId: payload, records: resJson })
  },
  async postRecord(_, { id, timestamp, value, comment }) {
    const res = await fetch(`/api/hobbits/${id}/records/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        timestamp,
        value,
        comment,
      }),
    })
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    // TODO: Put in store
  },
  async putRecord(_, { id: hobbitId, recordId, timestamp, value, comment }) {
    const res = await fetch(`/api/hobbits/${hobbitId}/records/${recordId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        timestamp,
        value,
        comment,
      }),
    })
    if (!res.ok) {
      throw new Error(res.statusText)
    }
  },
  async deleteRecord({ commit }, { hobbitId, recordId }) {
    const res = await fetch(`/api/hobbits/${hobbitId}/records/${recordId}`, {
      method: 'DELETE',
    })
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    commit('deleteRecordForHobbit', { hobbitId, recordId })
  },
  async postHobbit({ commit }, { name, description, image }) {
    const res = await fetch('/api/hobbits/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name,
        description,
        image,
      }),
    })
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    const resJson = await res.json()
    commit('setHobbit', resJson)
  },
  async putHobbit({ commit }, { id, name, description, image }) {
    const res = await fetch(`/api/hobbits/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name,
        description,
        image,
      }),
    })
    if (!res.ok) {
      throw new Error(res.statusText)
    }
    const resJson = await res.json()
    commit('setHobbit', resJson)
  },
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
