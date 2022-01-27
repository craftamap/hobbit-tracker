<template>
  <FormWrapper>
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
        <Button @click="submit" value="Login" type="primary"/>
      </div>
    </form>
  </FormWrapper>
</template>

<script>
import Button from '../components/form/Button.vue'
import FormWrapper from '../components/form/FormWrapper.vue'

export default {
  name: 'Login',
  components: {
    Button,
    FormWrapper,
  },
  methods: {
    submit() {
      this.$refs.form.submit()
    },
  },
}
</script>

<style lang="scss" scoped>
.welcome {
  font-size: 16pt;
  text-align: center;
  margin-bottom: 2rem;
  color: darkgray;

}
</style>
