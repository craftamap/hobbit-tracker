<template>
  <div>
    <div id="nav">
      <router-link to="/">Home</router-link> |
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
      <div id="dialog-target">
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import '@fontsource/lato/index.css' /* weight 400 */
import { defineComponent } from 'vue'
import DialogWrapper from '@/components/DialogWrapper.vue'

export default defineComponent({
  created() {
    this.dispatchFetchAuth()
  },
  components: {
    DialogWrapper,
  },
  computed: {
    isAuthenticated() {
      return this.$store.state.auth.authenticated
    },
  },
  methods: {
    dispatchFetchAuth() {
      this.$store.dispatch('fetchAuth')
    },
  },
})
</script>

<style lang="scss">
:root {
  --sky-blue-crayola: #66d2e3ff;
  --cadet-blue: #49a0abff;
  --ming: #076470ff;
  --midnight-green-eagle-green: #1a4b52ff;
  --dark-jungle-green: #111d1fff;
}

@media (prefers-color-scheme: dark) {
  body {
    background: var(--dark-jungle-green);
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
  color: var(--dark-jungle-green);
  margin: 0 auto;
  max-width: 1000px;

  @media (prefers-color-scheme: dark) {
    color: white;
  }
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: var(--dark-jungle-green);

    @media (prefers-color-scheme: dark) {
      color: white;
    }

    &.router-link-exact-active {
      color: var(--ming);
    }
  }
}

a {
  text-decoration: none;
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

.cursor-pointer {
  cursor: pointer;
}
</style>
