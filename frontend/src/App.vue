<template>
  <div>
    <div id="nav">
      <router-link to="/">Home</router-link> |
      <router-link to="/overview">Overview</router-link> |
      <template v-if="isAuthenticated">
        <router-link to="/profile/me">Profile</router-link> |
        <a href="/auth/logout">Logout</a>
      </template>
      <router-link v-if="!isAuthenticated" to="/login">Login</router-link>
    </div>
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
    <div>
      <div id="dialog-target"></div>
    </div>
  </div>
</template>

<script lang="ts">
import '@fontsource/lato/index.css' /* weight 400 */
import { defineComponent } from 'vue'
import DialogWrapper from './components/DialogWrapper.vue'
import { storeToRefs } from 'pinia'
import { useSocketStore } from './store/socket'
import { useAuthStore } from './store/auth'

export default defineComponent({
  components: {
    DialogWrapper,
  },
  setup() {
    const auth = useAuthStore()
    const socket = useSocketStore()

    const { authenticated: isAuthenticated } = storeToRefs(auth)

    auth.fetchAuthDetails()
    socket.createWebsocketConnection()

    return {
      isAuthenticated,
    }
  },
})
</script>

<style>
:root {
  --primary: #076470; /* ming */
  --primary-dark: #1a4b52ff;/* midnight-green-eagle-green */

  --primary-text: #111d1f;
  --secondary-text: #444444;
  --background: #ffffff;

  @media (prefers-color-scheme: dark) {
    --primary: #37A2B0;
    --primary-text: #ffffff;
    --secondary-text: #eeeeee;
    --background: #111d1f;
  }
}

@media (prefers-color-scheme: dark) {
  body {
    background: var(--background);
  }
}

body {
  margin: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

#app {
  font-family: "Lato", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--primary-text);
  margin: 0 auto;
  max-width: 1000px;
}

#nav {
  padding: 30px;

  a,
  a:link,
  :hover,
  a:active,
  a:visited {
    font-weight: bold;
    color: var(--primary-text);

    &.router-link-exact-active {
      color: var(--primary);
    }
  }
}

a {
  text-decoration: none;

  &:link,
  &:hover,
  &:active,
  &:visited {
    color: inherit;
  }
}

.h-16 {
  height: 16px;
}

.h-20 {
  height: 20px;
}

.h-24 {
  height: 24px;
}

.cursor-pointer {
  cursor: pointer;
}

.text-align-right {
  text-align: right;
}
</style>
