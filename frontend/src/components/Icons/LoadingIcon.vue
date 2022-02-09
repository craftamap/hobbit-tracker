<template>
  <div>
      <svg :width="width" :height="width" viewBox="0 0 44 44" xmlns="http://www.w3.org/2000/svg" :stroke="stroke">
          <g fill="none" fill-rule="evenodd" stroke-width="2">
              <circle cx="22" cy="22" r="1">
                  <animate attributeName="r"
                      begin="0s" dur="1.8s"
                      values="1; 20"
                      calcMode="spline"
                      keyTimes="0; 1"
                      keySplines="0.165, 0.84, 0.44, 1"
                      repeatCount="indefinite" />
                  <animate attributeName="stroke-opacity"
                      begin="0s" dur="1.8s"
                      values="1; 0"
                      calcMode="spline"
                      keyTimes="0; 1"
                      keySplines="0.3, 0.61, 0.355, 1"
                      repeatCount="indefinite" />
              </circle>
              <circle cx="22" cy="22" r="1">
                  <animate attributeName="r"
                      begin="-0.9s" dur="1.8s"
                      values="1; 20"
                      calcMode="spline"
                      keyTimes="0; 1"
                      keySplines="0.165, 0.84, 0.44, 1"
                      repeatCount="indefinite" />
                  <animate attributeName="stroke-opacity"
                      begin="-0.9s" dur="1.8s"
                      values="1; 0"
                      calcMode="spline"
                      keyTimes="0; 1"
                      keySplines="0.3, 0.61, 0.355, 1"
                      repeatCount="indefinite" />
              </circle>
          </g>
      </svg>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  props: {
    width: {
      type: [Number, String],
      default: 48,
    },
    height: {
      type: [Number, String],
      default: 48,
    },
    stroke: {
      type: String,
      default: '#000',
    },
  },
})
</script>

<style lang="scss" scoped>
div {
    padding: 1rem;

    svg {
       margin: 0 auto;
       display: block;
    }
}
</style>
