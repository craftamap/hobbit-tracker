<template>
  <div>
    <template v-if="!hobbit">
      <Loading />
    </template>
    <template v-if="hobbit">
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
          <router-link :to="`/hobbits/${$route.params.id}/add`" class="add-record">Add Record</router-link>
          <table>
            <thead>
              <tr>
                <td>id</td>
                <td>date</td>
                <td>value</td>
                <td>comment</td>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in (hobbit.records || []).slice().reverse()" :key="record.ID">
                <td>{{record.ID}}</td>
                <td>{{formatDate(record.timestamp)}}</td>
                <td>{{record.value}}</td>
                <td>{{record.comment}}</td>
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
import { Hobbit } from '@/models'
import Loading from '@/components/Loading.vue'
import moment from 'moment'

export default defineComponent({
  name: 'Hobbit',
  components: {
    Loading
  },
  computed: {
    hobbit (): Hobbit {
      return this.$store.getters.getHobbitById(Number(this.$route.params.id))
    }
  },
  created () {
    if (!this.hobbit) {
      this.dispatchFetchHobbits()
    }
  },
  methods: {
    dispatchFetchHobbits () {
      this.$store.dispatch('fetchHobbits').then(() => {
        return this.$store.dispatch('fetchRecords', Number(this.$route.params.id))
      })
    },
    formatDate (date: string) {
      return moment(date).format('YYYY-MM-DD HH:mm')
    }
  }
})
</script>

<style lang="scss" scoped>
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

.add-record {
  margin-bottom: 0.5rem;
  margin-top: 0.5rem;
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
  background: var(--ming);
  color: #fff;
  width: 100%;
  text-align: center;
}
</style>
