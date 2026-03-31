<template>
  <div class="studio-page">
    <AppHeader :show-logo="false">
      <template #left>
        <div class="page-title">
          <span class="page-kicker">Video Studio</span>
          <h1>视频镜头中心</h1>
          <p>围绕首帧、尾帧和多参考图重建视频页，让镜头控制和任务状态都落进同一套黑白组件里。</p>
        </div>
      </template>
      <template #right>
        <el-button type="primary" @click="showGenerateDialog = true">
          <el-icon><VideoPlay /></el-icon>
          <span>新建视频任务</span>
        </el-button>
      </template>
    </AppHeader>

    <section class="hero-grid">
      <BaseCard title="镜头连贯策略" subtitle="用同一批图片资产驱动视频单图、首尾帧和多参考三种模式。">
        <div class="hero-copy">
          <p>现在的视频弹窗已经支持首尾帧和多参考图输入。图片页里生成的成片会直接进入这个链路，减少人物漂移和镜头断裂。</p>
          <div class="workflow-list">
            <div class="workflow-item">
              <span class="workflow-index">A</span>
              <span>单图模式适合快速试镜头</span>
            </div>
            <div class="workflow-item">
              <span class="workflow-index">B</span>
              <span>首尾帧适合做动作和机位推进</span>
            </div>
            <div class="workflow-item">
              <span class="workflow-index">C</span>
              <span>多参考适合角色细节稳定场景</span>
            </div>
          </div>
        </div>
      </BaseCard>

      <div class="stat-grid">
        <StatCard label="总任务" :value="total" variant="compact" />
        <StatCard label="已完成" :value="completedCount" variant="compact" />
        <StatCard label="处理中" :value="processingCount" variant="compact" />
        <StatCard label="失败" :value="failedCount" variant="compact" />
      </div>
    </section>

    <BaseCard title="筛选与状态" subtitle="按剧目和任务状态过滤视频资产。">
      <div class="filter-shell">
        <el-form class="filter-form" label-position="top">
          <div class="filter-grid">
            <el-form-item label="剧目">
              <el-select v-model="filters.drama_id" placeholder="全部剧目" clearable>
                <el-option v-for="drama in dramas" :key="drama.id" :label="drama.title" :value="drama.id" />
              </el-select>
            </el-form-item>

            <el-form-item label="状态">
              <el-select v-model="filters.status" placeholder="全部状态" clearable>
                <el-option label="待开始" value="pending" />
                <el-option label="处理中" value="processing" />
                <el-option label="已完成" value="completed" />
                <el-option label="失败" value="failed" />
              </el-select>
            </el-form-item>

            <el-form-item label="当前剧目">
              <div class="active-summary">
                <strong>{{ activeDrama?.title || '未选择剧目' }}</strong>
                <span>{{ activeDrama?.style || '当前展示所有视频资产' }}</span>
              </div>
            </el-form-item>
          </div>
        </el-form>

        <div class="filter-actions">
          <el-button type="primary" @click="loadVideos">刷新列表</el-button>
          <el-button @click="resetFilters">重置筛选</el-button>
        </div>
      </div>
    </BaseCard>

    <section v-loading="loading" class="asset-section">
      <div v-if="videos.length" class="asset-grid">
        <BaseCard v-for="video in videos" :key="video.id" hoverable class="asset-card" no-padding>
          <div class="video-media">
            <video
              v-if="video.status === 'completed' && video.video_url"
              :src="video.video_url"
              class="video-player"
              controls
              :poster="video.first_frame_url"
            >
              当前浏览器不支持视频播放。
            </video>

            <div v-else-if="video.status === 'processing'" class="asset-placeholder processing">
              <el-icon class="spin"><Loading /></el-icon>
              <span>视频生成中</span>
              <small>系统会自动轮询任务状态</small>
            </div>

            <div v-else-if="video.status === 'failed'" class="asset-placeholder failed">
              <el-icon><CircleClose /></el-icon>
              <span>任务失败</span>
            </div>

            <div v-else class="asset-placeholder">
              <el-icon><VideoCamera /></el-icon>
              <span>等待生成</span>
            </div>

            <div class="asset-badges">
              <el-tag size="small" :type="getStatusType(video.status)">{{ getStatusText(video.status) }}</el-tag>
              <el-tag v-if="video.duration" size="small" effect="plain">{{ video.duration }}s</el-tag>
            </div>
          </div>

          <div class="asset-content">
            <div class="asset-header">
              <strong>{{ video.model || video.provider }}</strong>
              <span>#{{ video.id }}</span>
            </div>

            <p class="asset-prompt">{{ truncateText(video.prompt, 116) }}</p>

            <div class="asset-meta">
              <span>{{ video.provider }}</span>
              <span>{{ video.aspect_ratio || '自动比例' }}</span>
              <span>{{ formatTime(video.created_at) }}</span>
            </div>

            <div class="asset-extra">
              <span v-if="video.resolution">{{ video.resolution }}</span>
              <span v-if="video.camera_motion">{{ video.camera_motion }}</span>
              <span v-if="video.motion_level">运动等级 {{ video.motion_level }}</span>
            </div>

            <div class="asset-actions">
              <ActionButton :icon="View" tooltip="查看详情" @click="viewDetails(video)" />
              <ActionButton
                v-if="video.status === 'completed'"
                :icon="Download"
                tooltip="打开视频"
                @click="downloadVideo(video)"
              />
              <el-popconfirm title="确认删除这条视频记录？" @confirm="deleteVideo(video.id)">
                <template #reference>
                  <span>
                    <ActionButton :icon="Delete" variant="danger" tooltip="删除记录" />
                  </span>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </BaseCard>
      </div>

      <el-empty v-else-if="!loading" description="当前没有视频任务，先从图片资产挑一组参考开始生成。" class="studio-empty" />
    </section>

    <el-pagination
      v-if="total > 0"
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.page_size"
      :total="total"
      :page-sizes="[12, 24, 36, 48]"
      layout="total, sizes, prev, pager, next, jumper"
      class="pagination-shell"
      @current-change="loadVideos"
      @size-change="loadVideos"
    />

    <GenerateVideoDialog v-model="showGenerateDialog" :drama-id="filters.drama_id" @success="loadVideos" />
    <VideoDetailDialog v-model="showDetailDialog" :video="selectedVideo" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  CircleClose,
  Delete,
  Download,
  Loading,
  VideoCamera,
  VideoPlay,
  View
} from '@element-plus/icons-vue'
import { dramaAPI } from '@/api/drama'
import { videoAPI } from '@/api/video'
import { ActionButton, AppHeader, BaseCard, StatCard } from '@/components/common'
import type { Drama } from '@/types/drama'
import type { VideoGeneration, VideoStatus } from '@/types/video'
import GenerateVideoDialog from './components/GenerateVideoDialog.vue'
import VideoDetailDialog from './components/VideoDetailDialog.vue'

