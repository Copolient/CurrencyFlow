<template>
  <div class="auth-container">
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      class="auth-form"
      @submit.prevent="login"
    >
      <h2 class="form-title">登录</h2>
      <el-form-item label="用户名" prop="username" label-width="80px">
        <el-input v-model="form.username" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item label="密码" prop="password" label-width="80px">
        <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" native-type="submit" :loading="loading" style="width: 100%">
          登录
        </el-button>
      </el-form-item>
      <div class="form-footer">
        还没有账号？<router-link to="/register">立即注册</router-link>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../store/auth';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';

const formRef = ref<FormInstance>();
const form = ref({ username: '', password: '' });
const loading = ref(false);

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 3, message: '密码至少 3 个字符', trigger: 'blur' },
  ],
};

const login = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  loading.value = true;
  try {
    await authStore.login(form.value.username, form.value.password);
    ElMessage.success('登录成功');
    const redirect = (route.query.redirect as string) || '/news';
    router.push(redirect);
  } catch (err: any) {
    ElMessage.error(err.message || '登录失败');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 120px);
  padding: 20px;
}

.auth-form {
  width: 100%;
  max-width: 400px;
  padding: 32px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.form-title {
  text-align: center;
  margin-bottom: 24px;
  color: #303133;
}

.form-footer {
  text-align: center;
  font-size: 14px;
  color: #909399;
}

.form-footer a {
  color: #409eff;
  text-decoration: none;
}
</style>
