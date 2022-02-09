import { defineStore } from 'pinia'

export const useSocketStore = defineStore('socket', {
  state: () => {
    return {
      socket: (null as WebSocket | null),
    }
  },
  actions: {
    createWebsocketConnection() {
      const socket = new WebSocket(((window.location.protocol === 'https:') ? 'wss://' : 'ws://') + window.location.host + '/ws')
      socket.onmessage = (ev) => {
        const parsedEventData = JSON.parse(ev.data)
        this.handleWebSocketMessage(parsedEventData)
      }

      socket.onclose = (ev) => {
        console.debug('WebSocket close event:', ev)
        this.socket = null
      }

      socket.onerror = (ev) => {
        // TODO: add better error handling
        console.debug('WebSocket: received error event:', ev)
      }

      this.socket = socket
    },
    handleWebSocketMessage({ typus, optional_data: optionalData }: { typus: string; optional_data: Record<any, any> }) {
      console.debug('recieved WebSocketMessage of typus', typus, 'and optional data', optionalData)
      switch (typus) {
        case 'RecordDeleted':
        case 'RecordModified':
        case 'RecordCreated':
          // FIXME
          // dispatch('hobbits/fetchRecords', optionalData?.hobbit_id)
          break
        case 'HobbitCreated':
        case 'HobbitModified':
          // FIXME
          // dispatch('hobbits/fetchHobbit', { id: optionalData?.id })
          break
        case 'HobbitDeleted':
          // TODO:
          break
      }
    },
  },
})
