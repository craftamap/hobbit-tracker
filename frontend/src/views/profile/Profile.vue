<template>
  <div>
    <div class="welcome">
      <template v-if="isMe">Here there, </template>
      <template v-if="!isMe">This is </template>
      <span class="username">{{ user?.username }}</span>!
    </div>
    <div>
      <IconBar>
        <template v-slot:right>
          <span v-if="isMe">
            <CogIcon class="h-24 cursor-pointer" @click="goToAppPassword" />
          </span>
          <span v-if="!isMe">
            <UserAddIcon v-if="!follows" class="h-24 cursor-pointer" @click="follow" />
          </span>
          <span v-if="!isMe">
            <UserRemoveIcon v-if="follows" class="h-24 cursor-pointer" @click="unfollow" />
          </span>
        </template>
      </IconBar>
    </div>
    <div>Hobbits:</div>
    <SimpleHobbit v-for="hobbit in hobbitsOfUser" :key="`hobbit-${hobbit.id}`" :hobbit="hobbit" :withHeatmap="true" />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, watch } from 'vue'

import { CogIcon, UserPlusIcon as UserAddIcon, UserMinusIcon as UserRemoveIcon } from '@heroicons/vue/24/outline'

import SimpleHobbit from '../../components/SimpleHobbit.vue'
import IconBar from '../../components/IconBar.vue'
import { useAuthStore } from '../../store/auth'
import { useProfileStore } from '../../store/profile'
import { useHobbitsStore } from '../../store/hobbits'
import { useUsersStore } from '../../store/users'
import { useRoute, useRouter } from 'vue-router'

export default defineComponent({
  name: 'ProfileView',
  components: {
    SimpleHobbit,
    CogIcon,
    UserAddIcon,
    UserRemoveIcon,
    IconBar,
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const auth = useAuthStore()
    const users = useUsersStore()
    const hobbits = useHobbitsStore()
    const profiles = useProfileStore()

    const userId = computed(() => {
      if (!route.params.profileId) {
        return auth.userId
      } else {
        if (Array.isArray(route.params.profileId)) {
          return Number(route.params.profileId[0])
        }
        return Number(route.params.profileId)
      }
    })

    const user = computed(() => {
      if (userId.value) {
        return users.getUserById(userId.value)
      }
      return null
    })

    const isMe = computed(() => {
      // TODO: also check the id
      return !route.params.profileId
    })

    const hobbitsOfUser = computed(() => {
      if (userId.value) {
        return hobbits.getHobbitsByUser(userId.value)
      }
      return null
    })

    const follows = computed(() => {
      if (userId.value) {
        return profiles.followsUser(userId.value)
      }
      return false
    })

    const fetchHobbits = () => {
      if (userId.value) {
        hobbits.fetchHobbitsByUser(userId.value)
      }
    }
    const fetchFollow = () => {
      if (userId.value) {
        profiles.fetchFollow({ id: userId.value })
      }
    }

    const follow = () => {
      if (userId.value) {
        profiles.followUser({ id: userId.value })
      }
    }
    const unfollow = () => {
      if (userId.value) {
        profiles.unfollowUser({ id: userId.value })
      }
    }

    const goToAppPassword = () => {
      router.push('/profile/me/apppassword')
    }

    const deferredInit = () => {
      if (userId.value) {
        users.fetchUser({ id: userId.value })
        fetchFollow()
        fetchHobbits()
      }
    }

    watch(userId, () => {
      deferredInit()
    }, { immediate: true })

    return {
      user,
      isMe,
      hobbitsOfUser,
      follows,
      follow,
      unfollow,
      goToAppPassword,
    }
  },
})

</script>

<style scoped>
.welcome {
  font-size: 16pt;
  text-align: center;
  margin-bottom: 2rem;
  color: var(--secondary-text);

  .username {
    color: var(--primary);
  }
}
</style>
