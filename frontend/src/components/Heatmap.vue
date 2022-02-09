<template>
    <div :data-uid="uid">
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import CalendarHeatMap from 'calendar-heatmap-mini'

export default defineComponent({
  name: 'CalendarHeatmap',
  props: {
    data: Array,
  },
  watch: {
    data() {
      this.renderCalHeatMap()
    },
  },
  computed: {
    prefersDark() {
      return window.matchMedia('(prefers-color-scheme: dark)')
    },
  },
  mounted() {
    this.renderCalHeatMap()
  },
  data() {
    return {
      uid: uuidv4(),
    }
  },
  methods: {
    async renderCalHeatMap() {
      return new Promise((resolve) => {
        setTimeout(() => {
          const calHeatMap = new CalendarHeatMap()
            .data(this.data)
            .selector(`[data-uid='${this.uid}']`)
            .colorRange(['#e9f6f7', '#218380'])
            .tooltipEnabled(true)
          if (this.prefersDark.matches) {
            // replace with better color
            calHeatMap.colorRange(['#333333', '#218380'])
          }

          resolve(calHeatMap())
        }, 0)
      })
    },
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
