<template>
  <div class="mobile-tab-bar" v-if="isMobile">
    <router-link
      v-for="tab in tabs"
      :key="tab.path"
      :to="tab.path"
      class="tab-item"
      :class="{ active: isActive(tab.path) }"
    >
      <span class="tab-icon">{{ tab.icon }}</span>
      <span class="tab-label">{{ tab.label }}</span>
      <span v-if="isActive(tab.path)" class="tab-indicator"></span>
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '../store/auth';

const route = useRoute();
const authStore = useAuthStore();

const isMobile = ref(false);

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768;
};

onMounted(() => {
  checkMobile();
  window.addEventListener('resize', checkMobile);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile);
});

const tabs = computed(() => {
  const items = [
    { path: '/', icon: '◈', label: '首页' },
    { path: '/exchange', icon: '⇄', label: '兑换' },
    { path: '/chart', icon: '◇', label: '行情' },
    { path: '/ai', icon: '◎', label: 'AI' },
    { path: '/community', icon: '◉', label: '社区' },
  ];

  if (authStore.isAuthenticated) {
    items.push({ path: '/alerts', icon: '⬡', label: '预警' });
  }

  return items;
});

const isActive = (path: string) => {
  if (path === '/') return route.path === '/';
  return route.path.startsWith(path);
};
</script>

<style scoped>
.mobile-tab-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(5, 5, 5, 0.9);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-top: 1px solid var(--cf-border);
  display: flex;
  justify-content: space-around;
  padding: 6px 0;
  z-index: 1000;
  padding-bottom: calc(6px + env(safe-area-inset-bottom, 0px));
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 8px;
  text-decoration: none;
  position: relative;
  transition: color 0.2s ease;
  -webkit-tap-highlight-color: transparent;
}

.tab-item:active {
  transform: scale(0.92);
}

.tab-icon {
  font-size: 18px;
  color: var(--cf-text-muted);
  transition: color 0.2s ease;
}

.tab-label {
  font-size: 10px;
  font-weight: 500;
  color: var(--cf-text-muted);
  transition: color 0.2s ease;
}

.tab-item.active .tab-icon {
  color: var(--cf-accent);
}

.tab-item.active .tab-label {
  color: var(--cf-accent);
}

.tab-indicator {
  position: absolute;
  top: -1px;
  left: 50%;
  transform: translateX(-50%);
  width: 16px;
  height: 2px;
  background: var(--cf-accent);
  border-radius: 1px;
}
</style>
