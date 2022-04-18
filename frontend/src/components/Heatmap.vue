<template>
    <div :data-uid="uid">
    </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, toRefs, watch } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import CalendarHeatMap from 'calendar-heatmap-mini'

export default defineComponent({
  name: 'CalendarHeatmap',
  props: {
    data: Array,
  },
  setup(props) {
    const { data } = toRefs(props)
    const uid = ref(uuidv4())

    const prefersDark = computed(() => {
      return window.matchMedia('(prefers-color-scheme: dark)').matches
    })

    const renderCalHeatMap = () => {
      return new Promise((resolve) => {
        setTimeout(() => {
          const calHeatMap = new CalendarHeatMap()
            .data(data.value)
            .selector(`[data-uid='${uid.value}']`)
            .colorRange(['#e9f6f7', '#218380'])
            .tooltipEnabled(true)
          if (prefersDark.value) {
            // replace with better color
            calHeatMap.colorRange(['#333333', '#218380'])
          }

          resolve(calHeatMap())
        }, 0)
      })
    }
    watch(data, () => {
      renderCalHeatMap()
    })

    renderCalHeatMap()

    return {
      uid,
    }
  },
})
</script>

<style lang="scss">

[data-uid] {
  aspect-ratio: 5;

  text {
    fill: var(--primary-text);
  }
}

</style>