const route = useRoute()

const loading = ref(false)
const total = ref(0)
const videos = ref<VideoGeneration[]>([])
const dramas = ref<Drama[]>([])
const showGenerateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedVideo = ref<VideoGeneration>()
let pollTimer: number | null = null

const filters = reactive({
  drama_id: undefined as string | undefined,
  status: undefined as VideoStatus | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 12
})

const completedCount = computed(() => videos.value.filter((item) => item.status === 'completed').length)
const processingCount = computed(() => videos.value.filter((item) => item.status === 'processing').length)
const failedCount = computed(() => videos.value.filter((item) => item.status === 'failed').length)
const activeDrama = computed(() => dramas.value.find((item) => item.id === filters.drama_id))

const loadVideos = async () => {
  loading.value = true
  try {
    const result = await videoAPI.listVideos({
      drama_id: filters.drama_id,
      status: filters.status,
      page: pagination.page,
      page_size: pagination.page_size
    })
    videos.value = result.items
    total.value = result.pagination.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载视频列表失败')
  } finally {
    loading.value = false
  }
}

const loadDramas = async () => {
  try {
    const result = await dramaAPI.list({ page: 1, page_size: 100 })
    dramas.value = result.items
  } catch (error) {
    console.error('Failed to load dramas', error)
  }
}

const resetFilters = () => {
  filters.drama_id = undefined
  filters.status = undefined
  pagination.page = 1
  loadVideos()
}

const viewDetails = (video: VideoGeneration) => {
  selectedVideo.value = video
  showDetailDialog.value = true
}

