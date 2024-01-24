<template>
  <div class="grid">
    <div>
      <span class="icon-entry" @click="navigateAddHobbit">
        <Add /><span>Add Hobbit... </span>
      </span>
    </div>
    <div>
    </div>
    <div>
      <span class="icon-entry" @click="reload">
        <Reload />
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import Add from '../components/Icons/AddIcon.vue'
import { defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import Reload from './Icons/ReloadIcon.vue'

export default defineComponent({
  name: 'IconBar',
  components: {
    Add: Add,
    Reload: Reload,
  },
  setup(_, { emit }) {
    const router = useRouter()

    const navigateAddHobbit = () => {
      router.push('/hobbits/add')
    }
    const reload = () => {
      emit('reload')
    }
    return {
      navigateAddHobbit,
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

.icon-entry {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
}</style>
