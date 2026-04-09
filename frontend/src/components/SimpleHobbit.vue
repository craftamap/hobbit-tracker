<template>
  <div
    class="card"
    :data-id="hobbit?.id"
  >
    <div class="header">
      <div>
        <h1>
          <router-link :to="`/hobbits/${hobbit?.id}`">
            {{
              hobbit?.name
            }}
          </router-link>
        </h1>
        <div class="by">
          by {{ hobbit?.user.username }}
        </div>
        <div>
          {{ hobbit?.description }}
        </div>
      </div>
      <div>
        <img
          v-if="hobbit?.image"
          :src="hobbit?.image"
        >
      </div>
    </div>
    <div v-if="withHeatmap">
      <Loading v-if="loadingHeatmapData" />
      <Heatmap
        v-if="!loadingHeatmapData"
        :data="heatmapData"
        class="heatmap"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { Hobbit } from '../models/index'
import { computed, defineComponent, PropType, ref, toRefs } from 'vue'
import Loading from './Icons/LoadingIcon.vue'
import Heatmap from './Heatmap.vue'
import { useHobbitsStore } from '../store/hobbits'

export default defineComponent({
  components: {
    Loading,
    Heatmap,
  },
  props: {
    hobbit: {
      type: Object as PropType<Hobbit>,
      required: true,
    },
    withHeatmap: {
      type: Boolean as PropType<boolean>,
      default: false,
    },
  },
  setup(props) {
    const hobbitsStore = useHobbitsStore()

    const loadingHeatmapData = ref(true)
    const { hobbit } = toRefs(props)


    const heatmapData = computed(() => {
      if (!hobbit.value?.id) {
        return []
      }
      return hobbitsStore.getHeatmapByHobbitId(hobbit.value?.id)
    })

    const fetchHeatmapData = async() => {
      if (!hobbit.value?.id) {
        return;
      }

      if (!hobbitsStore.getRecordsByHobbitId(hobbit.value?.id)) {
        try {
          await hobbitsStore.fetchRecords(hobbit.value?.id);
        } catch (e) {
          console.log('failed to fetch records for heatmap', e);
        }
      }
      loadingHeatmapData.value = false;
    }

    fetchHeatmapData()

    return {
      loadingHeatmapData,
      heatmapData,
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
