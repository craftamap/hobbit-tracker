<template>
  <div>
    <SimpleHobbit v-for="hobbit in $store.state.hobbits" :key='hobbit.id' :hobbit="hobbit" />
  </div>
</template>

<script lang="ts">
import { useStore } from '../store'
import { defineComponent } from 'vue'
import SimpleHobbit from '../components/SimpleHobbit.vue'

export default defineComponent({
  name: 'Overview',
  components: { SimpleHobbit },
  created () {
    this.dispatchFetchHobbits()
  },
  methods: {
    dispatchFetchHobbits () {
      useStore().dispatch('fetchHobbits')
    }
  }
})
</script>
