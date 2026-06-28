<template>
  <div class="exchange-page">
    <div class="exchange-hero cf-animate-in">
      <h1 class="page-title">货币兑换</h1>
      <p class="page-desc">选择货币对，输入金额，即时计算兑换结果</p>
    </div>

    <div class="exchange-card cf-glass cf-animate-in cf-delay-1">
      <!-- From Currency -->
      <div class="field-group">
        <label class="field-label">从</label>
        <div class="currency-row">
          <el-select
            v-model="form.fromCurrency"
            placeholder="选择货币"
            class="currency-select"
            size="large"
          >
            <el-option
              v-for="currency in currencies"
              :key="currency"
              :label="currency"
              :value="currency"
            />
          </el-select>
          <el-input
            v-model.number="form.amount"
            type="number"
            placeholder="0.00"
            size="large"
            class="amount-input"
          />
        </div>
      </div>

      <!-- Swap Button -->
      <div class="swap-row">
        <button class="swap-btn" @click="swapCurrencies" :class="{ rotated: swapped }">
          ↕
        </button>
      </div>

      <!-- To Currency -->
      <div class="field-group">
        <label class="field-label">到</label>
        <div class="currency-row">
          <el-select
            v-model="form.toCurrency"
            placeholder="选择货币"
            class="currency-select"
            size="large"
          >
            <el-option
              v-for="currency in currencies"
              :key="currency"
              :label="currency"
              :value="currency"
            />
          </el-select>
          <div class="result-display" :class="{ active: result !== null }">
            <template v-if="result !== null">
              {{ result.toFixed(2) }}
            </template>
            <template v-else>
              <span class="result-placeholder">0.00</span>
            </template>
          </div>
        </div>
      </div>

      <!-- Exchange Button -->
      <button
        class="exchange-btn"
        @click="exchange"
        :disabled="loading || !form.fromCurrency || !form.toCurrency || !form.amount"
      >
        <span v-if="loading" class="btn-loading"></span>
        <span v-else>计算兑换</span>
      </button>

      <!-- Result Message -->
      <transition name="result">
        <div v-if="result !== null" class="result-banner">
          <div class="result-text">
            {{ form.amount }} {{ form.fromCurrency }} = <strong>{{ result.toFixed(4) }}</strong> {{ form.toCurrency }}
          </div>
          <div class="result-sub">基于当前实时汇率</div>
        </div>
      </transition>

      <div v-if="noRate" class="error-banner">
        未找到该货币对的汇率
      </div>
    </div>
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
const swapped = ref(false);

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
  } catch { /* */ } finally { loading.value = false; }
};

const swapCurrencies = () => {
  const temp = form.value.fromCurrency;
  form.value.fromCurrency = form.value.toCurrency;
  form.value.toCurrency = temp;
  swapped.value = !swapped.value;
  result.value = null;
};

const exchange = () => {
  noRate.value = false;
  result.value = null;

  if (!form.value.fromCurrency || !form.value.toCurrency || !form.value.amount) return;

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
.exchange-page {
  padding-top: 48px;
  max-width: 520px;
  margin: 0 auto;
}

.exchange-hero {
  text-align: center;
  margin-bottom: 40px;
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

/* ── Card ── */
.exchange-card {
  padding: 32px;
}

/* ── Field Groups ── */
.field-group {
  margin-bottom: 8px;
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

.currency-row {
  display: flex;
  gap: 12px;
}

.currency-select {
  width: 160px;
  flex-shrink: 0;
}

.amount-input {
  flex: 1;
}

.result-display {
  flex: 1;
  height: 40px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  background: var(--cf-surface);
  border: 1px solid var(--cf-border);
  border-radius: 8px;
  font-size: 18px;
  font-weight: 700;
  color: var(--cf-text);
  font-variant-numeric: tabular-nums;
  transition: border-color 0.3s ease;
}

.result-display.active {
  border-color: var(--cf-accent);
  box-shadow: 0 0 0 3px var(--cf-accent-glow);
}

.result-placeholder {
  color: var(--cf-text-muted);
  font-weight: 400;
}

/* ── Swap Button ── */
.swap-row {
  display: flex;
  justify-content: center;
  padding: 8px 0;
}

.swap-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--cf-surface-hover);
  border: 1px solid var(--cf-border);
  color: var(--cf-text-secondary);
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease, border-color 0.2s ease, color 0.2s ease;
}

.swap-btn:hover {
  border-color: var(--cf-accent);
  color: var(--cf-accent);
}

.swap-btn.rotated {
  transform: rotate(180deg);
}

/* ── Exchange Button ── */
.exchange-btn {
  width: 100%;
  height: 48px;
  margin-top: 24px;
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

.exchange-btn:hover:not(:disabled) {
  background: var(--cf-accent-hover);
}

.exchange-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.exchange-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-loading {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(0, 0, 0, 0.2);
  border-top-color: var(--cf-black);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Result Banner ── */
.result-banner {
  margin-top: 20px;
  padding: 16px 20px;
  background: var(--cf-accent-subtle);
  border: 1px solid rgba(16, 185, 129, 0.15);
  border-radius: 10px;
  text-align: center;
}

.result-text {
  font-size: 18px;
  font-weight: 600;
  color: var(--cf-text);
}

.result-text strong {
  color: var(--cf-accent);
  font-weight: 700;
}

.result-sub {
  font-size: 12px;
  color: var(--cf-text-muted);
  margin-top: 4px;
}

.error-banner {
  margin-top: 16px;
  padding: 12px 16px;
  background: rgba(244, 63, 94, 0.1);
  border: 1px solid rgba(244, 63, 94, 0.2);
  border-radius: 8px;
  font-size: 13px;
  color: var(--cf-rose);
  text-align: center;
}

/* ── Result Transition ── */
.result-enter-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.result-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.result-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.result-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* ── Responsive ── */
@media (max-width: 767px) {
  .exchange-page {
    padding-top: 24px;
  }

  .exchange-card {
    padding: 20px;
  }

  .currency-row {
    flex-direction: column;
    gap: 8px;
  }

  .currency-select {
    width: 100%;
  }
}
</style>
