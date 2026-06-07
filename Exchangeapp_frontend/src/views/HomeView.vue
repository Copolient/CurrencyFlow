<template>
  <div class="home-container">
    <!-- Dashboard for authenticated users -->
    <template v-if="authStore.isAuthenticated">
      <div class="dashboard-header">
        <h1>我的仪表盘</h1>
        <p class="subtitle">欢迎回来，{{ username }}</p>
      </div>

      <!-- Favorite Rates -->
      <el-card class="dashboard-section" v-loading="loading">
        <template #header>
          <div class="section-header">
            <span>💱 收藏汇率</span>
            <el-button text type="primary" @click="router.push('/exchange')">查看更多</el-button>
          </div>
        </template>

        <div v-if="favoriteRates.length === 0" class="empty-state">
          <el-empty description="还没有收藏汇率">
            <el-button type="primary" @click="router.push('/exchange')">去收藏</el-button>
          </el-empty>
        </div>

        <el-row :gutter="16" v-else>
          <el-col :span="6" v-for="rate in favoriteRates" :key="`${rate.fromCurrency}-${rate.toCurrency}`">
            <el-card shadow="hover" class="rate-card" @click="router.push('/chart')">
              <div class="rate-pair">{{ rate.fromCurrency }}/{{ rate.toCurrency }}</div>
              <div class="rate-value">{{ rate.rate.toFixed(4) }}</div>
              <div class="rate-time">{{ formatTime(rate.timestamp) }}</div>
            </el-card>
          </el-col>
        </el-row>
      </el-card>

      <!-- Quick Actions -->
      <el-card class="dashboard-section">
        <template #header>
          <span>⚡ 快捷操作</span>
        </template>
        <el-row :gutter="16">
          <el-col :span="6">
            <el-button class="action-btn" @click="router.push('/exchange')">
              💱 货币兑换
            </el-button>
          </el-col>
          <el-col :span="6">
            <el-button class="action-btn" @click="router.push('/chart')">
              📈 行情走势
            </el-button>
          </el-col>
          <el-col :span="6">
            <el-button class="action-btn" @click="router.push('/news')">
              📰 财经资讯
            </el-button>
          </el-col>
          <el-col :span="6">
            <el-button class="action-btn" @click="router.push('/alerts')">
              🔔 汇率预警
            </el-button>
          </el-col>
        </el-row>
      </el-card>
    </template>

    <!-- Landing page for guests -->
    <template v-else>
      <div class="hero">
        <h1 class="title">蓝鼠兑换</h1>
        <p class="subtitle">实时汇率查询 · AI 智能分析 · 社交交易社区</p>
        <div class="actions">
          <el-button type="primary" size="large" @click="router.push('/exchange')">
            开始兑换
          </el-button>
          <el-button size="large" @click="router.push('/chart')">
            行情走势
          </el-button>
        </div>
      </div>

      <!-- 实时汇率 -->
      <el-card class="guest-section" v-loading="ratesLoading">
        <template #header>
          <div class="section-header">
            <span>💱 实时汇率</span>
            <el-button text type="primary" @click="router.push('/chart')">查看详情 →</el-button>
          </div>
        </template>
        <el-row :gutter="16">
          <el-col :span="6" v-for="rate in latestRates" :key="`${rate.fromCurrency}-${rate.toCurrency}`">
            <el-card shadow="hover" class="rate-card" @click="router.push('/chart')">
              <div class="rate-pair">{{ rate.fromCurrency }}/{{ rate.toCurrency }}</div>
              <div class="rate-value">{{ rate.rate.toFixed(4) }}</div>
            </el-card>
          </el-col>
        </el-row>
      </el-card>

      <!-- AI 分析师 -->
      <el-card class="guest-section">
        <template #header>
          <div class="section-header">
            <span>🤖 AI 智能分析师</span>
            <el-button text type="primary" @click="router.push('/ai')">去分析 →</el-button>
          </div>
        </template>
        <div class="ai-preview">
          <p>基于历史数据的智能汇率分析，支持趋势判断、关键价位、风险提示。</p>
          <el-button type="primary" @click="router.push('/ai')">体验 AI 分析</el-button>
        </div>
      </el-card>

      <!-- 社区帖子 -->
      <el-card class="guest-section" v-loading="postsLoading">
        <template #header>
          <div class="section-header">
            <span>👥 交易社区</span>
            <el-button text type="primary" @click="router.push('/community')">进入社区 →</el-button>
          </div>
        </template>
        <div v-if="recentPosts.length === 0" class="empty-state">
          <el-empty description="暂无帖子" :image-size="60" />
        </div>
        <div v-for="post in recentPosts" :key="post.ID" class="post-preview">
          <div class="post-header">
            <span class="post-username">{{ post.username }}</span>
            <el-tag v-if="post.currency" size="small">{{ post.currency }}</el-tag>
          </div>
          <div class="post-content">{{ post.content }}</div>
          <div class="post-footer">❤️ {{ post.likes }}</div>
        </div>
      </el-card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import axios from '../axios';

