<template>
  <div>
    <div class="header">
      <h1>Add Hobbit</h1>
    </div>
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
          <Button value="Add Hobbit" @click="postHobbit()" type="primary" :loading="submitting"/>
          <Button value="Go back" @click="goBack()"/>
        </div>
      </form>
    </FormWrapper>
  </div>
</template>

<script lang="ts">
import FormWrapper from '@/components/form/FormWrapper.vue'
import { useHobbitsStore } from '@/store/hobbits'
import { defineComponent, ref } from 'vue'
import { useRouter } from 'vue-router'
import Button from '../../components/form/Button.vue'

export default defineComponent({
  name: 'AddHobbit',
  components: {
    Button,
    FormWrapper,
  },
  setup() {
    const router = useRouter()
    const hobbits = useHobbitsStore()

    const submitting = ref(false)
    const form = ref({ name: '', description: '', image: '' })

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

    const postHobbit = async() => {
      submitting.value = true
      return hobbits.postHobbit({
        name: form.value.name,
        description: form.value.description,
        image: form.value.image,
      }).then(() => {
        submitting.value = false
      })
    }

    const goBack = () => {
      router.push('/')
    }

    return {
      submitting,
      form,
      changeImage,
      goBack,
      postHobbit,
    }
  },
})
</script>
