<template>
  <div class="home-container">
    <div class="hero">
      <h1 class="title">蓝鼠兑换</h1>
      <p class="subtitle">实时汇率查询 · 财经资讯 · 一站掌握</p>
      <div class="actions">
        <el-button type="primary" size="large" @click="router.push('/exchange')">
          开始兑换
        </el-button>
        <el-button size="large" @click="router.push('/news')">
          浏览资讯
        </el-button>
      </div>
    </div>

    <div class="features">
      <el-row :gutter="24">
        <el-col :span="8">
          <el-card shadow="hover" class="feature-card">
            <div class="feature-icon">💱</div>
            <h3>实时汇率</h3>
            <p>支持多币种实时查询与兑换计算</p>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover" class="feature-card">
            <div class="feature-icon">📰</div>
            <h3>财经资讯</h3>
            <p>精选财经文章，洞察市场动态</p>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover" class="feature-card">
            <div class="feature-icon">🔒</div>
            <h3>安全可靠</h3>
            <p>JWT 鉴权 + HTTPS，数据安全有保障</p>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();
</script>

<style scoped>
.home-container {
  max-width: 960px;
  margin: 0 auto;
  padding: 60px 20px;
}

.hero {
  text-align: center;
  margin-bottom: 60px;
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
