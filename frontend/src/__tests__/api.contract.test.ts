import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

describe('API contract', () => {
  it('uses the versioned /api/v1 base path', () => {
    const env = readFileSync('.env', 'utf8');
    expect(env).toContain('VITE_API_BASE_URL=/api/v1');
    expect(env).not.toContain('VITE_API_BASE_URL=/api\n');
  });

  it('builds the WebSocket URL from the versioned base path', () => {
    const source = readFileSync('src/composables/useWebSocket.ts', 'utf8');
    expect(source).toContain('${base}/ws');
    expect(source).not.toContain('${base}/v1/ws');
  });
});
