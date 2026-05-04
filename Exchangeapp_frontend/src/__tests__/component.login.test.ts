import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import ElementPlus from 'element-plus';
import Login from '../components/Login.vue';

function createTestRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', name: 'Home', component: { template: '<div />' } },
      { path: '/login', name: 'Login', component: Login },
      { path: '/news', name: 'News', component: { template: '<div />' } },
    ],
  });
}

describe('Login Component', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
  });

  it('renders login form', () => {
    const router = createTestRouter();
    const wrapper = mount(Login, {
      global: { plugins: [router, ElementPlus] },
    });

    expect(wrapper.text()).toContain('登录');
    expect(wrapper.find('input').exists()).toBe(true);
  });

  it('has form fields for username and password', () => {
    const router = createTestRouter();
    const wrapper = mount(Login, {
      global: { plugins: [router, ElementPlus] },
    });

    const inputs = wrapper.findAll('input');
    expect(inputs.length).toBeGreaterThanOrEqual(2);
  });

  it('has a submit button', () => {
    const router = createTestRouter();
    const wrapper = mount(Login, {
      global: { plugins: [router, ElementPlus] },
    });

    const button = wrapper.find('button');
    expect(button.exists()).toBe(true);
    expect(button.text()).toContain('登录');
  });
});
