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
                <VButton
                  type="primary"
                  value="delete"
                  @click="deleteRecord"
                />
                <VButton
                  value="cancel"
                  @click="closeDeleteRecordDialog"
                />
              </div>
            </form>
          </FormWrapper>
        </DDialog>
      </Teleport>
      <div>
        <div class="header">
          <div>
            <h1>{{ hobbit.name }}</h1>
            <div class="by">
              by {{ hobbit.user.username }}
            </div>
            <div>{{ hobbit.description }}</div>
          </div>
          <div>
            <img :src="hobbit.image">
          </div>
        </div>
        <div>
          <IconBar>
            <template #left>
              <router-link
                v-if="isAuthenticated && userId === hobbit.user.id"
                v-slot="{ navigate }"
                :to="`/hobbits/${id}/records/add`"
                custom
              >
                <span @click="navigate">
                  <Add /><span>Add Record... </span>
                </span>
              </router-link>
            </template>
            <template #right>
              <router-link
                v-if="isAuthenticated && userId === hobbit.user.id"
                v-slot="{ navigate }"
                :to="`/hobbits/${id}/edit`"
                custom
              >
                <span @click="navigate">
                  <Pencil class="h-20 cursor-pointer" />
                </span>
              </router-link>
              <router-link
                v-if="isAuthenticated && userId === hobbit.user.id"
                v-slot="{ navigate }"
                :to="`/hobbits/${id}/delete`"
                custom
              >
                <span @click="navigate">
                  <Trash class="h-20 cursor-pointer" />
                </span>
              </router-link>
            </template>
          </IconBar>
          <Line
            :options="{ scales: { x: { type: 'time' }, y: { min: 0 } } }"
            :data="chartData"
          />
          <table>
            <thead>
              <tr>
                <td>value</td>
                <td>comment</td>
                <td>date</td>
                <td name="actions" />
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="record in (records || []).toReversed()"
                :key="`record-${record.id}`"
              >
                <td>{{ record.value }}</td>
                <td>{{ record.comment }}</td>
                <td>{{ formatDate(record.timestamp) }}</td>
                <td class="table-actions">
                  <Pencil
                    class="h-20 cursor-pointer"
                    tabindex="0"
                    @click="goToEditRecord($event, record)"
                  />
                  <Trash
                    class="h-20 cursor-pointer"
                    tabindex="0"
                    @click="openDeleteRecordDialog($event, record)"
                  />
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
import { computed, defineComponent, ref, watch } from 'vue'
import { NumericRecord } from '../../models'
import Loading from '../../components/Icons/LoadingIcon.vue'
import VButton from '../../components/form/Button.vue'
import DDialog from '../../components/Dialog.vue'
import FormWrapper from '../../components/form/FormWrapper.vue'
import { PencilIcon as Pencil, TrashIcon as Trash } from '@heroicons/vue/24/outline'
import { useAuthStore } from '../../store/auth'
import { storeToRefs } from 'pinia'
import { useHobbitsStore } from '../../store/hobbits'
import { useHobbitFromRoute } from '../../composables/hobbitFromRoute'
import { useRouter } from 'vue-router'
import IconBar from '../../components/IconBar.vue'
import Add from '../../components/Icons/AddIcon.vue'
import { Line } from 'vue-chartjs'
import {
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  TimeScale,
  Title,
  Tooltip,
  _adapters,
} from 'chart.js'
import temporalAdapter from 'chartjs-adapter-temporal'
import { formatDate } from '../../utils/date-utils'

_adapters._date.override(temporalAdapter);
ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement, TimeScale)


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
    Line,
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
    const records = computed(() => {
      return hobbits.getRecordsByHobbitId(id.value);
    });

    const { authenticated: isAuthenticated, userId } = storeToRefs(auth)

    watch(hobbit, () => {
      hobbits.fetchRecords(id.value);
    });

    const chartData = computed(() => {
      // cut off 1 year
      const yearAgo = Temporal.Now.zonedDateTimeISO().subtract({ years: 1 }).toInstant();
      return {
        datasets: [{
          label: hobbit.value.name,
          borderColor: getComputedStyle(document.body).getPropertyValue('--primary-dark'),
          backgroundColor: getComputedStyle(document.body).getPropertyValue('--primary'),
          data: records.value
            ?.filter(record => Temporal.Instant.from(record.timestamp).until(yearAgo).sign === -1)
            .toSorted((recordA, recordB) => Temporal.Instant.compare(Temporal.Instant.from(recordA.timestamp), Temporal.Instant.from(recordB.timestamp)))
            .map((record) => {
              return { x: record.timestamp as unknown as number, y: record.value }
            }),
        }],
      }
    })

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

    return {
      deleteDialog,
      hobbit,
      records,
      id,
      isAuthenticated,
      userId,
      chartData,
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
