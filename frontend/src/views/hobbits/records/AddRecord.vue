<template>
  <div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
      <div>
        <div class="header">
          <div>
            <h1>{{ hobbit.name }} - Add record</h1>
            <div class="by">by {{ hobbit.user.username }}</div>
            <div>
              {{ hobbit.description }}
            </div>
          </div>
          <div>
            <img :src="hobbit.image" />
          </div>
        </div>
        <FormWrapper>
          <form>
            <div>
              <label for="timestamp">When:</label>
              <input id="timestamp" name="timestamp" type="datetime-local" v-model="recordData.timestamp" />
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
              <Button value="Add record" @click="postRecord()" type="primary" :loading="submitting" />
              <Button value="Go back" @click="goBack()" />
            </div>
          </form>
        </FormWrapper>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import FormWrapper from '../../../components/form/FormWrapper.vue'
import { defineComponent, ref } from 'vue'
import Loading from '../../../components/Icons/LoadingIcon.vue'
import Button from '../../../components/form/Button.vue'
import { useHobbitsStore } from '../../../store/hobbits'
import { useHobbitFromRoute } from '../../../composables/hobbitFromRoute'
import { useRouter } from 'vue-router'
import { Temporal } from '@js-temporal/polyfill'

export default defineComponent({
  components: {
    Loading,
    Button,
    FormWrapper,
  },
  setup() {
    const hobbits = useHobbitsStore()
    const router = useRouter()

    const submitting = ref(false)
    const recordData = ref({
      timestamp: Temporal.Now.plainDateTimeISO().toString({fractionalSecondDigits: 0}),
      value: 10,
      comment: '',
    })

    const { id, hobbit } = useHobbitFromRoute()

    console.log(Temporal.Now.timeZoneId())
    const postRecord = () => {
      submitting.value = true
      hobbits.postRecord({
        id: id.value,
        timestamp: Temporal.PlainDateTime
          .from(recordData.value.timestamp)
          .toZonedDateTime(Temporal.Now.timeZoneId())
          .toInstant()
          .toString(),
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

    return {
      submitting,
      recordData,
      hobbit,
      postRecord,
      goBack,
    }
  },
})

</script>

<style scoped>
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
