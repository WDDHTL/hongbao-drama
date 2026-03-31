<template>
  <div class="studio-page">
    <AppHeader :show-logo="false">
      <template #left>
        <div class="page-title">
          <span class="page-kicker">Image Studio</span>
          <h1>图片资产中心</h1>
          <p>保留黑白基调，把角色图、场景图和连续参考图统一收进一套可复用组件系统。</p>
        </div>
      </template>
      <template #right>
        <el-button type="primary" @click="showGenerateDialog = true">
          <el-icon><Plus /></el-icon>
          <span>新建图片任务</span>
        </el-button>
      </template>
    </AppHeader>

    <section class="hero-grid">
      <BaseCard title="连续性工作流" subtitle="先做稳定视觉锚点，再把图片资产送入视频首尾帧和多参考流程。">
        <div class="hero-copy">
          <p>在图片弹窗里选择已完成图片作为连续性参考，可以锁定人物、服装、场景和镜头关系，减少后续视频漂移。</p>
          <div class="workflow-list">
            <div class="workflow-item">
              <span class="workflow-index">01</span>
              <span>先产出角色主视觉和场景母图</span>
            </div>
            <div class="workflow-item">
              <span class="workflow-index">02</span>
              <span>再用多张参考图固化细节</span>
            </div>
            <div class="workflow-item">
              <span class="workflow-index">03</span>
              <span>把完成图片直接送到视频首尾帧</span>
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

    <BaseCard title="筛选与状态" subtitle="按剧目和状态过滤图片资产。">
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
                <span>{{ activeDrama?.style || '当前展示所有图片资产' }}</span>
              </div>
            </el-form-item>
          </div>
        </el-form>

        <div class="filter-actions">
          <el-button type="primary" @click="loadImages">刷新列表</el-button>
          <el-button @click="resetFilters">重置筛选</el-button>
        </div>
      </div>
    </BaseCard>

    <section v-loading="loading" class="asset-section">
      <div v-if="images.length" class="asset-grid">
        <BaseCard v-for="image in images" :key="image.id" hoverable class="asset-card" no-padding>
          <div class="asset-media">
            <el-image
              v-if="image.status === 'completed' && image.image_url"
              :src="image.image_url"
              fit="cover"
              class="asset-image"
              :preview-src-list="[image.image_url]"
            >
              <template #error>
                <div class="asset-placeholder">
                  <el-icon><PictureFilled /></el-icon>
                  <span>图片加载失败</span>
                </div>
              </template>
            </el-image>

            <div v-else-if="image.status === 'processing'" class="asset-placeholder processing">
              <el-icon class="spin"><Loading /></el-icon>
              <span>生成中</span>
            </div>

            <div v-else-if="image.status === 'failed'" class="asset-placeholder failed">
              <el-icon><CircleClose /></el-icon>
              <span>任务失败</span>
            </div>

            <div v-else class="asset-placeholder">
              <el-icon><Picture /></el-icon>
              <span>等待生成</span>
            </div>

            <div class="asset-badges">
              <el-tag size="small" :type="getStatusType(image.status)">
                {{ getStatusText(image.status) }}
              </el-tag>
              <el-tag v-if="image.frame_type" size="small" effect="plain">{{ image.frame_type }}</el-tag>
            </div>
          </div>

          <div class="asset-content">
            <div class="asset-header">
              <strong>{{ image.model || image.provider }}</strong>
              <span>#{{ image.id }}</span>
            </div>

            <p class="asset-prompt">{{ truncateText(image.prompt, 96) }}</p>

            <div class="asset-meta">
              <span>{{ image.provider }}</span>
              <span>{{ image.size || '自动尺寸' }}</span>
              <span>{{ formatTime(image.created_at) }}</span>
            </div>

            <div class="asset-actions">
              <ActionButton :icon="View" tooltip="查看详情" @click="viewDetails(image)" />
              <ActionButton
                v-if="image.status === 'completed'"
                :icon="Download"
                tooltip="打开图片"
                @click="downloadImage(image)"
              />
              <el-popconfirm title="确认删除这条图片记录？" @confirm="deleteImage(image.id)">
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

      <el-empty v-else-if="!loading" description="当前没有图片资产，先创建一组角色图或场景图。" class="studio-empty" />
    </section>

    <el-pagination
      v-if="total > 0"
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.page_size"
      :total="total"
      :page-sizes="[12, 24, 36, 48]"
      layout="total, sizes, prev, pager, next, jumper"
      class="pagination-shell"
      @current-change="loadImages"
      @size-change="loadImages"
    />

    <GenerateImageDialog v-model="showGenerateDialog" :drama-id="filters.drama_id" @success="loadImages" />
    <ImageDetailDialog v-model="showDetailDialog" :image="selectedImage" />
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
  Picture,
  PictureFilled,
  Plus,
  View
} from '@element-plus/icons-vue'
import { dramaAPI } from '@/api/drama'
import { imageAPI } from '@/api/image'
import { ActionButton, AppHeader, BaseCard, StatCard } from '@/components/common'
import type { Drama } from '@/types/drama'
import type { ImageGeneration, ImageStatus } from '@/types/image'
import GenerateImageDialog from './components/GenerateImageDialog.vue'
import ImageDetailDialog from './components/ImageDetailDialog.vue'

