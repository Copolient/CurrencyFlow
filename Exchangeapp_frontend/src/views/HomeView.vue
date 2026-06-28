<template>
  <div class="home">
    <!-- ═══════════ GUEST LANDING ═══════════ -->
    <template v-if="!authStore.isAuthenticated">
      <!-- Hero -->
      <section class="hero">
        <div class="hero-glow"></div>
        <div class="hero-content cf-animate-in">
          <div class="hero-badge">
            <span class="badge-dot"></span>
            实时汇率 · 智能分析 · 社交交易
          </div>
          <h1 class="hero-title">
            <span class="title-line">全球汇率，</span>
            <span class="title-line accent">一目了然。</span>
          </h1>
          <p class="hero-desc">
            实时追踪多币种汇率，AI 驱动的趋势分析，<br class="br-hide" />
            加入交易社区分享你的洞察。
          </p>
          <div class="hero-actions">
            <router-link to="/exchange" class="btn-primary-lg">
              开始兑换
              <span class="btn-arrow">→</span>
            </router-link>
            <router-link to="/chart" class="btn-ghost-lg">
              查看行情
            </router-link>
          </div>
        </div>

        <!-- Floating rate ticker -->
        <div class="hero-ticker cf-animate-in cf-delay-2">
          <div class="ticker-track" ref="tickerTrack">
            <div
              v-for="(rate, i) in [...latestRates, ...latestRates]"
              :key="`ticker-${i}`"
              class="ticker-item"
            >
              <span class="ticker-pair">{{ rate.fromCurrency }}/{{ rate.toCurrency }}</span>
              <span class="ticker-value">{{ rate.rate.toFixed(4) }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Features -->
      <section class="features">
        <div class="feature-grid">
          <div
            v-for="(f, i) in features"
            :key="f.title"
            class="feature-card cf-glass cf-glass-hover cf-animate-in"
            :style="{ animationDelay: `${0.1 + i * 0.08}s` }"
          >
            <div class="feature-icon">{{ f.icon }}</div>
            <h3 class="feature-title">{{ f.title }}</h3>
            <p class="feature-desc">{{ f.desc }}</p>
          </div>
        </div>
      </section>

      <!-- Live Rates Section -->
      <section class="rates-section cf-animate-in cf-delay-3">
        <div class="section-head">
          <h2 class="section-title">实时汇率</h2>
          <router-link to="/chart" class="section-link">
            查看全部 →
          </router-link>
        </div>
        <div class="rates-grid" v-loading="ratesLoading">
          <div
            v-for="rate in latestRates.slice(0, 8)"
            :key="`${rate.fromCurrency}-${rate.toCurrency}`"
            class="rate-card cf-glass cf-glass-hover"
            @click="router.push('/chart')"
          >
            <div class="rate-pair">{{ rate.fromCurrency }}/{{ rate.toCurrency }}</div>
            <div class="rate-value">{{ rate.rate.toFixed(4) }}</div>
            <div class="rate-bar">
              <div class="rate-bar-fill" :style="{ width: `${Math.min(rate.rate * 10, 100)}%` }"></div>
            </div>
          </div>
        </div>
      </section>

      <!-- Community Preview -->
      <section class="community-section cf-animate-in cf-delay-4">
        <div class="section-head">
          <h2 class="section-title">社区动态</h2>
          <router-link to="/community" class="section-link">
            进入社区 →
          </router-link>
        </div>
        <div v-if="recentPosts.length === 0" class="empty-hint">
          暂无帖子，成为第一个分享观点的人
        </div>
        <div v-else class="posts-list">
          <div v-for="post in recentPosts.slice(0, 3)" :key="post.ID" class="post-card cf-glass">
            <div class="post-top">
              <div class="post-avatar">{{ post.username?.charAt(0)?.toUpperCase() }}</div>
              <div class="post-meta">
                <span class="post-name">{{ post.username }}</span>
                <span class="post-time">{{ formatTime(post.CreatedAt) }}</span>
              </div>
              <span v-if="post.currency" class="post-tag">{{ post.currency }}</span>
            </div>
            <p class="post-text">{{ post.content }}</p>
            <div class="post-stats">
              <span class="stat">♥ {{ post.likes }}</span>
            </div>
          </div>
        </div>
      </section>
    </template>

    <!-- ═══════════ AUTHENTICATED DASHBOARD ═══════════ -->
    <template v-else>
      <section class="dashboard">
        <!-- Welcome -->
        <div class="dash-welcome cf-animate-in">
          <div class="welcome-text">
            <h1 class="dash-greeting">你好，{{ username }}</h1>
            <p class="dash-subtitle">你的汇率仪表盘</p>
          </div>
          <div class="welcome-time">{{ currentTime }}</div>
        </div>

        <!-- Quick Actions -->
        <div class="dash-actions cf-animate-in cf-delay-1">
          <router-link
            v-for="action in quickActions"
            :key="action.to"
            :to="action.to"
            class="action-card cf-glass cf-glass-hover"
          >
            <span class="action-icon">{{ action.icon }}</span>
            <span class="action-label">{{ action.label }}</span>
          </router-link>
        </div>

        <!-- Favorites -->
        <div class="dash-section cf-animate-in cf-delay-2">
          <div class="section-head">
            <h2 class="section-title">收藏汇率</h2>
            <router-link to="/exchange" class="section-link">管理 →</router-link>
          </div>
          <div v-loading="loading">
            <div v-if="favoriteRates.length === 0" class="empty-hint">
              <router-link to="/exchange" class="empty-link">去收藏你关注的货币对</router-link>
            </div>
            <div v-else class="rates-grid">
              <div
                v-for="rate in favoriteRates"
                :key="`${rate.fromCurrency}-${rate.toCurrency}`"
                class="rate-card cf-glass cf-glass-hover"
                @click="router.push('/chart')"
              >
                <div class="rate-pair">{{ rate.fromCurrency }}/{{ rate.toCurrency }}</div>
                <div class="rate-value">{{ rate.rate.toFixed(4) }}</div>
                <div class="rate-time">{{ formatTime(rate.timestamp) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Community Feed -->
        <div class="dash-section cf-animate-in cf-delay-3">
          <div class="section-head">
            <h2 class="section-title">社区动态</h2>
            <router-link to="/community" class="section-link">更多 →</router-link>
          </div>
          <div v-if="recentPosts.length === 0" class="empty-hint">暂无动态</div>
          <div v-else class="posts-list">
            <div v-for="post in recentPosts.slice(0, 3)" :key="post.ID" class="post-card cf-glass">
              <div class="post-top">
                <div class="post-avatar">{{ post.username?.charAt(0)?.toUpperCase() }}</div>
                <div class="post-meta">
                  <span class="post-name">{{ post.username }}</span>
                  <span class="post-time">{{ formatTime(post.CreatedAt) }}</span>
                </div>
              </div>
              <p class="post-text">{{ post.content }}</p>
            </div>
          </div>
        </div>
      </section>
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

const currentTime = computed(() => {
  const now = new Date();
  const h = now.getHours();
  const greeting = h < 6 ? '凌晨好' : h < 12 ? '上午好' : h < 18 ? '下午好' : '晚上好';
  return greeting;
});

const features = [
  { icon: '💱', title: '多币种兑换', desc: '实时汇率查询，支持全球主流货币的即时兑换计算。' },
  { icon: '📈', title: '行情走势', desc: '交互式图表，支持多时间维度的趋势分析。' },
  { icon: '🤖', title: 'AI 分析师', desc: '基于历史数据的智能分析，趋势判断与关键价位。' },
  { icon: '🔔', title: '汇率预警', desc: '设置目标汇率，触发时实时推送通知。' },
  { icon: '👥', title: '交易社区', desc: '分享交易观点，关注其他用户，获取市场洞察。' },
  { icon: '⚡', title: '实时推送', desc: 'WebSocket 驱动的实时汇率变动，数据即时更新。' },
];

const quickActions = [
  { icon: '💱', label: '兑换', to: '/exchange' },
  { icon: '📈', label: '行情', to: '/chart' },
  { icon: '🤖', label: 'AI', to: '/ai' },
  { icon: '📰', label: '资讯', to: '/news' },
  { icon: '🔔', label: '预警', to: '/alerts' },
];

const fetchFavoriteRates = async () => {
  if (!authStore.isAuthenticated) return;
  loading.value = true;
  try {
    const favResp = await axios.get('/favorites');
    const favorites = favResp.data;
    if (favorites.length === 0) { favoriteRates.value = []; return; }
    const ratesResp = await axios.get<RateHistory[]>('/rates/latest');
    const latest = ratesResp.data;
    favoriteRates.value = favorites.map((fav: any) => {
      const match = latest.find((r) => r.fromCurrency === fav.fromCurrency && r.toCurrency === fav.toCurrency);
      return match || { _id: 0, fromCurrency: fav.fromCurrency, toCurrency: fav.toCurrency, rate: 0, timestamp: '' };
    });
  } catch { /* */ } finally { loading.value = false; }
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

const fetchLatestRates = async () => {
  ratesLoading.value = true;
  try {
    const resp = await axios.get<RateHistory[]>('/rates/latest');
    latestRates.value = resp.data;
  } catch { /* */ } finally { ratesLoading.value = false; }
};

const fetchRecentPosts = async () => {
  postsLoading.value = true;
  try {
    const resp = await axios.get<Post[]>('/posts', { params: { pageSize: 5 } });
    recentPosts.value = resp.data;
  } catch { /* */ } finally { postsLoading.value = false; }
};

watch(() => authStore.isAuthenticated, (val) => {
  if (val) fetchFavoriteRates();
});

onMounted(() => {
  if (authStore.isAuthenticated) fetchFavoriteRates();
  fetchLatestRates();
  fetchRecentPosts();
});
</script>

<style scoped>
.home {
  padding-top: 24px;
}

/* ═══════════ HERO ═══════════ */
.hero {
  position: relative;
  padding: 80px 0 40px;
  overflow: hidden;
}

.hero-glow {
  position: absolute;
  top: -200px;
  left: 50%;
  transform: translateX(-50%);
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, var(--cf-accent-glow) 0%, transparent 70%);
  pointer-events: none;
  opacity: 0.6;
}

.hero-content {
  position: relative;
  text-align: center;
  max-width: 640px;
  margin: 0 auto;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  font-size: 12px;
  font-weight: 500;
  color: var(--cf-accent);
  background: var(--cf-accent-subtle);
  border: 1px solid rgba(16, 185, 129, 0.15);
  border-radius: 100px;
  margin-bottom: 24px;
}

.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--cf-accent);
  animation: cf-glow-pulse 2s ease-in-out infinite;
}

.hero-title {
  font-size: 52px;
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.03em;
  margin: 0 0 20px;
}

.title-line {
  display: block;
  color: var(--cf-text);
}

.title-line.accent {
  background: linear-gradient(135deg, var(--cf-accent), var(--cf-indigo));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-desc {
  font-size: 16px;
  line-height: 1.7;
  color: var(--cf-text-secondary);
  margin: 0 0 32px;
}

.br-hide { display: none; }

.hero-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.btn-primary-lg {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 28px;
  font-size: 15px;
  font-weight: 600;
  color: var(--cf-black);
  background: var(--cf-accent);
  border-radius: 10px;
  text-decoration: none;
  transition: background 0.2s ease, transform 0.1s ease;
}

.btn-primary-lg:hover {
  background: var(--cf-accent-hover);
  color: var(--cf-black);
}

.btn-primary-lg:active {
  transform: scale(0.97);
}

.btn-arrow {
  transition: transform 0.2s ease;
}

.btn-primary-lg:hover .btn-arrow {
  transform: translateX(3px);
}

.btn-ghost-lg {
  padding: 12px 28px;
  font-size: 15px;
  font-weight: 500;
  color: var(--cf-text-secondary);
  border: 1px solid var(--cf-border);
  border-radius: 10px;
  text-decoration: none;
  transition: border-color 0.2s ease, color 0.2s ease;
}

.btn-ghost-lg:hover {
  border-color: var(--cf-border-hover);
  color: var(--cf-text);
}

/* ── Ticker ── */
.hero-ticker {
  margin-top: 48px;
  overflow: hidden;
  mask-image: linear-gradient(90deg, transparent, black 10%, black 90%, transparent);
  -webkit-mask-image: linear-gradient(90deg, transparent, black 10%, black 90%, transparent);
}

.ticker-track {
  display: flex;
  gap: 32px;
  animation: ticker-scroll 30s linear infinite;
  width: max-content;
}

@keyframes ticker-scroll {
  from { transform: translateX(0); }
  to { transform: translateX(-50%); }
}

.ticker-item {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.ticker-pair {
  font-size: 13px;
  font-weight: 600;
  color: var(--cf-text);
}

.ticker-value {
  font-size: 13px;
  font-weight: 500;
  color: var(--cf-accent);
  font-variant-numeric: tabular-nums;
}

/* ═══════════ FEATURES ═══════════ */
.features {
  padding: 40px 0 60px;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.feature-card {
  padding: 28px 24px;
  cursor: default;
}

.feature-icon {
  font-size: 28px;
  margin-bottom: 16px;
}

.feature-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--cf-text);
  margin: 0 0 8px;
}

.feature-desc {
  font-size: 13px;
  line-height: 1.6;
  color: var(--cf-text-secondary);
  margin: 0;
}

/* ═══════════ SECTIONS ═══════════ */
.section-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--cf-text);
  margin: 0;
  letter-spacing: -0.01em;
}

