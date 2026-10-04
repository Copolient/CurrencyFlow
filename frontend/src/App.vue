<template>
  <div class="app-shell">
    <!-- Desktop Navigation -->
    <header class="nav-bar" role="banner">
      <nav class="nav-inner" aria-label="主导航">
        <div class="nav-brand" @click="router.push('/')">
          <span class="brand-mark">◆</span>
          <span class="brand-text">CurrencyFlow</span>
        </div>

        <div class="nav-links">
          <router-link
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="nav-link"
            :class="{ active: isActive(link.match) }"
          >
            {{ link.label }}
          </router-link>
        </div>

        <div class="nav-actions">
          <NotificationCenter v-if="authStore.isAuthenticated" />

          <template v-if="!authStore.isAuthenticated">
            <router-link to="/login" class="nav-link nav-link--muted">登录</router-link>
            <router-link to="/register" class="btn-accent">注册</router-link>
          </template>

          <div v-else class="user-menu" @click="showUserMenu = !showUserMenu">
            <div class="user-avatar">{{ userInitial }}</div>
            <transition name="dropdown">
              <div v-if="showUserMenu" class="user-dropdown" @click.stop>
                <div class="dropdown-item" @click="handleLogout">退出登录</div>
              </div>
            </transition>
          </div>
        </div>

        <!-- Mobile menu toggle -->
        <button class="mobile-menu-btn" @click="mobileMenuOpen = !mobileMenuOpen" aria-label="菜单">
          <span :class="['hamburger', { open: mobileMenuOpen }]">
            <span></span><span></span><span></span>
          </span>
        </button>
      </nav>

      <!-- Mobile menu -->
      <transition name="slide-down">
        <div v-if="mobileMenuOpen" class="mobile-menu">
          <router-link
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="mobile-link"
            :class="{ active: isActive(link.match) }"
            @click="mobileMenuOpen = false"
          >
            {{ link.label }}
          </router-link>
          <div class="mobile-divider"></div>
          <template v-if="!authStore.isAuthenticated">
            <router-link to="/login" class="mobile-link" @click="mobileMenuOpen = false">登录</router-link>
            <router-link to="/register" class="mobile-link accent" @click="mobileMenuOpen = false">注册</router-link>
          </template>
          <div v-else class="mobile-link danger" @click="handleLogout">退出登录</div>
        </div>
      </transition>
    </header>

    <!-- Main Content -->
    <main class="app-main" role="main">
      <ErrorBoundary>
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </ErrorBoundary>
    </main>

    <!-- Mobile Tab Bar -->
    <MobileTabBar />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from './store/auth';
import ErrorBoundary from './components/ErrorBoundary.vue';
import NotificationCenter from './components/NotificationCenter.vue';
import MobileTabBar from './components/MobileTabBar.vue';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const showUserMenu = ref(false);
const mobileMenuOpen = ref(false);

const navLinks = computed(() => [
  { to: '/exchange', label: '兑换', match: '/exchange' },
  { to: '/chart', label: '行情', match: '/chart' },
  { to: '/ai', label: 'AI 分析', match: '/ai' },
  { to: '/community', label: '社区', match: '/community' },
  { to: '/news', label: '资讯', match: '/news' },
  ...(authStore.isAuthenticated ? [{ to: '/alerts', label: '预警', match: '/alerts' }] : []),
]);

const userInitial = computed(() => {
  const token = localStorage.getItem('token');
  if (!token) return '?';
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return (payload.username || '?').charAt(0).toUpperCase();
  } catch {
    return '?';
  }
});

const isActive = (match: string) => route.path.startsWith(match);

const handleLogout = () => {
  authStore.logout();
  showUserMenu.value = false;
  mobileMenuOpen.value = false;
  router.push('/');
};

// Close menus on route change
watch(route, () => {
  showUserMenu.value = false;
  mobileMenuOpen.value = false;
});

// Close user menu on outside click
const closeUserMenu = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.user-menu')) {
    showUserMenu.value = false;
  }
};
document.addEventListener('click', closeUserMenu);
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: var(--cf-bg);
}

/* ── Navigation Bar ── */
.nav-bar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(5, 5, 5, 0.8);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--cf-border);
}

.nav-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 56px;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ── Brand ── */
.nav-brand {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  margin-right: 32px;
  user-select: none;
}

.brand-mark {
  color: var(--cf-accent);
  font-size: 18px;
  line-height: 1;
}

.brand-text {
  font-size: 18px;
  font-weight: 700;
  color: var(--cf-text);
  letter-spacing: -0.02em;
}

