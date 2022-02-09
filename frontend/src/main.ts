import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { store } from './store'
import './registerServiceWorker'
import { createPinia } from 'pinia'

const pinia = createPinia()

createApp(App).use(store).use(pinia).use(router).mount('#app')
