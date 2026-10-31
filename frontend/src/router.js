import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from './stores/auth.js'
import AppShell from './components/layout/AppShell.vue'

// Hash history: main.go only serves "/", so deep links must not hit the server.
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/login', name: 'login', component: () => import('./views/LoginView.vue'), meta: { guest: true } },
    { path: '/register', name: 'register', component: () => import('./views/RegisterView.vue'), meta: { guest: true } },
    {
      path: '/',
      component: AppShell,
      meta: { auth: true },
      children: [
        { path: '', name: 'home', redirect: () => (useAuthStore().role === 'admin' ? '/admin' : '/customer') },
        { path: 'admin', name: 'admin', component: () => import('./views/admin/AdminDashboard.vue'), meta: { role: 'admin' } },
        { path: 'customer', name: 'customer', component: () => import('./views/customer/CustomerPortal.vue'), meta: { role: 'customer' } },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.guest && auth.isAuthenticated) return '/'
  if (to.matched.some((r) => r.meta.auth) && !auth.isAuthenticated) {
    auth.logout()
    return '/login'
  }
  if (to.meta.role && to.meta.role !== auth.role) return '/'
})

export default router
