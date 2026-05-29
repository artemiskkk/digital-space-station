import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../views/Dashboard.vue'),
      children: [
        { path: '', redirect: '/posts' },
        { path: 'posts', component: () => import('../views/Posts.vue') },
        { path: 'moments', component: () => import('../views/Moments.vue') },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

export default router
