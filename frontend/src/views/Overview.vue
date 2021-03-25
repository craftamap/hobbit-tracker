<template>
  <div>
    <div class="welcome" v-if="isAuthenticated">
      Here there, <span class="username">{{username}}</span>!
    </div>
    <SimpleHobbit v-for="hobbit in $store.state.hobbits" :key='hobbit.id' :hobbit="hobbit" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import SimpleHobbit from '../components/SimpleHobbit.vue'

export default defineComponent({
  name: 'Overview',
  components: { SimpleHobbit },
  created () {
    this.dispatchFetchHobbits()
  },
  computed: {
    username (): string {
      return this.$store.state.auth.username as string
    },
    isAuthenticated () {
      return this.$store.state.auth.authenticated
    }
  },
  methods: {
    dispatchFetchHobbits () {
      this.$store.dispatch('fetchHobbits')
    }
  }
})
</script>

<style lang="scss" scoped>
.welcome {
  font-size: 16pt;
  text-align: center;
  margin-bottom: 2rem;
  color: darkgray;

  .username {
    color: var(--ming);
  }
}
</style>
