<template>
  <div class="app-header-wrapper">
    <header class="app-header" :class="{ 'header-fixed': fixed }">
      <div class="header-shell">
        <div class="header-left">
          <router-link v-if="showLogo" to="/" class="brand-link">
            <span class="brand-mark" />
            <span class="brand-copy">
              <span class="brand-title font-display">HuoBao Drama</span>
              <span class="brand-subtitle">Monochrome Workflow System</span>
            </span>
          </router-link>
          <slot name="left" />
        </div>

        <div class="header-center">
          <slot name="center" />
        </div>

        <div class="header-right">
          <slot name="right" />
          <LanguageSwitcher v-if="showLanguage" />
          <ThemeToggle v-if="showTheme" />
          <el-button v-if="showAIConfig" class="config-button" @click="handleOpenAIConfig">
            <el-icon><Setting /></el-icon>
            <span class="button-text">{{ $t('drama.aiConfig') }}</span>
          </el-button>
        </div>
      </div>
    </header>

    <AIConfigDialog v-model="showConfigDialog" @config-updated="emit('config-updated')" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Setting } from '@element-plus/icons-vue'
import ThemeToggle from './ThemeToggle.vue'
import AIConfigDialog from './AIConfigDialog.vue'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'

interface Props {
  fixed?: boolean
  showLogo?: boolean
  showLanguage?: boolean
  showTheme?: boolean
  showAIConfig?: boolean
}

withDefaults(defineProps<Props>(), {
  fixed: true,
  showLogo: true,
  showLanguage: true,
  showTheme: true,
  showAIConfig: true
})

const emit = defineEmits<{
  (e: 'open-ai-config'): void
  (e: 'config-updated'): void
}>()

const showConfigDialog = ref(false)

const handleOpenAIConfig = () => {
  showConfigDialog.value = true
  emit('open-ai-config')
}

defineExpose({
  openAIConfig: () => {
    showConfigDialog.value = true
  }
})
</script>

<style scoped>
.app-header {
  z-index: 1000;
  background: transparent;
}

.header-fixed {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
}

.header-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 0.8fr) auto;
  align-items: center;
  gap: 16px;
  width: calc(100% - 24px);
  min-height: 74px;
  margin: 16px auto 0;
  padding: 14px 18px;
  border: 1px solid var(--border-primary);
  border-radius: 22px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.03), transparent 28%),
    var(--surface-overlay);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(18px);
}

.header-left,
.header-center,
.header-right {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.header-center {
  justify-content: center;
}

.header-right {
  justify-content: flex-end;
}

.brand-link {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  color: inherit;
  text-decoration: none;
}

.brand-mark {
  width: 38px;
  height: 38px;
  flex-shrink: 0;
  border-radius: 14px;
  border: 1px solid var(--border-primary);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.14), rgba(255, 255, 255, 0.02)),
    var(--bg-card);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.brand-copy {
  display: flex;
  min-width: 0;
  flex-direction: column;
}

.brand-title {
  font-size: 1.04rem;
  line-height: 1.1;
  color: var(--text-primary);
}

.brand-subtitle {
  color: var(--text-muted);
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.config-button {
  border-radius: 999px;
}

@media (max-width: 1024px) {
  .header-shell {
    grid-template-columns: 1fr;
  }

  .header-center,
  .header-right {
    justify-content: flex-start;
  }

  .header-right {
    flex-wrap: wrap;
  }
}

@media (max-width: 640px) {
  .header-shell {
    width: calc(100% - 16px);
    padding: 12px 14px;
    border-radius: 18px;
  }

  .button-text,
  .brand-subtitle {
    display: none;
  }
}
</style>
