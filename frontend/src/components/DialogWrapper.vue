<template>
  <div class="dialog-wrapper">
    <slot></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'DialogWrapper',
})
</script>

<style lang="scss" scoped>
.dialog-wrapper {
  position: fixed;
  left: 0;
  top: 0;
  z-index: 10;
  background: rgba(black, 0.75);
  height: 100vh;
  width: 100vw;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
