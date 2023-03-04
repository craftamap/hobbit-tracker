<template>
  <div :data-uid="uid">
  </div>
  <div class="svg-wrapper">
    <svg viewBox="-30 -20 1090 140" width="1090" height="160">
      <text x="-30" fill="currentColor" y="35">Mo</text>
      <text x="-30" fill="currentColor" y="75">We</text>
      <text x="-30" fill="currentColor" y="115">Fr</text>
      <template v-for="(timestamp, idx) in timestamps" :key="timestamp">
        <text v-if="timestamp.getDate() === 1" :x="Math.floor((idx / 7)) * 20" y="-10" fill="currentColor">{{
          timestamp.toLocaleString(undefined, {
            month:
              'short'
          }) }}</text>
        <rect width="16" height="16" :y="(idx % 7) * 20" :x="Math.floor((idx / 7)) * 20"
          :data-timestamp="timestamp.getTime()" ry='3'
          :style="`fill: ${getColor((mappedData?.get(timestamp.getTime()) || 0) / (max === 0 ? 1 : max))}`">
          <title>{{ timestamp }} - Count: {{ mappedData?.get(timestamp.getTime()) || 0 }}</title>
        </rect>
      </template>
    </svg>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, toRefs } from 'vue'
import { v4 as uuidv4 } from 'uuid'

export default defineComponent({
  name: 'CalendarHeatmap',
  props: {
    data: Array as PropType<{ count: number, date: Date }[]>,
  },
  setup(props) {
    const { data } = toRefs(props)
    const uid = ref(uuidv4())

    const tmp = computed(() => {
      console.error('foo')
      const map = new Map();
      let max = 0;
      for (const entry of data.value || []) {
        if (entry.count > max) {
          max = entry.count;
        }
        map.set(entry.date.setHours(0, 0, 0, 0), entry.count || 0)
      }
      console.error('foo done')
      return [map, max]
    })
    const mappedData = computed(() => tmp.value[0] as Map<number, number>)
    const max = computed(() => tmp.value[1] as number);
    console.log(mappedData.value)

    // monday, 52 weeks ago, midnight
    const startDate = new Date()
    startDate.setFullYear(startDate.getFullYear() - 1);
    startDate.setDate(startDate.getDate() - startDate.getDay());
    startDate.setHours(0, 0, 0, 0)
    console.log(startDate)

    // end of this week, midnight
    const endDate = new Date()
    endDate.setDate(endDate.getDate() - endDate.getDay());
    endDate.setDate(endDate.getDate() - endDate.getDay() + 6);
    startDate.setHours(0, 0, 0, 0)
    console.log(endDate)


    const timestamps = computed(() => {
      const timestamps = []
      const ptrDate = new Date(startDate)
      while (ptrDate <= endDate) {
        timestamps.push(new Date(ptrDate))
        ptrDate.setDate(ptrDate.getDate() + 1)
      }
      console.log(timestamps);
      return timestamps;
    })

    const prefersDark = computed(() => {
      return window.matchMedia('(prefers-color-scheme: dark)').matches
    })

    const getColor = (percentage: number) => {
      let start = [233, 246, 247];
      let end = [33, 131, 128];
      if (prefersDark.value) {
        start = [51, 51, 51];
        end = [33, 131, 128];
      }
      const r = Math.round(start[0] + (end[0] - start[0]) * percentage)
      const g = Math.round(start[1] + (end[1] - start[1]) * percentage)
      const b = Math.round(start[2] + (end[2] - start[2]) * percentage)
      return `rgb(${r}, ${g}, ${b})`;
    }

    // const renderCalHeatMap = () => {
    //   return new Promise((resolve) => {
    //     setTimeout(() => {
    //       const calHeatMap = new CalendarHeatMap()
    //         .data(data.value)
    //         .selector(`[data-uid='${uid.value}']`)
    //         .colorRange(['#e9f6f7', '#218380'])
    //         .tooltipEnabled(true)
    //       if (prefersDark.value) {
    //         // replace with better color
    //         calHeatMap.colorRange(['#333333', '#218380'])
    //       }

    //       resolve(calHeatMap())
    //     }, 0)
    //   })
    // }
    // watch(data, () => {
    //   renderCalHeatMap()
    // })

    // renderCalHeatMap()

    return {
      uid,
      mappedData,
      max,
      timestamps,
      getColor,
    }
  },
})
</script>

<style lang="scss" scoped>
[data-uid] {
  /* aspect-ratio: 5; */

  text {
    fill: var(--primary-text);
  }
}


svg {
  width: 100%;
}

.svg-wrapper {
  width: 100%;
}

@media (prefers-color-scheme: dark) {
  .svg-wrapper {
    color: black;
  }

}
</style>

