<template>
  <div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
      <Teleport to="#dialog-target">
        <!-- TODO: refactor this into it's own component -->
        <DDialog :shown="deleteDialog.shown">
          <FormWrapper>
            <form>
              <p>Do you really want to delete this record?</p>
              <div>
                <VButton type="primary" value="delete" @click="deleteRecord" />
                <VButton value="cancel" @click="closeDeleteRecordDialog" />
              </div>
            </form>
          </FormWrapper>
        </DDialog>
      </Teleport>
      <div>
        <div class="header">
          <div>
            <h1>{{ hobbit.name }}</h1>
            <div class="by">by {{ hobbit.user.username }}</div>
            <div>{{ hobbit.description }}</div>
          </div>
          <div>
            <img :src="hobbit.image" />
          </div>
        </div>
        <div>
          <IconBar>
            <template v-slot:left>
              <router-link v-if="isAuthenticated && userId === hobbit.user.id" :to="`/hobbits/${id}/records/add`" custom
                v-slot="{ navigate }">
                <span @click="navigate">
                  <Add /><span>Add Record... </span>
                </span>
              </router-link>
            </template>
            <template v-slot:right>
              <router-link v-if="isAuthenticated && userId === hobbit.user.id" :to="`/hobbits/${id}/edit`" custom
                v-slot="{ navigate }">
                <span @click="navigate">
                  <Pencil class="h-20 cursor-pointer" />
                </span>
              </router-link>
              <router-link v-if="isAuthenticated && userId === hobbit.user.id" :to="`/hobbits/${id}/delete`" custom
                v-slot="{ navigate }">
                <span @click="navigate">
                  <Trash class="h-20 cursor-pointer" />
                </span>
              </router-link>
            </template>
          </IconBar>
          <table>
            <thead>
              <tr>
                <td>value</td>
                <td>comment</td>
                <td>date</td>
                <td name="actions"></td>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in (hobbit.records || []).slice().reverse()" :key="`record-${record.id}`">
                <td>{{ record.value }}</td>
                <td>{{ record.comment }}</td>
                <td>{{ formatDate(record.timestamp) }}</td>
                <td class="table-actions">
                  <Pencil class="h-20 cursor-pointer" v-on:click="goToEditRecord($event, record)" tabindex="0" />
                  <Trash class="h-20 cursor-pointer" v-on:click="openDeleteRecordDialog($event, record)" tabindex="0" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { NumericRecord } from '../../models'
import Loading from '../../components/Icons/LoadingIcon.vue'
import VButton from '../../components/form/Button.vue'
import DDialog from '../../components/Dialog.vue'
import FormWrapper from '../../components/form/FormWrapper.vue'
import moment from 'moment'
import { TrashIcon as Trash, PencilIcon as Pencil } from '@heroicons/vue/24/outline'
import { useAuthStore } from '../../store/auth'
import { storeToRefs } from 'pinia'
import { useHobbitsStore } from '../../store/hobbits'
import { useHobbitFromRoute } from '../../composables/hobbitFromRoute'
import { useRouter } from 'vue-router'
import IconBar from '../../components/IconBar.vue'
import Add from '../../components/Icons/AddIcon.vue'

export default defineComponent({
  name: 'HobbitView',
  components: {
    Loading,
    VButton,
    Trash,
    Pencil,
    DDialog,
    FormWrapper,
    IconBar,
    Add,
  },
  setup() {
    const auth = useAuthStore()
    const hobbits = useHobbitsStore()
    const router = useRouter()

    const deleteDialog = ref({
      shown: false,
      record: null as NumericRecord | null,
    })

    const { hobbit, id } = useHobbitFromRoute()

    const { authenticated: isAuthenticated, userId } = storeToRefs(auth)

    hobbits.fetchRecords(id.value)

    const openDeleteRecordDialog = (_: Event, record: NumericRecord) => {
      deleteDialog.value.record = record
      deleteDialog.value.shown = true
    }
    const closeDeleteRecordDialog = () => {
      deleteDialog.value.shown = false
      deleteDialog.value.record = null
    }

    const goToEditRecord = (_: Event, record: NumericRecord) => {
      const recordId = record.id
      return router.push(`/hobbits/${id.value}/records/${recordId}/edit`)
    }

    const deleteRecord = () => {
      return hobbits.deleteRecord({
        hobbitId: Number(id.value),
        recordId: Number(deleteDialog.value.record?.id),
      })
    }

    const formatDate = (date: string) => {
      return moment(date).format('YYYY-MM-DD HH:mm')
    }

    return {
      deleteDialog,
      hobbit,
      id,
      isAuthenticated,
      userId,
      goToEditRecord,
      deleteRecord,
      openDeleteRecordDialog,
      closeDeleteRecordDialog,
      formatDate,
    }
  },
})
</script>

<style scoped>
table {
  width: 100%;

  thead {
    font-weight: bold;
  }
}

h1 {
  margin: 0;
  font-size: 16pt;
}

.by {
  color: gray;
}

.header {
  display: flex;
  justify-content: space-between;

  img {
    width: 2rem;
    height: 2rem;
  }
}

table {
  border-collapse: collapse;

  th,
  td {
    padding: 16px 0px;
  }

  tr {
    border-bottom: solid 1px lightgray;
  }
}

.table-actions {
  text-align: right;
  width: 0%;
}
</style>
