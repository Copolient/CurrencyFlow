<template>
  <div class="auth-container">
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      class="auth-form"
      @submit.prevent="register"
    >
      <h2 class="form-title">注册</h2>
      <el-form-item label="用户名" prop="username" label-width="80px">
        <el-input v-model="form.username" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item label="密码" prop="password" label-width="80px">
        <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword" label-width="80px">
        <el-input v-model="form.confirmPassword" type="password" placeholder="再次输入密码" show-password />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" native-type="submit" :loading="loading" style="width: 100%">
          注册
        </el-button>
      </el-form-item>
      <div class="form-footer">
        已有账号？<router-link to="/login">立即登录</router-link>
      </div>
    </el-form>
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
