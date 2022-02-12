<template>
  <div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
      <div>
        <div class="header">
          <div>
            <h1>{{ hobbit.name }} - Edit record {{ recordId }}</h1>
            <div class="by">by {{ hobbit.user.username }}</div>
            <div>{{ hobbit.description }}</div>
          </div>
          <div>
            <img :src="hobbit.image" />
          </div>
        </div>
        <FormWrapper>
          <form>
            <div>
              <label for="timestamp">When:</label>
              <input
                id="timestamp"
                name="timestamp"
                type="datetime-local"
                v-model="recordData.timestamp"
              />
            </div>
            <div>
              <label for="value">Value:</label>
              <input type="number" name="number" id="number" v-model="recordData.value" />
            </div>
            <div>
              <label for="comment">Comment:</label>
              <textarea name="comment" id="comment" rows="5" v-model="recordData.comment"></textarea>
            </div>
            <div>
              <Button
                value="Edit record"
                @click="putRecord()"
                type="primary"
                :loading="submitting"
              />
              <Button value="Go back" @click="goBack()" />
            </div>
          </form>
        </FormWrapper>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import FormWrapper from '@/components/form/FormWrapper.vue'
import { computed, defineComponent, ref } from 'vue'
import Loading from '@/components/Icons/LoadingIcon.vue'
import moment from 'moment'
import Button from '@/components/form/Button.vue'
import { useHobbitsStore } from '@/store/hobbits'
import { useHobbitFromRoute } from '@/composables/hobbitFromRoute'
import { useRoute, useRouter } from 'vue-router'

export default defineComponent({
  components: {
    Loading,
    Button,
    FormWrapper,
  },
  setup() {
    const hobbits = useHobbitsStore()
    const router = useRouter()
    const route = useRoute()

    const submitting = ref(false)
    const recordData = ref({
      timestamp: '',
      value: 10,
      comment: '',
    })

    const { id, hobbit } = useHobbitFromRoute()
    const recordId = computed(() => {
      const recordIdParam = route.params.recordId
      if (Array.isArray(recordIdParam)) {
        if (recordIdParam.length > 0) {
          return Number(recordIdParam[0])
        }
        return -1
      }
      return Number(recordIdParam)
    })

    const record = computed(() => {
      return hobbits.getRecordById(id.value, recordId.value)
    })

    const parseAndFormatDate = (date: string | undefined) => {
      if (date) {
        return moment(date).format('YYYY-MM-DDTHH:mm')
      }
      return ''
    }

    const putRecord = () => {
      submitting.value = true
      hobbits.putRecord({
        hobbitId: id.value,
        recordId: recordId.value,
        timestamp: moment(recordData.value.timestamp).toDate(),
        value: Number(recordData.value.value),
        comment: recordData.value.comment,
      }).then(() => {
        return Promise.all([
          hobbits.fetchRecords(id.value),
          hobbits.fetchHeatmapData(id.value),
        ])
      }).then(() => {
        submitting.value = false
      })
    }

    const goBack = () => {
      router.push(`/hobbits/${id.value}`)
    }

    hobbits.fetchRecords(id.value).then(() => {
      recordData.value = Object.assign({}, record.value, { timestamp: parseAndFormatDate(record.value?.timestamp) })
    })

    return {
      recordData,
      submitting,
      hobbit,
      recordId,
      goBack,
      putRecord,
    }
  },
})
</script>

<style lang="scss" scoped>
label {
  width: 8rem;
  display: inline-block;
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
</style>
