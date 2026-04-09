<template>
  <div :data-uid="uid" />
  <div class="svg-wrapper">
    <!-- note that svg is min-x min-y width height and NOT max-x max-y */ -->
    <svg viewBox="-30 -20 1106 160">
      <text
        x="-30"
        fill="currentColor"
        y="35"
      >Mo
      </text>
      <text
        x="-30"
        fill="currentColor"
        y="75"
      >We
      </text>
      <text
        x="-30"
        fill="currentColor"
        y="115"
      >Fr
      </text>
      <template
        v-for="(date, idx) in dates"
        :key="date.toString()"
      >
        <text
          v-if="date.day === 1"
          :x="Math.floor((idx / 7)) * 20"
          y="-10"
          fill="currentColor"
        >{{
          date.toLocaleString(undefined, {
            month: 'short'
          })
        }}
        </text>
        <rect
          width="16"
          height="16"
          :y="(idx % 7) * 20"
          :x="Math.floor((idx / 7)) * 20"
          :data-date="date.toString()"
          ry="3"
          :style="`fill: ${getColor((mappedData.map.get(date.toString()) || 0) / (mappedData.max === 0 ? 1 : mappedData.max))}`"
        >
          <title>{{ date.toString() }} - Count: {{ mappedData.map.get(date.toString()) || 0 }}</title>
        </rect>
      </template>
    </svg>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, toRefs } from 'vue'

export default defineComponent({
  name: 'CalendarHeatmap',
  props: {
    data: {
      type: Array as PropType<{ count: number, date: Temporal.PlainDate }[]>,
      required: true,
    },
  },
  setup(props) {
    const { data } = toRefs(props)
    const uid = ref(crypto.randomUUID())
    const tz = Temporal.Now.timeZoneId()

    const mappedData = computed(() => {
      const map = new Map<string, number>();
      let max = 0;
      for (const entry of data.value || []) {
        if (entry.count > max) {
          max = entry.count;
        }
        const dateStr = entry.date.toString();
        map.set(dateStr, entry.count || 0)
      }
      return { map, max }
    })

    const dates = computed(() => {
      const today = Temporal.Now.plainDateISO(tz)
      // End of this week (Sunday)
      const endDate = today.add({ days: 7 - today.dayOfWeek })
      // Start of the grid (Monday, 52 weeks ago)
      const startDate = endDate.subtract({ weeks: 52 }).subtract({ days: 6 })

      const dates = []
      let curr = startDate
      while (Temporal.PlainDate.compare(curr, endDate) <= 0) {
        dates.push(curr)
        curr = curr.add({ days: 1 })
      }
      return dates
    })

    const prefersDark = ref(window.matchMedia('(prefers-color-scheme: dark)').matches)

    window.matchMedia('(prefers-color-scheme: dark)').addEventListener("change", (pd) => {
      prefersDark.value = (pd.currentTarget as MediaQueryList).matches
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

    return {
      uid,
      mappedData,
      dates,
      getColor,
    }
  },
})
</script>

<style scoped>
svg {
  width: 100%;
}

.svg-wrapper {
  width: 100%;
}

@media (prefers-color-scheme: dark) {
  .svg-wrapper {
    color: white;
  }
}

text {
  fill: currentColor;
}
</style>
