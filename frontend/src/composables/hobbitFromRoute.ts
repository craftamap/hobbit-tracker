import { useHobbitsStore } from '../store/hobbits'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

// useHobbitFromRoute extracts the hobbitId from a route, and returns the associated hobbit from the store. If the hobbit is not in store yet, it's also fetched.
export const useHobbitFromRoute = () => {
  const route = useRoute()
  const hobbits = useHobbitsStore()

  const id = computed(() => {
    const hobbitIdParam = route.params.hobbitId
    if (Array.isArray(hobbitIdParam)) {
      if (hobbitIdParam.length > 0) {
        return Number(hobbitIdParam[0])
      }
      return -1
    }
    return Number(hobbitIdParam)
  })
  const hobbit = computed(() => {
    return hobbits.getHobbitById(id.value)
  })

  if (!hobbit.value) {
    hobbits.fetchHobbit(id.value)
  }

  return {
    id,
    hobbit,
  }
}
