<template>
  <div v-if="error" class="error-boundary">
    <el-result icon="error" title="页面出错了" :sub-title="error.message">
      <template #extra>
        <el-button type="primary" @click="reset">重试</el-button>
      </template>
    </el-result>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue';

const error = ref<Error | null>(null);

onErrorCaptured((err: Error) => {
  error.value = err;
  return false; // 阻止错误继续向上传播
});

const reset = () => {
  error.value = null;
};
</script>

<style scoped>
.error-boundary {
  padding: 40px 20px;
}
</style>
