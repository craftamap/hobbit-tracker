<template>
  <div class="card" :data-id="hobbit.id">
    <div class="header">
      <div>
        <h1>
          <router-link :to="`/hobbits/${hobbit.id}`">{{
            hobbit.name
          }}</router-link>
        </h1>
        <div class="by">by {{ hobbit.user.username }}</div>
        <div>
          {{ hobbit.description }}
        </div>
      </div>
      <div>
        <img :src="hobbit.image" v-if="hobbit.image" />
      </div>
    </div>
    <div>
      <Loading v-if="loading" />
      <Heatmap v-if="!loading" :data="getHeatmapData" class="heatmap" />
    </div>
  </div>
</template>

<script lang="ts">
import { Hobbit, NumericRecord } from '../models/index'
import { defineComponent, PropType } from 'vue'
import Loading from './Loading.vue'
import moment from 'moment'
import Heatmap from './Heatmap.vue'

export default defineComponent({
  props: {
    hobbit: Object as PropType<Hobbit>,
  },
  components: {
    Loading,
    Heatmap,
  },
  data() {
    return {
      loading: true,
      calHeatMap: undefined as any,
    }
  },
  created() {
    console.log('created', this.loading)
    this.dispatchFetchHeatmapData()
  },
  computed: {
    getRecords(): NumericRecord[] {
      return (this?.hobbit?.records as NumericRecord[]) || []
    },
    getHeatmapData(): object[] {
      const result: { [key: string]: { date: Date; count: number } } = {}
      return (
        (this?.hobbit?.heatmap as NumericRecord[]) || [
          {
            timestamp: new Date(),
            value: 0,
          },
        ]
      ).map((record) => {
        return {
          date: new Date(record.timestamp),
          count: record.value,
        }
      })
    },
  },
  methods: {
    dispatchFetchHeatmapData() {
      if (!this.$props.hobbit?.heatmap) {
        this.$store
          .dispatch('fetchHeatmapData', this.$props.hobbit?.id)
          .then(() => {
            this.$data.loading = false
          })
      } else {
        this.$data.loading = false
      }
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

  h1 {
    margin: 0;
    font-size: 16pt;
  }

  .by {
    color: var(--secondary-text);
  }

  .header {
    display: flex;
    justify-content: space-between;

    img {
      width: 2rem;
      height: 2rem;
    }
  }
  .heatmap {
    max-width: 600px;
    margin: 1rem auto;
  }
}
</style>
