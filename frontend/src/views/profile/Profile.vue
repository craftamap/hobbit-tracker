<template>
  <div>
    <div class="welcome">
      <template v-if="isMe">Here there, </template>
      <template v-if="!isMe">This is </template>
      <span class="username">{{ user?.username }}</span>!
    </div>
    <div>
      <div v-if="isMe" class="text-align-right">
        <CogIcon class="h-24 cursor-pointer" @click="navigateToAppPassword" />
      </div>
      <div v-if="!isMe" class="text-align-right">
        <UserAddIcon v-if="!follows" class="h-24 cursor-pointer" @click="follow" />
        <UserRemoveIcon v-if="follows" class="h-24 cursor-pointer" @click="unfollow" />
      </div>
      <div>Hobbits:</div>
      <SimpleHobbit
        v-for="hobbit in hobbitsOfUser"
        :key="`hobbit-${hobbit.id}`"
        :hobbit="hobbit"
        :withHeatmap="true"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import { CogIcon, UserAddIcon, UserRemoveIcon } from '@heroicons/vue/outline'

import { Hobbit, User } from '@/models/'
import SimpleHobbit from '@/components/SimpleHobbit.vue'
import { mapActions, mapState } from 'pinia'
import { useAuthStore } from '@/store/auth'
import { useProfileStore } from '@/store/profile'
import { useHobbitsStore } from '@/store/hobbits'
import { useUsersStore } from '@/store/users'

export default defineComponent({
  name: 'ProfileView',
  components: {
    SimpleHobbit,
    CogIcon,
    UserAddIcon,
    UserRemoveIcon,
  },
  created() {
    this.deferredInit()
  },
  watch: {
    $route() {
      this.deferredInit()
    },
    userId() {
      this.deferredInit()
    },
  },
  computed: {
    ...mapState(useUsersStore, ['getUserById']),
    ...mapState(useAuthStore, { myUserId: 'userId' }),
    ...mapState(useProfileStore, ['followsUser']),
    ...mapState(useHobbitsStore, {
      _hobbitsByUser: 'getHobbitsByUser',
    }),
    isMe(): boolean {
      return !this.$route.params.profileId
    },
    user(): User | undefined {
      console.log(this.userId)
      if (this.userId) {
        return this.getUserById(this.userId)
      }
      return undefined
    },
    userId(): number | null {
      if (!this.$route.params.profileId) {
        return this.myUserId
      }
      return Number(this.$route.params.profileId)
    },
    hobbitsOfUser(): Hobbit[] {
      console.log('hobbitsOfUser')
      if (this.userId) {
        return this._hobbitsByUser(this.userId)
      }
      return []
    },
    follows(): boolean {
      if (this.userId) {
        return this.followsUser(this.userId)
      }
      return false
    },
  },
  methods: {
    ...mapActions(useUsersStore, {
      fetchUser: 'fetchUser',
    }),
    ...mapActions(useProfileStore, {
      _fetchFollow: 'fetchFollow',
      followUser: 'followUser',
      unfollowUser: 'unfollowUser',
    }),
    ...mapActions(useHobbitsStore, {
      _fetchHobbitsByUser: 'fetchHobbitsByUser',
    }),
    fetchHobbits() {
      if (this.userId) {
        this._fetchHobbitsByUser(this.userId)
      }
    },
    navigateToAppPassword() {
      this.$router.push('/profile/me/apppassword')
    },
    fetchFollow() {
      if (this.userId) {
        this._fetchFollow({ id: this.userId })
      }
    },
    follow() {
      if (this.userId) {
        this.followUser({ id: this.userId })
      }
    },
    unfollow() {
      if (this.userId) {
        this.unfollowUser({ id: this.userId })
      }
    },
    deferredInit() {
      if (this.userId) {
        this.fetchUser({ id: this.userId })
        this.fetchFollow()
        this.fetchHobbits()
      }
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
