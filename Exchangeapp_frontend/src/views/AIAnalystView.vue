<template>
  <div class="ai-container">
    <el-card class="ai-card">
      <template #header>
        <div class="card-header">
          <span>🤖 AI 汇率分析师</span>
          <el-tag type="info" size="small">基于历史数据分析</el-tag>
        </div>
      </template>

      <el-form :model="form" label-width="80px" class="analysis-form">
        <el-form-item label="货币对">
          <div class="currency-input">
            <el-select v-model="form.from" placeholder="源货币" style="width: 120px">
              <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
            </el-select>
            <span class="arrow">→</span>
            <el-select v-model="form.to" placeholder="目标货币" style="width: 120px">
              <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
            </el-select>
          </div>
        </el-form-item>

        <el-form-item label="问题">
          <el-input
            v-model="form.question"
            type="textarea"
            :rows="2"
            placeholder="例如：美元还会继续涨吗？现在适合换汇吗？"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="analyze" :loading="loading" style="width: 100%">
            开始分析
          </el-button>
        </el-form-item>
      </el-form>

      <div v-if="result" class="analysis-result">
        <el-divider />

        <div class="result-header">
          <el-tag :type="trendTagType" size="large">
            {{ trendIcon }} {{ trendLabel }}
          </el-tag>
        </div>

        <div class="analysis-content" v-html="formatAnalysis(result.analysis)"></div>

        <el-card class="key-levels" shadow="never">
          <template #header>
            <span>关键价位</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="level-item">
                <span class="level-label">支撑位</span>
                <span class="level-value support">{{ result.keyLevels.support.toFixed(4) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="level-item">
                <span class="level-label">阻力位</span>
                <span class="level-value resistance">{{ result.keyLevels.resistance.toFixed(4) }}</span>
              </div>
            </el-col>
          </el-row>
        </el-card>

        <el-alert
          :title="result.riskWarning"
          type="warning"
          show-icon
          :closable="false"
          class="risk-warning"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import axios from '../axios';

interface AnalysisResult {
  analysis: string;
  trend: string;
  keyLevels: {
    support: number;
    resistance: number;
  };
  riskWarning: string;
}

const loading = ref(false);
const result = ref<AnalysisResult | null>(null);
const currencies = ref<string[]>([]);

const form = ref({
  from: 'USD',
  to: 'CNY',
  question: '',
});

const trendIcon = computed(() => {
  switch (result.value?.trend) {
    case 'bullish': return '📈';
    case 'bearish': return '📉';
    default: return '📊';
  }
});

const trendLabel = computed(() => {
  switch (result.value?.trend) {
    case 'bullish': return '看涨';
    case 'bearish': return '看跌';
    default: return '震荡';
  }
});

const trendTagType = computed(() => {
  switch (result.value?.trend) {
    case 'bullish': return 'success';
    case 'bearish': return 'danger';
    default: return 'info';
  }
});

const formatAnalysis = (text: string) => {
  return text
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/## (.*?)(<br>|$)/g, '<h3>$1</h3>');
};

const fetchCurrencies = async () => {
  try {
    const resp = await axios.get('/rates/latest');
    const set = new Set<string>();
    resp.data.forEach((r: any) => {
      set.add(r.fromCurrency);
      set.add(r.toCurrency);
    });
    currencies.value = Array.from(set).sort();
  } catch {
    // interceptor handles error
  }
};

const analyze = async () => {
  if (!form.value.from || !form.value.to) return;

  loading.value = true;
  result.value = null;

  try {
    const resp = await axios.post<AnalysisResult>('/ai/analyze', form.value);
    result.value = resp.data;
  } catch {
    // interceptor handles error
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchCurrencies();
});
</script>

<style scoped>
.ai-container {
  max-width: 720px;
  margin: 20px auto;
  padding: 0 20px;
}

.ai-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.analysis-form {
  margin-bottom: 0;
}

.currency-input {
  display: flex;
  align-items: center;
  gap: 8px;
}

.arrow {
  font-size: 18px;
  color: #909399;
}

.analysis-result {
  margin-top: 20px;
}

.result-header {
  text-align: center;
  margin-bottom: 20px;
}

.analysis-content {
  font-size: 14px;
  line-height: 1.8;
  color: #303133;
  margin-bottom: 20px;
}

.analysis-content :deep(h3) {
  font-size: 16px;
  margin: 16px 0 8px;
  color: #303133;
}

.key-levels {
  margin-bottom: 16px;
}

.level-item {
  text-align: center;
}

.level-label {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.level-value {
  font-size: 20px;
  font-weight: 600;
}

.level-value.support {
  color: #67c23a;
}

.level-value.resistance {
  color: #f56c6c;
}

.risk-warning {
  margin-top: 16px;
}
</style>
