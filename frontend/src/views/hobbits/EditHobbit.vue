<template>
  <div>
    <div class="header">
      <h1>Add Hobbit</h1>
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
              <Button value="Add Hobbit" @click="dispatchPutHobbit()" type="primary" :loading="submitting"/>
              <Button value="Go back" @click="goBack()"/>
            </div>
          </form>
        </FormWrapper>
      </template>
  </div>
</template>

<script lang="ts">
import { Hobbit } from '@/models'
import { defineComponent } from 'vue'
import Button from '../../components/form/Button.vue'
import Loading from '../../components/Loading.vue'
import FormWrapper from '@/components/form/FormWrapper.vue'

export default defineComponent({
  name: 'AddHobbit',
  components: {
    Button,
    Loading,
    FormWrapper,
  },
  created () {
    if (!this.hobbit) {
      this.$store.dispatch('fetchHobbit', { id: this.id })
    }
  },
  data () {
    return {
      submitting: false,
      form: {
        name: '',
        description: '',
        image: '',
      },
    }
  },
  computed: {
    id (): number {
      return Number(this.$route.params.id)
    },
    hobbit (): Hobbit {
      return this.$store.getters.getHobbitById(this.id)
    },
  },
  watch: {
    hobbit: {
      handler (newValue: Hobbit, oldValue) {
        console.log('WATCH', newValue, oldValue)
        if (newValue && newValue !== oldValue) {
          this.form.name = newValue.name
          this.form.description = newValue.description
          this.form.image = newValue.image
        }
      },
      immediate: true,
    },
  },
  methods: {
    goBack () {
      this.$router.push(`/hobbits/${this.id}`)
    },
    readUploadedFileAsDataURL (inputFile: File): Promise<string> {
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
    async changeImage (event: Event) {
      // TODO: Add validation
      const fileList = (event?.target as any).files as FileList
      const firstFile = fileList[0]
      this.form.image = await this.readUploadedFileAsDataURL(firstFile)
      console.log(this.form.image)
    },
    dispatchPutHobbit () {
      this.submitting = true
      this.$store.dispatch('putHobbit', {
        id: this.id,
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
