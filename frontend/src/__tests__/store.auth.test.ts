import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAuthStore } from '../store/auth';
import axios from '../axios';

vi.mock('../axios');

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
    vi.clearAllMocks();
  });

  it('starts unauthenticated', () => {
    const store = useAuthStore();
    expect(store.isAuthenticated).toBe(false);
    expect(store.token).toBeNull();
  });

  it('login sets token and persists to localStorage', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: { token: 'Bearer test-jwt-token' },
    });

    const store = useAuthStore();
    await store.login('alice', 'password123');

    expect(store.isAuthenticated).toBe(true);
    expect(store.token).toBe('Bearer test-jwt-token');
    expect(localStorage.getItem('token')).toBe('Bearer test-jwt-token');
  });

  it('login throws on failure', async () => {
    vi.mocked(axios.post).mockRejectedValueOnce({
      response: { data: { error: 'wrong credentials' } },
    });

    const store = useAuthStore();
    await expect(store.login('alice', 'wrong')).rejects.toThrow();
    expect(store.isAuthenticated).toBe(false);
  });

  it('register sets token on success', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: { token: 'Bearer new-user-token' },
    });

    const store = useAuthStore();
    await store.register('bob', 'secure123');

    expect(store.isAuthenticated).toBe(true);
    expect(store.token).toBe('Bearer new-user-token');
  });

  it('logout clears token and localStorage', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: { token: 'Bearer temp' },
    });

    const store = useAuthStore();
    await store.login('alice', 'password');
    store.logout();

    expect(store.isAuthenticated).toBe(false);
    expect(store.token).toBeNull();
    expect(localStorage.getItem('token')).toBeNull();
  });

  it('loading state resets after login', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: { token: 'Bearer t' },
    });

    const store = useAuthStore();
    expect(store.loading).toBe(false);

    const promise = store.login('a', 'b');
    // loading is true during the request
    expect(store.loading).toBe(true);

    await promise;
    expect(store.loading).toBe(false);
  });
});
