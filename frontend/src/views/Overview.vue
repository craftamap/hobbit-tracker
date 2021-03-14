<template>
  <div>
    <SimpleHobbit v-for="hobbit in $store.state.hobbits" :key='hobbit.id' :hobbit="hobbit" />
  </div>
</template>

<script>
import SimpleHobbit from '../components/SimpleHobbit'

export default {
  name: 'Overview',
  components: { SimpleHobbit },
  created () {
    this.dispatchFetchHobbits()
  },
  methods: {
    dispatchFetchHobbits () {
      this.$store.dispatch('fetchHobbits')
    }
  }
}
</script>
