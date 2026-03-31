<template>
  <el-dialog
    v-model="visible"
    title="AI 视频生成"
    width="920px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="dialog-shell">
      <div class="dialog-aside">
        <p class="eyebrow">Shot Continuity</p>
        <h3 class="font-display">Single frame, first-last frame, or multi-reference</h3>
        <p>
          First/last frame mode is the most stable option when you need the opening and ending composition to stay
          consistent. Multi-reference mode is better for character identity and costume continuity.
        </p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" class="dialog-form">
        <div class="form-grid">
          <el-form-item label="选择剧本" prop="drama_id">
            <el-select v-model="form.drama_id" placeholder="选择剧本" @change="onDramaChange">
              <el-option v-for="drama in dramas" :key="drama.id" :label="drama.title" :value="drama.id" />
            </el-select>
          </el-form-item>

          <el-form-item label="服务商">
            <el-select v-model="form.provider" placeholder="选择服务商">
              <el-option label="ChatFire" value="chatfire" />
              <el-option label="Volcengine" value="volces" />
              <el-option label="OpenAI" value="openai" />
              <el-option label="MiniMax" value="minimax" />
              <el-option label="OpenClaw" value="openclaw" />
              <el-option label="Custom Relay" value="custom_relay" />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="参考模式">
          <el-radio-group v-model="referenceMode">
            <el-radio-button label="single">单图</el-radio-button>
            <el-radio-button label="first_last">首尾帧</el-radio-button>
            <el-radio-button label="multiple">多参考图</el-radio-button>
            <el-radio-button label="none">纯文本</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="referenceMode === 'single'" label="参考图片">
          <el-select v-model="singleImageId" placeholder="选择一张已生成图片" clearable @change="syncPromptFromImage">
            <el-option
              v-for="image in imageCandidates"
              :key="image.id"
              :label="truncateText(image.prompt, 52)"
              :value="image.id"
            />
          </el-select>
          <div class="form-tip">适合 i2v 或让当前镜头承接上一张关键画面。</div>
        </el-form-item>

        <div v-if="referenceMode === 'first_last'" class="form-grid">
          <el-form-item label="首帧">
            <el-select v-model="firstFrameId" placeholder="选择首帧图片" clearable>
              <el-option
                v-for="image in imageCandidates"
                :key="`first-${image.id}`"
                :label="truncateText(image.prompt, 40)"
                :value="image.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="尾帧">
            <el-select v-model="lastFrameId" placeholder="选择尾帧图片" clearable>
              <el-option
                v-for="image in imageCandidates"
                :key="`last-${image.id}`"
                :label="truncateText(image.prompt, 40)"
                :value="image.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item v-if="referenceMode === 'multiple'" label="参考图片组">
          <el-select
            v-model="multipleImageIds"
            multiple
            filterable
            collapse-tags
            collapse-tags-tooltip
            placeholder="选择多张参考图片"
            style="width: 100%"
          >
            <el-option
              v-for="image in imageCandidates"
              :key="`multi-${image.id}`"
              :label="truncateText(image.prompt, 48)"
              :value="image.id"
            />
          </el-select>
          <div class="form-tip">适合锁定角色身份、服装、材质和场景资产。</div>
        </el-form-item>

        <div v-if="selectedPreviewImages.length > 0" class="reference-grid">
          <div v-for="image in selectedPreviewImages" :key="image.id" class="reference-card">
            <img v-if="image.image_url" :src="image.image_url" :alt="image.prompt" class="reference-thumb" />
            <div class="reference-copy">
              <strong>#{{ image.id }}</strong>
              <span>{{ truncateText(image.prompt, 56) }}</span>
            </div>
          </div>
        </div>

        <el-form-item label="视频提示词" prop="prompt">
          <el-input
            v-model="form.prompt"
            type="textarea"
            :rows="5"
            placeholder="描述动作、镜头运动、节奏、表演强度与氛围变化。"
            maxlength="2000"
            show-word-limit
          />
        </el-form-item>

        <div class="form-grid">
          <el-form-item label="视频时长">
            <el-slider v-model="form.duration" :min="3" :max="10" :marks="durationMarks" show-stops />
          </el-form-item>

          <el-form-item label="画幅比例">
            <el-radio-group v-model="form.aspect_ratio">
              <el-radio label="16:9">16:9</el-radio>
              <el-radio label="9:16">9:16</el-radio>
              <el-radio label="1:1">1:1</el-radio>
            </el-radio-group>
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="运动强度">
            <el-slider v-model="form.motion_level" :min="0" :max="100" :marks="motionMarks" />
          </el-form-item>

          <el-form-item label="镜头运动">
            <el-select v-model="form.camera_motion" placeholder="可选" clearable>
              <el-option label="Static" value="static" />
              <el-option label="Zoom In" value="zoom_in" />
              <el-option label="Zoom Out" value="zoom_out" />
              <el-option label="Pan Left" value="pan_left" />
              <el-option label="Pan Right" value="pan_right" />
              <el-option label="Tilt Up" value="tilt_up" />
              <el-option label="Tilt Down" value="tilt_down" />
              <el-option label="Orbit" value="orbit" />
            </el-select>
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="风格">
            <el-input v-model="form.style" placeholder="例如 cinematic, premium fashion, moody realism" />
          </el-form-item>

          <el-form-item label="随机种子">
            <el-input-number v-model="form.seed" :min="-1" placeholder="留空随机" style="width: 100%" />
          </el-form-item>
        </div>
      </el-form>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="generating" @click="handleGenerate">生成视频</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { dramaAPI } from '@/api/drama'
