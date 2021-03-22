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
  },
  {
    path: '/hobbits/:id',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/Hobbit.vue')
  },
  {
    path: '/hobbits/:id/add',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/Add.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
