<template>
  <div class="form-wrapper">
    <form ref="form" method="POST" action="/auth/login" v-on:keyup.enter="submit">
      <div class="welcome">
        Welcome back!
      </div>
      <div>
        <input
          type="text"
          id="username"
          name="username"
          placeholder="username"
        />
      </div>
      <div>
        <input
          type="password"
          id="password"
          name="password"
          placeholder="password"
        />
      </div>
      <div>
        <VButton @click="submit" value="Login" type="primary"/>
      </div>
    </form>
  </div>
</template>

<script>
import VButton from '../components/form/Button'

export default {
  name: 'Login',
  components: {
    VButton
  },
  methods: {
    submit () {
      this.$refs.form.submit()
    }
  }
}
</script>

<style lang="scss" scoped>
.welcome {
  font-size: 16pt;
  text-align: center;
  margin-bottom: 2rem;
  color: darkgray;

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

    input {
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
