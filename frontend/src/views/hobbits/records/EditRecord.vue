<template>
  <div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
      <div>
        <div class="header">
          <div>
            <h1>{{ hobbit.name }} - Edit record {{ id }}</h1>
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
              <input id="timestamp" name="timestamp" type="datetime-local" v-model="data.timestamp" />
            </div>
            <div>
              <label for="value">Value:</label>
              <input type="number" name="number" id="number" v-model="data.value" />
            </div>
            <div>
              <label for="comment">Comment:</label>
              <textarea name="comment" id="comment" rows="5" v-model="data.comment"></textarea>
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
import { defineComponent } from 'vue'
import { Hobbit, NumericRecord } from '@/models'
import Loading from '@/components/Icons/LoadingIcon.vue'
import moment from 'moment'
import Button from '@/components/form/Button.vue'
import { useHobbitsStore } from '@/store/hobbits'
import { mapActions, mapState } from 'pinia'

export default defineComponent({
  components: {
    Loading,
    Button,
    FormWrapper,
  },
  computed: {
    ...mapState(useHobbitsStore, {
      hobbitById: 'getHobbitById',
      recordById: 'getRecordById',
    }),
    id(): number {
      return Number(this.$route.params.hobbitId)
    },
    recordId(): number {
      return Number(this.$route.params.recordId)
    },
    hobbit(): Hobbit {
      return this.hobbitById(this.id)
    },
    record(): NumericRecord | undefined {
      if (this.id) {
        return this.recordById(this.id, this.recordId)
      }
      return undefined
    },
  },
  async created() {
    if (!this.hobbit) {
      await this.fetchHobbit()
    }
    if (!this.record) {
      await this.fetchRecords()
    }

    console.log(this.data, this.record)

    this.data = Object.assign({}, this.record, { timestamp: this.parseAndFormatDate(this.record?.timestamp) })
  },
  data() {
    return {
      submitting: false,
      data: {
        timestamp: '',
        value: 0,
        comment: '',
      },
    }
  },
  methods: {
    ...mapActions(useHobbitsStore, {
      _putRecord: 'putRecord',
      _fetchHobbit: 'fetchHobbit',
      _fetchRecords: 'fetchRecords',
      fetchHeatmapData: 'fetchHeatmapData',
    }),
    parseAndFormatDate(date: string | undefined) {
      if (date) {
        return moment(date).format('YYYY-MM-DDTHH:mm')
      }
      return ''
    },
    putRecord() {
      this.submitting = true
      this._putRecord({
        hobbitId: this.id,
        recordId: this.recordId,
        timestamp: moment(this.data.timestamp).toDate(),
        value: Number(this.data.value),
        comment: this.data.comment,
      }).then(() => {
        return Promise.all([
          this._fetchRecords(this.id),
          this.fetchHeatmapData(this.id),
        ])
      }).then(() => {
        this.submitting = false
      })
    },
    async fetchHobbit() {
      if (this.id) {
        return this._fetchHobbit(this.id)
      }
    },
    async fetchRecords() {
      return this._fetchRecords(Number(this.id))
    },
    goBack() {
      this.$router.push('/hobbits/' + this.id)
    },
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
