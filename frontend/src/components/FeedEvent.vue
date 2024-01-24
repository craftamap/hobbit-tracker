<template>
  <div class="card">
    <template v-if="isHobbitCreated">
      <div class="header">
        <router-link :to="`/profile/${hobbit?.user.id}`">{{ hobbit?.user.username }}</router-link> has created a new
        Hobbit.
      </div>
      <SimpleHobbit :hobbit="hobbit" />
    </template>
    <template v-if="isRecordCreated">
      <div class="header">
        <router-link :to="`/profile/${hobbit?.user.id}`">{{ hobbit?.user?.username }}</router-link> has created a new
        entry in
        <router-link :to="`/hobbits/${hobbit?.user.id}`">"{{ hobbit?.name }}".</router-link>
      </div>
      <h1>
        <router-link :to="`/hobbits/${hobbit?.id}`">{{ record?.value }}</router-link>
      </h1>
      <blockquote class="comment" v-if="!!record?.comment">{{ record?.comment }}</blockquote>
    </template>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, toRefs } from 'vue'
import { FeedEvent, FeedEventTypus, Hobbit, NumericRecord } from '../models'
import SimpleHobbit from '../components/SimpleHobbit.vue'

export default defineComponent({
  name: 'FeedEvent',
  components: {
    SimpleHobbit,
  },
  props: {
    feedEvent: Object as PropType<FeedEvent>,
  },
  setup(props) {
    const { feedEvent } = toRefs(props)

    const isHobbitCreated = computed(() => {
      return feedEvent.value?.FeedEventTypus === FeedEventTypus.HobbitCreated
    })
    const isRecordCreated = computed(() => {
      return feedEvent.value?.FeedEventTypus === FeedEventTypus.RecordCreated
    })

    const hobbit = computed(() => {
      if (isHobbitCreated.value) {
        return feedEvent.value?.Payload as Hobbit
      } else if (isRecordCreated.value) {
        return (feedEvent.value?.Payload as NumericRecord)?.hobbit
      }
      return undefined
    })

    const record = computed(() => {
      if (isRecordCreated.value) {
        return feedEvent.value?.Payload as NumericRecord
      }
      return undefined
    })

    return {
      isHobbitCreated,
      isRecordCreated,
      hobbit,
      record,
    }
  },
})
</script>

<style scoped>
.card {
  border-radius: 0.5rem;
  box-shadow: 0px 0px 5px -2px #000000;
  padding: 1rem;
  margin: 0.5rem 0.5rem;

  h2 {
    margin: 0;
    font-size: 1.2em;
  }

  h1 {
    margin: 0;
    font-size: 1.4em;
  }

  .comment {
    padding: 1em;
    border-left: solid 3px var(--primary);
  }
}
</style>
