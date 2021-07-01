<template>
  <div>
    <div class="welcome">
      <template v-if="isMe">
        Here there,
      </template>
      <template v-if="!isMe">
        This is
      </template>
      <span class="username">{{ user?.username }}</span
      >!
    </div>
    <div>
      <div v-if="isMe" class="text-align-right">
        <CogIcon class="h-24 cursor-pointer" @click="navigateToAppPassword" />
      </div>
      <div>Hobbits:</div>
      <SimpleHobbit v-for="hobbit in hobbitsOfUser" :key='hobbit.id' :hobbit="hobbit" :withHeatmap=true />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import { CogIcon } from '@heroicons/vue/outline'

import { Hobbit, User } from '@/models/'
import SimpleHobbit from '@/components/SimpleHobbit.vue'
import { createNamespacedHelpers } from 'vuex'

const { mapActions, mapGetters } = createNamespacedHelpers('users')

export default defineComponent({
  name: 'Profile',
  components: {
    SimpleHobbit,
    CogIcon,
  },
  created() {
    this.dispatchFetchHobbitsByUser()
    this.fetchUser({ id: this.userId })
  },
  watch: {
    $route() {
      this.dispatchFetchHobbitsByUser()
      this.fetchUser({ id: this.userId })
    },
  },
  computed: {
    isMe(): boolean {
      return !this.$route.params.id
    },
    ...mapGetters(['getUserById']),
    user(): User | undefined {
      console.log(this.userId)
      return this.getUserById(this.userId)
    },
    userId(): number {
      if (!this.$route.params.id) {
        return this.$store.state.auth.userId as number
      }
      return Number(this.$route.params.id)
    },
    hobbitsOfUser(): Hobbit[] {
      console.log('hobbitsOfUser')
      return this.$store.getters.getHobbitsByUser(this.userId)
    },
  },
  methods: {
    ...mapActions({
      fetchUser: 'fetchUser',
    }),
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
    color: var(--primary);
  }
}
</style>
