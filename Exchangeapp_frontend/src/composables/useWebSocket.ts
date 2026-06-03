import { ref, onMounted, onUnmounted } from 'vue';

export interface RateUpdate {
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  timestamp: string;
}

export function useWebSocket(pair?: string) {
  const ws = ref<WebSocket | null>(null);
  const connected = ref(false);
  const lastUpdate = ref<RateUpdate | null>(null);
  const reconnectAttempts = ref(0);
  const maxReconnectAttempts = 10;

  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  const getWsUrl = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const base = import.meta.env.VITE_API_BASE_URL || '/api';
    const host = window.location.host;
    let url = `${protocol}//${host}${base}/v1/ws`;
    if (pair) {
      url += `?pair=${encodeURIComponent(pair)}`;
    }
    return url;
  };

  const connect = () => {
    if (ws.value?.readyState === WebSocket.OPEN) return;

    try {
      ws.value = new WebSocket(getWsUrl());

      ws.value.onopen = () => {
        connected.value = true;
        reconnectAttempts.value = 0;
        console.log('WebSocket connected');
      };

      ws.value.onmessage = (event) => {
        try {
          const update: RateUpdate = JSON.parse(event.data);
          lastUpdate.value = update;
        } catch (e) {
          console.error('Failed to parse WebSocket message:', e);
        }
      };

      ws.value.onclose = () => {
        connected.value = false;
        console.log('WebSocket disconnected');
        scheduleReconnect();
      };

      ws.value.onerror = (error) => {
        console.error('WebSocket error:', error);
        ws.value?.close();
      };
    } catch (e) {
      console.error('Failed to create WebSocket:', e);
      scheduleReconnect();
    }
  };

  const scheduleReconnect = () => {
    if (reconnectAttempts.value >= maxReconnectAttempts) {
      console.warn('Max WebSocket reconnect attempts reached');
      return;
    }

    const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.value), 30000);
    reconnectAttempts.value++;

    reconnectTimer = setTimeout(() => {
      console.log(`WebSocket reconnecting (attempt ${reconnectAttempts.value})...`);
      connect();
    }, delay);
  };

  const disconnect = () => {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }
    connected.value = false;
  };

  onMounted(() => {
    connect();
  });

  onUnmounted(() => {
    disconnect();
  });

  return {
    connected,
    lastUpdate,
    disconnect,
    reconnect: connect,
  };
}
