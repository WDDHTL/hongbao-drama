<template>
  <el-dialog
    v-model="visible"
    width="1080px"
    :close-on-click-modal="false"
    destroy-on-close
    class="ai-config-dialog"
  >
    <template #header>
      <div class="dialog-header">
        <div>
          <p class="eyebrow">AI CONTROL CENTER</p>
          <h2 class="dialog-title">{{ $t('aiConfig.title') }}</h2>
          <p class="dialog-subtitle">在一个面板里管理直连模型、OpenClaw 和自定义中转站的完整路由配置。</p>
        </div>
        <div class="header-actions">
          <el-button class="soft-btn" @click="showQuickSetupDialog">
            <el-icon><MagicStick /></el-icon>
            <span>Quick Setup</span>
          </el-button>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            <span>{{ $t('aiConfig.addConfig') }}</span>
          </el-button>
        </div>
      </div>
    </template>

    <div class="overview-grid">
      <div class="overview-card">
        <span class="overview-label">Text</span>
        <strong>{{ countByType('text') }}</strong>
      </div>
      <div class="overview-card">
        <span class="overview-label">Image</span>
        <strong>{{ countByType('image') }}</strong>
      </div>
      <div class="overview-card">
        <span class="overview-label">Video</span>
        <strong>{{ countByType('video') }}</strong>
      </div>
      <div class="overview-card overview-accent">
        <span class="overview-label">Relay</span>
        <strong>{{ relayCount }}</strong>
      </div>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="config-tabs">
      <el-tab-pane :label="$t('aiConfig.tabs.text')" name="text">
        <ConfigList
          :configs="configs"
          :loading="loading"
          :show-test-button="true"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
          @test="handleTest"
        />
      </el-tab-pane>

      <el-tab-pane :label="$t('aiConfig.tabs.image')" name="image">
        <ConfigList
          :configs="configs"
          :loading="loading"
          :show-test-button="false"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
        />
      </el-tab-pane>

      <el-tab-pane :label="$t('aiConfig.tabs.video')" name="video">
        <ConfigList
          :configs="configs"
          :loading="loading"
          :show-test-button="false"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
        />
      </el-tab-pane>
    </el-tabs>

    <el-dialog
      v-model="quickSetupVisible"
      title="Quick Setup"
      width="520px"
      :close-on-click-modal="false"
      append-to-body
    >
      <div class="quick-setup-info">
        <p>This creates one ChatFire configuration for text, image, and video.</p>
        <ul>
          <li>Text: {{ getPresetModel('text', 'chatfire') }}</li>
          <li>Image: {{ getPresetModel('image', 'chatfire') }}</li>
          <li>Video: {{ getPresetModel('video', 'chatfire') }}</li>
        </ul>
        <p class="quick-setup-tip">Base URL: https://api.chatfire.site/v1</p>
      </div>
      <el-form label-width="80px">
        <el-form-item label="API Key" required>
          <el-input
            v-model="quickSetupApiKey"
            type="password"
            show-password
            placeholder="Enter ChatFire API key"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="quick-setup-footer">
          <a href="https://api.chatfire.site/login?inviteCode=C4453345" target="_blank" class="register-link">
            Need a ChatFire key?
          </a>
          <div class="footer-buttons">
            <el-button @click="quickSetupVisible = false">Cancel</el-button>
            <el-button type="primary" @click="handleQuickSetup" :loading="quickSetupLoading">
              Create presets
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      v-model="editDialogVisible"
      :title="isEdit ? $t('aiConfig.editConfig') : $t('aiConfig.addConfig')"
      width="720px"
      :close-on-click-modal="false"
      append-to-body
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="116px" class="config-form">
        <div class="form-grid">
          <el-form-item :label="$t('aiConfig.form.name')" prop="name">
            <el-input v-model="form.name" :placeholder="$t('aiConfig.form.namePlaceholder')" />
          </el-form-item>

          <el-form-item :label="$t('aiConfig.form.priority')" prop="priority">
            <el-input-number v-model="form.priority" :min="0" :max="100" :step="1" style="width: 100%" />
          </el-form-item>
        </div>

        <el-form-item :label="$t('aiConfig.form.provider')" prop="provider">
          <el-select v-model="form.provider" @change="handleProviderChange" style="width: 100%">
            <el-option
              v-for="provider in availableProviders"
              :key="provider.id"
              :label="provider.name"
              :value="provider.id"
            />
          </el-select>
          <div class="form-tip">{{ providerDescription }}</div>
        </el-form-item>

        <el-form-item :label="$t('aiConfig.form.model')" prop="model">
          <el-select
            v-model="form.model"
            multiple
            filterable
            allow-create
            default-first-option
            collapse-tags
            collapse-tags-tooltip
            style="width: 100%"
          >
            <el-option v-for="model in availableModels" :key="model" :label="model" :value="model" />
          </el-select>
          <div class="form-tip">Leave the preset list or type your own model name for OpenClaw / custom relay.</div>
        </el-form-item>

        <div class="form-grid">
          <el-form-item :label="$t('aiConfig.form.baseUrl')" prop="base_url">
            <el-input v-model="form.base_url" placeholder="https://your-relay.example/v1" />
          </el-form-item>

          <el-form-item :label="$t('aiConfig.form.apiKey')" prop="api_key">
            <el-input v-model="form.api_key" type="password" show-password placeholder="sk-..." />
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="Endpoint" prop="endpoint">
            <el-input v-model="form.endpoint" placeholder="/chat/completions" />
          </el-form-item>

          <el-form-item label="Query Endpoint" prop="query_endpoint">
            <el-input
              v-model="form.query_endpoint"
              :placeholder="form.service_type === 'video' ? '/videos/{taskId}' : 'Optional'"
            />
          </el-form-item>
        </div>

        <el-alert
          type="info"
          :closable="false"
          show-icon
          class="endpoint-alert"
          :title="`Resolved request path: ${fullEndpointExample}`"
        />

        <el-form-item v-if="isEdit" :label="$t('aiConfig.form.isActive')">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="quick-setup-footer">
          <span class="register-link">{{ providerDescription }}</span>
          <div class="footer-buttons">
            <el-button @click="editDialogVisible = false">{{ $t('common.cancel') }}</el-button>
            <el-button v-if="form.service_type === 'text'" @click="testConnection" :loading="testing">
              {{ $t('aiConfig.actions.test') }}
            </el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitting">
              {{ isEdit ? $t('common.save') : $t('common.create') }}
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { MagicStick, Plus } from '@element-plus/icons-vue'
import { aiAPI } from '@/api/ai'
import ConfigList from '@/views/settings/components/ConfigList.vue'
import type {
  AIServiceConfig,
  AIServiceType,
  CreateAIConfigRequest,
  UpdateAIConfigRequest
} from '@/types/ai'
import {
  buildConfigName,
  getProviderBaseUrl,
  getProviderConfig,
  getProviderConfigs,
  getProviderEndpoint,
  getProviderModels,
  getProviderQueryEndpoint
} from '@/utils/aiProviders'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'config-updated': []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const activeTab = ref<AIServiceType>('text')
const loading = ref(false)
const configs = ref<AIServiceConfig[]>([])
const editDialogVisible = ref(false)
const quickSetupVisible = ref(false)
const quickSetupApiKey = ref('')
const quickSetupLoading = ref(false)
const isEdit = ref(false)
const editingId = ref<number>()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const testing = ref(false)

