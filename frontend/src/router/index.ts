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
    path: '/hobbits/add',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/hobbits/AddHobbit.vue')
  },
  {
    path: '/hobbits/:id',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/hobbits/Hobbit.vue')
  },
  {
    path: '/hobbits/:id/records/add',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/hobbits/records/AddRecord.vue')
  },
  {
    path: '/profile',
    component: () => import(/* webpackChunkName: "profile" */ '../views/Profile.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
