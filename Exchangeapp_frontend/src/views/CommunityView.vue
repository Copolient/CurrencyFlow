<template>
  <div class="community-container">
    <el-card class="post-form-card" v-if="authStore.isAuthenticated">
      <el-input
        v-model="newPost"
        type="textarea"
        :rows="3"
        placeholder="分享你的交易观点..."
        maxlength="2000"
        show-word-limit
      />
      <div class="post-actions">
        <el-input v-model="newCurrency" placeholder="关联货币对（可选）" style="width: 150px" />
        <el-button type="primary" @click="createPost" :loading="posting">发布</el-button>
      </div>
    </el-card>

    <el-card class="feed-card">
      <template #header>
        <div class="feed-header">
          <el-radio-group v-model="feedType" @change="fetchPosts">
            <el-radio-button label="latest">最新</el-radio-button>
            <el-radio-button label="following" v-if="authStore.isAuthenticated">关注</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <div v-loading="loading">
        <div v-if="posts.length === 0" class="empty-state">
          <el-empty description="暂无帖子" />
        </div>

        <div v-for="post in posts" :key="post.id" class="post-item">
          <div class="post-header">
            <el-avatar :size="36" :src="post.avatar || undefined">
              {{ post.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <div class="post-meta">
              <span class="post-username" @click="viewProfile(post.userId)">{{ post.username }}</span>
              <span class="post-time">{{ formatTime(post.createdAt) }}</span>
            </div>
            <el-tag v-if="post.currency" size="small" class="post-currency">{{ post.currency }}</el-tag>
          </div>
          <div class="post-content">{{ post.content }}</div>
          <div class="post-footer">
            <el-button text @click="likePost(post.id)">
              ❤️ {{ post.likes }}
            </el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import axios from '../axios';

interface Post {
  id: number;
  userId: number;
  username: string;
  content: string;
  currency: string;
  likes: number;
  createdAt: string;
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
  } catch {
    // interceptor handles error
  } finally {
    loading.value = false;
  }
};

const createPost = async () => {
  if (!newPost.value.trim()) return;

  posting.value = true;
  try {
    await axios.post('/posts', {
      content: newPost.value,
      currency: newCurrency.value,
    });
    newPost.value = '';
    newCurrency.value = '';
    fetchPosts();
  } catch {
    // interceptor handles error
  } finally {
    posting.value = false;
  }
};

const likePost = async (id: number) => {
  try {
    await axios.post(`/posts/${id}/like`);
    const post = posts.value.find((p) => p.id === id);
    if (post) post.likes++;
  } catch {
    // interceptor handles error
  }
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

onMounted(() => {
  fetchPosts();
});
</script>

<style scoped>
.community-container {
  max-width: 640px;
  margin: 20px auto;
  padding: 0 20px;
}

.post-form-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.post-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.feed-card {
  border-radius: 8px;
}

.feed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-item {
  padding: 16px 0;
  border-bottom: 1px solid #ebeef5;
}

.post-item:last-child {
  border-bottom: none;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.post-meta {
  flex: 1;
}

.post-username {
  font-weight: 600;
  color: #303133;
  cursor: pointer;
}

.post-username:hover {
  color: #409eff;
}

.post-time {
  display: block;
  font-size: 12px;
  color: #909399;
}

.post-currency {
  margin-left: auto;
}

.post-content {
  font-size: 14px;
  color: #303133;
  line-height: 1.6;
  margin-bottom: 12px;
  white-space: pre-wrap;
}

.post-footer {
  display: flex;
  gap: 16px;
}

.empty-state {
  padding: 40px 0;
}
</style>