const form = reactive<
  CreateAIConfigRequest & { is_active?: boolean; provider?: string; endpoint?: string; query_endpoint?: string }
>({
  service_type: 'text',
  provider: 'chatfire',
  name: '',
  base_url: '',
  api_key: '',
  model: [],
  endpoint: '',
  query_endpoint: '',
  priority: 0,
  is_active: true
})

const rules: FormRules = {
  name: [{ required: true, message: 'Please enter a config name', trigger: 'blur' }],
  provider: [{ required: true, message: 'Please choose a provider', trigger: 'change' }],
  base_url: [{ required: true, message: 'Please enter a Base URL', trigger: 'blur' }],
  api_key: [{ required: true, message: 'Please enter an API key', trigger: 'blur' }],
  model: [
    {
      validator: (_rule, value, callback) => {
        if (Array.isArray(value) && value.length > 0) {
          callback()
          return
        }
        callback(new Error('Please select at least one model'))
      },
      trigger: 'change'
    }
  ]
}

const allConfigs = ref<AIServiceConfig[]>([])

const availableProviders = computed(() => getProviderConfigs(form.service_type))

const availableModels = computed(() => {
  const presetModels = getProviderModels(form.service_type, form.provider)
  const configuredModels = allConfigs.value
    .filter((config) => config.service_type === form.service_type && config.provider === form.provider)
    .flatMap((config) => config.model || [])

  return Array.from(new Set([...presetModels, ...configuredModels]))
})

const providerDescription = computed(
  () => getProviderConfig(form.service_type, form.provider)?.description || 'OpenAI-compatible relay.'
)

const fullEndpointExample = computed(() => {
  const baseUrl = form.base_url || 'https://relay.example/v1'
  const endpoint = form.endpoint || getProviderEndpoint(form.service_type, form.provider)
  return `${baseUrl}${endpoint || ''}`
})

