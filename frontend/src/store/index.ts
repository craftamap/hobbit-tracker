import { createStore } from 'vuex'
import feed from './modules/feed'
import users from './modules/users'
import hobbits from './modules/hobbits'

export interface State {
  socket?: WebSocket;
}

export const store = createStore<State>({
  modules: {
    feed: feed,
    users: users,
    hobbits: hobbits,
  },
})
