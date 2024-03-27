import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Overview from '../views/Overview.vue'
import Feed from '../views/Feed.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Feed',
    component: Feed,
  },
  {
    path: '/overview',
    name: 'Overview',
    component: Overview,
  },
  {
    path: '/login',
    name: 'Login',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/hobbits/add',
    component: () => import('../views/hobbits/AddHobbit.vue'),
  },
  {
    path: '/hobbits/:hobbitId',
    component: () => import('../views/hobbits/Hobbit.vue'),
  },
  {
    path: '/hobbits/:hobbitId/edit',
    component: () => import('../views/hobbits/EditHobbit.vue'),
  },
  {
    path: '/hobbits/:hobbitId/delete',
    component: () => import('../views/hobbits/DeleteHobbit.vue'),
  },
  {
    path: '/hobbits/:hobbitId/records/add',
    component: () => import('../views/hobbits/records/AddRecord.vue'),
  },
  {
    path: '/hobbits/:hobbitId/records/:recordId/edit',
    component: () => import('../views/hobbits/records/EditRecord.vue'),
  },
  {
    path: '/profile/:profileId',
    component: () => import('../views/profile/Profile.vue'),
  },
  {
    path: '/profile/me',
    component: () => import('../views/profile/Profile.vue'),
  },
  {
    path: '/profile/me/apppassword',
    component: () => import('../views/profile/AppPassword.vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