const relayCount = computed(
  () =>
    allConfigs.value.filter((config) => ['custom_relay', 'openclaw'].includes(config.provider || '') && config.is_active)
      .length
)

const countByType = (serviceType: AIServiceType) =>
  allConfigs.value.filter((config) => config.service_type === serviceType && config.is_active).length

const getPresetModel = (serviceType: AIServiceType, provider: string) =>
  getProviderModels(serviceType, provider)[0] || 'custom'

const applyProviderPreset = (serviceType: AIServiceType, providerId?: string) => {
  form.base_url = getProviderBaseUrl(serviceType, providerId)
  form.endpoint = getProviderEndpoint(serviceType, providerId)
  form.query_endpoint = getProviderQueryEndpoint(serviceType, providerId)
  const models = getProviderModels(serviceType, providerId)
  if (!isEdit.value) {
    form.model = models.length > 0 ? [models[0]] : []
    if (providerId) {
      form.name = buildConfigName(providerId, serviceType)
    }
  }
}

const loadConfigs = async () => {
  loading.value = true
  try {
    const [activeConfigs, textConfigs, imageConfigs, videoConfigs] = await Promise.all([
      aiAPI.list(activeTab.value),
      aiAPI.list('text'),
      aiAPI.list('image'),
      aiAPI.list('video')
    ])
    configs.value = activeConfigs
    allConfigs.value = [...textConfigs, ...imageConfigs, ...videoConfigs]
  } catch (error: any) {
    ElMessage.error(error.message || 'Failed to load AI configs')
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  Object.assign(form, {
    service_type: activeTab.value,
    provider: 'chatfire',
    name: '',
    base_url: '',
    api_key: '',
    model: [],
    endpoint: '',
    query_endpoint: '',
    priority: 0,
    is_active: true
  })
  formRef.value?.clearValidate()
  applyProviderPreset(activeTab.value, 'chatfire')
}

const showCreateDialog = () => {
  isEdit.value = false
  editingId.value = undefined
  resetForm()
  editDialogVisible.value = true
}

const handleEdit = (config: AIServiceConfig) => {
  isEdit.value = true
  editingId.value = config.id
  Object.assign(form, {
    service_type: config.service_type,
    provider: config.provider || 'chatfire',
    name: config.name,
    base_url: config.base_url,
    api_key: config.api_key,
    model: Array.isArray(config.model) ? config.model : [config.model],
    endpoint: config.endpoint || getProviderEndpoint(config.service_type, config.provider),
    query_endpoint: config.query_endpoint || '',
    priority: config.priority || 0,
    is_active: config.is_active
  })
  editDialogVisible.value = true
}

const handleDelete = async (config: AIServiceConfig) => {
  try {
    await ElMessageBox.confirm('Delete this AI config?', 'Warning', {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })

    await aiAPI.delete(config.id)
    ElMessage.success('Deleted')
    await loadConfigs()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || 'Delete failed')
    }
  }
}

const handleToggleActive = async (config: AIServiceConfig) => {
  try {
    await aiAPI.update(config.id, { is_active: !config.is_active })
    ElMessage.success(config.is_active ? 'Config disabled' : 'Config enabled')
    await loadConfigs()
  } catch (error: any) {
    ElMessage.error(error.message || 'Update failed')
  }
}

const testConnection = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  testing.value = true
  try {
    await aiAPI.testConnection({
      base_url: form.base_url,
      api_key: form.api_key,
      model: form.model,
      provider: form.provider,
      endpoint: form.endpoint,
      query_endpoint: form.query_endpoint
    })
    ElMessage.success('Connection test passed')
  } catch (error: any) {
    ElMessage.error(error.message || 'Connection test failed')
  } finally {
    testing.value = false
  }
}

const handleTest = async (config: AIServiceConfig) => {
  testing.value = true
  try {
    await aiAPI.testConnection({
      base_url: config.base_url,
      api_key: config.api_key,
      model: config.model,
      provider: config.provider,
      endpoint: config.endpoint,
      query_endpoint: config.query_endpoint
    })
    ElMessage.success('Connection test passed')
  } catch (error: any) {
    ElMessage.error(error.message || 'Connection test failed')
  } finally {
    testing.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value && editingId.value) {
        const updateData: UpdateAIConfigRequest = {
          name: form.name,
          provider: form.provider,
          base_url: form.base_url,
          api_key: form.api_key,
          model: form.model,
          endpoint: form.endpoint,
          query_endpoint: form.query_endpoint,
          priority: form.priority,
          is_active: form.is_active
        }
        await aiAPI.update(editingId.value, updateData)
      } else {
        await aiAPI.create(form)
      }

      ElMessage.success(isEdit.value ? 'Config updated' : 'Config created')
      editDialogVisible.value = false
      await loadConfigs()
      emit('config-updated')
    } catch (error: any) {
      ElMessage.error(error.message || 'Save failed')
    } finally {
      submitting.value = false
    }
  })
}

