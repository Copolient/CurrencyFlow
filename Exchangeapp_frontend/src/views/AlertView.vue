<template>
  <div class="alert-container">
    <el-card class="alert-card">
      <template #header>
        <div class="card-header">
          <span>🔔 汇率预警</span>
          <el-button type="primary" @click="showCreateDialog = true">创建预警</el-button>
        </div>
      </template>

      <div v-loading="loading">
        <div v-if="alerts.length === 0" class="empty-state">
          <el-empty description="还没有设置预警">
            <el-button type="primary" @click="showCreateDialog = true">创建第一个预警</el-button>
          </el-empty>
        </div>

        <el-table v-else :data="alerts" style="width: 100%">
          <el-table-column prop="fromCurrency" label="源货币" width="100" />
          <el-table-column prop="toCurrency" label="目标货币" width="100" />
          <el-table-column label="目标汇率" width="120">
            <template #default="{ row }">
              {{ row.targetRate.toFixed(4) }}
            </template>
          </el-table-column>
          <el-table-column label="方向" width="100">
            <template #default="{ row }">
              <el-tag :type="row.direction === 'above' ? 'success' : 'danger'">
                {{ row.direction === 'above' ? '高于' : '低于' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.triggered ? 'info' : 'warning'">
                {{ row.triggered ? '已触发' : '监控中' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button type="danger" text @click="deleteAlert(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- Create Alert Dialog -->
    <el-dialog v-model="showCreateDialog" title="创建汇率预警" width="400px">
      <el-form :model="newAlert" label-width="80px">
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
          <el-radio-group v-model="newAlert.direction">
            <el-radio label="above">高于目标</el-radio>
            <el-radio label="below">低于目标</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createAlert" :loading="creating">创建</el-button>
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
  } catch {
    // interceptor handles error
  } finally {
    loading.value = false;
  }
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
  } catch {
    // interceptor handles error
  } finally {
    creating.value = false;
  }
};

const deleteAlert = async (id: number) => {
  try {
    await axios.delete(`/alerts/${id}`);
    ElMessage.success('预警已删除');
    fetchAlerts();
  } catch {
    // interceptor handles error
  }
};

onMounted(() => {
  fetchAlerts();
  fetchCurrencies();
});
</script>

<style scoped>
.alert-container {
  max-width: 960px;
  margin: 20px auto;
  padding: 0 20px;
}

.alert-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-state {
  padding: 40px 0;
}
</style>
