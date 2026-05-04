<template>
  <div class="detail-container">
    <el-skeleton :loading="loading" animated :rows="6">
      <template #default>
        <el-card v-if="article" class="article-detail" shadow="never">
          <h1 class="article-title">{{ article.Title }}</h1>
          <el-divider />
          <div class="article-content">{{ article.Content }}</div>

          <div class="like-section">
            <el-button type="primary" :icon="StarFilled" @click="likeArticle" :loading="likeLoading">
              点赞
            </el-button>
            <span class="like-count">{{ likes }} 人点赞</span>
          </div>
        </el-card>
        <el-empty v-else description="文章不存在" />
      </template>
    </el-skeleton>

    <el-button text @click="router.back()" style="margin-top: 16px">
      ← 返回列表
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { StarFilled } from '@element-plus/icons-vue';
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
  } catch {
    // interceptor 已处理
  } finally {
    loading.value = false;
  }
};

const fetchLike = async () => {
  try {
    const res = await axios.get<LikeResponse>(`/articles/${id}/like`);
    likes.value = Number(res.data.likes) || 0;
  } catch {
    // 非关键错误，静默处理
  }
};

const likeArticle = async () => {
  likeLoading.value = true;
  try {
    await axios.post(`/articles/${id}/like`);
    await fetchLike();
  } catch {
    // interceptor 已处理
  } finally {
    likeLoading.value = false;
  }
};

onMounted(() => {
  fetchArticle();
  fetchLike();
});
</script>

<style scoped>
.detail-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.article-detail {
  border-radius: 8px;
}

.article-title {
  font-size: 28px;
  color: #303133;
  margin-bottom: 0;
}

.article-content {
  font-size: 16px;
  line-height: 1.8;
  color: #606266;
  white-space: pre-wrap;
  margin-bottom: 24px;
}

.like-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.like-count {
  color: #909399;
  font-size: 14px;
}
</style>
