<template>
  <el-dialog
    v-model="visible"
    title="AI 图片生成"
    width="820px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="dialog-shell">
      <div class="dialog-aside">
        <p class="eyebrow">Continuity</p>
        <h3 class="font-display">Use previous images as visual anchors</h3>
        <p>
          Select one or more completed images from the same project. They will be sent as continuity references to
          keep character identity, wardrobe, lighting, and framing stable.
        </p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" class="dialog-form">
        <div class="form-grid">
          <el-form-item label="选择剧本" prop="drama_id">
            <el-select v-model="form.drama_id" placeholder="选择剧本" @change="onDramaChange">
              <el-option v-for="drama in dramas" :key="drama.id" :label="drama.title" :value="drama.id" />
            </el-select>
          </el-form-item>

          <el-form-item label="选择场景" prop="scene_id">
            <el-select v-model="form.scene_id" placeholder="可选" clearable @change="onSceneChange">
              <el-option
                v-for="scene in scenes"
                :key="scene.id"
                :label="`${scene.storyboard_number || '-'} · ${scene.title}`"
                :value="scene.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="提示词" prop="prompt">
          <el-input
            v-model="form.prompt"
            type="textarea"
            :rows="6"
            placeholder="描述镜头主体、动作、景别、情绪、材质与光线。"
            maxlength="2000"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="反向提示词">
          <el-input
            v-model="form.negative_prompt"
            type="textarea"
            :rows="3"
            placeholder="例如：deformed face, duplicate limbs, blurry details"
            maxlength="1000"
            show-word-limit
          />
        </el-form-item>

        <div class="form-grid">
          <el-form-item label="服务商">
            <el-select v-model="form.provider" placeholder="选择服务商">
              <el-option label="ChatFire" value="chatfire" />
              <el-option label="OpenAI" value="openai" />
              <el-option label="Volcengine" value="volcengine" />
              <el-option label="Gemini" value="gemini" />
              <el-option label="OpenClaw" value="openclaw" />
              <el-option label="Custom Relay" value="custom_relay" />
            </el-select>
          </el-form-item>

          <el-form-item label="尺寸">
            <el-select v-model="form.size" placeholder="选择尺寸">
              <el-option label="1024x1024 · 方图" value="1024x1024" />
              <el-option label="1792x1024 · 横图" value="1792x1024" />
              <el-option label="1024x1792 · 竖图" value="1024x1792" />
            </el-select>
          </el-form-item>
        </div>

        <div class="form-grid" v-if="['openai', 'openclaw', 'custom_relay', 'chatfire'].includes(form.provider || '')">
          <el-form-item label="质量">
            <el-radio-group v-model="form.quality">
              <el-radio label="standard">Standard</el-radio>
              <el-radio label="hd">HD</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="风格">
            <el-radio-group v-model="form.style">
              <el-radio label="vivid">Vivid</el-radio>
              <el-radio label="natural">Natural</el-radio>
            </el-radio-group>
          </el-form-item>
        </div>

        <el-divider content-position="left">连续性参考</el-divider>

        <el-form-item label="参考图">
          <el-select
            v-model="referenceImageIds"
            multiple
            filterable
            collapse-tags
            collapse-tags-tooltip
            placeholder="选择已完成的图片作为连续性参考"
            style="width: 100%"
          >
            <el-option
              v-for="image in referenceCandidates"
              :key="image.id"
              :label="truncateText(image.prompt, 46)"
              :value="image.id"
            />
          </el-select>
          <div class="form-tip">最多建议选择 1-4 张。优先选择同角色、同场景、同镜头语言的图片。</div>
        </el-form-item>

        <div v-if="selectedReferenceImages.length > 0" class="reference-grid">
          <div v-for="image in selectedReferenceImages" :key="image.id" class="reference-card">
            <img v-if="image.image_url" :src="image.image_url" :alt="image.prompt" class="reference-thumb" />
            <div class="reference-copy">
              <strong>#{{ image.id }}</strong>
              <span>{{ truncateText(image.prompt, 56) }}</span>
            </div>
          </div>
        </div>
      </el-form>
    </div>

    <template #footer>
      <el-button @click="handleClose">{{ $t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="generating" @click="handleGenerate">生成图片</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { dramaAPI } from '@/api/drama'
import { imageAPI } from '@/api/image'
import type { Drama, Scene } from '@/types/drama'
import type { GenerateImageRequest, ImageGeneration } from '@/types/image'

interface Props {
  modelValue: boolean
  dramaId?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref<FormInstance>()
const generating = ref(false)
const dramas = ref<Drama[]>([])
const scenes = ref<Scene[]>([])
const referenceCandidates = ref<ImageGeneration[]>([])
const referenceImageIds = ref<number[]>([])

const form = reactive<GenerateImageRequest>({
  drama_id: props.dramaId || '',
  scene_id: undefined,
  prompt: '',
  negative_prompt: '',
  provider: 'chatfire',
  size: '1024x1024',
  quality: 'standard',
  style: 'vivid'
})

const rules: FormRules = {
  drama_id: [{ required: true, message: '请选择剧本', trigger: 'change' }],
  prompt: [
    { required: true, message: '请输入提示词', trigger: 'blur' },
    { min: 5, message: '提示词至少 5 个字符', trigger: 'blur' }
  ]
}

const selectedReferenceImages = computed(() =>
  referenceCandidates.value.filter((image) => referenceImageIds.value.includes(image.id))
)

watch(
  () => props.modelValue,
  (val) => {
    if (!val) return
    loadDramas()
    if (props.dramaId) {
      form.drama_id = props.dramaId
      loadScenes(props.dramaId)
      loadReferenceImages(props.dramaId)
    }
  }
)

const loadDramas = async () => {
  try {
    const result = await dramaAPI.list({ page: 1, page_size: 100 })
    dramas.value = result.items || []
  } catch (error) {
    console.error('Failed to load dramas', error)
  }
}

const loadScenes = async (dramaId: string) => {
  try {
    const drama = await dramaAPI.get(dramaId)
    const collectedScenes: Scene[] = []
    drama.episodes?.forEach((episode) => {
      if (episode.scenes) {
        collectedScenes.push(...episode.scenes)
      }
    })
    scenes.value = collectedScenes
  } catch (error) {
    console.error('Failed to load scenes', error)
  }
}

const loadReferenceImages = async (dramaId: string) => {
  try {
    const result = await imageAPI.listImages({
      drama_id: dramaId,
      status: 'completed',
      page: 1,
      page_size: 100
    })
    referenceCandidates.value = result.items || []
  } catch (error) {
    console.error('Failed to load reference images', error)
  }
}

const onDramaChange = (dramaId: string) => {
  form.scene_id = undefined
  referenceImageIds.value = []
  scenes.value = []
  referenceCandidates.value = []
  if (dramaId) {
    loadScenes(dramaId)
    loadReferenceImages(dramaId)
  }
}

const onSceneChange = (sceneId?: number) => {
  if (!sceneId) return
  const scene = scenes.value.find((item) => item.id === sceneId)
  if (scene?.prompt && !form.prompt) {
    form.prompt = scene.prompt
  }
}

const truncateText = (text: string, length: number) => {
  if (!text || text.length <= length) return text
  return `${text.slice(0, length)}...`
}

const handleGenerate = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    generating.value = true
    try {
      const params: GenerateImageRequest = {
        drama_id: form.drama_id,
        scene_id: form.scene_id,
        prompt: form.prompt,
        negative_prompt: form.negative_prompt,
        provider: form.provider,
        size: form.size,
        quality: form.quality,
        style: form.style,
        reference_images: selectedReferenceImages.value
          .map((image) => image.local_path || image.image_url)
          .filter(Boolean) as string[]
      }

      await imageAPI.generateImage(params)
      ElMessage.success('图片生成任务已提交')
      emit('success')
      handleClose()
    } catch (error: any) {
      ElMessage.error(error.message || '生成失败')
    } finally {
      generating.value = false
    }
  })
}

