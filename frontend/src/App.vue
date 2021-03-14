<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
@import "@fontsource/lato/index.css"; /* weight 400 */

:root {
  --sky-blue-crayola: #66d2e3ff;
  --cadet-blue: #49a0abff;
  --ming: #076470ff;
  --midnight-green-eagle-green: #1a4b52ff;
  --dark-jungle-green: #111d1fff;
}

body {
  margin: 0;
}

#app {
  font-family: 'Lato', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--dark-jungle-green);
}
</style>
