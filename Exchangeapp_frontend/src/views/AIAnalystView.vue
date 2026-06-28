<template>
  <div class="ai-page">
    <div class="page-header cf-animate-in">
      <h1 class="page-title">AI 汇率分析师</h1>
      <p class="page-desc">基于历史数据的智能汇率分析，趋势判断与关键价位</p>
    </div>

    <!-- Input Form -->
    <div class="ai-form cf-glass cf-animate-in cf-delay-1">
      <div class="form-row">
        <div class="pair-input">
          <label class="field-label">货币对</label>
          <div class="pair-selectors">
            <el-select v-model="form.from" placeholder="源货币" size="large" class="pair-sel">
              <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
            </el-select>
            <span class="pair-arrow">→</span>
            <el-select v-model="form.to" placeholder="目标货币" size="large" class="pair-sel">
              <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
            </el-select>
          </div>
        </div>
      </div>

      <div class="form-row">
        <label class="field-label">你的问题</label>
        <el-input
          v-model="form.question"
          type="textarea"
          :rows="2"
          placeholder="例如：美元还会继续涨吗？现在适合换汇吗？"
        />
      </div>

      <button
        class="analyze-btn"
        @click="analyze"
        :disabled="loading || !form.from || !form.to"
      >
        <span v-if="loading" class="btn-loading">
          <span class="spinner"></span>
          分析中...
        </span>
        <span v-else>开始分析</span>
      </button>
    </div>

    <!-- Result -->
    <transition name="result">
      <div v-if="result" class="ai-result cf-animate-in">
        <!-- Trend Badge -->
        <div class="trend-banner" :class="result.trend">
          <span class="trend-icon">{{ trendIcon }}</span>
          <span class="trend-label">{{ trendLabel }}</span>
        </div>

        <!-- Analysis -->
        <div class="analysis-card cf-glass">
          <div class="analysis-content" v-html="formatAnalysis(result.analysis)"></div>
        </div>

        <!-- Key Levels -->
        <div class="levels-row">
          <div class="level-card cf-glass">
            <span class="level-label">支撑位</span>
            <span class="level-value support">{{ result.keyLevels.support.toFixed(4) }}</span>
          </div>
          <div class="level-card cf-glass">
            <span class="level-label">阻力位</span>
            <span class="level-value resistance">{{ result.keyLevels.resistance.toFixed(4) }}</span>
          </div>
        </div>

        <!-- Risk Warning -->
        <div class="risk-card">
          <span class="risk-icon">⚠</span>
          <span class="risk-text">{{ result.riskWarning }}</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import axios from '../axios';

interface AnalysisResult {
  analysis: string;
  trend: string;
  keyLevels: { support: number; resistance: number };
  riskWarning: string;
}

const loading = ref(false);
const result = ref<AnalysisResult | null>(null);
const currencies = ref<string[]>([]);

const form = ref({ from: 'USD', to: 'CNY', question: '' });

const trendIcon = computed(() => {
  switch (result.value?.trend) {
    case 'bullish': return '↑';
    case 'bearish': return '↓';
    default: return '→';
  }
});

const trendLabel = computed(() => {
  switch (result.value?.trend) {
    case 'bullish': return '看涨趋势';
    case 'bearish': return '看跌趋势';
    default: return '震荡整理';
  }
});

const escapeHtml = (str: string) =>
  str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');

const formatAnalysis = (text: string) => {
  const safe = escapeHtml(text);
  return safe
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/## (.*?)(<br>|$)/g, '<h3>$1</h3>');
};

const fetchCurrencies = async () => {
  try {
    const resp = await axios.get('/rates/latest');
    const set = new Set<string>();
    resp.data.forEach((r: any) => { set.add(r.fromCurrency); set.add(r.toCurrency); });
    currencies.value = Array.from(set).sort();
  } catch { /* */ }
};

