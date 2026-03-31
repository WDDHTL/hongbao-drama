<template>
  <div class="auto-workflow-page">
    <AppHeader :fixed="false" :show-logo="false">
      <template #left>
        <el-button text class="back-btn" @click="router.push(`/dramas/${dramaId}`)">
          <el-icon><ArrowLeft /></el-icon>
          <span>返回项目</span>
        </el-button>
        <div class="page-title">
          <h1>{{ drama?.title || '自动工作流' }}</h1>
          <span class="subtitle">输入剧本后，从角色、场景、分镜、图片、视频到合成的自动流水线</span>
        </div>
      </template>
      <template #right>
        <el-button :loading="starting" type="primary" @click="startWorkflow">
          开始自动生成
        </el-button>
        <el-button v-if="canResume" :loading="resuming" @click="resumeWorkflow">
          继续执行
        </el-button>
      </template>
    </AppHeader>

    <div class="workflow-shell" v-loading="loading">
      <section class="hero-panel">
        <div class="hero-copy">
          <div class="hero-kicker">Project Automation</div>
          <h2>整集自动生产控制台</h2>
          <p>
            当前工作流会自动准备剧集、提取角色/场景/道具、生成角色基准图，并逐集完成分镜、图片、视频和最终合成。
          </p>
        </div>
        <div class="hero-metrics">
          <div class="metric-card">
            <span>当前状态</span>
            <strong>{{ getRunStatusText(status.run?.status) }}</strong>
          </div>
          <div class="metric-card">
            <span>当前阶段</span>
            <strong>{{ formatStage(status.run?.current_stage) }}</strong>
          </div>
          <div class="metric-card">
            <span>总进度</span>
            <strong>{{ status.run?.progress ?? 0 }}%</strong>
          </div>
        </div>
      </section>

      <el-alert
        v-if="status.run?.error_msg"
        class="status-alert"
        type="error"
        :closable="false"
        :title="status.run.error_msg"
      />

      <section class="run-overview">
        <div class="section-head">
          <h3>工作流总览</h3>
          <span v-if="status.run">最近一次执行：{{ formatTime(status.run.updated_at) }}</span>
        </div>
        <el-progress
          :percentage="status.run?.progress ?? 0"
          :status="progressStatus"
          :stroke-width="16"
        />
        <div class="stage-strip">
          <div
            v-for="item in stageOrder"
            :key="item.key"
            class="stage-pill"
            :class="getStageClass(item.key)"
          >
            <span class="stage-name">{{ item.label }}</span>
            <small>{{ getStageStatusText(item.key) }}</small>
          </div>
        </div>
      </section>

      <section class="episode-panel">
        <div class="section-head">
          <h3>分集执行情况</h3>
          <span>{{ status.episodes.length }} 集</span>
        </div>
        <div v-if="status.episodes.length === 0" class="empty-card">
          当前项目还没有可执行的分集。工作流会优先从项目描述或 `metadata.full_script` 自动创建第 1 集。
        </div>
        <div v-else class="episode-grid">
          <article
            v-for="episode in status.episodes"
            :key="episode.episode_id"
            class="episode-card"
            :class="{ warning: episode.needs_manual_review }"
          >
            <div class="episode-head">
              <div>
                <div class="episode-kicker">第 {{ episode.episode_number }} 集</div>
                <h4>{{ episode.title }}</h4>
              </div>
              <el-tag :type="episode.needs_manual_review ? 'danger' : 'success'" effect="plain">
                {{ episode.needs_manual_review ? '需人工处理' : '可继续自动化' }}
              </el-tag>
            </div>

            <div class="episode-stats">
              <div>
                <span>分镜</span>
                <strong>{{ episode.storyboard_count }}</strong>
              </div>
              <div>
                <span>图片完成</span>
                <strong>{{ episode.completed_images }}</strong>
              </div>
              <div>
                <span>视频完成</span>
                <strong>{{ episode.completed_videos }}</strong>
              </div>
            </div>

            <div class="episode-foot">
              <span>合成状态：{{ formatMergeStatus(episode.merge_status) }}</span>
              <a
                v-if="episode.final_video_url"
                :href="resolveStaticUrl(episode.final_video_url)"
                target="_blank"
                rel="noreferrer"
              >
                查看成片
              </a>
            </div>
          </article>
        </div>
      </section>

      <section class="steps-panel">
        <div class="section-head">
          <h3>执行日志</h3>
          <span>{{ status.steps.length }} 条</span>
        </div>
        <div v-if="status.steps.length === 0" class="empty-card">
          还没有执行记录。点击“开始自动生成”即可启动全流程。
        </div>
        <div v-else class="step-list">
          <div
            v-for="step in status.steps"
            :key="step.id"
            class="step-row"
            :class="step.status"
          >
            <div class="step-status-dot" />
            <div class="step-main">
              <div class="step-title">
                <strong>{{ formatStage(step.stage) }}</strong>
                <span v-if="step.episode_id">Episode #{{ findEpisodeNumber(step.episode_id) }}</span>
              </div>
              <p>{{ step.message || '等待执行' }}</p>
              <small v-if="step.error_msg">{{ step.error_msg }}</small>
            </div>
            <div class="step-side">
              <span>{{ step.progress }}%</span>
              <em>{{ getRunStatusText(step.status) }}</em>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import AppHeader from '@/components/common/AppHeader.vue'
