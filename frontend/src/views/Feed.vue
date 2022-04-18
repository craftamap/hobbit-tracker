<template>
  <div class="feed">
    <div class="greeting">
      <div class="welcome" v-if="authenticated">
        Here there,
        <span class="username">{{ username }}</span>!
      </div>
    </div>
    <div class="events">
      <h1>Your personal feed:</h1>
      <template v-for="(feedEvent, idx) in feedEvents" v-bind:key="`feedEvent-${idx}`">
        <FeedEvent :feedEvent="feedEvent" />
      </template>
    </div>
    <div class="sidebar">
      <h1>Your hobbits:</h1>
      <div>
        <span class="icon-entry" @click="goToAddHobbit">
          <PlusIcon class="w-24 h-24" />
          <span>Add Hobbit...</span>
        </span>
      </div>
      <SimpleHobbit v-for="hobbit in hobbitsOfUser" :key="`hobbit-${hobbit.id}`" :hobbit="hobbit" />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, toRefs, watch } from 'vue'
import FeedEvent from '@/components/FeedEvent.vue'
import SimpleHobbit from '@/components/SimpleHobbit.vue'
import { PlusIcon } from '@heroicons/vue/outline'
import { useAuthStore } from '@/store/auth'
import { mapActions, storeToRefs } from 'pinia'
import { useFeedStore } from '@/store/feed'
import { useHobbitsStore } from '@/store/hobbits'
import { useRouter } from 'vue-router'

export default defineComponent({
  name: 'FeedView',
  components: {
    FeedEvent,
    SimpleHobbit,
    PlusIcon,
  },
  setup() {
    const hobbits = useHobbitsStore()
    const router = useRouter()

    const feed = useFeedStore()
    const { feedEvents } = storeToRefs(feed)

    const auth = useAuthStore()
    const { userId, authenticated, username } = storeToRefs(auth)

    const hobbitsOfUser = computed(() => {
      if (userId) {
        return hobbits.getHobbitsByUser(userId.value!)
      }
      return []
    })

    feed.fetchFeed()

    hobbits.fetchHobbitsByUser()

    const goToAddHobbit = () => {
      router.push('/hobbits/add')
    }

    return {
      feedEvents,
      authenticated,
      username,
      hobbitsOfUser,
      goToAddHobbit,
    }
  },
})
</script>

<style lang="scss" scoped>
.welcome {
  font-size: 16pt;
  text-align: center;
  margin-bottom: 2rem;
  color: var(--secondary-text);

  .username {
    color: var(--primary);
  }
}

.feed {
  display: grid;
  gap: 0px 0px;
  grid-template-columns: 1fr 1fr 1fr;
  grid-template-areas:
    "greeting greeting greeting"
    "events events sidebar";
  justify-items: stretch;
  align-items: stretch;
  @media (max-width: 1000px) {
    grid-template-columns: 1fr;
    grid-template-areas:
      "greeting"
      "sidebar"
      "events";
  }

  h1 {
    font-size: 1.2em;
  }
}

.greeting {
  grid-area: greeting;
}
.sidebar {
  grid-area: sidebar;
}
.events {
  grid-area: events;
}
.icon-entry {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
}
</style>
