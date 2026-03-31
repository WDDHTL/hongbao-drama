<template>
  <section
    :class="[
      'base-card',
      `variant-${variant}`,
      { 'is-hoverable': hoverable, 'is-clickable': clickable, 'has-header': $slots.header || title }
    ]"
    :tabindex="clickable ? 0 : undefined"
    @click="clickable ? $emit('click') : undefined"
    @keydown.enter="clickable ? $emit('click') : undefined"
  >
    <header v-if="$slots.header || title" class="card-header">
      <slot name="header">
        <div class="header-main">
          <div v-if="icon" class="header-icon">
            <el-icon :size="iconSize">
              <component :is="icon" />
            </el-icon>
          </div>
          <div class="header-copy">
            <h3 class="card-title">{{ title }}</h3>
            <p v-if="subtitle" class="card-subtitle">{{ subtitle }}</p>
          </div>
        </div>
        <div v-if="$slots.headerActions" class="header-actions">
          <slot name="headerActions" />
        </div>
      </slot>
    </header>

    <div :class="['card-body', { 'no-padding': noPadding }]">
      <slot />
    </div>

    <footer v-if="$slots.footer" class="card-footer">
      <slot name="footer" />
    </footer>
  </section>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

withDefaults(
  defineProps<{
    title?: string
    subtitle?: string
    icon?: Component
    iconSize?: number
    variant?: 'default' | 'elevated' | 'outlined' | 'ghost'
    hoverable?: boolean
    clickable?: boolean
    noPadding?: boolean
  }>(),
  {
    iconSize: 18,
    variant: 'default',
    hoverable: false,
    clickable: false,
    noPadding: false
  }
)

defineEmits<{
  click: []
}>()
</script>

<style scoped>
.base-card {
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-primary);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.028), transparent 30%),
    var(--bg-card);
  box-shadow: var(--shadow-card);
  transition:
    transform var(--transition-fast),
    border-color var(--transition-fast),
    box-shadow var(--transition-fast),
    background-color var(--transition-fast);
}

.base-card::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.035), transparent 22%);
  opacity: 0.7;
}

.variant-elevated {
  border-color: transparent;
  box-shadow: var(--shadow-lg);
}

.variant-outlined {
  background: transparent;
  box-shadow: none;
}

.variant-ghost {
  border-color: transparent;
  background: transparent;
  box-shadow: none;
}

.is-hoverable:hover,
.is-clickable:hover {
  transform: translateY(-2px);
  border-color: var(--border-secondary);
  box-shadow: var(--shadow-card-hover);
}

.is-clickable {
  cursor: pointer;
}

.is-clickable:focus-visible {
  outline: 2px solid var(--border-focus);
  outline-offset: 2px;
}

.card-header,
.card-body,
.card-footer {
  position: relative;
  z-index: 1;
}

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 20px 22px 0;
}

.header-main {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  min-width: 0;
}

.header-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  border-radius: 12px;
  border: 1px solid var(--border-primary);
  background: var(--accent-light);
  color: var(--text-primary);
}

.header-copy {
  min-width: 0;
}

.card-title {
  margin: 0;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1rem;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.card-subtitle {
  margin: 6px 0 0;
  color: var(--text-muted);
  line-height: 1.6;
  font-size: 0.84rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-body {
  flex: 1;
  padding: 22px;
}

.has-header .card-body {
  padding-top: 18px;
}

.card-body.no-padding {
  padding: 0;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 22px 20px;
  border-top: 1px solid var(--border-primary);
}
</style>
