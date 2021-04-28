<template>
  <div>
    <div class="welcome" v-if="isAuthenticated">
      Here there, <span class="username">{{username}}</span>!
    </div>
    <IconBar @reload="reload" />
    <SimpleHobbit v-for="hobbit in hobbits" :key='hobbit.id' :hobbit="hobbit" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import SimpleHobbit from '../components/SimpleHobbit.vue'
import IconBar from '@/components/IconBar.vue'
import { Hobbit } from '@/models'

export default defineComponent({
  name: 'Overview',
  components: { SimpleHobbit, IconBar },
  created () {
    this.dispatchFetchHobbits()
  },
  computed: {
    username (): string {
      return this.$store.state.auth.username as string
    },
    isAuthenticated () {
      return this.$store.state.auth.authenticated
    },
    hobbits (): Hobbit[] {
      return this.$store.getters.getHobbits()
    },
  },
  methods: {
    dispatchFetchHobbits () {
      if (!this.$store.state.hobbits.initialLoaded) {
        this.$store.dispatch('fetchHobbits')
      }
    },
    reload () {
      this.$store.dispatch('fetchHobbits')
    },
  },
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
