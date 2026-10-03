import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: () => import('../views/HomeView.vue') },
  { path: '/exchange', name: 'CurrencyExchange', component: () => import('../views/CurrencyExchangeView.vue') },
  { path: '/chart', name: 'Chart', component: () => import('../views/ChartView.vue') },
  { path: '/alerts', name: 'Alerts', component: () => import('../views/AlertView.vue'), meta: { requiresAuth: true } },
  { path: '/ai', name: 'AIAnalyst', component: () => import('../views/AIAnalystView.vue'), meta: { requiresAuth: true } },
  { path: '/community', name: 'Community', component: () => import('../views/CommunityView.vue') },
  { path: '/user/:id', name: 'UserProfile', component: () => import('../views/UserProfileView.vue') },
  {
    path: '/news',
    name: 'News',
    component: () => import('../views/NewsView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/news/:id',
    name: 'NewsDetail',
    component: () => import('../views/NewsDetailView.vue'),
    meta: { requiresAuth: true },
  },
  { path: '/login', name: 'Login', component: () => import('../components/Login.vue') },
  { path: '/register', name: 'Register', component: () => import('../components/Register.vue') },
  { path: '/:pathMatch(.*)*', redirect: '/' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Global auth guard
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token');

  if (to.meta.requiresAuth && token === null) {
    next({ name: 'Login', query: { redirect: to.fullPath } });
  } else if ((to.name === 'Login' || to.name === 'Register') && token) {
    next({ name: 'Home' });
  } else {
    next();
  }
});

export default router;