const handleClose = () => {
  visible.value = false
  referenceImageIds.value = []
  formRef.value?.resetFields()
}
</script>

<style scoped>
.dialog-shell {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: 22px;
}

.dialog-aside {
  padding: 18px;
  border-radius: 18px;
  border: 1px solid rgba(199, 164, 91, 0.12);
  background: rgba(255, 255, 255, 0.025);
  height: fit-content;
}

.eyebrow {
  margin: 0 0 8px;
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.16em;
  color: var(--accent);
}

.dialog-aside h3 {
  margin: 0 0 12px;
  color: var(--text-primary);
}

.dialog-aside p {
  color: var(--text-secondary);
  line-height: 1.7;
}

.dialog-form {
  min-width: 0;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.form-tip {
  margin-top: 6px;
  color: var(--text-muted);
}

.reference-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.reference-card {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 16px;
  border: 1px solid rgba(199, 164, 91, 0.1);
  background: rgba(255, 255, 255, 0.025);
}

.reference-thumb {
  width: 88px;
  height: 88px;
  border-radius: 12px;
  object-fit: cover;
}

.reference-copy {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--text-secondary);
}

.reference-copy strong {
  color: var(--text-primary);
}

@media (max-width: 900px) {
  .dialog-shell,
  .form-grid,
  .reference-grid {
    grid-template-columns: 1fr;
  }
}
</style>
