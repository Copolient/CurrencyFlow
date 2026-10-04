import { describe, it, expect, beforeEach } from 'vitest';
import router from '../router/index';

describe('Router Guards', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('redirects unauthenticated users from /news to /login', async () => {
    await router.push('/news');
    expect(router.currentRoute.value.name).toBe('Login');
    expect(router.currentRoute.value.query.redirect).toBe('/news');
  });

  it('allows authenticated users to access /news', async () => {
    localStorage.setItem('token', 'Bearer test');
    await router.push('/news');
    expect(router.currentRoute.value.name).toBe('News');
  });

  it('redirects authenticated users away from /login', async () => {
    localStorage.setItem('token', 'Bearer test');
    await router.push('/login');
    expect(router.currentRoute.value.name).toBe('Home');
  });

  it('redirects unknown routes to home', async () => {
    await router.push('/nonexistent');
    expect(router.currentRoute.value.name).toBe('Home');
  });

  it('allows unauthenticated users to access /chart', async () => {
    await router.push('/chart');
    expect(router.currentRoute.value.name).toBe('Chart');
  });
});
