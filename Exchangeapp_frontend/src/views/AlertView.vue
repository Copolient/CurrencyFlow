<template>
  <div class="alert-page">
    <div class="page-header cf-animate-in">
      <div class="header-row">
        <div>
          <h1 class="page-title">汇率预警</h1>
          <p class="page-desc">设置目标汇率，达到时实时通知你</p>
        </div>
        <button class="create-btn" @click="showCreateDialog = true">
          + 创建预警
        </button>
      </div>
    </div>

    <div class="alert-content cf-glass cf-animate-in cf-delay-1" v-loading="loading">
      <div v-if="alerts.length === 0" class="empty-state">
        <div class="empty-icon">🔔</div>
        <p class="empty-text">还没有设置预警</p>
        <button class="empty-action" @click="showCreateDialog = true">创建第一个预警</button>
      </div>

      <div v-else class="alerts-list">
        <div v-for="alert in alerts" :key="alert.id" class="alert-item">
          <div class="alert-pair">
            <span class="pair-from">{{ alert.fromCurrency }}</span>
            <span class="pair-arrow">→</span>
            <span class="pair-to">{{ alert.toCurrency }}</span>
          </div>
          <div class="alert-target">
            <span class="target-label">目标</span>
            <span class="target-value">{{ alert.targetRate.toFixed(4) }}</span>
          </div>
          <div class="alert-direction">
            <span
              class="direction-tag"
              :class="alert.direction === 'above' ? 'above' : 'below'"
            >
              {{ alert.direction === 'above' ? '↑ 高于' : '↓ 低于' }}
            </span>
          </div>
          <div class="alert-status">
            <span class="status-dot" :class="{ triggered: alert.triggered }"></span>
            <span class="status-text">{{ alert.triggered ? '已触发' : '监控中' }}</span>
          </div>
          <button class="delete-btn" @click="deleteAlert(alert.id)">删除</button>
        </div>
      </div>
    </div>

    <!-- Create Dialog -->
    <el-dialog v-model="showCreateDialog" title="创建汇率预警" width="420px" :close-on-click-modal="false">
      <el-form :model="newAlert" label-position="top">
        <el-form-item label="源货币">
          <el-select v-model="newAlert.fromCurrency" placeholder="选择货币" style="width: 100%">
            <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标货币">
          <el-select v-model="newAlert.toCurrency" placeholder="选择货币" style="width: 100%">
            <el-option v-for="c in currencies" :key="c" :label="c" :value="c" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标汇率">
          <el-input v-model.number="newAlert.targetRate" type="number" placeholder="输入目标汇率" />
        </el-form-item>
        <el-form-item label="触发方向">
          <div class="direction-picker">
            <button
              class="dir-btn"
              :class="{ active: newAlert.direction === 'above', above: newAlert.direction === 'above' }"
              @click="newAlert.direction = 'above'"
            >
              ↑ 高于目标时通知
            </button>
            <button
              class="dir-btn"
              :class="{ active: newAlert.direction === 'below', below: newAlert.direction === 'below' }"
              @click="newAlert.direction = 'below'"
            >
              ↓ 低于目标时通知
            </button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <button class="cancel-btn" @click="showCreateDialog = false">取消</button>
          <button
            class="confirm-btn"
            @click="createAlert"
            :disabled="creating || !newAlert.fromCurrency || !newAlert.toCurrency || !newAlert.targetRate"
          >
            <span v-if="creating" class="btn-spinner"></span>
            <span v-else>创建</span>
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import axios from '../axios';

interface RateAlert {
  id: number;
  fromCurrency: string;
  toCurrency: string;
  targetRate: number;
  direction: string;
  triggered: boolean;
}

const loading = ref(false);
const creating = ref(false);
const alerts = ref<RateAlert[]>([]);
const currencies = ref<string[]>([]);
const showCreateDialog = ref(false);

const newAlert = ref({
  fromCurrency: '',
  toCurrency: '',
  targetRate: 0,
  direction: 'above',
});

const fetchAlerts = async () => {
  loading.value = true;
  try {
    const resp = await axios.get<RateAlert[]>('/alerts');
    alerts.value = resp.data;
  } catch { /* */ } finally { loading.value = false; }
};

const fetchCurrencies = async () => {
  try {
    const resp = await axios.get('/rates/latest');
    const set = new Set<string>();
    resp.data.forEach((r: any) => { set.add(r.fromCurrency); set.add(r.toCurrency); });
    currencies.value = Array.from(set).sort();
  } catch { /* */ }
};

const createAlert = async () => {
  if (!newAlert.value.fromCurrency || !newAlert.value.toCurrency || !newAlert.value.targetRate) {
    ElMessage.warning('请填写完整信息');
    return;
  }
  creating.value = true;
  try {
    await axios.post('/alerts', newAlert.value);
    ElMessage.success('预警创建成功');
    showCreateDialog.value = false;
    newAlert.value = { fromCurrency: '', toCurrency: '', targetRate: 0, direction: 'above' };
    fetchAlerts();
  } catch { /* */ } finally { creating.value = false; }
};

