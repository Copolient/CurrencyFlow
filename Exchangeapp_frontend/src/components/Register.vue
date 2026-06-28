<template>
  <div class="auth-page">
    <!-- Left: Brand -->
    <div class="auth-brand">
      <div class="brand-content cf-animate-in">
        <div class="brand-logo">
          <span class="logo-mark">◆</span>
          <span class="logo-text">CurrencyFlow</span>
        </div>
        <h1 class="brand-headline">加入社区，<br />洞察市场。</h1>
        <p class="brand-sub">实时汇率 · AI 分析 · 交易社区</p>
      </div>
      <div class="brand-glow"></div>
    </div>

    <!-- Right: Form -->
    <div class="auth-form-side">
      <div class="auth-form-wrap cf-animate-in cf-delay-1">
        <h2 class="form-title">创建账户</h2>
        <p class="form-subtitle">注册你的 CurrencyFlow 账户</p>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          class="auth-form"
          @submit.prevent="register"
          label-position="top"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="form.username"
              placeholder="2-20 个字符"
              size="large"
              clearable
            />
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="至少 6 个字符"
              size="large"
              show-password
            />
          </el-form-item>

          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="再次输入密码"
              size="large"
              show-password
            />
          </el-form-item>

          <el-form-item>
            <button
              type="submit"
              class="submit-btn"
              :disabled="loading"
            >
              <span v-if="loading" class="btn-spinner"></span>
              <span v-else>注册</span>
            </button>
          </el-form-item>
        </el-form>

        <div class="form-footer">
          已有账号？<router-link to="/login" class="footer-link">立即登录</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';

const formRef = ref<FormInstance>();
const form = ref({ username: '', password: '', confirmPassword: '' });
const loading = ref(false);

const authStore = useAuthStore();
const router = useRouter();

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名 2-20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: Function) => {
        if (value !== form.value.password) {
          callback(new Error('两次密码不一致'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    },
  ],
};

const register = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  loading.value = true;
  try {
    await authStore.register(form.value.username, form.value.password);
    ElMessage.success('注册成功');
    router.push('/news');
  } catch (err: any) {
    ElMessage.error(err.message || '注册失败');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.auth-page {
  display: flex;
  min-height: calc(100vh - 56px);
  margin: 0 -24px;
}

.auth-brand {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--cf-surface);
  border-right: 1px solid var(--cf-border);
}

.brand-content {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 40px;
}

.brand-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 40px;
}

.logo-mark {
  font-size: 24px;
  color: var(--cf-accent);
}

.logo-text {
  font-size: 22px;
  font-weight: 700;
  color: var(--cf-text);
  letter-spacing: -0.02em;
}

.brand-headline {
  font-size: 36px;
  font-weight: 800;
  line-height: 1.2;
  color: var(--cf-text);
  margin: 0 0 16px;
  letter-spacing: -0.03em;
}

.brand-sub {
  font-size: 14px;
  color: var(--cf-text-muted);
  margin: 0;
}

.brand-glow {
  position: absolute;
  bottom: -100px;
  left: 50%;
  transform: translateX(-50%);
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.12) 0%, transparent 70%);
  pointer-events: none;
}

.auth-form-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.auth-form-wrap {
  width: 100%;
  max-width: 380px;
}

.form-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--cf-text);
  margin: 0 0 4px;
}

.form-subtitle {
  font-size: 14px;
  color: var(--cf-text-muted);
  margin: 0 0 32px;
}

.auth-form {
  margin-bottom: 24px;
}

.submit-btn {
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

.submit-btn:hover:not(:disabled) {
  background: var(--cf-accent-hover);
}

.submit-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-spinner {
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

.form-footer {
  text-align: center;
  font-size: 14px;
  color: var(--cf-text-muted);
}

.footer-link {
  color: var(--cf-accent);
  font-weight: 500;
}

@media (max-width: 767px) {
  .auth-page {
    flex-direction: column;
    min-height: auto;
  }

  .auth-brand {
    border-right: none;
    border-bottom: 1px solid var(--cf-border);
    padding: 40px 24px;
  }

  .brand-headline {
    font-size: 24px;
  }

  .auth-form-side {
    padding: 32px 24px;
  }
}
</style>
