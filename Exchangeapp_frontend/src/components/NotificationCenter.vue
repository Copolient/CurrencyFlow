<template>
  <el-popover placement="bottom" :width="320" trigger="click">
    <template #reference>
      <div class="bell-btn">
        <span class="bell-icon">⬡</span>
        <span v-if="unreadCount > 0" class="bell-badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
      </div>
    </template>

    <template #default>
      <div class="notif-panel">
        <div class="notif-header">
          <span class="notif-title">通知</span>
          <button class="mark-all-btn" @click="markAllRead" v-if="unreadCount > 0">
            全部已读
          </button>
        </div>

        <el-scrollbar max-height="300px">
          <div v-if="notifications.length === 0" class="notif-empty">
            暂无通知
          </div>
          <div
            v-for="notif in notifications"
            :key="notif.id"
            class="notif-item"
            :class="{ unread: !notif.read }"
            @click="markRead(notif.id)"
          >
            <div class="notif-item-title">{{ notif.title }}</div>
            <div class="notif-item-content">{{ notif.content }}</div>
            <div class="notif-item-time">{{ formatTime(notif.createdAt) }}</div>
          </div>
        </el-scrollbar>
      </div>
    </template>
  </el-popover>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from '../axios';

interface Notification {
  id: number;
  type: string;
  title: string;
  content: string;
  read: boolean;
  createdAt: string;
}

const notifications = ref<Notification[]>([]);
const unreadCount = ref(0);

const fetchNotifications = async () => {
  try {
    const resp = await axios.get<Notification[]>('/notifications');
    notifications.value = resp.data;
  } catch { /* */ }
};

const fetchUnreadCount = async () => {
  try {
    const resp = await axios.get('/notifications/unread-count');
    unreadCount.value = resp.data.count;
  } catch { /* */ }
};

const markRead = async (id: number) => {
  try {
    await axios.put(`/notifications/${id}/read`);
    const notif = notifications.value.find((n) => n.id === id);
    if (notif) notif.read = true;
    unreadCount.value = Math.max(0, unreadCount.value - 1);
  } catch { /* */ }
};

const markAllRead = async () => {
  try {
    await axios.put('/notifications/read-all');
    notifications.value.forEach((n) => (n.read = true));
    unreadCount.value = 0;
  } catch { /* */ }
};

const formatTime = (timestamp: string) => {
  if (!timestamp) return '';
  const d = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - d.getTime();
  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`;
  return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`;
};

onMounted(() => { fetchNotifications(); fetchUnreadCount(); });
</script>

<style scoped>
.bell-btn {
  position: relative;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.2s ease;
}

.bell-btn:hover {
  background: var(--cf-surface-hover);
}

.bell-icon {
  font-size: 18px;
  color: var(--cf-text-secondary);
}

.bell-badge {
  position: absolute;
  top: -2px;
  right: -2px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  background: var(--cf-accent);
  color: var(--cf-black);
  font-size: 10px;
  font-weight: 700;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

/* ── Panel ── */
.notif-panel {
  padding: 4px 0;
}

.notif-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px 12px;
  border-bottom: 1px solid var(--cf-border);
  margin-bottom: 4px;
}

.notif-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--cf-text);
}

.mark-all-btn {
  background: none;
  border: none;
  font-size: 12px;
  font-weight: 500;
  font-family: var(--cf-font);
  color: var(--cf-accent);
  cursor: pointer;
  padding: 0;
}

.notif-empty {
  text-align: center;
  padding: 24px 0;
  font-size: 13px;
  color: var(--cf-text-muted);
}

.notif-item {
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.15s ease;
  margin: 2px 4px;
}

.notif-item:hover {
  background: var(--cf-surface-hover);
}

.notif-item.unread {
  background: var(--cf-accent-subtle);
}

.notif-item-title {
  font-weight: 600;
  font-size: 13px;
  color: var(--cf-text);
  margin-bottom: 3px;
}

.notif-item-content {
  font-size: 12px;
  color: var(--cf-text-secondary);
  margin-bottom: 3px;
}

.notif-item-time {
  font-size: 11px;
  color: var(--cf-text-muted);
}
</style>
