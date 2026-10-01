import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn } from './api'
import PublicSite from './components/PublicSite.vue'
import Login from './components/Login.vue'
import AdminDashboard from './components/AdminDashboard.vue'

const routes = [
  { path: '/', name: 'public', component: PublicSite },
  { path: '/admin/login', name: 'login', component: Login },
  { path: '/admin', name: 'admin', component: AdminDashboard, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !isLoggedIn()) {
    return { name: 'login' }
  }
  if (to.name === 'login' && isLoggedIn()) {
    return { name: 'admin' }
  }
})

export default router
