import axios from 'axios';
import { ElMessage } from 'element-plus';

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
});

// 请求拦截：自动附加 Token
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = token;
  }
  return config;
});

// 响应拦截：统一错误处理
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;
    const message = error.response?.data?.error || '请求失败';

    if (status === 401) {
      localStorage.removeItem('token');
      // 避免在登录页循环跳转
      if (window.location.pathname !== '/login') {
        ElMessage.error('登录已过期，请重新登录');
        window.location.href = '/login';
      }
    } else if (status === 403) {
      ElMessage.error('没有权限执行此操作');
    } else if (status === 500) {
      ElMessage.error('服务器内部错误，请稍后重试');
    } else if (!error.response) {
      ElMessage.error('网络连接失败，请检查网络');
    }

    return Promise.reject(error);
  },
);

export default instance;
