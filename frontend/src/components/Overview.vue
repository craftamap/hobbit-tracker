<template>
  <div>

  </div>
</template>

<script>
export default {
  name: 'Overview',
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
