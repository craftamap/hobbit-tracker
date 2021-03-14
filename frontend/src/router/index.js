import Vue from 'vue'
import Router from 'vue-router'
import Login from '../components/Login'
import Overview from '../components/Overview.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Overview',
      component: Overview
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    }
  ]
})
