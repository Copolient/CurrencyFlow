import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import CurrencyExchangeView from '../views/CurrencyExchangeView.vue';
import ChartView from '../views/ChartView.vue';
import AlertView from '../views/AlertView.vue';
import AIAnalystView from '../views/AIAnalystView.vue';
import CommunityView from '../views/CommunityView.vue';
import UserProfileView from '../views/UserProfileView.vue';
import NewsView from '../views/NewsView.vue';
import NewsDetailView from '../views/NewsDetailView.vue';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/exchange', name: 'CurrencyExchange', component: CurrencyExchangeView },
  { path: '/chart', name: 'Chart', component: ChartView },
  { path: '/alerts', name: 'Alerts', component: AlertView, meta: { requiresAuth: true } },
  { path: '/ai', name: 'AIAnalyst', component: AIAnalystView },
  { path: '/community', name: 'Community', component: CommunityView },
  { path: '/user/:id', name: 'UserProfile', component: UserProfileView },
  {
    path: '/news',
    name: 'News',
    component: NewsView,
    meta: { requiresAuth: true },
  },
  {
    path: '/news/:id',
    name: 'NewsDetail',
    component: NewsDetailView,
    meta: { requiresAuth: true },
  },
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  { path: '/:pathMatch(.*)*', redirect: '/' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 全局前置守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token');

  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } });
  } else if ((to.name === 'Login' || to.name === 'Register') && token) {
    next({ name: 'Home' });
  } else {
    next();
  }
});

export default router;
