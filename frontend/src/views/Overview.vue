<template>
  <div>
    <div class="welcome" v-if="isAuthenticated">
      Here there, <span class="username">{{username}}</span>!
    </div>
    <IconBar @reload="reload" />
    <SimpleHobbit v-for="hobbit in hobbits" :key='`hobbit-${hobbit?.id}`' :hobbit="hobbit" :withHeatmap=true  />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import SimpleHobbit from '../components/SimpleHobbit.vue'
import IconBar from '@/components/IconBar.vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/store/auth'
import { useHobbitsStore } from '@/store/hobbits'

export default defineComponent({
  name: 'OverviewView',
  components: { SimpleHobbit, IconBar },
  setup() {
    const auth = useAuthStore()
    const hobbits = useHobbitsStore()

    const { getHobbits } = storeToRefs(hobbits)
    const { authenticated: isAuthenticated, username } = storeToRefs(auth)

    const reload = () => {
      hobbits.fetchHobbits()
    }

    hobbits.fetchHobbits()

    return {
      hobbits: getHobbits,
      reload,
      username,
      isAuthenticated,
    }
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
