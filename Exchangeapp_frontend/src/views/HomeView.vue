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
        <p class="subtitle">实时汇率查询 · 财经资讯 · 一站掌握</p>
        <div class="actions">
          <el-button type="primary" size="large" @click="router.push('/exchange')">
            开始兑换
          </el-button>
          <el-button size="large" @click="router.push('/chart')">
            行情走势
          </el-button>
        </div>
      </div>

      <div class="features">
        <el-row :gutter="24">
          <el-col :span="6">
            <el-card shadow="hover" class="feature-card">
              <div class="feature-icon">💱</div>
              <h3>实时汇率</h3>
              <p>支持多币种实时查询与兑换计算</p>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="feature-card">
              <div class="feature-icon">📈</div>
              <h3>行情走势</h3>
              <p>交互式图表，支持多时间范围</p>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="feature-card">
              <div class="feature-icon">📰</div>
              <h3>财经资讯</h3>
              <p>精选财经文章，洞察市场动态</p>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="feature-card">
              <div class="feature-icon">🔒</div>
              <h3>安全可靠</h3>
              <p>JWT 鉴权 + HTTPS，数据安全有保障</p>
            </el-card>
          </el-col>
        </el-row>
      </div>
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

interface RateHistory {
  _id: number;
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  timestamp: string;
}

const favoriteRates = ref<RateHistory[]>([]);

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

watch(() => authStore.isAuthenticated, (val) => {
  if (val) fetchFavoriteRates();
});

onMounted(() => {
  if (authStore.isAuthenticated) {
    fetchFavoriteRates();
  }
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
</style>
