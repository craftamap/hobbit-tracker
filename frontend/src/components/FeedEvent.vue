<template>
  <div class="card">
    <div class="header">
      <h2>{{title}}</h2>
    </div>
    <div class="subcard">
      <template v-if="isHobbitCreated">
          <h1>
            <router-link :to="`/hobbits/${hobbit.id}`">{{
              hobbit.name
            }}</router-link>
          </h1>
          <div class="by">by {{ hobbit.user.username }}</div>
          <div>
            {{ hobbit.description }}
          </div>
          <div>
            <img :src="hobbit.image" v-if="hobbit.image" />
          </div>
      </template >
      <template v-if="isRecordCreated">
          <h1>
            <router-link :to="`/hobbits/${hobbit.id}`">{{
              record.value
            }}</router-link>
          </h1>
          <div class="by">
            <router-link :to="`/hobbits/${hobbit.id}`">"{{
              hobbit.name
            }}"</router-link>
            by {{ hobbit.user.username }}
          </div>
          <div>
            <img :src="hobbit.image" v-if="hobbit.image" />
          </div>
      </template >
  </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { FeedEvent, FeedEventTypus } from '@/store/modules/feed'
import { Hobbit, NumericRecord } from '@/models'

export default defineComponent({
  name: 'FeedEvent',
  props: {
    feedEvent: Object as PropType<FeedEvent>,
  },
  computed: {
    title(): string {
      if (this.feedEvent?.FeedEventTypus === FeedEventTypus.HobbitCreated) {
        const hobbit = this.feedEvent.Payload as Hobbit
        return `${(hobbit.user as any).username} has created a new Hobbit.`
      } else if (this.feedEvent?.FeedEventTypus === FeedEventTypus.RecordCreated) {
        const record = this.feedEvent.Payload as NumericRecord
        return `${((record as any).hobbit.user as any).username} has created a new entry in "${((record as any).hobbit).name}."`
      }
      return ''
    },
    isHobbitCreated(): boolean {
      return this.feedEvent?.FeedEventTypus === FeedEventTypus.HobbitCreated
    },
    isRecordCreated(): boolean {
      return this.feedEvent?.FeedEventTypus === FeedEventTypus.RecordCreated
    },
    hobbit(): Hobbit | null {
      if (this.isHobbitCreated) {
        return this?.feedEvent?.Payload as Hobbit
      } else if (this.isRecordCreated) {
        return (this?.feedEvent?.Payload as any).hobbit as Hobbit
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
    font-size: 14pt;
  }
  .subcard {
    border-radius: 0.5rem;
    box-shadow: 0px 0px 5px -2px #000000;
    padding: 1rem;
    margin: 0.5rem 0.5rem;
    display: grid;

    h1 {
      margin: 0;
      font-size: 16pt;
    }
    .by {
      color: var(--secondary-text);
    }
  }
}
</style>
