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
          <div class="buttons">
            <router-link :to="`/hobbits/${$route.params.id}/add`"
              custom
              v-slot="{ navigate }">
              <VButton value="Add Record" type="primary" @click="navigate" />
            </router-link>
            <router-link :to="`/hobbits/${$route.params.id}/edit`"
              custom
              v-slot="{ navigate }">
              <VButton value="Edit" @click="navigate" />
            </router-link>
          </div>
          <table>
            <thead>
              <tr>
                <td>date</td>
                <td>value</td>
                <td>comment</td>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in (hobbit.records || []).slice().reverse()" :key="record.ID">
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
import VButton from '@/components/form/Button.vue'
import moment from 'moment'

export default defineComponent({
  name: 'Hobbit',
  components: {
    Loading,
    VButton
  },
  computed: {
    hobbit (): Hobbit {
      return this.$store.getters.getHobbitById(Number(this.$route.params.id))
    }
  },
  created () {
    if (!this.hobbit) {
      this.dispatchFetchHobbits().then(() => {
        this.dispathFetchRecords()
      })
    } else {
      this.dispathFetchRecords()
    }
  },
  methods: {
    dispatchFetchHobbits () {
      return this.$store.dispatch('fetchHobbits')
    },
    dispathFetchRecords () {
      return this.$store.dispatch('fetchRecords', Number(this.$route.params.id))
    },
    formatDate (date: string) {
      return moment(date).format('YYYY-MM-DD HH:mm')
    }
  }
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
</style>
