<template>
  <div>
    <div class="welcome" v-if="isAuthenticated">
      Here there, <span class="username">{{ username }}</span>!
    </div>
    <IconBar @reload="reload">
      <template v-slot:left>
        <span @click="navigateAddHobbit">
          <Add /><span>Add Hobbit... </span>
        </span>
      </template>
      <template v-slot:right>
        <span @click="reload">
          <Reload />
        </span>
      </template>
    </IconBar>
    <SimpleHobbit v-for="hobbit in hobbits" :key='`hobbit-${hobbit?.id}`' :hobbit="hobbit" :withHeatmap=true />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import SimpleHobbit from '../components/SimpleHobbit.vue'
import IconBar from '../components/IconBar.vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '../store/auth'
import { useHobbitsStore } from '../store/hobbits'
import { useRouter } from 'vue-router'
import Add from '../components/Icons/AddIcon.vue'
import Reload from '../components/Icons/ReloadIcon.vue'

export default defineComponent({
  name: 'OverviewView',
  components: { SimpleHobbit, IconBar, Add, Reload },
  setup() {
    const auth = useAuthStore()
    const hobbits = useHobbitsStore()

    const { getHobbits } = storeToRefs(hobbits)
    const { authenticated: isAuthenticated, username } = storeToRefs(auth)

    const reload = () => {
      hobbits.fetchHobbits()
    }

    hobbits.fetchHobbits()

    const router = useRouter()

    const navigateAddHobbit = () => {
      router.push('/hobbits/add')
    }

    return {
      hobbits: getHobbits,
      reload,
      username,
      isAuthenticated,
      navigateAddHobbit,
    }
  },
})
</script>

<style scoped>
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