import { imageAPI } from '@/api/image'
import { videoAPI } from '@/api/video'
import type { Drama } from '@/types/drama'
import type { ImageGeneration } from '@/types/image'
import type { GenerateVideoRequest } from '@/types/video'

interface Props {
  modelValue: boolean
  dramaId?: string
}

type ReferenceMode = 'single' | 'first_last' | 'multiple' | 'none'

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
const imageCandidates = ref<ImageGeneration[]>([])

const referenceMode = ref<ReferenceMode>('single')
const singleImageId = ref<number>()
const firstFrameId = ref<number>()
const lastFrameId = ref<number>()
const multipleImageIds = ref<number[]>([])

const form = reactive<GenerateVideoRequest>({
  drama_id: props.dramaId || '',
  prompt: '',
  provider: 'chatfire',
  duration: 5,
  aspect_ratio: '16:9',
  motion_level: 45,
  camera_motion: undefined,
  style: '',
  seed: undefined
})

const rules: FormRules = {
  drama_id: [{ required: true, message: '请选择剧本', trigger: 'change' }],
  prompt: [
    { required: true, message: '请输入视频提示词', trigger: 'blur' },
    { min: 5, message: '提示词至少 5 个字符', trigger: 'blur' }
  ]
}

const durationMarks = {
  3: '3s',
  5: '5s',
  7: '7s',
  10: '10s'
}

const motionMarks = {
  0: '静',
  50: '中',
  100: '强'
}

const selectedPreviewImages = computed(() => {
  const ids =
    referenceMode.value === 'single'
      ? [singleImageId.value]
      : referenceMode.value === 'first_last'
        ? [firstFrameId.value, lastFrameId.value]
        : multipleImageIds.value

  return imageCandidates.value.filter((image) => ids.filter(Boolean).includes(image.id))
})

watch(
  () => props.modelValue,
  (val) => {
    if (!val) return
    loadDramas()
    if (props.dramaId) {
      form.drama_id = props.dramaId
      loadImages(props.dramaId)
    }
  }
)

const loadDramas = async () => {
  try {
    const result = await dramaAPI.list({ page: 1, page_size: 100 })
    dramas.value = result.items
  } catch (error) {
    console.error('Failed to load dramas', error)
  }
}

const loadImages = async (dramaId: string) => {
  try {
    const result = await imageAPI.listImages({
      drama_id: dramaId,
      status: 'completed',
      page: 1,
      page_size: 100
    })
    imageCandidates.value = result.items
  } catch (error) {
    console.error('Failed to load images', error)
  }
}

const onDramaChange = (dramaId: string) => {
  singleImageId.value = undefined
  firstFrameId.value = undefined
  lastFrameId.value = undefined
  multipleImageIds.value = []
  imageCandidates.value = []
  if (dramaId) {
    loadImages(dramaId)
  }
}

const getImageById = (id?: number) => imageCandidates.value.find((image) => image.id === id)

const syncPromptFromImage = (imageId?: number) => {
  const image = getImageById(imageId)
  if (image && !form.prompt) {
    form.prompt = image.prompt
  }
}

const getImagePath = (image?: ImageGeneration) => image?.local_path || image?.image_url || ''

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
      const params: GenerateVideoRequest = {
        drama_id: form.drama_id,
        prompt: form.prompt,
        provider: form.provider,
        duration: form.duration,
        aspect_ratio: form.aspect_ratio,
        motion_level: form.motion_level,
        camera_motion: form.camera_motion,
        style: form.style,
        seed: form.seed
      }

      if (referenceMode.value === 'single') {
        const image = getImageById(singleImageId.value)
        params.reference_mode = 'single'
        params.image_gen_id = image?.id
        params.image_url = image?.image_url
        params.image_local_path = image?.local_path
      } else if (referenceMode.value === 'first_last') {
        const first = getImageById(firstFrameId.value)
        const last = getImageById(lastFrameId.value)
        params.reference_mode = 'first_last'
        params.first_frame_url = first?.image_url
        params.first_frame_local_path = first?.local_path
        params.last_frame_url = last?.image_url
        params.last_frame_local_path = last?.local_path
      } else if (referenceMode.value === 'multiple') {
        params.reference_mode = 'multiple'
        params.reference_image_urls = multipleImageIds.value
          .map((id) => getImagePath(getImageById(id)))
          .filter(Boolean)
      } else {
        params.reference_mode = 'none'
      }

      await videoAPI.generateVideo(params)
      ElMessage.success('视频生成任务已提交')
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
  singleImageId.value = undefined
  firstFrameId.value = undefined
  lastFrameId.value = undefined
  multipleImageIds.value = []
  referenceMode.value = 'single'
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
  margin-bottom: 16px;
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
