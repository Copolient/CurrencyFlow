<template>
  <div class="chart-page">
    <!-- Header -->
    <div class="chart-top cf-animate-in">
      <div class="chart-title-row">
        <h1 class="page-title">行情走势</h1>
        <div class="ws-badge" :class="{ online: wsConnected }">
          <span class="ws-dot"></span>
          {{ wsConnected ? '实时' : '离线' }}
        </div>
      </div>
      <p class="page-desc">交互式汇率图表，支持多时间维度分析</p>
    </div>

    <!-- Controls -->
    <div class="chart-controls cf-glass cf-animate-in cf-delay-1">
      <div class="pair-selector">
        <el-select v-model="fromCurrency" placeholder="源货币" size="large" class="pair-select">
          <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
        </el-select>
        <span class="pair-arrow">→</span>
        <el-select v-model="toCurrency" placeholder="目标货币" size="large" class="pair-select">
          <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
        </el-select>
      </div>
      <div class="range-tabs">
        <button
          v-for="range in ranges"
          :key="range.value"
          class="range-btn"
          :class="{ active: timeRange === range.value }"
          @click="timeRange = range.value; fetchHistory()"
        >
          {{ range.label }}
        </button>
      </div>
    </div>

    <!-- Chart -->
    <div class="chart-area cf-glass cf-animate-in cf-delay-2" v-loading="loading">
      <v-chart :option="chartOption" autoresize class="chart-canvas" />
    </div>

    <!-- Latest Rates Ticker -->
    <div class="latest-section cf-animate-in cf-delay-3" v-if="latestRates.length > 0">
      <h2 class="section-title">最新汇率</h2>
      <div class="latest-grid">
        <div
          v-for="rate in latestRates"
          :key="rate._id"
          class="latest-card cf-glass cf-glass-hover"
        >
          <div class="latest-pair">{{ rate.fromCurrency }}/{{ rate.toCurrency }}</div>
          <div class="latest-value">{{ rate.rate.toFixed(4) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import VChart from 'vue-echarts';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart } from 'echarts/charts';
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DataZoomComponent,
} from 'echarts/components';
import axios from '../axios';
import { useWebSocket } from '../composables/useWebSocket';

use([CanvasRenderer, LineChart, TitleComponent, TooltipComponent, GridComponent, DataZoomComponent]);

interface RateHistory {
  _id: number;
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  timestamp: string;
}

const fromCurrency = ref('USD');
const toCurrency = ref('CNY');
const timeRange = ref('1M');
const loading = ref(false);
const history = ref<RateHistory[]>([]);
const latestRates = ref<RateHistory[]>([]);
const currencies = ref<string[]>([]);

const ranges = [
  { value: '1D', label: '1天' },
  { value: '1W', label: '1周' },
  { value: '1M', label: '1月' },
  { value: '3M', label: '3月' },
  { value: '1Y', label: '1年' },
];

const { connected: wsConnected, lastUpdate } = useWebSocket();

watch(lastUpdate, (update) => {
  if (!update) return;
  const idx = latestRates.value.findIndex(
    (r) => r.fromCurrency === update.fromCurrency && r.toCurrency === update.toCurrency
  );
  if (idx >= 0) {
    latestRates.value[idx] = { ...latestRates.value[idx], rate: update.rate, timestamp: update.timestamp };
  }
  if (update.fromCurrency === fromCurrency.value && update.toCurrency === toCurrency.value) {
    history.value.push({ _id: Date.now(), ...update });
  }
});

