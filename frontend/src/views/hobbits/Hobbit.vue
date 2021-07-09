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
            <h1>{{hobbit.name}}</h1>
            <div class="by">by {{hobbit.user.username}}</div>
            <div>
            {{hobbit.description}}
            </div>
          </div>
          <div>
            <img :src="hobbit.image" />
          </div>
        </div>
        <div>
          <div class="buttons" v-if="isAuthenticated && userId === hobbit.user.id">
            <router-link :to="`/hobbits/${id}/records/add`"
              custom
              v-slot="{ navigate }">
              <VButton value="Add Record" type="primary" @click="navigate" />
            </router-link>
            <router-link :to="`/hobbits/${id}/edit`"
              custom
              v-slot="{ navigate }">
              <VButton value="Edit" @click="navigate" />
            </router-link>
          </div>
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
              <tr v-for="record in (hobbit.records || []).slice().reverse()" :key="record.id">
                <td>{{record.value}}</td>
                <td>{{record.comment}}</td>
                <td>{{formatDate(record.timestamp)}}</td>
                <td class="table-actions">
                  <Pencil class="h-20 cursor-pointer" v-on:click="editRecord($event, record)" tabindex="0" />
                  <Trash class="h-20 cursor-pointer"  v-on:click="openDeleteRecordDialog($event, record)" tabindex="0" />
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
import { defineComponent } from 'vue'
import { Hobbit, NumericRecord } from '@/models'
import Loading from '@/components/Loading.vue'
import VButton from '@/components/form/Button.vue'
import DDialog from '@/components/Dialog.vue'
import FormWrapper from '@/components/form/FormWrapper.vue'
import moment from 'moment'
import { TrashIcon as Trash, PencilIcon as Pencil } from '@heroicons/vue/outline'
import { createNamespacedHelpers } from 'vuex'
import { AuthenticationState } from '@/store/modules/auth'

const { mapState: authMapState } = createNamespacedHelpers('auth')
const { mapActions: mapHobbitsActions, mapGetters: mapHobbitsGetters } = createNamespacedHelpers('hobbits')

export default defineComponent({
  name: 'Hobbit',
  components: {
    Loading,
    VButton,
    Trash,
    Pencil,
    DDialog,
    FormWrapper,
  },
  data(): {
    deleteDialog: {
      shown: boolean;
      record: NumericRecord | null;
    };
    } {
    return {
      deleteDialog: {
        shown: false,
        record: null,
      },
    }
  },
  computed: {
    id(): number {
      return Number(this.$route.params.hobbitId)
    },
    ...mapHobbitsGetters({
      hobbitById: 'getHobbitById',
    }),
    hobbit(): Hobbit {
      return this.hobbitById(this.id)
    },
    ...authMapState({
      isAuthenticated: state => (state as AuthenticationState).authenticated,
      userId: state => (state as AuthenticationState).userId,
    }),
  },
  created() {
    if (!this.hobbit) {
      this.fetchHobbit().then(() => {
        this.fetchRecords()
      })
    } else {
      this.fetchRecords()
    }
  },
  methods: {
    ...mapHobbitsActions({
      _putRecord: 'putRecord',
      _fetchRecords: 'fetchRecords',
      _fetchHobbit: 'fetchHobbit',
      _deleteRecord: 'deleteRecord',
    }),
    fetchHobbit() {
      return this._fetchHobbit({ id: Number(this.id) })
    },
    fetchRecords() {
      return this._fetchRecords(Number(this.id))
    },
    formatDate(date: string) {
      return moment(date).format('YYYY-MM-DD HH:mm')
    },
    editRecord(_: Event, record: NumericRecord) {
      const recordId = record.id
      return this.$router.push(`/hobbits/${this.id}/records/${recordId}/edit`)
    },
    deleteRecord() {
      return this._deleteRecord({
        hobbitId: Number(this.id),
        recordId: Number(this.deleteDialog.record?.id),
      })
    },
    openDeleteRecordDialog(_: Event, record: NumericRecord) {
      this.deleteDialog.record = record
      this.deleteDialog.shown = true
    },
    closeDeleteRecordDialog() {
      this.deleteDialog.shown = false
      this.deleteDialog.record = null
    },
  },
})
</script>

<style lang="scss" scoped>

  .buttons {
    display: flex;
  }

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
    th, td {
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
