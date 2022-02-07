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
})
