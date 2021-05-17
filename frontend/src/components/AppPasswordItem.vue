<template>
  <div class="app-password">
    <div class="description-wrapper">
      <span class="description">{{appPassword.description}}</span> · <span class="id">{{appPassword.id}}</span>
    </div>
    <div class="last-used-wrapper">
      <span class="last-used-at-label">Last Used At: </span>{{ formatDate(appPassword.last_used_at) }}
    </div>
    <div class="icons-wrapper">
      <span tabindex="0"><Trash class="h-16 cursor-pointer"  @click="emitDelete($event)"/></span>
    </div>
  </div>
</template>

<script lang="ts">
import moment from 'moment'
import { defineComponent, PropType } from 'vue'
import { AppPassword } from '../models'
import { TrashIcon as Trash } from '@heroicons/vue/outline'

export default defineComponent({
  name: 'AppPasswordItem',
  components: {
    Trash,
  },
  props: {
    appPassword: {
      type: Object as PropType<AppPassword>,
    },
  },
  methods: {
    formatDate(date: string) {
      return moment(date).format('YYYY-MM-DD HH:mm')
    },
    emitDelete() {
      this.$emit('delete', {
        id: this?.appPassword?.id,
      })
    },
  },
})
</script>

<style lang="scss" scoped>
.app-password {
  display: grid;
  grid-template-columns: 3fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 0px 0px;
  grid-template-areas:
    "description icons"
    "last-used icons";
  border-radius: 0.5em;
  padding: 1rem;
  margin: 0.5rem;
  box-shadow: 0px 0px 5px -2px #000000;

  .icons-wrapper {
    grid-area: icons;
    text-align: right;
  }

  .description-wrapper { grid-area: description; }

  .last-used-wrapper { grid-area: last-used; }

  .description {
    font-weight: bold;
  }

  .id, .last-used-at-label {
    color: var(--secondary-text);
  }
}
</style>
