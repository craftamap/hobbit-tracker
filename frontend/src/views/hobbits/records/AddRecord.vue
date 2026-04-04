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
            <div class="by">
              by {{ hobbit.user.username }}
            </div>
            <div>
              {{ hobbit.description }}
            </div>
          </div>
          <div>
            <img :src="hobbit.image">
          </div>
        </div>
        <FormWrapper>
          <form>
            <div>
              <label for="timestamp">When:</label>
              <input id="timestamp" v-model="recordData.timestamp" name="timestamp" type="datetime-local">
            </div>
            <div>
              <label for="value">Value:</label>
              <input id="number" v-model="recordData.value" type="number" name="number">
            </div>
            <div>
              <label for="comment">Comment:</label>
              <textarea id="comment" v-model="recordData.comment" name="comment" rows="5" />
            </div>
            <div>
              <Button value="Add record" type="primary" :loading="submitting" @click="postRecord()" />
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
import { fromDateTimeLocal, toDateTimeLocal } from '../../../utils/date-utils'

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
      timestamp: toDateTimeLocal(Temporal.Now.instant()),
      value: 10,
      comment: '',
    })

    const { id, hobbit } = useHobbitFromRoute()

    const postRecord = () => {
      submitting.value = true
      hobbits.postRecord({
        id: id.value,
        timestamp: fromDateTimeLocal(recordData.value.timestamp),
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
