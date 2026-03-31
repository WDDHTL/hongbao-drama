<template>
  <div v-loading="loading" class="config-list">
    <el-empty v-if="!loading && configs.length === 0" :description="$t('aiConfig.empty')" />

    <BaseCard v-for="config in configs" :key="config.id" hoverable class="config-card">
      <template #header>
        <div class="config-header">
          <div class="config-header-main">
            <div class="config-tags">
              <el-tag size="small" effect="plain">{{ config.provider || 'relay' }}</el-tag>
              <el-tag v-if="config.is_active" size="small" type="success">{{ $t('aiConfig.enabled') }}</el-tag>
              <el-tag v-else size="small" type="info">{{ $t('aiConfig.disabled') }}</el-tag>
            </div>
            <div>
              <h3 class="config-name">{{ config.name }}</h3>
              <p class="config-url">{{ config.base_url }}</p>
            </div>
          </div>

          <div class="config-actions">
            <el-button v-if="showTestButton" text @click="$emit('test', config)" :icon="Connection">
              {{ $t('aiConfig.actions.test') }}
            </el-button>
            <el-button text @click="$emit('edit', config)" :icon="Edit">{{ $t('common.edit') }}</el-button>
            <el-button text :type="config.is_active ? 'warning' : 'success'" @click="$emit('toggleActive', config)">
              {{ config.is_active ? $t('aiConfig.disable') : $t('aiConfig.enable') }}
            </el-button>
            <el-popconfirm :title="$t('aiConfig.messages.deleteConfirm')" @confirm="$emit('delete', config)">
              <template #reference>
                <el-button text type="danger" :icon="Delete">{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>
      </template>

      <div class="metric-grid">
        <div class="metric-card">
          <label>Endpoint</label>
          <span>{{ config.endpoint || '/chat/completions' }}</span>
        </div>

        <div class="metric-card">
          <label>Query Endpoint</label>
          <span>{{ config.query_endpoint || 'Not required' }}</span>
        </div>

        <div class="metric-card">
          <label>Priority</label>
          <el-tag size="small" :type="priorityType(config.priority || 0)">
            {{ config.priority || 0 }}
          </el-tag>
        </div>

        <div class="metric-card">
          <label>API Key</label>
          <span class="font-mono">{{ maskApiKey(config.api_key) }}</span>
        </div>

        <div class="metric-card models-card">
          <label>Models</label>
          <div class="model-list">
            <el-tag v-for="model in config.model" :key="model" size="small" effect="plain">
              {{ model }}
            </el-tag>
          </div>
        </div>

        <div class="metric-card">
          <label>Created</label>
          <span>{{ formatDate(config.created_at) }}</span>
        </div>
      </div>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { Connection, Delete, Edit } from '@element-plus/icons-vue'
import { BaseCard } from '@/components/common'
import type { AIServiceConfig } from '@/types/ai'

defineProps<{
  configs: AIServiceConfig[]
  loading: boolean
  showTestButton?: boolean
}>()

defineEmits<{
  edit: [config: AIServiceConfig]
  delete: [config: AIServiceConfig]
  toggleActive: [config: AIServiceConfig]
  test: [config: AIServiceConfig]
}>()

const maskApiKey = (key: string) => {
  if (!key) return ''
  if (key.length <= 8) return '***'
  return `${key.slice(0, 4)}***${key.slice(-4)}`
}

const formatDate = (value: string) => new Date(value).toLocaleString('zh-CN')

const priorityType = (priority: number) => {
  if (priority >= 50) return 'danger'
  if (priority >= 20) return 'warning'
  return 'info'
}
</script>

<style scoped>
.config-list {
  display: grid;
  gap: 14px;
  min-height: 320px;
}

.config-card {
  overflow: hidden;
}

.config-header {
  display: flex;
  justify-content: space-between;
  gap: 18px;
  width: 100%;
}

.config-header-main {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.config-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.config-name {
  margin: 0;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1.05rem;
  letter-spacing: -0.03em;
}

.config-url {
  margin: 6px 0 0;
  color: var(--text-muted);
  line-height: 1.6;
}

.config-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 8px;
  justify-content: flex-end;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.metric-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  background: rgba(255, 255, 255, 0.02);
}

.metric-card label {
  color: var(--text-muted);
  font-size: 0.74rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.metric-card span {
  color: var(--text-secondary);
  word-break: break-word;
  line-height: 1.6;
}

.models-card {
  grid-column: 1 / -1;
}

.model-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media (max-width: 900px) {
  .config-header {
    flex-direction: column;
  }

  .config-actions {
    justify-content: flex-start;
  }

  .metric-grid {
    grid-template-columns: 1fr;
  }
}
</style>
