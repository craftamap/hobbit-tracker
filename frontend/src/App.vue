<template>
  <div>
    <div id="nav">
      <router-link to="/">Home</router-link> |
      <template v-if="isAuthenticated">
        <router-link to="/profile">Profile</router-link> |
        <a href="/auth/logout">Logout</a>
      </template>
      <router-link v-if="!isAuthenticated" to="/login">Login</router-link>
    </div>
    <router-view/>
  </div>
</template>

<script lang="ts">
import '@fontsource/lato/index.css' /* weight 400 */

import { defineComponent } from 'vue'
export default defineComponent({
  created () {
    this.dispatchFetchAuth()
  },
  computed: {
    isAuthenticated () {
      return this.$store.state.auth.authenticated
    }
  },
  methods: {
    dispatchFetchAuth () {
      this.$store.dispatch('fetchAuth')
    }
  }
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

body {
  margin: 0;
}

#app {
  font-family: 'Lato', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--dark-jungle-green);
  margin: 0 auto;
  max-width: 1000px;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: var(--dark-jungle-green);

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
</style>
