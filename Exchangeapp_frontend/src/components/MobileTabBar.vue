<template>
  <div class="mobile-tab-bar" v-if="isMobile">
    <div
      v-for="tab in tabs"
      :key="tab.path"
      class="tab-item"
      :class="{ active: isActive(tab.path) }"
      @click="router.push(tab.path)"
    >
      <span class="tab-icon">{{ tab.icon }}</span>
      <span class="tab-label">{{ tab.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../store/auth';

const router = useRouter();
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
    { path: '/', icon: '🏠', label: '首页' },
    { path: '/exchange', icon: '💱', label: '兑换' },
    { path: '/chart', icon: '📈', label: '行情' },
    { path: '/ai', icon: '🤖', label: 'AI' },
    { path: '/community', icon: '👥', label: '社区' },
  ];

  if (authStore.isAuthenticated) {
    items.push({ path: '/alerts', icon: '🔔', label: '预警' });
  }

  return items;
});

const isActive = (path: string) => {
  if (path === '/') {
    return route.path === '/';
  }
  return route.path.startsWith(path);
};
</script>

<style scoped>
.mobile-tab-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  display: flex;
  justify-content: space-around;
  padding: 8px 0;
  z-index: 1000;
  /* Safe area for iOS */
  padding-bottom: calc(8px + env(safe-area-inset-bottom, 0px));
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: all 0.2s;
}

.tab-item:active {
  transform: scale(0.95);
}

.tab-item.active {
  color: #409eff;
}

.tab-icon {
  font-size: 20px;
}

.tab-label {
  font-size: 10px;
  color: #909399;
}

.tab-item.active .tab-label {
  color: #409eff;
}
</style>
