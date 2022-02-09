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
    <div v-if="withHeatmap">
      <Loading v-if="loadingHeatmapData" />
      <Heatmap v-if="!loadingHeatmapData" :data="heatmapData" class="heatmap" />
    </div>
  </div>
</template>

<script lang="ts">
import { Hobbit, NumericRecord } from '../models/index'
import { defineComponent, PropType } from 'vue'
import Loading from './Icons/LoadingIcon.vue'
import Heatmap from './Heatmap.vue'
import { createNamespacedHelpers } from 'vuex'

const { mapState: mapHobbitsState, mapGetters: mapHobbitsGetters, mapActions: mapHobbitsActions } = createNamespacedHelpers('hobbits')

export default defineComponent({
  props: {
    hobbit: Object as PropType<Hobbit>,
    withHeatmap: {
      type: Boolean as PropType<boolean>,
      default: false,
    },
  },
  components: {
    Loading,
    Heatmap,
  },
  data() {
    return {
      loadingHeatmapData: true,
    }
  },
  created() {
    if (this.withHeatmap) {
      this.fetchHeatmapData()
    }
  },
  computed: {
    getRecords(): NumericRecord[] {
      return (this?.hobbit?.records as NumericRecord[]) || []
    },
    heatmapData(): object[] {
      const retVal = (
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
      return retVal
    },
  },
  methods: {
    ...mapHobbitsActions({
      _fetchHeatmapData: 'fetchHeatmapData',
    }),
    fetchHeatmapData() {
      if (!this.$props.hobbit?.heatmap) {
        this._fetchHeatmapData(this.$props.hobbit?.id)
          .then(() => {
            this.$data.loadingHeatmapData = false
          })
      } else {
        this.$data.loadingHeatmapData = false
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
