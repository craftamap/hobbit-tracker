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
          <Button value="Add Hobbit" @click="dispatchPostHobbit()" type="primary" :loading="submitting"/>
          <Button value="Go back" @click="goBack()"/>
        </div>
      </form>
    </FormWrapper>
  </div>
</template>

<script lang="ts">
import FormWrapper from '@/components/form/FormWrapper.vue'
import { useHobbitsStore } from '@/store/hobbits'
import { mapActions } from 'pinia'
import { defineComponent } from 'vue'
import Button from '../../components/form/Button.vue'

export default defineComponent({
  name: 'AddHobbit',
  components: {
    Button,
    FormWrapper,
  },
  data() {
    return {
      submitting: false,
      form: {
        name: '',
        description: '',
        image: '',
      },
    }
  },
  methods: {
    ...mapActions(useHobbitsStore, {
      _postHobbit: 'postHobbit',
    }),
    goBack() {
      this.$router.push('/')
    },
    readUploadedFileAsDataURL(inputFile: File): Promise<string> {
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
    },
    async changeImage(event: Event) {
      // TODO: Add validation
      const fileList = (event?.target as any).files as FileList
      const firstFile = fileList[0]
      this.form.image = await this.readUploadedFileAsDataURL(firstFile)
      console.log(this.form.image)
    },
    async dispatchPostHobbit() {
      this.submitting = true
      return this._postHobbit({
        name: this.form.name,
        description: this.form.description,
        image: this.form.image,
      }).then(() => {
        this.submitting = false
      })
    },
  },
})
</script>
