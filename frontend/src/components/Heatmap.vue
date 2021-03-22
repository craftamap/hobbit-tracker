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
    data: Array
  },
  watch: {
    data (newVal) {
      const calHeatMap = new CalendarHeatMap()
        .data(newVal)
        .selector(`[data-uid='${this.uid}']`)
        .colorRange(['#D8E6E7', '#218380'])
        .tooltipEnabled(true)
        .onClick(function (data: any) {
          console.log('onClick callback. Data:', data)
        })
      calHeatMap()
      console.log('draw')
    }
  },
  mounted () {
    const calHeatMap = new CalendarHeatMap()
      .data(this.data)
      .selector(`[data-uid='${this.uid}']`)
      .colorRange(['#D8E6E7', '#218380'])
      .tooltipEnabled(true)
      .onClick(function (data: any) {
        console.log('onClick callback. Data:', data)
      })
    calHeatMap()
    console.log('draw')
  },
  data () {
    return {
      uid: uuidv4()
    }
  }
})
</script>