const downloadVideo = (video: VideoGeneration) => {
  if (!video.video_url) return
  window.open(video.video_url, '_blank')
}

const deleteVideo = async (id: number) => {
  try {
    await videoAPI.deleteVideo(id)
    ElMessage.success('视频记录已删除')
    loadVideos()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

const getStatusType = (status: VideoStatus) => {
  const map: Record<VideoStatus, 'info' | 'warning' | 'success' | 'danger'> = {
    pending: 'info',
    processing: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return map[status]
}

const getStatusText = (status: VideoStatus) => {
  const map: Record<VideoStatus, string> = {
    pending: '待开始',
    processing: '处理中',
    completed: '已完成',
    failed: '失败'
  }
  return map[status]
}

const truncateText = (text: string, length: number) => (text.length <= length ? text : `${text.slice(0, length)}...`)

const formatTime = (value: string) => {
  const date = new Date(value)
  const diff = Date.now() - date.getTime()
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  return date.toLocaleDateString('zh-CN')
}

const startPolling = () => {
  stopPolling()
  pollTimer = window.setInterval(() => {
    if (videos.value.some((item) => item.status === 'processing')) {
      loadVideos()
    }
  }, 10000)
}

const stopPolling = () => {
  if (pollTimer) {
    window.clearInterval(pollTimer)
    pollTimer = null
  }
}

onMounted(() => {
  const dramaId = route.query.drama_id as string | undefined
  if (dramaId) filters.drama_id = dramaId
  loadDramas()
  loadVideos()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.studio-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding-bottom: 32px;
}

.page-title {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.page-title h1 {
  margin: 0;
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: clamp(1.8rem, 3vw, 2.6rem);
  letter-spacing: -0.06em;
}

.page-title p {
  max-width: 760px;
  margin: 0;
  color: var(--text-muted);
  line-height: 1.7;
}

.page-kicker {
  color: var(--text-muted);
  font-size: 0.76rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.hero-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(300px, 0.9fr);
  gap: 18px;
}

.hero-copy {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-copy p {
  color: var(--text-secondary);
  line-height: 1.75;
}

.workflow-list {
  display: grid;
  gap: 10px;
}

.workflow-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  color: var(--text-secondary);
}

.workflow-index {
  display: inline-flex;
  width: 28px;
  height: 28px;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: var(--accent-light);
  color: var(--text-primary);
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 0.8rem;
  font-weight: 700;
}

.stat-grid {
  display: grid;
  gap: 12px;
}

.filter-shell {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.filter-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.active-summary {
  display: flex;
  min-height: 54px;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  padding: 10px 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  background: rgba(255, 255, 255, 0.02);
}

.active-summary strong {
  font-weight: 700;
}

.active-summary span {
  color: var(--text-muted);
  font-size: 0.84rem;
}

.filter-actions {
  display: flex;
  gap: 10px;
}

.asset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 18px;
}

.asset-card {
  overflow: hidden;
}

.video-media {
  position: relative;
  aspect-ratio: 16 / 9;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-secondary);
}

.video-player {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.asset-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-muted);
}

.asset-placeholder .el-icon {
  font-size: 42px;
}

.asset-placeholder.processing {
  color: var(--warning);
}

.asset-placeholder.failed {
  color: var(--error);
}

.asset-placeholder small {
  color: var(--text-muted);
}

.spin {
  animation: spin 1s linear infinite;
}

.asset-badges {
  position: absolute;
  top: 14px;
  left: 14px;
  display: flex;
  gap: 8px;
}

.asset-content {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 18px;
}

.asset-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.asset-header strong {
  font-family: 'Space Grotesk', 'Noto Sans SC', sans-serif;
  font-size: 1rem;
  letter-spacing: -0.03em;
}

.asset-header span,
.asset-meta,
.asset-extra {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.asset-prompt {
  min-height: 72px;
  color: var(--text-secondary);
  line-height: 1.75;
}

.asset-meta,
.asset-extra {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.asset-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.studio-empty {
  padding: 60px 0 20px;
}

.pagination-shell {
  display: flex;
  justify-content: center;
  margin-top: 6px;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 980px) {
  .hero-grid,
  .filter-grid {
    grid-template-columns: 1fr;
  }
}
</style>
