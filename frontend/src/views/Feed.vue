<template>
  <div class="feed">
    <div class="greeting">
      <div class="welcome" v-if="isAuthenticated">
        Here there, <span class="username">{{username}}</span>!
      </div>
    </div>
    <div class="sidebar">
      <h1>Your hobbits:</h1>
      <SimpleHobbit  v-for="hobbit in hobbitsOfUser" :key="hobbit.id" :hobbit="hobbit"/>
    </div>
    <div class="events">
      <h1>Your personal feed:</h1>
      <div v-for="feedEvent, idx in feedEvents" :key="idx">
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
import { Hobbit } from '@/models'

const { mapState, mapActions } = createNamespacedHelpers('feed')

export default defineComponent({
  name: 'Feed',
  components: {
    FeedEvent,
    SimpleHobbit,
  },
  created() {
    this.fetchFeed()
    this.fetchHobbitsByUser()
  },
  computed: {
    ...mapState({
      feedEvents: state => (state as FeedState).feedEvents,
    }),
    isAuthenticated(): boolean {
      return this.$store.state.auth.authenticated
    },
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
    ...mapActions({
      fetchFeed: 'fetchFeed',
    }),
    fetchHobbitsByUser() {
      this.$store.dispatch('fetchHobbitsByUser', { userId: this.userId })
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
</style>