const chartOption = computed(() => {
  if (history.value.length === 0) {
    return {
      title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#52525b', fontSize: 14, fontWeight: 500 } },
      backgroundColor: 'transparent',
    };
  }

  const dates = history.value.map((h) => {
    const d = new Date(h.timestamp);
    return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`;
  });
  const rates = history.value.map((h) => h.rate);

  return {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 15, 15, 0.95)',
      borderColor: 'rgba(255, 255, 255, 0.08)',
      textStyle: { color: '#fafafa', fontSize: 13 },
      formatter: (params: any) => {
        const p = params[0];
        return `<span style="color:#a1a1aa">${p.axisValue}</span><br/><span style="color:#10b981;font-weight:700">${fromCurrency.value}/${toCurrency.value}: ${p.value.toFixed(4)}</span>`;
      },
    },
    grid: { left: '3%', right: '4%', bottom: '15%', top: '5%', containLabel: true },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: { rotate: 30, color: '#52525b', fontSize: 11 },
      axisLine: { lineStyle: { color: 'rgba(255,255,255,0.06)' } },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      scale: true,
      axisLabel: { formatter: (v: number) => v.toFixed(4), color: '#52525b', fontSize: 11 },
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.04)' } },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    dataZoom: [
      { type: 'inside', start: 0, end: 100 },
      {
        type: 'slider',
        start: 0,
        end: 100,
        height: 20,
        bottom: 8,
        borderColor: 'transparent',
        backgroundColor: 'rgba(255,255,255,0.02)',
        fillerColor: 'rgba(16, 185, 129, 0.08)',
        handleStyle: { color: '#10b981', borderColor: '#10b981' },
        textStyle: { color: '#52525b' },
      },
    ],
    series: [
      {
        name: `${fromCurrency.value}/${toCurrency.value}`,
        type: 'line',
        data: rates,
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: '#10b981' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(16, 185, 129, 0.25)' },
              { offset: 1, color: 'rgba(16, 185, 129, 0.02)' },
            ],
          },
        },
      },
    ],
  };
});

const fetchHistory = async () => {
  if (!fromCurrency.value || !toCurrency.value) return;
  loading.value = true;
  try {
    const resp = await axios.get<RateHistory[]>('/rates/history', {
      params: { from: fromCurrency.value, to: toCurrency.value, range: timeRange.value },
    });
    history.value = resp.data;
  } catch { /* */ } finally { loading.value = false; }
};

const fetchLatest = async () => {
  try {
    const resp = await axios.get<RateHistory[]>('/rates/latest');
    latestRates.value = resp.data;
    const set = new Set<string>();
    resp.data.forEach((r) => { set.add(r.fromCurrency); set.add(r.toCurrency); });
    currencies.value = Array.from(set).sort();
  } catch { /* */ }
};

watch([fromCurrency, toCurrency], () => { fetchHistory(); });

onMounted(() => { fetchLatest(); fetchHistory(); });
</script>

<style scoped>
.chart-page {
  padding-top: 48px;
}

/* ── Header ── */
.chart-top {
  margin-bottom: 28px;
}

.chart-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.page-title {
  font-size: 28px;
  font-weight: 800;
  color: var(--cf-text);
  margin: 0;
  letter-spacing: -0.02em;
}

.page-desc {
  font-size: 14px;
  color: var(--cf-text-muted);
  margin: 0;
}

.ws-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  font-size: 11px;
  font-weight: 600;
  color: var(--cf-text-muted);
  background: var(--cf-surface);
  border: 1px solid var(--cf-border);
  border-radius: 100px;
}

.ws-badge.online {
  color: var(--cf-accent);
  border-color: rgba(16, 185, 129, 0.2);
  background: var(--cf-accent-subtle);
}

.ws-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--cf-text-muted);
}

.ws-badge.online .ws-dot {
  background: var(--cf-accent);
  animation: cf-glow-pulse 2s ease-in-out infinite;
}

/* ── Controls ── */
.chart-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  margin-bottom: 16px;
  gap: 16px;
  flex-wrap: wrap;
}

.pair-selector {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pair-select {
  width: 130px;
}

.pair-arrow {
  font-size: 18px;
  color: var(--cf-text-muted);
  font-weight: 300;
}

.range-tabs {
  display: flex;
  gap: 4px;
  background: var(--cf-surface);
  border-radius: 8px;
  padding: 3px;
}

.range-btn {
  padding: 6px 14px;
  font-size: 12px;
  font-weight: 500;
  font-family: var(--cf-font);
  color: var(--cf-text-muted);
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: color 0.2s ease, background 0.2s ease;
}

.range-btn:hover {
  color: var(--cf-text-secondary);
}

.range-btn.active {
  color: var(--cf-black);
  background: var(--cf-accent);
}

/* ── Chart ── */
.chart-area {
  padding: 20px;
  min-height: 420px;
}

.chart-canvas {
  width: 100%;
  height: 400px;
}

/* ── Latest Rates ── */
.latest-section {
  margin-top: 32px;
  padding-bottom: 32px;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--cf-text);
  margin: 0 0 16px;
}

.latest-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 10px;
}

.latest-card {
  padding: 16px;
  text-align: center;
}

.latest-pair {
  font-size: 11px;
  font-weight: 500;
  color: var(--cf-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 6px;
}

.latest-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--cf-text);
  font-variant-numeric: tabular-nums;
}

/* ── Responsive ── */
@media (max-width: 767px) {
  .chart-page {
    padding-top: 24px;
  }

  .chart-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .pair-selector {
    flex-wrap: wrap;
  }

  .pair-select {
    flex: 1;
    min-width: 0;
  }

  .range-tabs {
    justify-content: center;
  }

  .chart-canvas {
    height: 300px;
  }
}
</style>
