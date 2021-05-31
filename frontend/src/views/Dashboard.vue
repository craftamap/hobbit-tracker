<template>
  <div v-for="feedEvent, idx in feedEvents" :key="idx">
    <div v-if="feedEvent.FeedEventTypus == 'HobbitCreated'">
      <h3>{{feedEvent.Payload.user.username}} created a new Hobbit: {{feedEvent.Payload.name}}</h3>
    </div>
    <div v-if="feedEvent.FeedEventTypus == 'RecordCreated'" >
      <h3>{{feedEvent.Payload.hobbit.user.username}} added a Record to {{feedEvent.Payload.hobbit.name}}:</h3>
      <h2>{{feedEvent.Payload.value}}</h2>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { createNamespacedHelpers } from 'vuex'
import { FeedState } from '../store/modules/feed'

const { mapState, mapActions } = createNamespacedHelpers('feed')

export default defineComponent({
  name: 'Dashboard',
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
