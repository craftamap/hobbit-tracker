<template>
  <div class="form-wrapper">
    <slot></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'FormWrapper',
})
</script>

<style>
.form-wrapper {
  display: flex;
  justify-content: center;
  justify-items: center;

  form {
    border-radius: 0.5rem;
    padding: 2rem;
    background: #eee;
    width: 300px;

    @media (prefers-color-scheme: dark) {
      background: var(--primary-dark);
    }

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
        background: var(--primary);
        color: #fff;
        width: 100%;
      }
    }
  }
}
</style>