/* ── Nav Links ── */
.nav-links {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
}

.nav-link {
  padding: 6px 14px;
  font-size: 14px;
  font-weight: 500;
  color: var(--cf-text-secondary);
  text-decoration: none;
  border-radius: 8px;
  transition: color 0.2s ease, background 0.2s ease;
  position: relative;
}

.nav-link:hover {
  color: var(--cf-text);
  background: var(--cf-surface-hover);
}

.nav-link.active {
  color: var(--cf-text);
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 14px;
  right: 14px;
  height: 2px;
  background: var(--cf-accent);
  border-radius: 1px;
}

.nav-link--muted {
  color: var(--cf-text-muted);
}

/* ── Actions ── */
.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.btn-accent {
  padding: 6px 16px;
  font-size: 13px;
  font-weight: 600;
  color: var(--cf-black);
  background: var(--cf-accent);
  border-radius: 8px;
  text-decoration: none;
  transition: background 0.2s ease, transform 0.1s ease;
}

.btn-accent:hover {
  background: var(--cf-accent-hover);
  color: var(--cf-black);
}

.btn-accent:active {
  transform: scale(0.97);
}

/* ── User Avatar ── */
.user-menu {
  position: relative;
  cursor: pointer;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--cf-surface-hover);
  border: 1px solid var(--cf-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  color: var(--cf-accent);
  transition: border-color 0.2s ease;
}

.user-avatar:hover {
  border-color: var(--cf-accent);
}

.user-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 140px;
  background: var(--cf-surface);
  border: 1px solid var(--cf-border);
  border-radius: var(--cf-radius);
  padding: 4px;
  box-shadow: var(--cf-shadow-lg);
}

.dropdown-item {
  padding: 8px 12px;
  font-size: 13px;
  color: var(--cf-text-secondary);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.dropdown-item:hover {
  background: var(--cf-surface-hover);
  color: var(--cf-text);
}

/* ── Mobile Menu Button ── */
.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  margin-left: auto;
}

.hamburger {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 20px;
}

.hamburger span {
  display: block;
  height: 1.5px;
  background: var(--cf-text-secondary);
  border-radius: 1px;
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.hamburger.open span:nth-child(1) {
  transform: translateY(5.5px) rotate(45deg);
}

.hamburger.open span:nth-child(2) {
  opacity: 0;
}

.hamburger.open span:nth-child(3) {
  transform: translateY(-5.5px) rotate(-45deg);
}

/* ── Mobile Menu ── */
.mobile-menu {
  display: none;
  flex-direction: column;
  padding: 8px 24px 16px;
  gap: 2px;
}

.mobile-link {
  padding: 10px 12px;
  font-size: 15px;
  font-weight: 500;
  color: var(--cf-text-secondary);
  text-decoration: none;
  border-radius: 8px;
  transition: background 0.15s ease, color 0.15s ease;
}

.mobile-link:hover,
.mobile-link.active {
  color: var(--cf-text);
  background: var(--cf-surface-hover);
}

.mobile-link.accent {
  color: var(--cf-accent);
}

.mobile-link.danger {
  color: var(--cf-rose);
  cursor: pointer;
}

.mobile-divider {
  height: 1px;
  background: var(--cf-border);
  margin: 8px 0;
}

/* ── Main Content ── */
.app-main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  min-height: calc(100vh - 56px);
  padding-bottom: 80px;
}

/* ── Transitions ── */
.page-enter-active {
  animation: cf-fade-in 0.3s ease-out;
}

.page-leave-active {
  animation: cf-fade-in 0.2s ease-in reverse;
}

.dropdown-enter-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.dropdown-leave-active {
  transition: opacity 0.1s ease, transform 0.1s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px) scale(0.96);
}

.slide-down-enter-active {
  transition: opacity 0.2s ease, max-height 0.3s ease;
}

.slide-down-leave-active {
  transition: opacity 0.15s ease, max-height 0.2s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  max-height: 0;
}

.slide-down-enter-to,
.slide-down-leave-from {
  max-height: 400px;
}

/* ── Responsive ── */
@media (max-width: 767px) {
  .nav-links,
  .nav-actions .nav-link--muted,
  .nav-actions .btn-accent {
    display: none;
  }

  .nav-actions {
    display: none;
  }

  .mobile-menu-btn {
    display: block;
    margin-left: auto;
  }

  .mobile-menu {
    display: flex;
  }

  .app-main {
    padding: 16px 16px 80px;
  }
}
</style>