import { workflowAPI } from '@/api/workflow'
import { dramaAPI } from '@/api/drama'
import type { Drama } from '@/types/drama'
import type { WorkflowStatusResponse } from '@/types/workflow'

const route = useRoute()
const router = useRouter()
const dramaId = route.params.id as string

const drama = ref<Drama | null>(null)
const loading = ref(false)
const starting = ref(false)
const resuming = ref(false)
const pollTimer = ref<number | null>(null)

const status = ref<WorkflowStatusResponse>({
  run: null,
  steps: [],
  episodes: [],
})

const stageOrder = [
  { key: 'prepare_episodes', label: '准备剧集' },
  { key: 'character_extract', label: '提取角色' },
  { key: 'scene_extract', label: '提取场景' },
  { key: 'prop_extract', label: '提取道具' },
  { key: 'character_baseline_generate', label: '角色基准图' },
  { key: 'storyboard_generate', label: '分镜生成' },
  { key: 'storyboard_image_generate', label: '镜头图片' },
  { key: 'storyboard_video_generate', label: '镜头视频' },
  { key: 'episode_merge', label: '自动合成' },
]

const canResume = computed(
  () => status.value.run?.status === 'failed' || status.value.run?.status === 'paused',
)

const progressStatus = computed(() => {
  if (status.value.run?.status === 'failed') return 'exception'
  if (status.value.run?.status === 'completed') return 'success'
  return undefined
})

const normalizeStage = (stage?: string) => stage?.replace(/_episode_\d+$/, '') || ''

const loadDrama = async () => {
  drama.value = await dramaAPI.get(dramaId)
}

const loadWorkflowStatus = async (silent = false) => {
  if (!silent) loading.value = true
  try {
    status.value = await workflowAPI.getProjectWorkflowStatus(dramaId)
  } finally {
    loading.value = false
  }
}

const stopPolling = () => {
  if (pollTimer.value) {
    window.clearInterval(pollTimer.value)
    pollTimer.value = null
  }
}

const startPolling = () => {
  stopPolling()
  pollTimer.value = window.setInterval(async () => {
    await loadWorkflowStatus(true)
    if (!status.value.run || ['completed', 'failed', 'paused'].includes(status.value.run.status)) {
      stopPolling()
    }
  }, 5000)
}

const startWorkflow = async () => {
  starting.value = true
  try {
    await workflowAPI.startProjectWorkflow(dramaId)
    await loadWorkflowStatus(true)
    startPolling()
    ElMessage.success('自动工作流已启动')
  } catch (error: any) {
    ElMessage.error(error?.message || '启动失败')
  } finally {
    starting.value = false
  }
}

const resumeWorkflow = async () => {
  resuming.value = true
  try {
    await workflowAPI.resumeProjectWorkflow(dramaId)
    await loadWorkflowStatus(true)
    startPolling()
    ElMessage.success('自动工作流已继续执行')
  } catch (error: any) {
    ElMessage.error(error?.message || '继续执行失败')
  } finally {
    resuming.value = false
  }
}

const formatStage = (stage?: string) => {
  const normalized = normalizeStage(stage)
  if (!normalized) return '未开始'
  const found = stageOrder.find((item) => item.key === normalized)
  return found?.label || normalized.replace(/_/g, ' / ')
}

const getRunStatusText = (value?: string) => {
  const map: Record<string, string> = {
    pending: '排队中',
    processing: '执行中',
    completed: '已完成',
    failed: '失败',
    paused: '已暂停',
  }
  return map[value || ''] || '未开始'
}

const getStageStatusText = (stage: string) => {
  const matched = status.value.steps
    .filter((step) => normalizeStage(step.stage) === stage)
    .sort((a, b) => a.id - b.id)
    .at(-1)
  return getRunStatusText(matched?.status)
}

const getStageClass = (stage: string) => {
  const matched = status.value.steps
    .filter((step) => normalizeStage(step.stage) === stage)
    .sort((a, b) => a.id - b.id)
    .at(-1)
  return matched?.status || 'pending'
}

const findEpisodeNumber = (episodeId?: number) => {
  return (
    status.value.episodes.find((episode) => episode.episode_id === episodeId)?.episode_number ??
    episodeId ??
    '-'
  )
}

const formatMergeStatus = (statusValue?: string) => {
  if (!statusValue) return '未开始'
  const map: Record<string, string> = {
    pending: '排队中',
    processing: '合成中',
    completed: '已完成',
    failed: '失败',
  }
  return map[statusValue] || statusValue
}

const resolveStaticUrl = (value: string) => {
  if (!value) return ''
  const normalized = value.replace(/\\/g, '/')
  if (normalized.startsWith('http')) return normalized
  if (normalized.startsWith('/static/')) return normalized
  return `/static/${normalized.replace(/^\/+/, '')}`
}

