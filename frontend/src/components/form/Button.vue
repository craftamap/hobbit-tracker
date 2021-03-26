<template>
  <button :class="[type, loading ? 'loading' : '']" type="button">
    <span>
      {{value}}
    </span>
    <Loading height=16 width=16 class="loading-icon" stroke="#fff" />
  </button>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import Loading from '../Loading.vue'

export default defineComponent({
  name: 'VButton',
  components: {
    Loading
  },
  props: {
    type: String,
    value: String,
    loading: Boolean as PropType<boolean>
  }
})
</script>

<style lang="scss" scoped>
button {
  margin-bottom: 0.25rem;
  margin-top: 0.25rem;
  margin-left: auto;
  margin-right: auto;
  appearance: none;
  border: none;
  line-height: 2em;
  padding: 5px;
  border-radius: 3px;
  display: block;
  width: 95%;
  position: relative;
  background: lightgray;

  &:focus {
    outline: none;
  }

  &.primary {
    background: var(--ming);
    color: #fff;

    &:active {
      background: var(--midnight-green-eagle-green);
    }
  }

  &.loading {

    span {
      visibility: hidden;
    }
    .loading-icon {
      visibility: inherit;
    }
  }
  .loading-icon {
    padding: 0;

    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    visibility: hidden;
  }

}
</style>