const route = useRoute()

const loading = ref(false)
const total = ref(0)
const images = ref<ImageGeneration[]>([])
const dramas = ref<Drama[]>([])
const showGenerateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedImage = ref<ImageGeneration>()
let pollTimer: number | null = null

const filters = reactive({
  drama_id: undefined as string | undefined,
  status: undefined as ImageStatus | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 12
})

const completedCount = computed(() => images.value.filter((item) => item.status === 'completed').length)
const processingCount = computed(() => images.value.filter((item) => item.status === 'processing').length)
const failedCount = computed(() => images.value.filter((item) => item.status === 'failed').length)
const activeDrama = computed(() => dramas.value.find((item) => item.id === filters.drama_id))

const loadImages = async () => {
  loading.value = true
  try {
    const result = await imageAPI.listImages({
      drama_id: filters.drama_id,
      status: filters.status,
      page: pagination.page,
      page_size: pagination.page_size
    })
    images.value = result.items
    total.value = result.pagination.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载图片列表失败')
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
  loadImages()
}

const viewDetails = (image: ImageGeneration) => {
  selectedImage.value = image
  showDetailDialog.value = true
}

const downloadImage = (image: ImageGeneration) => {
  if (!image.image_url) return
  window.open(image.image_url, '_blank')
}

const deleteImage = async (id: number) => {
  try {
    await imageAPI.deleteImage(id)
    ElMessage.success('图片记录已删除')
    loadImages()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

const truncateText = (text: string, length: number) => (text.length <= length ? text : `${text.slice(0, length)}...`)

const getStatusType = (status: ImageStatus) => {
  const map: Record<ImageStatus, 'info' | 'warning' | 'success' | 'danger'> = {
    pending: 'info',
    processing: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return map[status]
}

const getStatusText = (status: ImageStatus) => {
  const map: Record<ImageStatus, string> = {
    pending: '待开始',
    processing: '处理中',
    completed: '已完成',
    failed: '失败'
  }
  return map[status]
}

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
    if (images.value.some((item) => item.status === 'processing')) {
      loadImages()
    }
  }, 5000)
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
  loadImages()
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
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 18px;
}

.asset-card {
  overflow: hidden;
}

.asset-media {
  position: relative;
  aspect-ratio: 1 / 1;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-secondary);
}

.asset-image,
.asset-image :deep(img) {
  width: 100%;
  height: 100%;
}

.asset-image :deep(img) {
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
.asset-meta {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.asset-prompt {
  min-height: 68px;
  color: var(--text-secondary);
  line-height: 1.75;
}

.asset-meta {
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
