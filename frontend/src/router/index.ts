import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Overview from '../views/Overview.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Overview',
    component: Overview
  },
  {
    path: '/login',
    name: 'Login',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "login" */ '../views/Login.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
