<template>
  <el-container class="app-container">
    <el-header class="app-header" role="banner">
      <nav aria-label="主导航">
        <el-menu
          :default-active="activeIndex"
          mode="horizontal"
          :ellipsis="false"
          @select="handleSelect"
          class="nav-menu"
          role="menubar"
        >
          <el-menu-item index="home" class="logo-item" aria-label="首页">蓝鼠兑换</el-menu-item>
          <div class="flex-grow" />
          <el-menu-item index="exchange" role="menuitem">兑换货币</el-menu-item>
          <el-menu-item index="news" role="menuitem">查看资讯</el-menu-item>
          <el-menu-item index="login" v-if="!authStore.isAuthenticated" role="menuitem">登录</el-menu-item>
          <el-menu-item index="register" v-if="!authStore.isAuthenticated" role="menuitem">注册</el-menu-item>
          <el-sub-menu index="user" v-if="authStore.isAuthenticated">
            <template #title>用户</template>
            <el-menu-item index="logout" role="menuitem">退出登录</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </nav>
    </el-header>
    <el-main class="app-main" role="main">
      <ErrorBoundary>
        <router-view />
      </ErrorBoundary>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from './store/auth';
import ErrorBoundary from './components/ErrorBoundary.vue';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const routeNameToIndex: Record<string, string> = {
  Home: 'home',
  CurrencyExchange: 'exchange',
  News: 'news',
  NewsDetail: 'news',
  Login: 'login',
  Register: 'register',
};

const activeIndex = ref(routeNameToIndex[route.name as string] || 'home');

watch(route, (newRoute) => {
  activeIndex.value = routeNameToIndex[newRoute.name as string] || 'home';
});

const handleSelect = (key: string) => {
  if (key === 'logout') {
    authStore.logout();
    router.push({ name: 'Home' });
    return;
  }

  const indexToRoute: Record<string, string> = {
    home: 'Home',
    exchange: 'CurrencyExchange',
    news: 'News',
    login: 'Login',
    register: 'Register',
  };

  const routeName = indexToRoute[key];
  if (routeName) {
    router.push({ name: routeName });
  }
};
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.app-header {
  padding: 0;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  background: #fff;
}

.nav-menu {
  max-width: 960px;
  margin: 0 auto;
}

.logo-item {
  font-weight: 700;
  font-size: 16px;
}

.flex-grow {
  flex-grow: 1;
}

.app-main {
  max-width: 960px;
  margin: 0 auto;
  width: 100%;
}
</style>
