<template>
  <div>
    <form method="POST" action="/auth/login">
      <input type="text" id="username" name="username" />
      <input type="password" id="password" name="password" />
      <input type="submit" />
    </form>
  </div>
</template>

<script>
export default {
  name: 'Login'
}
</script>
