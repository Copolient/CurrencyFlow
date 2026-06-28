<template>
  <div class="news-page">
    <div class="page-header cf-animate-in">
      <div class="header-row">
        <div>
          <h1 class="page-title">财经资讯</h1>
          <p class="page-desc">浏览最新的财经文章与市场分析</p>
        </div>
        <button class="refresh-btn" @click="fetchArticles" :disabled="loading">
          <span :class="['refresh-icon', { spinning: loading }]">↻</span>
        </button>
      </div>
    </div>

    <div class="news-content cf-animate-in cf-delay-1">
      <el-skeleton :loading="loading" animated :rows="3" :count="3">
        <template #default>
          <div v-if="articles.length" class="articles-list">
            <div
              v-for="article in articles"
              :key="article.ID"
              class="article-card cf-glass cf-glass-hover"
              @click="viewDetail(article.ID)"
            >
              <h3 class="article-title">{{ article.Title }}</h3>
              <p class="article-preview">{{ article.Preview }}</p>
              <span class="article-link">阅读更多 →</span>
            </div>
          </div>
          <div v-else class="empty-state">暂无文章</div>
        </template>
      </el-skeleton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from '../axios';
import type { Article } from '../types/Article';

const articles = ref<Article[]>([]);
const router = useRouter();
const loading = ref(false);

const fetchArticles = async () => {
  loading.value = true;
  try {
    const response = await axios.get<Article[]>('/articles');
    articles.value = response.data;
  } catch { /* */ } finally { loading.value = false; }
};

const viewDetail = (id: string) => {
  router.push({ name: 'NewsDetail', params: { id } });
};

onMounted(fetchArticles);
</script>

<style scoped>
.news-page {
  padding-top: 48px;
  max-width: 720px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 28px;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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

.refresh-btn {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--cf-surface);
  border: 1px solid var(--cf-border);
  color: var(--cf-text-secondary);
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s ease, color 0.2s ease;
}

.refresh-btn:hover {
  border-color: var(--cf-accent);
  color: var(--cf-accent);
}

.refresh-icon {
  display: inline-block;
  transition: transform 0.3s ease;
}

.refresh-icon.spinning {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Articles ── */
.articles-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.article-card {
  padding: 24px;
  cursor: pointer;
}

.article-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--cf-text);
  margin: 0 0 10px;
  line-height: 1.4;
}

.article-preview {
  font-size: 14px;
  line-height: 1.6;
  color: var(--cf-text-secondary);
  margin: 0 0 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-link {
  font-size: 13px;
  font-weight: 500;
  color: var(--cf-accent);
  transition: color 0.2s ease;
}

.article-card:hover .article-link {
  color: var(--cf-accent-hover);
}

.empty-state {
  text-align: center;
  padding: 80px 0;
  font-size: 14px;
  color: var(--cf-text-muted);
}

@media (max-width: 767px) {
  .news-page {
    padding-top: 24px;
  }
}
</style>
