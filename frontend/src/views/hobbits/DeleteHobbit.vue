<template>
  <div>
    <div class="header">
      <h1>Delete Hobbit</h1>
    </div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
      <FormWrapper>
        <form>
          <div>
            <p>Id: {{hobbit.id}}</p>
            <p>Name: {{hobbit.name}}</p>
            <p>Description: {{hobbit.description}}</p>
            <Button value="Delete Hobbit" @click="deleteHobbit()" type="primary" :loading="submitting" />
            <Button value="Go back" @click="goBack()" />
          </div>
        </form>
      </FormWrapper>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import Button from '../../components/form/Button.vue'
import Loading from '../../components/Icons/LoadingIcon.vue'
import FormWrapper from '../../components/form/FormWrapper.vue'
import { useHobbitsStore } from '../../store/hobbits'
import { useHobbitFromRoute } from '../../composables/hobbitFromRoute'
import { useRouter } from 'vue-router'

export default defineComponent({
  name: 'DeleteHobbit',
  components: {
    Button,
    Loading,
    FormWrapper,
  },
  setup() {
    const hobbits = useHobbitsStore()
    const router = useRouter()

    const submitting = ref(false)

    const { hobbit, id } = useHobbitFromRoute()

    const goBack = () => {
      router.push(`/hobbits/${id.value}`)
    }

    const goHome = () => {
      router.push('/')
    }

    const deleteHobbit = async() => {
      submitting.value = true
      await hobbits.deleteHobbit({
        id: id.value,
      })
      submitting.value = false
      goHome()
    }

    return {
      submitting,
      hobbit,
      goBack,
      deleteHobbit,
    }
  },
})
</script>
