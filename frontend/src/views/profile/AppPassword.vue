<template>
  <div>
    <Teleport to="#dialog-target">
      <!-- TODO: refactor this into it's own component -->
      <Dialog :shown="deleteDialog.shown">
        <FormWrapper>
          <form>
            <p>Do you really want to delete this app password?</p>
            <p>Description: {{ deleteDialog.appPassword?.description }}</p>
            <div>
              <Button type="primary" value="delete" @click="deleteAppPassword" :loading="deleteDialog.loading" />
              <Button value="cancel" @click="closeDeleteDialog" />
            </div>
          </form>
        </FormWrapper>
      </Dialog>
      <Dialog :shown="addDialog.shown">
        <FormWrapper>
          <form>
            <p>Create a new app password</p>
            <div>
              <label for="description">Description:</label>
              <input id="description" type="text" v-model="addDialog.description" />
            </div>
            <template v-if="addDialog.password">
              <p>Your new password is:</p>
              <input type="text" readonly :value="addDialog.password" />
            </template>
            <div>
              <Button type="primary" value="add" @click="addAppPassword" :loading="addDialog.loading"
                v-if="!addDialog.password" />
              <Button value="cancel" @click="closeAddDialog" />
            </div>
          </form>
        </FormWrapper>
      </Dialog>
    </Teleport>
    <h1>App Passwords</h1>
    <div>
      <div class="add-app-password cursor-pointer" @click="openAddDialog" tabindex="0">
        <Add class="h-24" />Add App Password...
      </div>
    </div>
    <AppPasswordItem v-for="appPassword in appPasswords" :key="`appPassword-${appPassword.id}`" :appPassword="appPassword"
      @delete="openDeleteDialog($event)" />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import AppPasswordItem from '../../components/profile/AppPasswordItem.vue'
import Dialog from '../../components/Dialog.vue'
import Button from '../../components/form/Button.vue'
import FormWrapper from '../../components/form/FormWrapper.vue'
import { AppPassword } from '../../models'
import { PlusIcon as Add } from '@heroicons/vue/24/outline'
import { useAppPasswordStore } from '../../store/profile'
import { storeToRefs } from 'pinia'

export default defineComponent({
  components: {
    AppPasswordItem,
    Dialog,
    Button,
    FormWrapper,
    Add,
  },
  setup() {
    const appPasswordStore = useAppPasswordStore()

    const deleteDialog = ref({
      appPassword: undefined as AppPassword | undefined,
      shown: false,
      loading: false,
    })

    const addDialog = ref({
      description: '',
      shown: false,
      loading: false,
      password: '',
    })

    const { appPasswords } = storeToRefs(appPasswordStore)

    const openDeleteDialog = ({ id }: { id: string }) => {
      deleteDialog.value.appPassword = appPasswords.value.find((appPassword) => {
        return appPassword.id === id
      })
      deleteDialog.value.shown = true
    }

    const closeDeleteDialog = () => {
      deleteDialog.value.shown = false
      deleteDialog.value.appPassword = undefined
    }

    const deleteAppPassword = async() => {
      deleteDialog.value.loading = true
      if (deleteDialog.value.appPassword?.id) {
        await appPasswordStore.deleteAppPassword({ id: deleteDialog.value.appPassword?.id })
      }
      deleteDialog.value.loading = false
      closeDeleteDialog()
    }

    const addAppPassword = async() => {
      addDialog.value.loading = true
      const newPassword = await appPasswordStore.postAppPassword({
        description: addDialog.value.description,
      })
      addDialog.value.loading = true
      addDialog.value.password = newPassword
    }

    const openAddDialog = () => {
      addDialog.value.description = ''
      addDialog.value.password = ''
      addDialog.value.shown = true
    }
    const closeAddDialog = () => {
      addDialog.value.shown = false
      addDialog.value.description = ''
      addDialog.value.password = ''
    }

    appPasswordStore.fetchAppPasswords()

    return {
      deleteDialog,
      addDialog,
      appPasswords,
      openDeleteDialog,
      closeDeleteDialog,
      deleteAppPassword,
      addAppPassword,
      openAddDialog,
      closeAddDialog,
    }
  },
})
</script>

<style scoped>
.add-app-password {
  display: flex;
  align-items: center;
}
</style>
