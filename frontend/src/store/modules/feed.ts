import { Module, ActionTree, GetterTree, MutationTree } from 'vuex'

import { State as rootState } from '@/store/index'
import { Hobbit, NumericRecord } from '@/models'

enum FeedEventTypus {
  HobbitCreated = 'HobbitCreated',
  RecordCreated = 'RecordCreated',
}

interface FeedEvent {
  FeedEventTypus: FeedEventTypus;
  CreatedAt: string;
  Payload: Hobbit | NumericRecord;
}

export interface FeedState {
  feedEvents: FeedEvent[];
}

export const mutations: MutationTree<FeedState> = {
  setFeedEvents(state, feedEvents: FeedEvent[]) {
    state.feedEvents = feedEvents
  },
}

export const getters: GetterTree<FeedState, rootState> = {

}

export const actions: ActionTree<FeedState, rootState> = {
  async fetchFeed({ commit }) {
    await fetch('/api/profile/me/feed')
      .then(res => {
        return res.json()
      }).then(json => {
        commit(mutations.setFeedEvents.name, json)
      })
  },
}

export const feedModule: Module<FeedState, rootState> = {
  namespaced: true,
  state: {
    feedEvents: [],
  },
  actions: actions,
  getters: getters,
  mutations: mutations,
}

export default feedModule
