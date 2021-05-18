<template>
  <div>
    <Teleport to="#dialog-target">
      <!-- TODO: refactor this into it's own component -->
      <Dialog :shown="deleteDialog.shown">
        <FormWrapper>
          <form>
          <p>Do you really want to delete this app password?</p>
          <p>Description: {{deleteDialog.appPassword.description}}</p>
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
            <Button type="primary" value="add" @click="addAppPassword" :loading="addDialog.loading" v-if="!addDialog.password" />
            <Button value="cancel" @click="closeAddDialog" />
          </div>
          </form>
        </FormWrapper>
      </Dialog>
    </Teleport>
    <h1>App Passwords</h1>
    <div>
      <div class="add-app-password cursor-pointer" @click="openAddDialog" tabindex="0">
        <Add class="h-24" />
        Add App Password...
      </div>
    </div>
    <AppPasswordItem v-for="appPassword in appPasswords" :key="appPassword.ID" :appPassword="appPassword" @delete="openDeleteDialog($event)" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { createNamespacedHelpers } from 'vuex'
import { ProfileState } from '@/store/modules/profile'
import AppPasswordItem from '@/components/profile/AppPasswordItem.vue'
import Dialog from '@/components/Dialog.vue'
import Button from '@/components/form/Button.vue'
import FormWrapper from '@/components/form/FormWrapper.vue'
import { AppPassword } from '@/models'
import { PlusIcon as Add } from '@heroicons/vue/outline'

const { mapState, mapActions } = createNamespacedHelpers('profile')

export default defineComponent({
  data(): {
    deleteDialog: {
      appPassword: AppPassword | undefined;
      shown: boolean;
      loading: boolean;
    };
    addDialog: {
      description: string;
      shown: boolean;
      loading: boolean;
      password: string;
    };
    } {
    return {
      deleteDialog: {
        appPassword: undefined,
        shown: false,
        loading: false,
      },
      addDialog: {
        description: '',
        shown: false,
        loading: false,
        password: '',
      },
    }
  },
  created() {
    this.fetchAppPasswords()
  },
  components: {
    AppPasswordItem,
    Dialog,
    Button,
    FormWrapper,
    Add,
  },
  computed: {
    ...mapState({
      appPasswords: state => (state as ProfileState).apppassword.apppasswords,
    }),
  },
  methods: {
    openDeleteDialog({ id }: {id: string}) {
      this.deleteDialog.appPassword = this.appPasswords.find((appPassword) => {
        return appPassword.id === id
      })
      this.deleteDialog.shown = true
    },
    closeDeleteDialog() {
      this.deleteDialog.shown = false
      this.deleteDialog.appPassword = undefined
    },
    async deleteAppPassword() {
      this.deleteDialog.loading = true
      await this._deleteAppPasswords({ id: this.deleteDialog.appPassword?.id })
      this.deleteDialog.loading = false
      this.closeDeleteDialog()
    },
    async addAppPassword() {
      this.addDialog.loading = true
      const newPassword = await this._addAppPassword({
        description: this.addDialog.description,
      })
      this.addDialog.loading = false
      this.addDialog.password = newPassword
    },
    openAddDialog() {
      this.addDialog.description = ''
      this.addDialog.password = ''
      this.addDialog.shown = true
    },
    closeAddDialog() {
      this.addDialog.shown = false
      this.addDialog.password = ''
      this.addDialog.description = ''
    },
    ...mapActions({
      _deleteAppPasswords: 'deleteAppPassword',
      _addAppPassword: 'postAppPassword',
      fetchAppPasswords: 'fetchAppPasswords',
    }),
  },
})
</script>

<style lang="scss" scoped>
.add-app-password {
  display: flex;
  align-items: center;
}
</style>
