<template>
  <div class="app-password">
    <div class="description-wrapper">
      <span class="description">{{ appPassword?.description }}</span> ·
      <span class="id">{{ appPassword?.id }}</span>
    </div>
    <div class="last-used-wrapper">
      <span class="last-used-at-label">Last Used At:</span>
      {{ formatDate(appPassword?.last_used_at) }} ({{ formatDateAgo(appPassword?.last_used_at) }})
    </div>
    <div class="icons-wrapper">
      <span tabindex="0">
        <Trash
          class="h-16 cursor-pointer"
          @click="emitDelete()"
        />
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, toRefs } from 'vue'
import { AppPassword } from '../../models'
import { TrashIcon as Trash } from '@heroicons/vue/24/outline'
import { formatDate, formatDateAgo } from '../../utils/date-utils'


export default defineComponent({
  name: 'AppPasswordItem',
  components: {
    Trash,
  },
  props: {
    appPassword: {
      type: Object as PropType<AppPassword>,
      required: true,
    },
  },
  emits: ['delete'],
  setup(props, { emit }) {
    const { appPassword } = toRefs(props)


    const emitDelete = () => {
      emit('delete', {
        id: appPassword.value?.id,
      })
    }

    return {
      formatDate,
      formatDateAgo,
      emitDelete,
    }
  },
})
</script>

<style scoped>
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

  .description-wrapper {
    grid-area: description;
  }

  .last-used-wrapper {
    grid-area: last-used;
  }

  .description {
    font-weight: bold;
  }

  .id,
  .last-used-at-label {
    color: var(--secondary-text);
  }
}
</style>
