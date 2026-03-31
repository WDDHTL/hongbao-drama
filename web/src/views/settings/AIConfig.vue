<template>
  <div class="page-container">
    <div class="content-wrapper animate-fade-in ai-config-page">
      <AppHeader :fixed="false" :show-ai-config="false">
        <template #left>
          <div class="page-title">
            <span class="page-kicker">AI Control</span>
            <h1>模型与中转配置中心</h1>
            <p>统一维护直连模型、OpenClaw 和自定义中转站，把文本、图片、视频的路由放到同一套组件里。</p>
          </div>
        </template>
        <template #right>
          <el-button type="primary" @click="dialogVisible = true">
            <el-icon><Setting /></el-icon>
            <span>打开配置面板</span>
          </el-button>
        </template>
      </AppHeader>

      <section class="hero-grid">
        <BaseCard
          title="黑白系 AI 控制台"
          subtitle="保留项目现有暗色基调，但把信息结构、卡片层次和组件语法统一下来。"
        >
          <div class="hero-copy">
            <p>现在可以在一个面板里管理直连服务、OpenClaw、以及自定义中转站的 Base URL、生成接口和查询接口。</p>
            <div class="hero-tags">
              <el-tag effect="plain">OpenClaw</el-tag>
              <el-tag effect="plain">Custom Relay</el-tag>
              <el-tag effect="plain">Endpoint Mapping</el-tag>
              <el-tag effect="plain">Video Polling</el-tag>
            </div>
          </div>
        </BaseCard>

        <BaseCard title="组件目标" subtitle="后续页面都能直接复用这一套信息密度和交互模式。">
          <div class="check-list">
            <div class="check-item">
              <span class="status-dot success" />
              <span>统一卡片、按钮、页头和指标组件</span>
            </div>
            <div class="check-item">
              <span class="status-dot success" />
              <span>减少花哨强调色，回到纯黑白灰层级</span>
            </div>
            <div class="check-item">
              <span class="status-dot success" />
              <span>保留 OpenClaw 和自定义中转能力</span>
            </div>
          </div>
        </BaseCard>
      </section>
    </div>

    <AIConfigDialog v-model="dialogVisible" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Setting } from '@element-plus/icons-vue'
import { AppHeader, BaseCard } from '@/components/common'
import AIConfigDialog from '@/components/common/AIConfigDialog.vue'

const dialogVisible = ref(true)
</script>

<style scoped>
.ai-config-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.page-title {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.page-title h1 {
  margin: 0;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: clamp(1.8rem, 3vw, 2.6rem);
  letter-spacing: -0.06em;
}

.page-title p {
  max-width: 760px;
  margin: 0;
  color: var(--text-muted);
  line-height: 1.7;
}

.page-kicker {
  color: var(--text-muted);
  font-size: 0.76rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.hero-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(320px, 0.8fr);
  gap: 18px;
}

.hero-copy {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hero-copy p {
  color: var(--text-secondary);
  line-height: 1.75;
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.check-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.check-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-secondary);
  line-height: 1.6;
}

@media (max-width: 900px) {
  .hero-grid {
    grid-template-columns: 1fr;
  }
}
</style>