const analyze = async () => {
  if (!form.value.from || !form.value.to) return;
  loading.value = true;
  result.value = null;
  try {
    const resp = await axios.post<AnalysisResult>('/ai/analyze', form.value);
    result.value = resp.data;
  } catch { /* */ } finally { loading.value = false; }
};

onMounted(() => { fetchCurrencies(); });
</script>

<style scoped>
.ai-page {
  padding-top: 48px;
  max-width: 720px;
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

/* ── Form ── */
.ai-form {
  padding: 24px;
  margin-bottom: 24px;
}

.form-row {
  margin-bottom: 20px;
}

.form-row:last-of-type {
  margin-bottom: 24px;
}

.field-label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--cf-text-muted);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.pair-selectors {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pair-sel {
  flex: 1;
}

.pair-arrow {
  font-size: 18px;
  color: var(--cf-text-muted);
  font-weight: 300;
}

.analyze-btn {
  width: 100%;
  height: 48px;
  background: var(--cf-accent);
  color: var(--cf-black);
  border: none;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  font-family: var(--cf-font);
  cursor: pointer;
  transition: background 0.2s ease, transform 0.1s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.analyze-btn:hover:not(:disabled) {
  background: var(--cf-accent-hover);
}

.analyze-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.analyze-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-loading {
  display: flex;
  align-items: center;
  gap: 8px;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(0, 0, 0, 0.2);
  border-top-color: var(--cf-black);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Result ── */
.trend-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 16px;
  border-radius: var(--cf-radius);
  margin-bottom: 16px;
  font-weight: 600;
}

.trend-banner.bullish {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.2);
  color: #4ade80;
}

.trend-banner.bearish {
  background: rgba(244, 63, 94, 0.1);
  border: 1px solid rgba(244, 63, 94, 0.2);
  color: #fb7185;
}

.trend-banner.neutral {
  background: var(--cf-accent-subtle);
  border: 1px solid rgba(16, 185, 129, 0.15);
  color: var(--cf-accent);
}

.trend-icon {
  font-size: 20px;
}

.trend-label {
  font-size: 16px;
}

.analysis-card {
  padding: 24px;
  margin-bottom: 16px;
}

.analysis-content {
  font-size: 14px;
  line-height: 1.8;
  color: var(--cf-text-secondary);
}

.analysis-content :deep(strong) {
  color: var(--cf-text);
  font-weight: 600;
}

.analysis-content :deep(h3) {
  font-size: 16px;
  font-weight: 600;
  color: var(--cf-text);
  margin: 20px 0 8px;
}

/* ── Key Levels ── */
.levels-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}

.level-card {
  padding: 20px;
  text-align: center;
}

.level-label {
  display: block;
  font-size: 11px;
  font-weight: 500;
  color: var(--cf-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 8px;
}

.level-value {
  font-size: 24px;
  font-weight: 800;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.01em;
}

.level-value.support {
  color: #4ade80;
}

.level-value.resistance {
  color: #fb7185;
}

/* ── Risk Warning ── */
.risk-card {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 14px 18px;
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.15);
  border-radius: 10px;
}

.risk-icon {
  font-size: 16px;
  color: var(--cf-amber);
  flex-shrink: 0;
  line-height: 1.4;
}

.risk-text {
  font-size: 13px;
  line-height: 1.6;
  color: var(--cf-text-secondary);
}

/* ── Transitions ── */
.result-enter-active {
  transition: opacity 0.4s ease, transform 0.4s ease;
}

.result-leave-active {
  transition: opacity 0.2s ease;
}

.result-enter-from {
  opacity: 0;
  transform: translateY(16px);
}

.result-leave-to {
  opacity: 0;
}

/* ── Responsive ── */
@media (max-width: 767px) {
  .ai-page {
    padding-top: 24px;
  }

  .pair-selectors {
    flex-direction: column;
    gap: 8px;
  }

  .pair-sel {
    width: 100%;
  }

  .levels-row {
    grid-template-columns: 1fr;
  }
}
</style>
