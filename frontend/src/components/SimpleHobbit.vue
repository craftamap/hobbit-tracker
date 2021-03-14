<template>
  <div class="card">
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
  </div>
</template>

<script lang="ts">
import { Hobbit } from '../models/index'
import { defineComponent, PropType } from 'vue'

export default defineComponent({
  props: {
    hobbit: Object as PropType<Hobbit>
  }
})
</script>

<style lang="scss" scoped>

.card {
  border-radius: 0.5rem;
  box-shadow: 0px 0px 5px -2px #000000;
  padding: 1rem;
  margin: 0.5rem 0;

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
}

</style>
