<template>
  <div>
    <div class="welcome">
      Here there, <span class="username">{{ username }}</span
      >!
    </div>
    <div>
      <div class="text-align-right">
        <CogIcon class="h-24 cursor-pointer" @click="navigateToAppPassword" />
      </div>
      <div>Your hobbits:</div>
      <SimpleHobbit v-for="hobbit in hobbitsOfUser" :key='hobbit.id' :hobbit="hobbit" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import { CogIcon } from '@heroicons/vue/outline'

import { Hobbit } from '@/models/'
import SimpleHobbit from '@/components/SimpleHobbit.vue'

export default defineComponent({
  name: 'Profile',
  components: {
    SimpleHobbit,
    CogIcon,
  },
  created() {
    this.dispatchFetchHobbitsByUser()
  },
  computed: {
    username(): string {
      return this.$store.state.auth.username as string
    },
    userId(): number {
      return this.$store.state.auth.userId as number
    },
    hobbitsOfUser(): Hobbit[] {
      return this.$store.getters.getHobbitsByUser(this.userId)
    },
  },
  methods: {
    dispatchFetchHobbitsByUser() {
      this.$store.dispatch('fetchHobbitsByUser', { userId: this.userId })
    },
    navigateToAppPassword() {
      this.$router.push('/profile/me/apppassword')
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