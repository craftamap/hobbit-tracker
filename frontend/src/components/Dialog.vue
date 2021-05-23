<template>
  <DialogWrapper v-if="shown">
    <div class="dialog">
      <slot></slot>
    </div>
  </DialogWrapper>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import DialogWrapper from './DialogWrapper.vue'

export default defineComponent({
  components: {
    DialogWrapper,
  },
  props: {
    shown: {
      type: Boolean,
      default: false,
    },
  },
})
</script>

<style lang="scss" scoped>
.dialog {
  background: var(--background);
  padding: 1.5em;
  border-radius: 0.75em;
}
</style>
