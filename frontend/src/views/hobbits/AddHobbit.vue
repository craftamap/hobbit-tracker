<template>
  <div>
    <div class="header">
      <h1>Add Hobbit</h1>
    </div>
        <div class="form-wrapper">
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
        </div>
        </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Button from '../../components/form/Button.vue'

export default defineComponent({
  name: 'AddHobbit',
  components: {
    Button
  },
  data () {
    return {
      submitting: false,
      form: {
        name: '',
        description: '',
        image: ''
      }
    }
  },
  methods: {
    goBack () {
      this.$router.push('/')
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
    dispatchPostHobbit () {
      this.submitting = true
      this.$store.dispatch('postHobbit', {
        name: this.form.name,
        description: this.form.description,
        image: this.form.image
      }).then(() => {
        this.submitting = false
      })
    }
  }
})
</script>

<style lang="scss" scoped>
.form-wrapper {
  display: flex;
  justify-content: center;
  justify-items: center;

  form {
    border-radius: 0.5rem;
    padding: 2rem;
    background: #eee;
    width: 300px;

    input, textarea {
      margin-bottom: 0.25rem;
      appearance: none;
      &:focus {
        outline: none;
      }
      border: none;
      line-height: 2em;
      padding: 5px;
      border-radius: 3px;
      display: block;
      margin-left: auto;
      margin-right: auto;
      width: 95%;

      &.submit {
        background: var(--ming);
        color: #fff;
        width: 100%;
      }
    }
  }
}
</style>
