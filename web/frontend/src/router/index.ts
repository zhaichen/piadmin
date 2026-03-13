import { createRouter, createWebHistory } from 'vue-router'
import { hasToken } from '@/api/client'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/',
      component: () => import('@/views/Layout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'dashboard', component: () => import('@/views/Dashboard.vue') },
        { path: 'processes', name: 'processes', component: () => import('@/views/Processes.vue') },
        { path: 'services', name: 'services', component: () => import('@/views/Services.vue') },
        { path: 'network', name: 'network', component: () => import('@/views/Network.vue') },
        { path: 'terminal', name: 'terminal', component: () => import('@/views/Terminal.vue') },
        { path: 'files', name: 'files', component: () => import('@/views/Files.vue') },
        { path: 'gpio', name: 'gpio', component: () => import('@/views/Gpio.vue') },
      ],
    },
  ],
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !hasToken()) {
    return { name: 'login' }
  }
})

export default router
