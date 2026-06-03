<template>
  <div class="chart-container">
    <el-card class="chart-card">
      <template #header>
        <div class="chart-header">
          <span>汇率走势</span>
          <el-tag :type="wsConnected ? 'success' : 'danger'" size="small">
            {{ wsConnected ? '实时' : '离线' }}
          </el-tag>
          <div class="chart-controls">
            <el-select v-model="fromCurrency" placeholder="源货币" style="width: 100px">
              <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
            </el-select>
            <span class="arrow">→</span>
            <el-select v-model="toCurrency" placeholder="目标货币" style="width: 100px">
              <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
            </el-select>
          </div>
        </div>
      </template>

      <div class="range-tabs">
        <el-radio-group v-model="timeRange" @change="fetchHistory">
          <el-radio-button label="1D">1天</el-radio-button>
          <el-radio-button label="1W">1周</el-radio-button>
          <el-radio-button label="1M">1月</el-radio-button>
          <el-radio-button label="3M">3月</el-radio-button>
          <el-radio-button label="1Y">1年</el-radio-button>
        </el-radio-group>
      </div>

      <div class="chart-wrapper" v-loading="loading">
        <v-chart :option="chartOption" autoresize style="height: 400px" />
      </div>

      <div class="latest-rates" v-if="latestRates.length > 0">
        <h3>最新汇率</h3>
        <el-row :gutter="16">
          <el-col :span="6" v-for="rate in latestRates" :key="rate._id">
            <el-card shadow="hover" class="rate-mini-card">
              <div class="rate-pair">{{ rate.fromCurrency }}/{{ rate.toCurrency }}</div>
              <div class="rate-value">{{ rate.rate.toFixed(4) }}</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
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

// WebSocket for real-time updates
const { connected: wsConnected, lastUpdate } = useWebSocket();

// Watch for real-time updates
watch(lastUpdate, (update) => {
  if (!update) return;

  // Update latest rates
  const idx = latestRates.value.findIndex(
    (r) => r.fromCurrency === update.fromCurrency && r.toCurrency === update.toCurrency
  );
  if (idx >= 0) {
    latestRates.value[idx] = {
      ...latestRates.value[idx],
      rate: update.rate,
      timestamp: update.timestamp,
    };
  }

  // Append to history if viewing the same pair
  if (update.fromCurrency === fromCurrency.value && update.toCurrency === toCurrency.value) {
    history.value.push({
      _id: Date.now(),
      fromCurrency: update.fromCurrency,
      toCurrency: update.toCurrency,
      rate: update.rate,
      timestamp: update.timestamp,
    });
  }
});

const chartOption = computed(() => {
  if (history.value.length === 0) {
    return {
      title: { text: '暂无数据', left: 'center', top: 'center' },
    };
  }

  const dates = history.value.map((h) => {
    const d = new Date(h.timestamp);
    return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`;
  });
  const rates = history.value.map((h) => h.rate);

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const p = params[0];
        return `${p.axisValue}<br/>${fromCurrency.value}/${toCurrency.value}: <b>${p.value.toFixed(4)}</b>`;
      },
    },
    grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: { rotate: 30 },
    },
    yAxis: {
      type: 'value',
      scale: true,
      axisLabel: { formatter: (v: number) => v.toFixed(4) },
    },
    dataZoom: [
      { type: 'inside', start: 0, end: 100 },
      { type: 'slider', start: 0, end: 100 },
    ],
    series: [
      {
        name: `${fromCurrency.value}/${toCurrency.value}`,
        type: 'line',
        data: rates,
        smooth: true,
        lineStyle: { width: 2 },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.05)' },
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
  } catch {
    // interceptor handles error
  } finally {
    loading.value = false;
  }
};

const fetchLatest = async () => {
  try {
    const resp = await axios.get<RateHistory[]>('/rates/latest');
    latestRates.value = resp.data;
    // Extract unique currencies
    const set = new Set<string>();
    resp.data.forEach((r) => {
      set.add(r.fromCurrency);
      set.add(r.toCurrency);
    });
    currencies.value = Array.from(set).sort();
  } catch {
    // interceptor handles error
  }
};

watch([fromCurrency, toCurrency], () => {
  fetchHistory();
});

onMounted(() => {
  fetchLatest();
  fetchHistory();
});
</script>

<style scoped>
.chart-container {
  max-width: 960px;
  margin: 20px auto;
  padding: 0 20px;
}

.chart-card {
  border-radius: 8px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.arrow {
  font-size: 18px;
  color: #909399;
}

.range-tabs {
  margin-bottom: 16px;
  text-align: center;
}

.chart-wrapper {
  min-height: 400px;
}

.latest-rates {
  margin-top: 24px;
}

.latest-rates h3 {
  margin-bottom: 12px;
  color: #303133;
}

.rate-mini-card {
  text-align: center;
  margin-bottom: 8px;
}

.rate-pair {
  font-size: 12px;
  color: #909399;
}

.rate-value {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-top: 4px;
  transition: color 0.3s ease;
}

.rate-value.up {
  color: #67c23a;
}

.rate-value.down {
  color: #f56c6c;
}
</style>
