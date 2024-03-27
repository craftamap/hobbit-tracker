<template>
  <div class="grid">
    <div>
      <slot name="left">
      </slot>
    </div>
    <div>
      <slot name="center">
      </slot>
    </div>
    <div>
      <slot name="right">
      </slot>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'IconBar',
  setup(_, { emit }) {
    const reload = () => {
      emit('reload')
    }
    return {
      reload,
    }
  },
})
</script>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;

  :nth-child(2) {
    justify-self: center;
  }

  :nth-child(3) {
    justify-self: right;
  }
}

.grid :deep(> div > span) {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
}

</style>
