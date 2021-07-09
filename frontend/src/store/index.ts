import { createStore } from 'vuex'
import { Hobbit, NumericRecord } from '../models/index'
import profile from './modules/profile'
import feed from './modules/feed'
import users from './modules/users'
import auth from './modules/auth'
import hobbits from './modules/hobbits'

export interface State {
  socket?: WebSocket;
}

export const store = createStore<State>({
  modules: {
    auth: auth,
    profile: profile,
    feed: feed,
    users: users,
    hobbits: hobbits,
  },
  state: {
    socket: undefined,
  },
  mutations: {
    setWebsocket(state, { socket }: { socket: WebSocket }) {
      state.socket = socket
    },
  },
  getters: {
  },
  actions: {
    async createWebSocketConnection({ commit, dispatch }) {
      const socket = new WebSocket(((window.location.protocol === 'https:') ? 'wss://' : 'ws://') + window.location.host + '/ws')
      socket.onmessage = (ev) => {
        const parsedEventData = JSON.parse(ev.data)
        dispatch('recieveWebSocketMessage', parsedEventData)
      }

      socket.onclose = (ev) => {
        console.debug('WebSocket close event:', ev)
        dispatch('recieveWebSocketMessage', { socket: undefined })
      }

      socket.onerror = (ev) => {
        // TODO: add better error handling
        console.debug('WebSocket: recieved error event:', ev)
      }

      commit('setWebsocket', { socket })
    },
    async recieveWebSocketMessage({ dispatch }, { typus, optional_data: optionalData }) {
      console.debug('recieved WebSocketMessage of typus', typus, 'and optional data', optionalData)
      switch (typus) {
        case 'RecordDeleted':
        case 'RecordModified':
        case 'RecordCreated':
          dispatch('hobbits/fetchRecords', optionalData?.hobbit_id)
          break
        case 'HobbitCreated':
        case 'HobbitModified':
          dispatch('hobbits/fetchHobbit', { id: optionalData?.id })
          break
        case 'HobbitDeleted':
          // TODO:
          break
      }
    },
  },
})
