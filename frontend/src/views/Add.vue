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
        <div class="form-wrapper">
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
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Hobbit } from '@/models'
import Loading from '@/components/Loading.vue'
import moment from 'moment'
import Button from '@/components/form/Button.vue'

export default defineComponent({
  components: {
    Loading,
    Button
  },
  computed: {
    id (): number {
      return Number(this.$route.params.id)
    },
    hobbit (): Hobbit {
      return this.$store.getters.getHobbitById(Number(this.$route.params.id))
    }
  },
  data () {
    return {
      submitting: false,
      data: {
        timestamp: this.getToday(),
        value: 10,
        comment: ''
      }
    }
  },
  methods: {
    getToday () {
      return moment().format('YYYY-MM-DDTHH:mm')
    },
    dispatchPostRecord () {
      this.submitting = true
      this.$store.dispatch('postRecord', {
        id: this.id,
        timestamp: moment(this.data.timestamp).toDate(),
        value: Number(this.data.value),
        comment: this.data.comment
      }).then(() => {
        return Promise.all([
          this.$store.dispatch('fetchRecords', this.id),
          this.$store.dispatch('fetchHeatmapData', this.id)
        ])
      }).then(() => {
        this.submitting = false
      })
    },
    goBack () {
      this.$router.push('/hobbits/' + this.id)
    }
  }
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
.form-wrapper {
  display: flex;
  justify-content: center;
  justify-items: center;

  form {
    border-radius: 0.5rem;
    padding: 2rem;
    background: #eee;
    width: 300px;

    input, textarea {
      margin-bottom: 0.25rem;
      appearance: none;
      &:focus {
        outline: none;
      }
      border: none;
      line-height: 2em;
      padding: 5px;
      border-radius: 3px;
      display: block;
      margin-left: auto;
      margin-right: auto;
      width: 95%;

      &.submit {
        background: var(--ming);
        color: #fff;
        width: 100%;
      }
    }
  }
}
</style>