const handleTabChange = (tabName: string | number) => {
  activeTab.value = tabName as AIServiceType
  resetForm()
  loadConfigs()
}

const handleProviderChange = () => {
  applyProviderPreset(form.service_type, form.provider)
}

const showQuickSetupDialog = () => {
  quickSetupApiKey.value = ''
  quickSetupVisible.value = true
}

const handleQuickSetup = async () => {
  if (!quickSetupApiKey.value.trim()) {
    ElMessage.warning('Please enter an API key')
    return
  }

  quickSetupLoading.value = true
  try {
    const baseUrl = 'https://api.chatfire.site/v1'
    const apiKey = quickSetupApiKey.value.trim()
    const serviceTypes: AIServiceType[] = ['text', 'image', 'video']

    for (const serviceType of serviceTypes) {
      const existing = allConfigs.value.find(
        (config) => config.service_type === serviceType && config.provider === 'chatfire' && config.base_url === baseUrl
      )

      if (existing) {
        continue
      }

      await aiAPI.create({
        service_type: serviceType,
        provider: 'chatfire',
        name: buildConfigName('chatfire', serviceType),
        base_url: baseUrl,
        api_key: apiKey,
        model: [getPresetModel(serviceType, 'chatfire')],
        endpoint: getProviderEndpoint(serviceType, 'chatfire'),
        query_endpoint: getProviderQueryEndpoint(serviceType, 'chatfire'),
        priority: 0
      })
    }

    ElMessage.success('ChatFire presets created')
    quickSetupVisible.value = false
    await loadConfigs()
    emit('config-updated')
  } catch (error: any) {
    ElMessage.error(error.message || 'Quick setup failed')
  } finally {
    quickSetupLoading.value = false
  }
}

watch(visible, (val) => {
  if (val) {
    resetForm()
    loadConfigs()
  }
})
</script>

<style scoped>
.dialog-header {
  display: flex;
  width: 100%;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
}

.eyebrow {
  margin: 0 0 6px;
  color: var(--text-muted);
  font-size: 0.72rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.dialog-title {
  margin: 0;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1.42rem;
  letter-spacing: -0.05em;
  line-height: 1.05;
  color: var(--text-primary);
}

.dialog-subtitle {
  margin: 8px 0 0;
  color: var(--text-muted);
  line-height: 1.65;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.soft-btn {
  border-radius: 999px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}

.overview-card {
  padding: 16px 18px;
  border: 1px solid var(--border-primary);
  border-radius: 18px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.03), transparent 30%),
    var(--bg-card);
}

.overview-card strong {
  display: block;
  margin-top: 8px;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1.7rem;
  letter-spacing: -0.05em;
  color: var(--text-primary);
}

.overview-accent {
  border-color: var(--border-secondary);
}

.overview-label {
  color: var(--text-muted);
  font-size: 0.76rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.quick-setup-info {
  margin-bottom: 16px;
  padding: 14px 16px;
  border: 1px solid var(--border-primary);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.02);
  color: var(--text-secondary);
  line-height: 1.7;
}

.quick-setup-info ul {
  margin: 12px 0 0;
  padding-left: 18px;
}

.quick-setup-tip {
  margin-top: 14px;
  color: var(--text-muted);
}

.config-form {
  margin-top: 6px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.form-tip {
  margin-top: 6px;
  color: var(--text-muted);
  line-height: 1.55;
}

.endpoint-alert {
  margin-bottom: 18px;
}

.quick-setup-footer {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.register-link {
  color: var(--text-muted);
  font-size: 0.78rem;
  text-decoration: none;
}

.register-link:hover {
  color: var(--text-primary);
}

.footer-buttons {
  display: flex;
  gap: 8px;
}

.ai-config-dialog :deep(.el-dialog__body) {
  max-height: 70vh;
  overflow-y: auto;
}

@media (max-width: 900px) {
  .overview-grid,
  .form-grid {
    grid-template-columns: 1fr;
  }

  .dialog-header,
  .quick-setup-footer {
    flex-direction: column;
    align-items: stretch;
  }

  .header-actions,
  .footer-buttons {
    width: 100%;
  }
}
</style>
