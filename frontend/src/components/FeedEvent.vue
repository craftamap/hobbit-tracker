<template>
  <div class="card">
    <div class="header">
      <h2>{{title}}</h2>
    </div>
    <template v-if="isHobbitCreated">
      <SimpleHobbit :hobbit="hobbit" />
    </template >
    <template v-if="isRecordCreated">
        <h1>
          <router-link :to="`/hobbits/${hobbit.id}`">{{
            record.value
          }}</router-link>
        </h1>
        <blockquote class="comment" v-if="!!record.comment">
          {{  record.comment  }}
        </blockquote>
      <SimpleHobbit :hobbit="hobbit" />
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { FeedEvent, FeedEventTypus } from '@/store/modules/feed'
import { Hobbit, NumericRecord } from '@/models'
import SimpleHobbit from '@/components/SimpleHobbit.vue'

export default defineComponent({
  name: 'FeedEvent',
  components: {
    SimpleHobbit,
  },
  props: {
    feedEvent: Object as PropType<FeedEvent>,
  },
  computed: {
    title(): string {
      if (this.feedEvent?.FeedEventTypus === FeedEventTypus.HobbitCreated) {
        const hobbit = this.feedEvent.Payload as Hobbit
        return `${hobbit.user.username} has created a new Hobbit.`
      } else if (this.feedEvent?.FeedEventTypus === FeedEventTypus.RecordCreated) {
        const record = this.feedEvent.Payload as NumericRecord
        return `${record?.hobbit?.user?.username} has created a new entry in "${record?.hobbit?.name}."`
      }
      return ''
    },
    isHobbitCreated(): boolean {
      return this.feedEvent?.FeedEventTypus === FeedEventTypus.HobbitCreated
    },
    isRecordCreated(): boolean {
      return this.feedEvent?.FeedEventTypus === FeedEventTypus.RecordCreated
    },
    hobbit(): Hobbit | null | undefined {
      if (this.isHobbitCreated) {
        return this?.feedEvent?.Payload as Hobbit
      } else if (this.isRecordCreated) {
        return (this?.feedEvent?.Payload as NumericRecord)?.hobbit
      }
      return null
    },
    record(): NumericRecord | null {
      if (this.isRecordCreated) {
        return this?.feedEvent?.Payload as NumericRecord
      }
      return null
    },
  },
})
</script>

<style lang="scss" scoped>
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
