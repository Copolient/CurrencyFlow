<template>
  <div class="community-page">
    <div class="page-header cf-animate-in">
      <h1 class="page-title">交易社区</h1>
      <p class="page-desc">分享你的交易观点，关注市场动态</p>
    </div>

    <!-- Post Form -->
    <div class="post-form cf-glass cf-animate-in cf-delay-1" v-if="authStore.isAuthenticated">
      <el-input
        v-model="newPost"
        type="textarea"
        :rows="3"
        placeholder="分享你的交易观点..."
        maxlength="2000"
        show-word-limit
      />
      <div class="post-form-actions">
        <el-input v-model="newCurrency" placeholder="关联货币对（可选）" class="currency-input" />
        <button class="post-btn" @click="createPost" :disabled="posting || !newPost.trim()">
          <span v-if="posting" class="btn-spinner"></span>
          <span v-else>发布</span>
        </button>
      </div>
    </div>

    <!-- Feed -->
    <div class="feed-section cf-animate-in cf-delay-2">
      <div class="feed-tabs">
        <button
          class="feed-tab"
          :class="{ active: feedType === 'latest' }"
          @click="feedType = 'latest'; fetchPosts()"
        >
          最新
        </button>
        <button
          v-if="authStore.isAuthenticated"
          class="feed-tab"
          :class="{ active: feedType === 'following' }"
          @click="feedType = 'following'; fetchPosts()"
        >
          关注
        </button>
      </div>

      <div v-loading="loading">
        <div v-if="posts.length === 0" class="empty-state">
          暂无帖子
        </div>

        <div v-for="post in posts" :key="post.ID" class="post-card cf-glass">
          <div class="post-header">
            <div class="post-avatar">{{ post.username?.charAt(0)?.toUpperCase() }}</div>
            <div class="post-meta">
              <span class="post-username" @click="viewProfile(post.userId)">{{ post.username }}</span>
              <span class="post-time">{{ formatTime(post.CreatedAt) }}</span>
            </div>
            <span v-if="post.currency" class="post-currency-tag">{{ post.currency }}</span>
          </div>
          <p class="post-content">{{ post.content }}</p>
          <div class="post-footer">
            <button class="like-btn" @click="likePost(post.ID)">
              <span class="like-icon">♥</span>
              <span class="like-count">{{ post.likes }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import axios from '../axios';

interface Post {
  ID: number;
  userId: number;
  username: string;
  content: string;
  currency: string;
  likes: number;
  CreatedAt: string;
  avatar?: string;
}

const router = useRouter();
const authStore = useAuthStore();
const loading = ref(false);
const posting = ref(false);
const posts = ref<Post[]>([]);
const feedType = ref('latest');
const newPost = ref('');
const newCurrency = ref('');

const fetchPosts = async () => {
  loading.value = true;
  try {
    const resp = await axios.get<Post[]>('/posts', {
      params: { type: feedType.value, pageSize: 50 },
    });
    posts.value = resp.data;
  } catch { /* */ } finally { loading.value = false; }
};

const createPost = async () => {
  if (!newPost.value.trim()) return;
  posting.value = true;
  try {
    await axios.post('/posts', { content: newPost.value, currency: newCurrency.value });
    newPost.value = '';
    newCurrency.value = '';
    fetchPosts();
  } catch { /* */ } finally { posting.value = false; }
};

const likePost = async (id: number) => {
  try {
    await axios.post(`/posts/${id}/like`);
    const post = posts.value.find((p) => p.ID === id);
    if (post) post.likes++;
  } catch { /* */ }
};

const viewProfile = (userId: number) => {
  router.push(`/user/${userId}`);
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

onMounted(() => { fetchPosts(); });
</script>

<style scoped>
.community-page {
  padding-top: 48px;
  max-width: 640px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 28px;
}

.page-title {
  font-size: 28px;
  font-weight: 800;
  color: var(--cf-text);
  margin: 0 0 8px;
  letter-spacing: -0.02em;
}

.page-desc {
  font-size: 14px;
  color: var(--cf-text-muted);
  margin: 0;
}

/* ── Post Form ── */
.post-form {
  padding: 20px;
  margin-bottom: 24px;
}

.post-form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
  gap: 12px;
}

.currency-input {
  width: 180px;
}

.post-btn {
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
  display: flex;
  align-items: center;
  gap: 6px;
}

.post-btn:hover:not(:disabled) {
  background: var(--cf-accent-hover);
}

.post-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(0, 0, 0, 0.2);
  border-top-color: var(--cf-black);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Feed Tabs ── */
.feed-tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 20px;
  background: var(--cf-surface);
  border-radius: 8px;
  padding: 3px;
  width: fit-content;
}

.feed-tab {
  padding: 6px 16px;
  font-size: 13px;
  font-weight: 500;
  font-family: var(--cf-font);
  color: var(--cf-text-muted);
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: color 0.2s ease, background 0.2s ease;
}

.feed-tab:hover {
  color: var(--cf-text-secondary);
}

.feed-tab.active {
  color: var(--cf-black);
  background: var(--cf-accent);
}

/* ── Post Cards ── */
.post-card {
  padding: 20px;
  margin-bottom: 12px;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.post-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--cf-surface-hover);
  border: 1px solid var(--cf-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: var(--cf-accent);
  flex-shrink: 0;
}

.post-meta {
  flex: 1;
  min-width: 0;
}

.post-username {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--cf-text);
  cursor: pointer;
  transition: color 0.2s ease;
}

.post-username:hover {
  color: var(--cf-accent);
}

.post-time {
  display: block;
  font-size: 11px;
  color: var(--cf-text-muted);
}

.post-currency-tag {
  font-size: 11px;
  font-weight: 500;
  color: var(--cf-accent);
  background: var(--cf-accent-subtle);
  padding: 3px 10px;
  border-radius: 6px;
}

.post-content {
  font-size: 14px;
  line-height: 1.7;
  color: var(--cf-text-secondary);
  margin: 0;
  white-space: pre-wrap;
}

.post-footer {
  margin-top: 12px;
}

.like-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: none;
  border: none;
  cursor: pointer;
  font-family: var(--cf-font);
  transition: color 0.2s ease;
}

.like-icon {
  font-size: 14px;
  color: var(--cf-text-muted);
  transition: color 0.2s ease;
}

.like-btn:hover .like-icon {
  color: var(--cf-rose);
}

.like-count {
  font-size: 12px;
  color: var(--cf-text-muted);
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  font-size: 14px;
  color: var(--cf-text-muted);
}

/* ── Responsive ── */
@media (max-width: 767px) {
  .community-page {
    padding-top: 24px;
  }

  .post-form-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .currency-input {
    width: 100%;
  }
}
</style>