.section-link {
  font-size: 13px;
  font-weight: 500;
  color: var(--cf-accent);
  text-decoration: none;
  transition: color 0.2s ease;
}

.section-link:hover {
  color: var(--cf-accent-hover);
}

.rates-section,
.community-section {
  padding-bottom: 48px;
}

/* ── Rate Cards ── */
.rates-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.rate-card {
  padding: 20px;
  text-align: center;
  cursor: pointer;
}

.rate-pair {
  font-size: 12px;
  font-weight: 500;
  color: var(--cf-text-muted);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.rate-value {
  font-size: 22px;
  font-weight: 700;
  color: var(--cf-text);
  font-variant-numeric: tabular-nums;
}

.rate-time {
  font-size: 11px;
  color: var(--cf-text-muted);
  margin-top: 6px;
}

.rate-bar {
  height: 3px;
  background: var(--cf-border);
  border-radius: 2px;
  margin-top: 12px;
  overflow: hidden;
}

.rate-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--cf-accent), var(--cf-indigo));
  border-radius: 2px;
  transition: width 0.6s ease;
}

/* ── Posts ── */
.posts-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.post-card {
  padding: 16px 20px;
}

.post-top {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.post-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--cf-surface-hover);
  border: 1px solid var(--cf-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  color: var(--cf-accent);
  flex-shrink: 0;
}

