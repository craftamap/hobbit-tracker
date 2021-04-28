<template>
  <div class="grid">
    <div>
      <span class="icon-entry" @click="navigateAddHobbit"><Add/><span>Add Hobbit... </span></span>
    </div>
    <div>
    </div>
    <div>
      <span class="icon-entry" @click="reload"><Reload/></span>
    </div>
  </div>
</template>

<script lang="ts">
import Add from '@/components/Icons/Add.vue'
import { defineComponent } from 'vue'
import Reload from './Icons/Reload.vue'

export default defineComponent({
  components: {
    Add: Add,
    Reload: Reload,
  },
  methods: {
    navigateAddHobbit() {
      this.$router.push('/hobbits/add')
    },
    reload() {
      this.$emit('reload')
    },
  },
})
</script>

<style lang="scss" scoped>
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

.icon-entry {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
}
</style>
