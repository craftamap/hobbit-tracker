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
      <div v-if="!isMe" class="text-align-right">
        <UserAddIcon v-if="!follows" class="h-24 cursor-pointer" @click="follow" />
        <UserRemoveIcon v-if="follows" class="h-24 cursor-pointer" @click="unfollow" />
      </div>
      <div>Hobbits:</div>
      <SimpleHobbit v-for="hobbit in hobbitsOfUser" :key='`hobbit-${hobbit.id}`' :hobbit="hobbit" :withHeatmap=true />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import { CogIcon, UserAddIcon, UserRemoveIcon } from '@heroicons/vue/outline'

import { Hobbit, User } from '@/models/'
import SimpleHobbit from '@/components/SimpleHobbit.vue'
import { createNamespacedHelpers } from 'vuex'
import { AuthenticationState } from '@/store/modules/auth'

const { mapActions: mapUsersActions, mapGetters: mapUsersGetters } = createNamespacedHelpers('users')
const { mapActions: mapProfileActions, mapGetters: mapProfileGetters } = createNamespacedHelpers('profile')
const { mapState: mapAuthState } = createNamespacedHelpers('auth')
const { mapActions: mapHobbitsActions, mapGetters: mapHobbitsGetters } = createNamespacedHelpers('hobbits')

export default defineComponent({
  name: 'Profile',
  components: {
    SimpleHobbit,
    CogIcon,
    UserAddIcon,
    UserRemoveIcon,
  },
  created() {
    this.defferedInit()
  },
  watch: {
    $route() {
      this.defferedInit()
    },
    userId() {
      this.defferedInit()
    },
  },
  computed: {
    ...mapUsersGetters(['getUserById']),
    ...mapProfileGetters(['followsUser']),
    ...mapAuthState({
      myUserId: state => (state as AuthenticationState).userId,
    }),
    ...mapHobbitsGetters({
      _hobbitsByUser: 'getHobbitsByUser',
    }),
    isMe(): boolean {
      return !this.$route.params.profileId
    },
    user(): User | undefined {
      console.log(this.userId)
      return this.getUserById(this.userId)
    },
    userId(): number {
      if (!this.$route.params.profileId) {
        return this.myUserId as number
      }
      return Number(this.$route.params.profileId)
    },
    hobbitsOfUser(): Hobbit[] {
      console.log('hobbitsOfUser')
      return this._hobbitsByUser(this.userId)
    },
    follows(): boolean {
      return this.followsUser(this.userId)
    },
  },
  methods: {
    ...mapUsersActions({
      fetchUser: 'fetchUser',
    }),
    ...mapProfileActions({
      _fetchFollow: 'fetchFollow',
      followUser: 'followUser',
      unfollowUser: 'unfollowUser',
    }),
    ...mapHobbitsActions({
      _fetchHobbitsByUser: 'fetchHobbitsByUser',
    }),
    fetchHobbits() {
      this._fetchHobbitsByUser({ userId: this.userId })
    },
    navigateToAppPassword() {
      this.$router.push('/profile/me/apppassword')
    },
    fetchFollow() {
      this._fetchFollow({ id: this.userId })
    },
    follow() {
      this.followUser({ id: this.userId })
    },
    unfollow() {
      this.unfollowUser({ id: this.userId })
    },
    defferedInit() {
      if (this.userId) {
        this.fetchUser({ id: this.userId })
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
