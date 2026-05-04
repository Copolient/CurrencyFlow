<template>
  <div class="exchange-container">
    <el-card class="exchange-card">
      <template #header>
        <span>货币兑换计算器</span>
      </template>

      <el-form :model="form" label-width="80px">
        <el-form-item label="从货币">
          <el-select v-model="form.fromCurrency" placeholder="选择货币" style="width: 100%">
            <el-option
              v-for="currency in currencies"
              :key="currency"
              :label="currency"
              :value="currency"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="到货币">
          <el-select v-model="form.toCurrency" placeholder="选择货币" style="width: 100%">
            <el-option
              v-for="currency in currencies"
              :key="currency"
              :label="currency"
              :value="currency"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="金额">
          <el-input v-model.number="form.amount" type="number" placeholder="输入金额" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="exchange" :loading="loading" style="width: 100%">
            计算兑换
          </el-button>
        </el-form-item>
      </el-form>

      <el-result
        v-if="result !== null"
        icon="success"
        :title="`${form.amount} ${form.fromCurrency} = ${result.toFixed(2)} ${form.toCurrency}`"
        sub-title="兑换结果基于当前汇率"
      />

      <el-alert
        v-if="noRate"
        title="未找到该货币对的汇率"
        type="warning"
        show-icon
        :closable="false"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from '../axios';
import type { ExchangeRate } from '../types/Article';

const form = ref({ fromCurrency: '', toCurrency: '', amount: 0 });
const result = ref<number | null>(null);
const noRate = ref(false);
const currencies = ref<string[]>([]);
const rates = ref<ExchangeRate[]>([]);
const loading = ref(false);

const fetchCurrencies = async () => {
  loading.value = true;
  try {
    const response = await axios.get<ExchangeRate[]>('/exchangeRates');
    rates.value = response.data;
    const currencySet = new Set<string>();
    response.data.forEach((r) => {
      currencySet.add(r.fromCurrency);
      currencySet.add(r.toCurrency);
    });
    currencies.value = Array.from(currencySet).sort();
  } catch {
    // interceptor 已处理错误提示
  } finally {
    loading.value = false;
  }
};

const exchange = () => {
  noRate.value = false;
  result.value = null;

  if (!form.value.fromCurrency || !form.value.toCurrency || !form.value.amount) {
    return;
  }

  const rate = rates.value.find(
    (r) => r.fromCurrency === form.value.fromCurrency && r.toCurrency === form.value.toCurrency,
  );

  if (rate) {
    result.value = form.value.amount * rate.rate;
  } else {
    noRate.value = true;
  }
};

onMounted(fetchCurrencies);
</script>

<style scoped>
.exchange-container {
  max-width: 600px;
  margin: 40px auto;
  padding: 0 20px;
}

.exchange-card {
  border-radius: 8px;
}
</style>
