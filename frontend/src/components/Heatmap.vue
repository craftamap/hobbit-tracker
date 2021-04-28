<template>
    <div :data-uid="uid">
    </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import CalendarHeatMap from 'calendar-heatmap-mini'

export default defineComponent({
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
    renderCalHeatMap() {
      const calHeatMap = new CalendarHeatMap()
        .data(this.data)
        .selector(`[data-uid='${this.uid}']`)
        .colorRange(['#e9f6f7', '#218380'])
        .tooltipEnabled(true)
        .onClick(function(data: any) {
          console.log('onClick callback. Data:', data)
        })
      if (this.prefersDark.matches) {
        // replace with better color
        calHeatMap.colorRange(['#333333', '#218380'])
        console.log('a')
      }
      calHeatMap()
      console.log('draw')
    },
  },
})
</script>
