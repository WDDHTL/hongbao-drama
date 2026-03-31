<template>
  <el-dialog
    v-model="visible"
    :title="$t('drama.createNew')"
    width="980px"
    :close-on-click-modal="false"
    class="create-dialog"
    @closed="handleClosed"
  >
    <div class="dialog-desc">{{ $t('drama.createDesc') }}</div>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      class="create-form"
      @submit.prevent="handleSubmit"
    >
      <div class="form-grid">
        <el-form-item :label="$t('drama.projectName')" prop="title" required>
          <el-input
            v-model="form.title"
            :placeholder="$t('drama.projectNamePlaceholder')"
            size="large"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>

        <el-form-item :label="$t('drama.projectDesc')" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            :placeholder="$t('drama.projectDescPlaceholder')"
            maxlength="500"
            show-word-limit
            resize="none"
          />
        </el-form-item>
      </div>

      <el-form-item :label="$t('drama.videoLanguage')" prop="video_language" required>
        <div class="language-shell" role="radiogroup" :aria-label="t('drama.videoLanguage')">
          <button
            v-for="item in languageOptions"
            :key="item.value"
            type="button"
            :class="['language-card', { selected: form.video_language === item.value }]"
            @click="form.video_language = item.value"
          >
            <strong>{{ t(item.labelKey) }}</strong>
            <span>{{ t(item.noteKey) }}</span>
          </button>
        </div>
        <div class="style-note">{{ $t('drama.videoLanguageHint') }}</div>
      </el-form-item>

      <el-form-item :label="$t('drama.style')" prop="style" required>
        <div class="style-shell">
          <div class="style-grid" role="radiogroup" aria-label="Drama styles">
            <button
              v-for="item in styleOptions"
              :key="item.value"
              type="button"
              :class="['style-card', { selected: form.style === item.value }]"
              @click="form.style = item.value"
            >
              <img :src="item.previewImage" :alt="t(item.labelKey)" class="style-preview" />
              <div class="style-copy">
                <strong>{{ t(item.labelKey) }}</strong>
                <span>{{ item.summary }}</span>
              </div>
            </button>
          </div>

          <div v-if="selectedStyle" class="style-inspector">
            <div class="inspector-head">
              <div>
                <p class="inspector-kicker">风格预览</p>
                <h3>{{ t(selectedStyle.labelKey) }}</h3>
              </div>
              <el-button text @click="copySelectedPrompt">复制提示词</el-button>
            </div>

            <img :src="selectedStyle.previewImage" :alt="t(selectedStyle.labelKey)" class="inspector-image" />

            <div class="inspector-meta">
              <div class="meta-block">
                <label>当前占位图路径</label>
                <code>{{ selectedStyle.assetPath }}</code>
              </div>
              <div class="meta-block">
                <label>替换说明</label>
                <span>如果你后面给我真实风格图，我直接替换这个位置或改映射即可。</span>
              </div>
              <div class="meta-block prompt-block">
                <label>对应生图提示词</label>
                <p>{{ selectedStyle.prompt }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="style-note">
          当前预览图是我放进去的占位风格图。真实图片资源统一映射在
          <code>/web/public/style-previews/</code>。
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button size="large" @click="handleClose">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">
          <el-icon v-if="!loading"><Plus /></el-icon>
          {{ $t('drama.createNew') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { dramaAPI } from '@/api/drama'
import { dramaStyleGallery, dramaStyleGalleryMap } from '@/constants/dramaStyleGallery'
import { videoLanguageOptions } from '@/constants/videoLanguages'
import type { CreateDramaRequest } from '@/types/drama'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  created: [id: string]
}>()

const { t } = useI18n()
const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const visible = ref(props.modelValue)
watch(
  () => props.modelValue,
  (val) => {
    visible.value = val
  }
)
watch(visible, (val) => {
  emit('update:modelValue', val)
})

const form = reactive<CreateDramaRequest>({
  title: '',
  description: '',
  style: 'ghibli',
  video_language: 'zh-CN'
})

const styleOptions = dramaStyleGallery
const languageOptions = videoLanguageOptions
const selectedStyle = computed(() => dramaStyleGalleryMap.get(form.style || 'ghibli'))

const rules: FormRules = {
  title: [
    { required: true, message: '请输入项目标题', trigger: 'blur' },
    { min: 1, max: 100, message: '标题长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  video_language: [{ required: true, message: '请选择视频语言', trigger: 'change' }],
  style: [{ required: true, message: '请选择风格', trigger: 'change' }]
}

const handleClosed = () => {
  form.title = ''
  form.description = ''
  form.video_language = 'zh-CN'
  formRef.value?.resetFields()
}

const handleClose = () => {
  visible.value = false
}

const copySelectedPrompt = async () => {
  if (!selectedStyle.value?.prompt) return
  try {
    await navigator.clipboard.writeText(selectedStyle.value.prompt)
    ElMessage.success('风格提示词已复制')
  } catch {
    ElMessage.warning('复制失败，请手动复制提示词')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const drama = await dramaAPI.create(form)
      ElMessage.success('创建成功')
      visible.value = false
      emit('created', drama.id)
      router.push(`/dramas/${drama.id}`)
    } catch (error: any) {
      ElMessage.error(error.message || '创建失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.create-dialog :deep(.el-dialog) {
  border-radius: var(--radius-2xl);
}

.create-dialog :deep(.el-dialog__header) {
  margin-right: 0;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-primary);
}

.create-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.dialog-desc {
  margin-bottom: 20px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
}

.create-form :deep(.el-form-item) {
  margin-bottom: 20px;
}

.style-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(300px, 0.8fr);
  gap: 16px;
  width: 100%;
}

.language-shell {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.language-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-height: 82px;
  padding: 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  background: var(--bg-card);
  cursor: pointer;
  text-align: left;
  transition:
    transform var(--transition-fast),
    border-color var(--transition-fast),
    box-shadow var(--transition-fast);
}

.language-card:hover {
  transform: translateY(-2px);
  border-color: var(--border-secondary);
  box-shadow: var(--shadow-card);
}

.language-card.selected {
  border-color: var(--border-focus);
  box-shadow: var(--shadow-glow);
}

.language-card strong {
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 0.95rem;
  color: var(--text-primary);
}

.language-card span {
  color: var(--text-muted);
  font-size: 0.8rem;
  line-height: 1.55;
}

.style-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.style-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  background: var(--bg-card);
  cursor: pointer;
  text-align: left;
  transition:
    transform var(--transition-fast),
    border-color var(--transition-fast),
    box-shadow var(--transition-fast);
}

.style-card:hover {
  transform: translateY(-2px);
  border-color: var(--border-secondary);
  box-shadow: var(--shadow-card);
}

.style-card.selected {
  border-color: var(--border-focus);
  box-shadow: var(--shadow-glow);
}

.style-preview {
  width: 100%;
  aspect-ratio: 4 / 3;
  max-height: 92px;
  object-fit: cover;
  border-radius: 12px;
  background: var(--bg-secondary);
}

.style-copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.style-copy strong {
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 0.88rem;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.style-copy span {
  color: var(--text-muted);
  font-size: 0.74rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.style-inspector {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  background: rgba(255, 255, 255, 0.02);
}

.inspector-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.inspector-kicker {
  margin: 0 0 6px;
  color: var(--text-muted);
  font-size: 0.72rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.inspector-head h3 {
  margin: 0;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1.18rem;
  letter-spacing: -0.04em;
}

.inspector-image {
  width: 100%;
  aspect-ratio: 16 / 10;
  object-fit: cover;
  border-radius: 16px;
  background: var(--bg-secondary);
}

.inspector-meta {
  display: grid;
  gap: 12px;
}

.meta-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.meta-block label {
  color: var(--text-muted);
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.meta-block span,
.meta-block code,
.meta-block p {
  color: var(--text-secondary);
  line-height: 1.65;
}

.meta-block code {
  padding: 8px 10px;
  border-radius: 10px;
  background: var(--accent-light);
  word-break: break-all;
}

.prompt-block p {
  margin: 0;
}

.style-note {
  margin-top: 10px;
  color: var(--text-muted);
  font-size: 0.8rem;
  line-height: 1.6;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.dialog-footer .el-button {
  min-width: 110px;
}

@media (max-width: 900px) {
  .form-grid,
  .style-shell,
  .style-grid,
  .language-shell {
    grid-template-columns: 1fr;
  }
}
</style>
