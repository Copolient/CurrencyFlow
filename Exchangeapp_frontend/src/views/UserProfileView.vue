<template>
  <div class="profile-page" v-loading="loading">
    <!-- Profile Header -->
    <div class="profile-card cf-glass cf-animate-in">
      <div class="profile-header">
        <div class="profile-avatar">{{ profile.username?.charAt(0)?.toUpperCase() || '?' }}</div>
        <div class="profile-info">
          <h2 class="profile-name">{{ profile.username }}</h2>
          <p class="profile-bio">{{ profile.bio || '这个人很懒，什么都没写' }}</p>
          <div class="profile-stats">
            <span class="stat"><strong>{{ profile.followersCount }}</strong> 粉丝</span>
            <span class="stat"><strong>{{ profile.followingCount }}</strong> 关注</span>
          </div>
        </div>
        <button
          v-if="authStore.isAuthenticated && !isSelf"
          class="follow-btn"
          :class="{ following: isFollowing }"
          @click="toggleFollow"
          :disabled="followLoading"
        >
          {{ isFollowing ? '已关注' : '关注' }}
        </button>
      </div>
    </div>

    <!-- User Posts -->
    <div class="user-posts cf-animate-in cf-delay-1">
      <h3 class="section-title">TA 的帖子</h3>
      <div v-if="userPosts.length === 0" class="empty-state">暂无帖子</div>
      <div v-for="post in userPosts" :key="post.id" class="post-card cf-glass">
        <p class="post-content">{{ post.content }}</p>
        <div class="post-meta">
          <span v-if="post.currency" class="post-currency-tag">{{ post.currency }}</span>
          <span class="post-time">{{ formatTime(post.createdAt) }}</span>
          <span class="post-likes">♥ {{ post.likes }}</span>
        </div>
      </div>
    </div>
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
  id: 0, username: '', avatar: '', bio: '', followersCount: 0, followingCount: 0, createdAt: '',
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
  } catch { return false; }
});

const fetchProfile = async () => {
  loading.value = true;
  try {
    const resp = await axios.get<UserProfile>(`/users/${userId.value}`);
    profile.value = resp.data;
  } catch { /* */ } finally { loading.value = false; }
};

const fetchPosts = async () => {
  try {
    const resp = await axios.get<Post[]>('/posts', {
      params: { type: 'user', userId: userId.value, pageSize: 50 },
    });
    userPosts.value = resp.data;
  } catch { /* */ }
};

const checkFollowing = async () => {
  if (!authStore.isAuthenticated || isSelf.value) return;
  try {
    const resp = await axios.get(`/users/${userId.value}/following`);
    isFollowing.value = resp.data.following;
  } catch { /* */ }
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
  } catch { /* */ } finally { followLoading.value = false; }
};

const formatTime = (timestamp: string) => {
  if (!timestamp) return '';
  const d = new Date(timestamp);
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`;
};

onMounted(() => { fetchProfile(); fetchPosts(); checkFollowing(); });
</script>

<style scoped>
.profile-page {
  padding-top: 48px;
  max-width: 640px;
  margin: 0 auto;
}

/* ── Profile Card ── */
.profile-card {
  padding: 28px;
  margin-bottom: 24px;
}

.profile-header {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.profile-avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: var(--cf-surface-hover);
  border: 2px solid var(--cf-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 700;
  color: var(--cf-accent);
  flex-shrink: 0;
}

.profile-info {
  flex: 1;
  min-width: 0;
}

.profile-name {
  font-size: 22px;
  font-weight: 700;
  color: var(--cf-text);
  margin: 0 0 6px;
}

.profile-bio {
  font-size: 14px;
  color: var(--cf-text-secondary);
  margin: 0 0 12px;
}

.profile-stats {
  display: flex;
  gap: 20px;
}

.stat {
  font-size: 13px;
  color: var(--cf-text-muted);
}

.stat strong {
  color: var(--cf-text);
  font-weight: 600;
}

.follow-btn {
  padding: 8px 20px;
  background: var(--cf-accent);
  color: var(--cf-black);
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  font-family: var(--cf-font);
  cursor: pointer;
  transition: background 0.2s ease;
  flex-shrink: 0;
}

.follow-btn:hover:not(:disabled) {
  background: var(--cf-accent-hover);
}

.follow-btn.following {
  background: var(--cf-surface-hover);
  color: var(--cf-text-secondary);
  border: 1px solid var(--cf-border);
}

.follow-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ── Posts ── */
.section-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--cf-text);
  margin: 0 0 16px;
}

.post-card {
  padding: 16px 20px;
  margin-bottom: 10px;
}

.post-content {
  font-size: 14px;
  line-height: 1.7;
  color: var(--cf-text-secondary);
  margin: 0 0 10px;
  white-space: pre-wrap;
}

.post-meta {
  display: flex;
  gap: 12px;
  align-items: center;
  font-size: 12px;
  color: var(--cf-text-muted);
}

.post-currency-tag {
  font-size: 11px;
  font-weight: 500;
  color: var(--cf-accent);
  background: var(--cf-accent-subtle);
  padding: 2px 8px;
  border-radius: 6px;
}

.post-likes {
  margin-left: auto;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  font-size: 14px;
  color: var(--cf-text-muted);
}

@media (max-width: 767px) {
  .profile-page {
    padding-top: 24px;
  }

  .profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .profile-stats {
    justify-content: center;
  }

  .follow-btn {
    width: 100%;
  }
}
</style>