const deleteAlert = async (id: number) => {
  try {
    await axios.delete(`/alerts/${id}`);
    ElMessage.success('预警已删除');
    fetchAlerts();
  } catch { /* */ }
};

onMounted(() => { fetchAlerts(); fetchCurrencies(); });
</script>

<style scoped>
.alert-page {
  padding-top: 48px;
}

.page-header {
  margin-bottom: 28px;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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

.create-btn {
  padding: 10px 20px;
  background: var(--cf-accent);
  color: var(--cf-black);
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  font-family: var(--cf-font);
  cursor: pointer;
  transition: background 0.2s ease;
}

.create-btn:hover {
  background: var(--cf-accent-hover);
}

/* ── Content ── */
.alert-content {
  padding: 24px;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
}

.empty-icon {
  font-size: 40px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 14px;
  color: var(--cf-text-muted);
  margin: 0 0 20px;
}

.empty-action {
  padding: 10px 20px;
  background: var(--cf-accent);
  color: var(--cf-black);
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  font-family: var(--cf-font);
  cursor: pointer;
}

/* ── Alert List ── */
.alerts-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.alert-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 16px 20px;
  background: var(--cf-surface);
  border: 1px solid var(--cf-border);
  border-radius: 10px;
  transition: border-color 0.2s ease;
}

.alert-item:hover {
  border-color: var(--cf-border-hover);
}

.alert-pair {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
}

.pair-from,
.pair-to {
  font-size: 15px;
  font-weight: 600;
  color: var(--cf-text);
}

.pair-arrow {
  color: var(--cf-text-muted);
  font-size: 14px;
}

.alert-target {
  min-width: 120px;
}

.target-label {
  display: block;
  font-size: 10px;
  color: var(--cf-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.target-value {
  font-size: 16px;
  font-weight: 700;
  color: var(--cf-text);
  font-variant-numeric: tabular-nums;
}

.alert-direction {
  min-width: 80px;
}

.direction-tag {
  font-size: 12px;
  font-weight: 500;
  padding: 3px 10px;
  border-radius: 6px;
}

.direction-tag.above {
  color: #4ade80;
  background: rgba(34, 197, 94, 0.12);
}

.direction-tag.below {
  color: #fb7185;
  background: rgba(244, 63, 94, 0.12);
}

.alert-status {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 80px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--cf-amber);
}

.status-dot.triggered {
  background: var(--cf-text-muted);
}

.status-text {
  font-size: 12px;
  color: var(--cf-text-secondary);
}

.delete-btn {
  margin-left: auto;
  padding: 6px 12px;
  background: none;
  border: 1px solid var(--cf-border);
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  font-family: var(--cf-font);
  color: var(--cf-text-muted);
  cursor: pointer;
  transition: color 0.2s ease, border-color 0.2s ease;
}

.delete-btn:hover {
  color: var(--cf-rose);
  border-color: rgba(244, 63, 94, 0.3);
}

/* ── Dialog ── */
.direction-picker {
  display: flex;
  gap: 8px;
}

.dir-btn {
  flex: 1;
  padding: 10px 16px;
  background: var(--cf-surface);
  border: 1px solid var(--cf-border);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  font-family: var(--cf-font);
  color: var(--cf-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.dir-btn:hover {
  border-color: var(--cf-border-hover);
}

.dir-btn.active.above {
  color: #4ade80;
  border-color: rgba(34, 197, 94, 0.3);
  background: rgba(34, 197, 94, 0.08);
}

.dir-btn.active.below {
  color: #fb7185;
  border-color: rgba(244, 63, 94, 0.3);
  background: rgba(244, 63, 94, 0.08);
}

.dialog-footer {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.cancel-btn {
  padding: 8px 20px;
  background: var(--cf-surface-hover);
  border: 1px solid var(--cf-border);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  font-family: var(--cf-font);
  color: var(--cf-text-secondary);
  cursor: pointer;
}

.confirm-btn {
  padding: 8px 20px;
  background: var(--cf-accent);
  color: var(--cf-black);
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  font-family: var(--cf-font);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
}

.confirm-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(0, 0, 0, 0.2);
  border-top-color: var(--cf-black);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Responsive ── */
@media (max-width: 767px) {
  .alert-page {
    padding-top: 24px;
  }

  .header-row {
    flex-direction: column;
    gap: 16px;
  }

  .alert-item {
    flex-wrap: wrap;
    gap: 12px;
  }

  .alert-pair,
  .alert-target,
  .alert-direction,
  .alert-status {
    min-width: auto;
  }

  .delete-btn {
    margin-left: 0;
    width: 100%;
  }

  .direction-picker {
    flex-direction: column;
  }
}
</style>
