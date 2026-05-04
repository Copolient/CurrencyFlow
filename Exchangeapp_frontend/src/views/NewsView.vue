<template>
  <div class="news-container">
    <div class="news-header">
      <h2>财经资讯</h2>
      <el-button :icon="Refresh" circle @click="fetchArticles" :loading="loading" />
    </div>

    <el-skeleton :loading="loading" animated :rows="3" :count="3">
      <template #default>
        <div v-if="articles.length">
          <el-card
            v-for="article in articles"
            :key="article.ID"
            class="article-card"
            shadow="hover"
          >
            <h3>{{ article.Title }}</h3>
            <p class="preview">{{ article.Preview }}</p>
            <el-button text type="primary" @click="viewDetail(article.ID)">
              阅读更多 →
            </el-button>
          </el-card>
        </div>
        <el-empty v-else description="暂无文章" />
      </template>
    </el-skeleton>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Refresh } from '@element-plus/icons-vue';
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
  } catch {
    // interceptor 已处理
  } finally {
    loading.value = false;
  }
};

const viewDetail = (id: string) => {
  router.push({ name: 'NewsDetail', params: { id } });
};

onMounted(fetchArticles);
</script>

<style scoped>
.news-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.news-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.article-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.article-card h3 {
  margin-bottom: 8px;
  color: #303133;
}

.preview {
  color: #909399;
  font-size: 14px;
  margin-bottom: 12px;
}
</style>
