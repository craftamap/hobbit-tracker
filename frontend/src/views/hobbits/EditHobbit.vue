<template>
  <div>
    <div class="header">
      <h1>Edit Hobbit</h1>
    </div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
      <FormWrapper>
        <form>
          <div>
            <label for="name">Hobbit name:</label>
            <input id="name" name="name" type="text" v-model="form.name" />
          </div>
          <div>
            <label for="description">Description:</label>
            <textarea name="description" id="description" rows="5" v-model="form.description"></textarea>
          </div>
          <div>
            <label for="image">Image:</label>
            <input id="image" name="image" type="file" @change="changeImage" />
          </div>
          <div>
            <Button value="Edit Hobbit" @click="putHobbit()" type="primary" :loading="submitting" />
            <Button value="Go back" @click="goBack()" />
          </div>
        </form>
      </FormWrapper>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue'
import Button from '../../components/form/Button.vue'
import Loading from '../../components/Icons/LoadingIcon.vue'
import FormWrapper from '../../components/form/FormWrapper.vue'
import { useHobbitsStore } from '../../store/hobbits'
import { useHobbitFromRoute } from '../../composables/hobbitFromRoute'
import { useRouter } from 'vue-router'

export default defineComponent({
  name: 'AddHobbit',
  components: {
    Button,
    Loading,
    FormWrapper,
  },
  setup() {
    const hobbits = useHobbitsStore()
    const router = useRouter()

    const submitting = ref(false)
    const form = ref({ name: '', description: '', image: '' })

    const { hobbit, id } = useHobbitFromRoute()

    watch(hobbit, (newValue, oldValue) => {
      console.log('foo 1')
      if (newValue && newValue !== oldValue) {
        console.log('foo 2')
        form.value.name = newValue.name
        form.value.description = newValue.description
        form.value.image = newValue.image
      }
    }, { immediate: true })

    const readUploadedFileAsDataURL = (inputFile: File): Promise<string> => {
      const temporaryFileReader = new FileReader()

      return new Promise((resolve, reject) => {
        temporaryFileReader.onerror = () => {
          temporaryFileReader.abort()
          reject(new DOMException('Problem parsing input file.'))
        }

        temporaryFileReader.onload = () => {
          resolve(temporaryFileReader.result as string)
        }
        temporaryFileReader.readAsDataURL(inputFile)
      })
    }
    const changeImage = async(event: Event) => {
      // TODO: Add validation
      const fileList = (event?.target as HTMLInputElement).files!
      const firstFile = fileList[0]
      form.value.image = await readUploadedFileAsDataURL(firstFile)
      console.log(form.value.image)
    }

    const goBack = () => {
      router.push(`/hobbits/${id.value}`)
    }

    const putHobbit = async() => {
      submitting.value = true
      await hobbits.putHobbit({
        id: id.value,
        name: form.value.name,
        description: form.value.description,
        image: form.value.image,
      })
      submitting.value = false
    }

    return {
      changeImage,
      submitting,
      form,
      hobbit,
      goBack,
      putHobbit,
    }
  },
})
</script>
