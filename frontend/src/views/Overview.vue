<template>
  <div>
    <div class="welcome" v-if="isAuthenticated">
      Here there, <span class="username">{{username}}</span>!
    </div>
    <IconBar @reload="reload" />
    <SimpleHobbit v-for="hobbit in hobbits" :key='`hobbit-${hobbit.id}`' :hobbit="hobbit" :withHeatmap=true  />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import SimpleHobbit from '../components/SimpleHobbit.vue'
import IconBar from '@/components/IconBar.vue'
import { createNamespacedHelpers } from 'vuex'
import { HobbitsState } from '@/store/modules/hobbits'
import { mapState } from 'pinia'
import { useAuthStore } from '@/store/auth'

const { mapState: mapHobbitsState, mapGetters: mapHobbitsGetters, mapActions: mapHobbitsActions } = createNamespacedHelpers('hobbits')

export default defineComponent({
  name: 'OverviewView',
  components: { SimpleHobbit, IconBar },
  created() {
    this.dispatchFetchHobbits()
  },
  computed: {
    ...mapState(useAuthStore, { isAuthenticated: 'authenticated', username: 'username' }),
    ...mapHobbitsState({
      initialLoaded: (state) => (state as HobbitsState).initialLoaded,
    }),
    ...mapHobbitsGetters({
      hobbits: 'getHobbits',
    }),
  },
  methods: {
    ...mapHobbitsActions({
      _fetchHobbits: 'fetchHobbits',
    }),
    dispatchFetchHobbits() {
      if (!this.initialLoaded) {
        this._fetchHobbits()
      }
    },
    reload() {
      this._fetchHobbits()
    },
  },
})
</script>

<style lang="scss" scoped>
.welcome {
  font-size: 16pt;
  text-align: center;
  margin-bottom: 2rem;
  color: var(--secondary-text);

  .username {
    color: var(--primary);
  }
}
</style>
