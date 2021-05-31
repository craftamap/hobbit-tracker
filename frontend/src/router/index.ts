import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Overview from '../views/Overview.vue'
import Dashboard from '../views/Dashboard.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Overview',
    component: Overview,
  },
  {
    path: '/dashboard',
    name: 'Feed',
    component: Dashboard,
  },
  {
    path: '/login',
    name: 'Login',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "login" */ '../views/Login.vue'),
  },
  {
    path: '/hobbits/add',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/hobbits/AddHobbit.vue'),
  },
  {
    path: '/hobbits/:id',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/hobbits/Hobbit.vue'),
  },
  {
    path: '/hobbits/:id/edit',
    component: () => import(/* webpackChunkName: "hobbits" */ '../views/hobbits/EditHobbit.vue'),
  },
  {
    path: '/hobbits/:id/records/add',
    component: () => import(/* webpackChunkName: "records" */ '../views/hobbits/records/AddRecord.vue'),
  },
  {
    path: '/hobbits/:id/records/:recordId/edit',
    component: () => import(/* webpackChunkName: "records" */ '../views/hobbits/records/EditRecord.vue'),
  },
  {
    path: '/profile/me',
    component: () => import(/* webpackChunkName: "profile" */ '../views/profile/Profile.vue'),
  },
  {
    path: '/profile/me/apppassword',
    component: () => import(/* webpackChunkName: "profile" */ '../views/profile/AppPassword.vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
