<template>
  <div :class="['stat-card', `variant-${variant}`]">
    <div v-if="icon" class="stat-icon" :style="{ background: iconBg, color: iconColor }">
      <el-icon :size="20">
        <component :is="icon" />
      </el-icon>
    </div>

    <div class="stat-content">
      <span class="stat-label">{{ label }}</span>
      <div class="stat-value-row">
        <span class="stat-value" :style="{ color: valueColor }">{{ formattedValue }}</span>
        <span v-if="suffix" class="stat-suffix">{{ suffix }}</span>
      </div>
      <span v-if="description" class="stat-description">{{ description }}</span>
    </div>

    <div v-if="trend !== undefined" :class="['stat-trend', trendDirection]">
      <el-icon :size="14">
        <component :is="trendIcon" />
      </el-icon>
      <span>{{ Math.abs(trend) }}%</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import { ArrowDown, ArrowUp } from '@element-plus/icons-vue'

const props = withDefaults(
  defineProps<{
    label: string
    value: number | string
    icon?: Component
    iconColor?: string
    iconBg?: string
    valueColor?: string
    suffix?: string
    description?: string
    trend?: number
    variant?: 'default' | 'compact'
  }>(),
  {
    iconColor: 'var(--text-primary)',
    iconBg: 'var(--accent-light)',
    valueColor: 'var(--text-primary)',
    variant: 'default'
  }
)

const formattedValue = computed(() => {
  if (typeof props.value === 'string') return props.value
  if (props.value >= 1000000) return `${(props.value / 1000000).toFixed(1)}M`
  if (props.value >= 1000) return `${(props.value / 1000).toFixed(1)}K`
  return props.value.toString()
})

const trendDirection = computed(() => {
  if (props.trend === undefined) return ''
  return props.trend >= 0 ? 'up' : 'down'
})

const trendIcon = computed(() => (props.trend !== undefined && props.trend < 0 ? ArrowDown : ArrowUp))
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 16px 18px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-primary);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.025), transparent 36%),
    var(--bg-card);
  transition:
    border-color var(--transition-fast),
    box-shadow var(--transition-fast),
    transform var(--transition-fast);
}

.stat-card:hover {
  border-color: var(--border-secondary);
  box-shadow: var(--shadow-sm);
  transform: translateY(-1px);
}

.variant-compact {
  padding: 14px;
  gap: 12px;
}

.stat-icon {
  display: inline-flex;
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  border: 1px solid var(--border-primary);
}

.stat-content {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  color: var(--text-muted);
  font-size: 0.78rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.stat-value-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.stat-value {
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1.9rem;
  font-weight: 700;
  letter-spacing: -0.05em;
  line-height: 1;
}

.variant-compact .stat-value {
  font-size: 1.5rem;
}

.stat-suffix {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.stat-description {
  color: var(--text-muted);
  font-size: 0.78rem;
  line-height: 1.5;
}

.stat-trend {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 700;
}

.stat-trend.up {
  background: var(--success-light);
  color: var(--success);
}

.stat-trend.down {
  background: var(--error-light);
  color: var(--error);
}
</style>