const formatTime = (value?: string) => {
  if (!value) return '--'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString()
}

onMounted(async () => {
  await Promise.all([loadDrama(), loadWorkflowStatus()])
  if (status.value.run && ['pending', 'processing'].includes(status.value.run.status)) {
    startPolling()
  }
})

onBeforeUnmount(() => {
  stopPolling()
})
</script>

<style scoped lang="scss">
.auto-workflow-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(255, 255, 255, 0.06), transparent 28%),
    linear-gradient(180deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
  color: var(--text-primary);
}

.workflow-shell {
  max-width: 1360px;
  margin: 0 auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.hero-panel,
.run-overview,
.episode-panel,
.steps-panel {
  border: 1px solid var(--border-primary);
  border-radius: 24px;
  background: color-mix(in srgb, var(--bg-card) 92%, transparent);
  box-shadow: var(--shadow-sm);
}

.hero-panel {
  display: grid;
  grid-template-columns: 1.4fr 1fr;
  gap: 18px;
  padding: 24px;
}

.hero-kicker,
.episode-kicker {
  font-size: 0.72rem;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.hero-copy h2,
.section-head h3,
.episode-head h4 {
  margin: 0;
}

.hero-copy p {
  margin: 10px 0 0;
  line-height: 1.7;
  color: var(--text-secondary);
}

.hero-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.metric-card {
  padding: 16px;
  border: 1px solid var(--border-primary);
  border-radius: 18px;
  background: color-mix(in srgb, var(--bg-secondary) 88%, transparent);
}

.metric-card span {
  display: block;
  margin-bottom: 10px;
  color: var(--text-muted);
  font-size: 0.8rem;
}

.metric-card strong {
  font-size: 1rem;
}

.status-alert {
  margin-top: -4px;
}

.run-overview,
.episode-panel,
.steps-panel {
  padding: 22px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.section-head span {
  color: var(--text-muted);
  font-size: 0.86rem;
}

.stage-strip {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.stage-pill {
  padding: 14px;
  border-radius: 18px;
  border: 1px solid var(--border-primary);
  background: color-mix(in srgb, var(--bg-secondary) 90%, transparent);
}

.stage-pill small,
.episode-foot span,
.step-main p,
.step-main small {
  color: var(--text-muted);
}

.stage-pill.processing {
  border-color: var(--border-focus);
}

.stage-pill.completed {
  border-color: color-mix(in srgb, var(--success) 30%, var(--border-primary));
}

.stage-pill.failed {
  border-color: color-mix(in srgb, var(--danger) 40%, var(--border-primary));
}

.stage-name {
  display: block;
  margin-bottom: 6px;
  font-weight: 700;
}

.empty-card {
  padding: 22px;
  border: 1px dashed var(--border-primary);
  border-radius: 20px;
  color: var(--text-muted);
  line-height: 1.7;
}

.episode-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.episode-card {
  padding: 18px;
  border: 1px solid var(--border-primary);
  border-radius: 20px;
  background: color-mix(in srgb, var(--bg-secondary) 90%, transparent);
}

.episode-card.warning {
  border-color: color-mix(in srgb, var(--danger) 32%, var(--border-primary));
}

.episode-head,
.episode-foot,
.step-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.episode-stats {
  margin: 16px 0;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.episode-stats div {
  padding: 12px;
  border-radius: 16px;
  background: color-mix(in srgb, var(--bg-card) 92%, transparent);
}

.episode-stats span {
  display: block;
  margin-bottom: 6px;
  color: var(--text-muted);
  font-size: 0.78rem;
}

.episode-stats strong {
  font-size: 1.1rem;
}

.episode-foot a {
  color: var(--accent);
  text-decoration: none;
}

.step-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.step-row {
  display: grid;
  grid-template-columns: 12px minmax(0, 1fr) auto;
  gap: 14px;
  padding: 16px;
  border: 1px solid var(--border-primary);
  border-radius: 18px;
  background: color-mix(in srgb, var(--bg-secondary) 90%, transparent);
}

.step-status-dot {
  width: 12px;
  height: 12px;
  margin-top: 6px;
  border-radius: 999px;
  background: var(--border-secondary);
}

.step-row.processing .step-status-dot {
  background: var(--warning);
}

.step-row.completed .step-status-dot {
  background: var(--success);
}

.step-row.failed .step-status-dot {
  background: var(--danger);
}

.step-title strong {
  font-size: 0.95rem;
}

.step-title span {
  color: var(--text-muted);
  font-size: 0.8rem;
}

.step-main p,
.step-main small {
  margin: 8px 0 0;
  line-height: 1.6;
}

.step-side {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
  white-space: nowrap;
}

.page-title h1 {
  margin: 0;
}

.subtitle {
  color: var(--text-muted);
}

.back-btn {
  margin-right: 12px;
}

@media (max-width: 960px) {
  .hero-panel {
    grid-template-columns: 1fr;
  }

  .hero-metrics,
  .stage-strip,
  .episode-stats {
    grid-template-columns: 1fr;
  }

  .workflow-shell {
    padding: 16px;
  }
}
</style>
