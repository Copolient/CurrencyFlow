import { ref, onMounted, onUnmounted } from 'vue';

export interface RateUpdate {
  type?: string;
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  timestamp: string;
}

export interface NotificationMessage {
  type: string;
  id: number;
  notificationType: string;
  title: string;
  content: string;
  read: boolean;
  createdAt: string;
}

export interface UseWebSocketOptions {
  pair?: string;
  token?: string;
  autoConnect?: boolean;
  onMessage?: (message: RateUpdate | NotificationMessage) => void;
}

export function useWebSocket(pairOrOptions: string | UseWebSocketOptions = {}) {
  const options: UseWebSocketOptions =
    typeof pairOrOptions === 'string' ? { pair: pairOrOptions } : pairOrOptions;

  const ws = ref<WebSocket | null>(null);
  const connected = ref(false);
  const lastUpdate = ref<RateUpdate | null>(null);
  const lastNotification = ref<NotificationMessage | null>(null);
  const reconnectAttempts = ref(0);
  const maxReconnectAttempts = 10;

  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let stopped = false;

  const getWsUrl = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const base = import.meta.env.VITE_API_BASE_URL || '/api/v1';
    const host = window.location.host;
    const params = new URLSearchParams();
    if (options.pair) params.set('pair', options.pair);
    if (options.token) params.set('token', options.token);
    const query = params.toString();
    return `${protocol}//${host}${base}/ws${query ? `?${query}` : ''}`;
  };

  const connect = () => {
    if (stopped || ws.value?.readyState === WebSocket.OPEN) return;

    try {
      const socket = new WebSocket(getWsUrl());
      ws.value = socket;

      socket.onopen = () => {
        connected.value = true;
        reconnectAttempts.value = 0;
      };

      socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as RateUpdate | NotificationMessage;
          if (message.type === 'notification') {
            lastNotification.value = message as NotificationMessage;
          } else {
            lastUpdate.value = message as RateUpdate;
          }
          options.onMessage?.(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      socket.onclose = () => {
        connected.value = false;
        scheduleReconnect();
      };

      socket.onerror = () => {
        socket.close();
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      scheduleReconnect();
    }
  };

  const scheduleReconnect = () => {
    if (stopped || reconnectAttempts.value >= maxReconnectAttempts) return;

    const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.value), 30000);
    reconnectAttempts.value += 1;

    reconnectTimer = setTimeout(() => {
      connect();
    }, delay);
  };

  const disconnect = () => {
    stopped = true;
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

  const reconnect = () => {
    stopped = false;
    reconnectAttempts.value = 0;
    connect();
  };

  onMounted(() => {
    if (options.autoConnect === false) return;
    connect();
  });

  onUnmounted(() => {
    disconnect();
  });

  return {
    connected,
    lastUpdate,
    lastNotification,
    disconnect,
    reconnect,
  };
}
