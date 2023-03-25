<script lang="ts">
import { useAuthStore } from '@/store/auth'
import { useHobbitsStore } from '@/store/hobbits'
import { computed, defineComponent, ref } from 'vue'
import { useRoute } from 'vue-router'
import Loading from '@/components/Icons/LoadingIcon.vue'
import Button from '@/components/form/Button.vue'

export default defineComponent({
  components: {
    Loading,
    Button,
  },
  setup() {
    const route = useRoute()
    const auth = useAuthStore()
    const hobbits = useHobbitsStore()

    const sharedFileId = route.params.sharedFileId
    console.log(sharedFileId)

    const hobbitsByUser = computed(() => {
      if (auth.userId) {
        return hobbits.getHobbitsByUser(auth.userId)
      }
      return []
    })

    hobbits.fetchHobbitsByUser()

    const shareData = ref(null);
    (async() => {
      const shareFetch = await fetch(`/share/${sharedFileId}`)
      shareData.value = await shareFetch.json()
    })()

    const isLoading = computed(() => {
      return shareData.value == null || hobbitsByUser.value == null
    })

    return {
      sharedFileId,
      hobbitsByUser,
      shareData,
      isLoading,
    }
  },
})
</script>

<template>
  <h1>Create Record from Shared File...</h1>
  <Loading v-if="isLoading"/>
  <template v-if="!isLoading">
  <label for="hobbit">Hobbit</label>
  <select name="hobbit" id="hobbit">
    <option v-for="hobbit in hobbitsByUser" :key="hobbit.id" :value="hobbit.id" >{{ hobbit.name }}</option>
  </select>
      <p><strong>Time: </strong>{{shareData?.Time}}</p>
      <p><strong>Moving Time: </strong>{{shareData?.MovingTime / 1000 / 1000 / 1000 }}s</p>
      <p><strong>Distance: </strong>{{shareData?.Distance}}m</p>
    <Button type="primary" value="Create Record"/>
  </template>
</template>

<style lang="scss" scoped>
label {
  display: block;
}
select {
  width: 100%;
  height: 3em;
}
</style>