const router = useRouter();
const authStore = useAuthStore();
const loading = ref(false);
const ratesLoading = ref(false);
const postsLoading = ref(false);

interface RateHistory {
  _id: number;
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  timestamp: string;
}

interface Post {
  ID: number;
  userId: number;
  username: string;
  content: string;
  currency: string;
  likes: number;
  CreatedAt: string;
}

const favoriteRates = ref<RateHistory[]>([]);
const latestRates = ref<RateHistory[]>([]);
const recentPosts = ref<Post[]>([]);

const username = computed(() => {
  const token = localStorage.getItem('token');
  if (!token) return '';
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.username || '';
  } catch {
    return '';
  }
});

const fetchFavoriteRates = async () => {
  if (!authStore.isAuthenticated) return;

  loading.value = true;
  try {
    // Get favorites
    const favResp = await axios.get('/favorites');
    const favorites = favResp.data;

    if (favorites.length === 0) {
      favoriteRates.value = [];
      return;
    }

    // Get latest rates
    const ratesResp = await axios.get<RateHistory[]>('/rates/latest');
    const latestRates = ratesResp.data;

    // Match favorites with latest rates
    favoriteRates.value = favorites.map((fav: any) => {
      const match = latestRates.find(
        (r) => r.fromCurrency === fav.fromCurrency && r.toCurrency === fav.toCurrency
      );
      return match || {
        _id: 0,
        fromCurrency: fav.fromCurrency,
        toCurrency: fav.toCurrency,
        rate: 0,
        timestamp: '',
      };
    });
  } catch {
    // interceptor handles error
  } finally {
    loading.value = false;
  }
};

const formatTime = (timestamp: string) => {
  if (!timestamp) return '';
  const d = new Date(timestamp);
  return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`;
};

const fetchLatestRates = async () => {
  ratesLoading.value = true;
  try {
    const resp = await axios.get<RateHistory[]>('/rates/latest');
    latestRates.value = resp.data.slice(0, 8);
  } catch {
    // interceptor handles error
  } finally {
    ratesLoading.value = false;
  }
};

const fetchRecentPosts = async () => {
  postsLoading.value = true;
  try {
    const resp = await axios.get<Post[]>('/posts', { params: { pageSize: 5 } });
    recentPosts.value = resp.data;
  } catch {
    // interceptor handles error
  } finally {
    postsLoading.value = false;
  }
};

watch(() => authStore.isAuthenticated, (val) => {
  if (val) fetchFavoriteRates();
});

onMounted(() => {
  if (authStore.isAuthenticated) {
    fetchFavoriteRates();
  }
  fetchLatestRates();
  fetchRecentPosts();
});
</script>

<style scoped>
.home-container {
  max-width: 960px;
  margin: 0 auto;
  padding: 20px;
}

.dashboard-header {
  margin-bottom: 24px;
}

.dashboard-header h1 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 8px;
}

.dashboard-header .subtitle {
  color: #909399;
  font-size: 14px;
}

.dashboard-section {
  margin-bottom: 20px;
  border-radius: 8px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.rate-card {
  text-align: center;
  cursor: pointer;
  transition: transform 0.2s;
}

.rate-card:hover {
  transform: translateY(-2px);
}

.rate-pair {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.rate-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.rate-time {
  font-size: 12px;
  color: #c0c4cc;
  margin-top: 4px;
}

.empty-state {
  padding: 40px 0;
}

.action-btn {
  width: 100%;
  height: 80px;
  font-size: 16px;
}

/* Landing page styles */
.hero {
  text-align: center;
  margin-bottom: 60px;
  padding-top: 60px;
}

.title {
  font-size: 48px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 12px;
}

.subtitle {
  font-size: 18px;
  color: #909399;
  margin-bottom: 32px;
}

.actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.features {
  margin-top: 40px;
}

.feature-card {
  text-align: center;
  padding: 20px 0;
}

.feature-icon {
  font-size: 40px;
  margin-bottom: 12px;
}

.feature-card h3 {
  margin-bottom: 8px;
  color: #303133;
}

.feature-card p {
  color: #909399;
  font-size: 14px;
}

.guest-section {
  margin-bottom: 20px;
  border-radius: 8px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.rate-card {
  text-align: center;
  cursor: pointer;
  transition: transform 0.2s;
}

.rate-card:hover {
  transform: translateY(-2px);
}

.rate-pair {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.rate-value {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.ai-preview {
  text-align: center;
  padding: 20px 0;
}

.ai-preview p {
  color: #606266;
  margin-bottom: 16px;
}

.post-preview {
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
}

.post-preview:last-child {
  border-bottom: none;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.post-username {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.post-content {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-footer {
  font-size: 12px;
  color: #909399;
}
</style>
