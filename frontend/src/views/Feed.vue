<template>
  <div class="feed">
    <div class="greeting">
      <div class="welcome" v-if="isAuthenticated">
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
        <span class="icon-entry" @click="navigateAddHobbit">
          <PlusIcon class="w-24 h-24" />
          <span>Add Hobbit...</span>
        </span>
      </div>
      <SimpleHobbit v-for="hobbit in hobbitsOfUser" :key="`hobbit-${hobbit.id}`" :hobbit="hobbit" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import FeedEvent from '@/components/FeedEvent.vue'
import SimpleHobbit from '@/components/SimpleHobbit.vue'
import { PlusIcon } from '@heroicons/vue/outline'
import { Hobbit } from '@/models'
import { useAuthStore } from '@/store/auth'
import { mapActions, mapState } from 'pinia'
import { useFeedStore } from '@/store/feed'
import { useHobbitsStore } from '@/store/hobbits'

export default defineComponent({
  name: 'FeedView',
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
    ...mapState(useFeedStore, ['feedEvents']),
    ...mapState(useAuthStore, { isAuthenticated: 'authenticated', username: 'username', userId: 'userId' }),
    ...mapState(useHobbitsStore, { _hobbitsByUser: 'getHobbitsByUser' }),
    hobbitsOfUser(): Hobbit[] {
      if (this.userId) {
        return this._hobbitsByUser(this.userId)
      }
      return []
    },
  },
  methods: {
    ...mapActions(useFeedStore, [
      'fetchFeed',
    ]),
    ...mapActions(useHobbitsStore, {
      _fetchHobbitsByUser: 'fetchHobbitsByUser',
    }),
    fetchHobbitsByUser() {
      this._fetchHobbitsByUser()
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