.post-meta {
  flex: 1;
  min-width: 0;
}

.post-name {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--cf-text);
}

.post-time {
  display: block;
  font-size: 11px;
  color: var(--cf-text-muted);
}

.post-tag {
  font-size: 11px;
  font-weight: 500;
  color: var(--cf-accent);
  background: var(--cf-accent-subtle);
  padding: 2px 8px;
  border-radius: 6px;
}

.post-text {
  font-size: 14px;
  line-height: 1.6;
  color: var(--cf-text-secondary);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-stats {
  margin-top: 10px;
}

.stat {
  font-size: 12px;
  color: var(--cf-text-muted);
}

.empty-hint {
  text-align: center;
  padding: 40px 0;
  font-size: 14px;
  color: var(--cf-text-muted);
}

.empty-link {
  color: var(--cf-accent);
  text-decoration: none;
}

/* ═══════════ DASHBOARD ═══════════ */
.dashboard {
  padding-top: 32px;
}

.dash-welcome {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--cf-border);
}

.dash-greeting {
  font-size: 28px;
  font-weight: 800;
  color: var(--cf-text);
  margin: 0 0 4px;
  letter-spacing: -0.02em;
}

.dash-subtitle {
  font-size: 14px;
  color: var(--cf-text-muted);
  margin: 0;
}

.welcome-time {
  font-size: 13px;
  color: var(--cf-text-muted);
}

/* Quick Actions */
.dash-actions {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
  margin-bottom: 40px;
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px 12px;
  text-decoration: none;
  cursor: pointer;
}

.action-icon {
  font-size: 24px;
}

.action-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--cf-text-secondary);
}

/* Dashboard Sections */
.dash-section {
  margin-bottom: 40px;
}

/* ═══════════ RESPONSIVE ═══════════ */
@media (max-width: 767px) {
  .hero {
    padding: 48px 0 24px;
  }

  .hero-title {
    font-size: 32px;
  }

  .hero-desc {
    font-size: 14px;
  }

  .br-hide {
    display: none;
  }

  .hero-actions {
    flex-direction: column;
  }

  .btn-primary-lg,
  .btn-ghost-lg {
    width: 100%;
    justify-content: center;
  }

  .feature-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }

  .rates-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }

  .dash-actions {
    grid-template-columns: repeat(3, 1fr);
  }

  .dash-welcome {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}

@media (min-width: 768px) {
  .br-hide {
    display: inline;
  }
}
</style>
