<template>
  <header :class="['page-header', { 'with-border': showBorder }]">
    <div class="header-main">
      <div class="header-start">
        <button v-if="showBack" class="back-button" @click="handleBack">
          <el-icon><ArrowLeft /></el-icon>
          <span>{{ backText }}</span>
        </button>

        <div class="title-block">
          <div class="title-row">
            <div v-if="$slots.icon" class="title-icon">
              <slot name="icon" />
            </div>
            <div>
              <h1 class="header-title">{{ title }}</h1>
              <p v-if="subtitle" class="header-subtitle">{{ subtitle }}</p>
            </div>
          </div>
          <slot name="badge" />
        </div>
      </div>

      <div class="header-actions">
        <slot name="actions" />
      </div>
    </div>

    <div v-if="$slots.extra" class="header-extra">
      <slot name="extra" />
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

const props = withDefaults(
  defineProps<{
    title: string
    subtitle?: string
    showBack?: boolean
    backText?: string
    showBorder?: boolean
  }>(),
  {
    showBack: false,
    backText: '返回',
    showBorder: true
  }
)

const emit = defineEmits<{
  back: []
}>()

const router = useRouter()

const handleBack = () => {
  emit('back')
  router.back()
}
</script>

<style scoped>
.page-header {
  margin-bottom: var(--space-4);
}

.page-header.with-border {
  padding-bottom: var(--space-4);
  border-bottom: 1px solid var(--border-primary);
}

.header-main {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header-start {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 14px;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 38px;
  padding: 0 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-full);
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition:
    border-color var(--transition-fast),
    background-color var(--transition-fast),
    color var(--transition-fast);
}

.back-button:hover {
  background: var(--accent-light);
  border-color: var(--border-secondary);
  color: var(--text-primary);
}

.back-button:focus-visible {
  outline: 2px solid var(--border-focus);
  outline-offset: 2px;
}

.title-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 14px;
  border: 1px solid var(--border-primary);
  background: var(--accent-light);
  color: var(--text-primary);
}

.header-title {
  margin: 0;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1.6rem;
  font-weight: 700;
  letter-spacing: -0.05em;
  color: var(--text-primary);
}

.header-subtitle {
  margin: 4px 0 0;
  color: var(--text-muted);
  line-height: 1.6;
}

.header-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.header-extra {
  margin-top: var(--space-4);
}

@media (min-width: 900px) {
  .header-main {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
