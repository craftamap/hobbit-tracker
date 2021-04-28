<template>
  <div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
      <div>
      <div class="header">
          <div>
            <h1>{{hobbit.name}} - Add record</h1>
            <div class="by">by {{hobbit.user.username}}</div>
            <div>
            {{hobbit.description}}
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
              <Button value="Add record" @click="dispatchPostRecord()" type="primary" :loading="submitting"/>
              <Button value="Go back" @click="goBack()"/>
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
import { Hobbit } from '@/models'
import Loading from '@/components/Loading.vue'
import moment from 'moment'
import Button from '@/components/form/Button.vue'

export default defineComponent({
  components: {
    Loading,
    Button,
    FormWrapper,
  },
  computed: {
    id(): number {
      return Number(this.$route.params.id)
    },
    hobbit(): Hobbit {
      return this.$store.getters.getHobbitById(Number(this.$route.params.id))
    },
  },
  created() {
    if (!this.hobbit) {
      this.dispatchFetchHobbit()
    }
  },
  data() {
    return {
      submitting: false,
      data: {
        timestamp: this.getToday(),
        value: 10,
        comment: '',
      },
    }
  },
  methods: {
    getToday() {
      return moment().format('YYYY-MM-DDTHH:mm')
    },
    dispatchPostRecord() {
      this.submitting = true
      this.$store.dispatch('postRecord', {
        id: this.id,
        timestamp: moment(this.data.timestamp).toDate(),
        value: Number(this.data.value),
        comment: this.data.comment,
      }).then(() => {
        return Promise.all([
          this.$store.dispatch('fetchRecords', this.id),
          this.$store.dispatch('fetchHeatmapData', this.id),
        ])
      }).then(() => {
        this.submitting = false
      })
    },
    dispatchFetchHobbit() {
      this.$store.dispatch('fetchHobbit', { id: this.id })
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
