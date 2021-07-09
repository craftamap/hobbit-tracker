<template>
  <div class="feed">
    <div class="greeting">
      <div class="welcome" v-if="isAuthenticated">
        Here there, <span class="username">{{username}}</span>!
      </div>
    </div>
    <div class="sidebar">
      <h1>Your hobbits:</h1>
      <div>
        <span class="icon-entry" @click="navigateAddHobbit"><PlusIcon class="w-24 h-24"/><span>Add Hobbit... </span></span>
      </div>
      <SimpleHobbit  v-for="hobbit in hobbitsOfUser" :key="hobbit.id" :hobbit="hobbit"/>
    </div>
    <div class="events">
      <h1>Your personal feed:</h1>
      <div v-for="(feedEvent, idx) in feedEvents" v-bind:key="idx">
        <FeedEvent :feedEvent="feedEvent" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { createNamespacedHelpers } from 'vuex'
import { FeedState } from '../store/modules/feed'
import FeedEvent from '@/components/FeedEvent.vue'
import SimpleHobbit from '@/components/SimpleHobbit.vue'
import { PlusIcon } from '@heroicons/vue/outline'
import { AuthenticationState } from '@/store/modules/auth'
import { Hobbit } from '@/models'

const { mapState: feedMapState, mapActions: feedMapActions } = createNamespacedHelpers('feed')
const { mapState: authMapState } = createNamespacedHelpers('auth')
const { mapActions: mapHobbitsActions, mapGetters: mapHobbitsGetters } = createNamespacedHelpers('hobbits')

export default defineComponent({
  name: 'Feed',
  components: {
    FeedEvent,
    SimpleHobbit,
    PlusIcon,
  },
  created() {
    this.fetchFeed()
    this.fetchHobbitsByUser()
  },
  computed: {
    ...feedMapState({
      feedEvents: state => (state as FeedState).feedEvents,
    }),
    ...authMapState({
      isAuthenticated: state => (state as AuthenticationState).authenticated,
      username: state => (state as AuthenticationState).username,
      userId: state => (state as AuthenticationState).userId,
    }),
    ...mapHobbitsGetters({
      _hobbitsByUser: 'getHobbitsByUser',
    }),
    hobbitsOfUser(): Hobbit[] {
      return this._hobbitsByUser(this.userId)
    },
  },
  methods: {
    ...feedMapActions({
      fetchFeed: 'fetchFeed',
    }),
    ...mapHobbitsActions({
      _fetchHobbitsByUser: 'fetchHobbitsByUser',
    }),
    fetchHobbitsByUser() {
      this._fetchHobbitsByUser({ userId: this.userId })
    },
    navigateAddHobbit() {
      this.$router.push('/hobbits/add')
    },
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
    "sidebar events events";
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
