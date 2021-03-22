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
        <div>
          <form>
            <div>
              <label for="timestamp">When:</label>
              <input id="timestamp" name="timestamp" type="datetime-local" :value="getToday()" />
            </div>
            <div>
              <label for="value">Value:</label>
              <input type="number" name="number" id="number" value="10" />
            </div>
            <div>
              <label for="comment">Comment:</label>
              <textarea name="comment" id="comment" rows="5"></textarea>
            </div>
          </form>
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useStore } from '../store'
import { Hobbit } from '@/models'
import Loading from '@/components/Loading.vue'
import moment from 'moment'

export default defineComponent({
  components: {
    Loading
  },
  computed: {
    hobbit (): Hobbit {
      return useStore().getters.getHobbitById(Number(this.$route.params.id))
    }
  },
  created () {
    if (!this.hobbit) {
      this.dispatchFetchHobbits()
    }
  },
  methods: {
    dispatchFetchHobbits () {
      const store = useStore()
      store.dispatch('fetchHobbits')
    },
    getToday () {
      return moment().format('YYYY-MM-DDTHH:mm')
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
</style>
