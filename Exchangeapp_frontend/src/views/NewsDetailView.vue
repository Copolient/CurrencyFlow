<template>
  <div class="detail-page">
    <el-skeleton :loading="loading" animated :rows="6">
      <template #default>
        <div v-if="article" class="article-detail cf-glass cf-animate-in">
          <button class="back-btn" @click="router.back()">← 返回列表</button>
          <h1 class="article-title">{{ article.Title }}</h1>
          <div class="article-divider"></div>
          <div class="article-content">{{ article.Content }}</div>

          <div class="like-section">
            <button class="like-btn" @click="likeArticle" :disabled="likeLoading">
              <span class="like-icon">♥</span>
              <span>点赞</span>
            </button>
            <span class="like-count">{{ likes }} 人点赞</span>
          </div>
        </div>
        <div v-else class="empty-state">文章不存在</div>
      </template>
    </el-skeleton>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from '../axios';
import type { Article, LikeResponse } from '../types/Article';

const article = ref<Article | null>(null);
const route = useRoute();
const router = useRouter();
const likes = ref(0);
const loading = ref(false);
const likeLoading = ref(false);

const id = route.params.id as string;

const fetchArticle = async () => {
  loading.value = true;
  try {
    const response = await axios.get<Article>(`/articles/${id}`);
    article.value = response.data;
  } catch { /* */ } finally { loading.value = false; }
};

const fetchLike = async () => {
  try {
    const res = await axios.get<LikeResponse>(`/articles/${id}/like`);
    likes.value = Number(res.data.likes) || 0;
  } catch { /* */ }
};

const likeArticle = async () => {
  likeLoading.value = true;
  try {
    await axios.post(`/articles/${id}/like`);
    await fetchLike();
  } catch { /* */ } finally { likeLoading.value = false; }
};

onMounted(() => { fetchArticle(); fetchLike(); });
</script>

<style scoped>
.detail-page {
  padding-top: 48px;
  max-width: 720px;
  margin: 0 auto;
}

.article-detail {
  padding: 32px;
}

.back-btn {
  background: none;
  border: none;
  font-size: 13px;
  font-weight: 500;
  font-family: var(--cf-font);
  color: var(--cf-accent);
  cursor: pointer;
  padding: 0;
  margin-bottom: 24px;
  transition: color 0.2s ease;
}

.back-btn:hover {
  color: var(--cf-accent-hover);
}

.article-title {
  font-size: 28px;
  font-weight: 800;
  color: var(--cf-text);
  margin: 0 0 20px;
  letter-spacing: -0.02em;
  line-height: 1.3;
}

.article-divider {
  height: 1px;
  background: var(--cf-border);
  margin-bottom: 24px;
}

.article-content {
  font-size: 15px;
  line-height: 1.8;
  color: var(--cf-text-secondary);
  white-space: pre-wrap;
  margin-bottom: 32px;
}

.like-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid var(--cf-border);
}

.like-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 18px;
  background: var(--cf-accent);
  color: var(--cf-black);
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  font-family: var(--cf-font);
  cursor: pointer;
  transition: background 0.2s ease;
}

.like-btn:hover:not(:disabled) {
  background: var(--cf-accent-hover);
}

.like-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.like-icon {
  font-size: 14px;
}

.like-count {
  font-size: 13px;
  color: var(--cf-text-muted);
}

.empty-state {
  text-align: center;
  padding: 80px 0;
  font-size: 14px;
  color: var(--cf-text-muted);
}

@media (max-width: 767px) {
  .detail-page {
    padding-top: 24px;
  }

  .article-detail {
    padding: 20px;
  }

  .article-title {
    font-size: 22px;
  }
}
</style>
