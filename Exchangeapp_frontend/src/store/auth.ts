import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from '../axios';

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const loading = ref(false);

  const isAuthenticated = computed(() => !!token.value);

  const login = async (username: string, password: string) => {
    loading.value = true;
    try {
      const response = await axios.post('/auth/login', { username, password });
      token.value = response.data.token;
      localStorage.setItem('token', token.value || '');
    } catch (error: any) {
      const msg = error.response?.data?.error || '登录失败';
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  };

  const register = async (username: string, password: string) => {
    loading.value = true;
    try {
      const response = await axios.post('/auth/register', { username, password });
      token.value = response.data.token;
      localStorage.setItem('token', token.value || '');
    } catch (error: any) {
      const msg = error.response?.data?.error || '注册失败';
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  };

  const logout = () => {
    token.value = null;
    localStorage.removeItem('token');
  };

  return { token, loading, isAuthenticated, login, register, logout };
});
