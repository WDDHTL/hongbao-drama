<template>
  <el-tooltip v-if="tooltip" :content="tooltip" placement="top" :show-after="400">
    <button :class="buttonClass" :disabled="disabled" @click="$emit('click')">
      <el-icon :size="size">
        <component :is="icon" />
      </el-icon>
    </button>
  </el-tooltip>
  <button v-else :class="buttonClass" :disabled="disabled" @click="$emit('click')">
    <el-icon :size="size">
      <component :is="icon" />
    </el-icon>
  </button>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'

const props = withDefaults(
  defineProps<{
    icon: Component
    tooltip?: string
    variant?: 'default' | 'primary' | 'danger'
    size?: number
    disabled?: boolean
  }>(),
  {
    variant: 'default',
    size: 16,
    disabled: false
  }
)

defineEmits<{
  click: []
}>()

const buttonClass = computed(() => ['action-button', `variant-${props.variant}`, { disabled: props.disabled }])
</script>

<style scoped>
.action-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 1px solid transparent;
  border-radius: 12px;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition:
    background-color var(--transition-fast),
    border-color var(--transition-fast),
    color var(--transition-fast),
    transform var(--transition-fast);
}

.action-button:hover {
  background: var(--accent-light);
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.action-button:focus-visible {
  outline: 2px solid var(--border-focus);
  outline-offset: 2px;
}

.variant-primary:hover {
  color: var(--text-primary);
}

.variant-danger:hover {
  background: var(--error-light);
  border-color: rgba(255, 122, 122, 0.18);
  color: var(--error);
}

.disabled {
  opacity: 0.42;
  cursor: not-allowed;
}

.disabled:hover {
  background: transparent;
  border-color: transparent;
  color: var(--text-muted);
  transform: none;
}
</style>
