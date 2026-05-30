<template>
  <div class="litestar-badge">
    <a-tooltip title="基于 LiteStar 框架构建">
      <div class="framework-info">
        <span class="framework-icon">🚀</span>
        <span class="framework-name">LiteStar</span>
        <span class="framework-version">{{ version }}</span>
      </div>
    </a-tooltip>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// 从全局配置获取版本信息
const version = ref(window.KNIFE4J_FRAMEWORK_VERSION || '2.17.0')
</script>

<style scoped>
.litestar-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  font-size: 12px;
  margin-left: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.litestar-badge:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.framework-info {
  display: flex;
  align-items: center;
  gap: 4px;
}

.framework-icon {
  font-size: 14px;
}

.framework-name {
  font-weight: 600;
  letter-spacing: 0.5px;
}

.framework-version {
  opacity: 0.8;
  font-size: 10px;
  background: rgba(255, 255, 255, 0.2);
  padding: 1px 4px;
  border-radius: 4px;
}
</style>
