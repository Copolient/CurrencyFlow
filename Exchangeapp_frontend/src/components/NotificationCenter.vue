<template>
  <el-popover placement="bottom" :width="320" trigger="click">
    <template #reference>
      <el-badge :value="unreadCount" :hidden="unreadCount === 0" :max="99">
        <el-button :icon="Bell" circle />
      </el-badge>
    </template>

    <template #default>
      <div class="notification-header">
        <span>通知</span>
        <el-button text type="primary" size="small" @click="markAllRead" v-if="unreadCount > 0">
          全部已读
        </el-button>
      </div>

      <el-scrollbar max-height="300px">
        <div v-if="notifications.length === 0" class="empty-notifications">
          暂无通知
        </div>
        <div
          v-for="notif in notifications"
          :key="notif.id"
          class="notification-item"
          :class="{ unread: !notif.read }"
          @click="markRead(notif.id)"
        >
          <div class="notification-title">{{ notif.title }}</div>
          <div class="notification-content">{{ notif.content }}</div>
          <div class="notification-time">{{ formatTime(notif.createdAt) }}</div>
        </div>
      </el-scrollbar>
    </template>
  </el-popover>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Bell } from '@element-plus/icons-vue';
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
  } catch {
    // interceptor handles error
  }
};

const fetchUnreadCount = async () => {
  try {
    const resp = await axios.get('/notifications/unread-count');
    unreadCount.value = resp.data.count;
  } catch {
    // interceptor handles error
  }
};

const markRead = async (id: number) => {
  try {
    await axios.put(`/notifications/${id}/read`);
    const notif = notifications.value.find((n) => n.id === id);
    if (notif) notif.read = true;
    unreadCount.value = Math.max(0, unreadCount.value - 1);
  } catch {
    // interceptor handles error
  }
};

const markAllRead = async () => {
  try {
    await axios.put('/notifications/read-all');
    notifications.value.forEach((n) => (n.read = true));
    unreadCount.value = 0;
  } catch {
    // interceptor handles error
  }
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

onMounted(() => {
  fetchNotifications();
  fetchUnreadCount();
});
</script>

<style scoped>
.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 8px;
}

.notification-item {
  padding: 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.unread {
  background-color: #ecf5ff;
}

.notification-title {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}

.notification-content {
  font-size: 12px;
  color: #606266;
  margin-bottom: 4px;
}

.notification-time {
  font-size: 11px;
  color: #909399;
}

.empty-notifications {
  text-align: center;
  color: #909399;
  padding: 20px 0;
}
</style>
