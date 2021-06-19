<template>
  <div>
    <div v-for="feedEvent, idx in feedEvents" :key="idx">
      <FeedEvent :feedEvent="feedEvent" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { createNamespacedHelpers } from 'vuex'
import { FeedState } from '../store/modules/feed'
import FeedEvent from '@/components/FeedEvent.vue'

const { mapState, mapActions } = createNamespacedHelpers('feed')

export default defineComponent({
  name: 'Feed',
  components: {
    FeedEvent,
  },
  created() {
    this.fetchFeed()
  },
  computed: {
    ...mapState({
      feedEvents: state => (state as FeedState).feedEvents,
    }),
  },
  methods: {
    ...mapActions({
      fetchFeed: 'fetchFeed',
    }),
  },
})
</script>
