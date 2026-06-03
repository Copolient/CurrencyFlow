<template>
  <div class="profile-container" v-loading="loading">
    <el-card class="profile-card">
      <div class="profile-header">
        <el-avatar :size="80" :src="profile.avatar || undefined">
          {{ profile.username?.charAt(0)?.toUpperCase() }}
        </el-avatar>
        <div class="profile-info">
          <h2>{{ profile.username }}</h2>
          <p class="bio">{{ profile.bio || '这个人很懒，什么都没写' }}</p>
          <div class="stats">
            <span><strong>{{ profile.followersCount }}</strong> 粉丝</span>
            <span><strong>{{ profile.followingCount }}</strong> 关注</span>
          </div>
        </div>
        <div class="profile-actions" v-if="authStore.isAuthenticated && !isSelf">
          <el-button
            :type="isFollowing ? 'default' : 'primary'"
            @click="toggleFollow"
            :loading="followLoading"
          >
            {{ isFollowing ? '已关注' : '关注' }}
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card class="posts-card">
      <template #header>
        <span>TA 的帖子</span>
      </template>

      <div v-if="userPosts.length === 0" class="empty-state">
        <el-empty description="暂无帖子" />
      </div>

      <div v-for="post in userPosts" :key="post.id" class="post-item">
        <div class="post-content">{{ post.content }}</div>
        <div class="post-meta">
          <span v-if="post.currency" class="post-currency">{{ post.currency }}</span>
          <span class="post-time">{{ formatTime(post.createdAt) }}</span>
          <span class="post-likes">❤️ {{ post.likes }}</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '../store/auth';
import axios from '../axios';

interface UserProfile {
  id: number;
  username: string;
  avatar: string;
  bio: string;
  followersCount: number;
  followingCount: number;
  createdAt: string;
}

interface Post {
  id: number;
  content: string;
  currency: string;
  likes: number;
  createdAt: string;
}

const route = useRoute();
const authStore = useAuthStore();
const loading = ref(false);
const followLoading = ref(false);
const profile = ref<UserProfile>({
  id: 0,
  username: '',
  avatar: '',
  bio: '',
  followersCount: 0,
  followingCount: 0,
  createdAt: '',
});
const userPosts = ref<Post[]>([]);
const isFollowing = ref(false);

const userId = computed(() => Number(route.params.id));

const isSelf = computed(() => {
  const token = localStorage.getItem('token');
  if (!token) return false;
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.username === profile.value.username;
  } catch {
    return false;
  }
});

const fetchProfile = async () => {
  loading.value = true;
  try {
    const resp = await axios.get<UserProfile>(`/users/${userId.value}`);
    profile.value = resp.data;
  } catch {
    // interceptor handles error
  } finally {
    loading.value = false;
  }
};

const fetchPosts = async () => {
  try {
    const resp = await axios.get<Post[]>('/posts', {
      params: { type: 'user', userId: userId.value, pageSize: 50 },
    });
    userPosts.value = resp.data;
  } catch {
    // interceptor handles error
  }
};

const checkFollowing = async () => {
  if (!authStore.isAuthenticated || isSelf.value) return;
  try {
    const resp = await axios.get(`/users/${userId.value}/following`);
    isFollowing.value = resp.data.following;
  } catch {
    // interceptor handles error
  }
};

const toggleFollow = async () => {
  followLoading.value = true;
  try {
    if (isFollowing.value) {
      await axios.delete(`/users/${userId.value}/follow`);
      isFollowing.value = false;
      profile.value.followersCount--;
    } else {
      await axios.post(`/users/${userId.value}/follow`);
      isFollowing.value = true;
      profile.value.followersCount++;
    }
  } catch {
    // interceptor handles error
  } finally {
    followLoading.value = false;
  }
};

const formatTime = (timestamp: string) => {
  if (!timestamp) return '';
  const d = new Date(timestamp);
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`;
};

onMounted(() => {
  fetchProfile();
  fetchPosts();
  checkFollowing();
});
</script>

<style scoped>
.profile-container {
  max-width: 640px;
  margin: 20px auto;
  padding: 0 20px;
}

.profile-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.profile-header {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.profile-info {
  flex: 1;
}

.profile-info h2 {
  margin: 0 0 8px 0;
  color: #303133;
}

.bio {
  color: #606266;
  font-size: 14px;
  margin-bottom: 12px;
}

.stats {
  display: flex;
  gap: 20px;
  color: #909399;
  font-size: 14px;
}

.stats strong {
  color: #303133;
}

.posts-card {
  border-radius: 8px;
}

.post-item {
  padding: 16px 0;
  border-bottom: 1px solid #ebeef5;
}

.post-item:last-child {
  border-bottom: none;
}

.post-content {
  font-size: 14px;
  color: #303133;
  line-height: 1.6;
  margin-bottom: 8px;
  white-space: pre-wrap;
}

.post-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #909399;
}

.post-currency {
  color: #409eff;
}

.empty-state {
  padding: 40px 0;
}
</style>
