<template>
  <div class="page-container">
    <div class="content-wrapper animate-fade-in">
      <!-- Page Header / 页面头部 -->
      <AppHeader :fixed="false" :show-logo="false">
        <template #left>
          <el-button text @click="goBack" class="back-btn">
            <el-icon><ArrowLeft /></el-icon>
            <span>{{ $t("common.back") }}</span>
          </el-button>
          <div class="page-title">
            <h1>{{ drama?.title || "" }}</h1>
            <span class="subtitle">{{
              drama?.description || $t("drama.management.overview")
            }}</span>
          </div>
        </template>
        <template #right>
          <el-button type="primary" class="action-btn" @click="goToVideoIdeaTab">
            <el-icon><VideoCamera /></el-icon>
            <span>{{ $t("drama.management.videoIdea") }}</span>
          </el-button>
          <el-button type="primary" plain @click="goToAutoWorkflow">
            自动工作流
          </el-button>
        </template>
      </AppHeader>

      <!-- Tabs / 标签页 -->
      <div class="tabs-wrapper">
        <el-tabs v-model="activeTab" class="management-tabs">
          <!-- 项目概览 -->
          <el-tab-pane :label="$t('drama.management.overview')" name="overview">
            <div class="stats-grid">
              <StatCard
                :label="$t('drama.management.episodeStats')"
                :value="episodesCount"
                :icon="Document"
                icon-color="var(--accent)"
                icon-bg="var(--accent-light)"
                value-color="var(--accent)"
                :description="$t('drama.management.episodesCreated')"
              />
              <StatCard
                :label="$t('drama.management.characterStats')"
                :value="charactersCount"
                :icon="User"
                icon-color="var(--success)"
                icon-bg="var(--success-light)"
                value-color="var(--success)"
                :description="$t('drama.management.charactersCreated')"
              />
              <StatCard
                :label="$t('drama.management.sceneStats')"
                :value="scenesCount"
                :icon="Picture"
                icon-color="var(--warning)"
                icon-bg="var(--warning-light)"
                value-color="var(--warning)"
                :description="$t('drama.management.sceneLibraryCount')"
              />
              <StatCard
                :label="$t('drama.management.propStats')"
                :value="propsCount"
                :icon="Box"
                icon-color="var(--primary)"
                icon-bg="var(--primary-light)"
                value-color="var(--primary)"
                :description="$t('drama.management.propsCreated')"
              />
            </div>

            <!-- 引导卡片：无章节时显示 -->
            <el-alert
              v-if="episodesCount === 0"
              :title="$t('drama.management.startFirstEpisode')"
              type="info"
              :closable="false"
              style="margin-top: 20px"
            >
              <template #default>
                <p style="margin: 8px 0">
                  {{ $t("drama.management.noEpisodesYet") }}
                </p>
                <el-button
                  type="primary"
                  :icon="Plus"
                  @click="createNewEpisode"
                  style="margin-top: 8px"
                >
                  {{ $t("drama.management.createFirstEpisode") }}
                </el-button>
              </template>
            </el-alert>

            <el-card shadow="never" class="project-info-card">
              <template #header>
                <div class="card-header">
                  <h3 class="card-title">
                    {{ $t("drama.management.projectInfo") }}
                  </h3>
                  <el-tag :type="getStatusType(drama?.status)" size="small">{{
                    getStatusText(drama?.status)
                  }}</el-tag>
                </div>
              </template>
              <el-descriptions :column="2" border class="project-descriptions">
                <el-descriptions-item
                  :label="$t('drama.management.projectName')"
                >
                  <span class="info-value">{{ drama?.title }}</span>
                </el-descriptions-item>
                <el-descriptions-item :label="$t('common.createdAt')">
                  <span class="info-value">{{
                    formatDate(drama?.created_at)
                  }}</span>
                </el-descriptions-item>
                <el-descriptions-item
                  :label="$t('drama.management.projectDesc')"
                  :span="2"
                >
                  <span class="info-desc">{{
                    drama?.description || $t("drama.management.noDescription")
                  }}</span>
                </el-descriptions-item>
              </el-descriptions>
            </el-card>
          </el-tab-pane>

          <!-- 对话生成视频 -->
          <el-tab-pane :label="$t('drama.management.videoIdea')" name="videoIdea">
            <div class="video-idea-wrapper">
              <div class="video-idea-hero">
                <div>
                  <div class="video-idea-kicker">{{ $t("drama.management.videoIdeaKicker") }}</div>
                  <h2 class="video-idea-title">{{ $t("drama.management.videoIdeaTitle") }}</h2>
                  <p class="video-idea-desc">
                    {{ $t("drama.management.videoIdeaGuideDesc") }}
                  </p>
                </div>
                <div class="video-idea-tags">
                  <el-tag type="success" effect="dark">{{ $t("drama.management.videoIdeaTagScript") }}</el-tag>
                  <el-tag type="warning" effect="dark">{{ $t("drama.management.videoIdeaTagAuto") }}</el-tag>
                  <el-tag type="info" effect="dark">{{ $t("drama.management.videoIdeaTagWorkflow") }}</el-tag>
                </div>
              </div>

              <el-card shadow="never" class="video-idea-card">
                <template #header>
                  <div class="card-header">
                    <h3 class="card-title">{{ $t("drama.management.videoIdeaActionTitle") }}</h3>
                    <span class="video-idea-subtitle">{{ $t("drama.management.videoIdeaActionSubTitle") }}</span>
                  </div>
                </template>
                <div class="video-idea-form">
                  <el-input
                    v-model="videoIdeaInput"
                    type="textarea"
                    :rows="8"
                    :placeholder="$t('drama.management.videoIdeaPlaceholder')"
                  />
                  <div class="video-idea-actions">
                    <el-button
                      type="primary"
                      size="large"
                      :loading="videoIdeaSubmitting"
                      @click="submitVideoIdea"
                    >
                      {{ $t("drama.management.videoIdeaSubmit") }}
                    </el-button>
                    <el-button link type="primary" @click="goToAutoWorkflow">
                      {{ $t("drama.management.videoIdeaWorkflowLink") }}
                    </el-button>
                  </div>
                  <div class="video-idea-tip">
                    <el-icon><VideoCamera /></el-icon>
                    <span>{{ $t("drama.management.videoIdeaTip") }}</span>
                  </div>
                  <div v-if="videoIdeaResponse" class="video-idea-response">
                    <el-alert
                      :title="videoIdeaResponse"
                      :type="videoIdeaResponse === '做不到' ? 'warning' : 'success'"
                      :closable="false"
                      show-icon
                    />
                  </div>
                </div>
              </el-card>
            </div>
          </el-tab-pane>

          <!-- 章节管理 -->
          <el-tab-pane :label="$t('drama.management.episodes')" name="episodes">
            <div class="tab-header">
              <h2>{{ $t("drama.management.episodeList") }}</h2>
              <el-button
                type="primary"
                :icon="Plus"
                @click="createNewEpisode"
                >{{ $t("drama.management.createNewEpisode") }}</el-button
              >
            </div>

            <!-- 空状态引导 -->
            <el-empty
              v-if="episodesCount === 0"
              :description="$t('drama.management.noEpisodes')"
              style="margin-top: 40px"
            >
              <template #image>
                <el-icon :size="80" class="empty-icon"><Document /></el-icon>
              </template>
              <el-button type="primary" :icon="Plus" @click="createNewEpisode">
                {{ $t("drama.management.createFirstEpisode") }}
              </el-button>
            </el-empty>

            <el-table
              v-else
              :data="sortedEpisodes"
              border
              stripe
              style="margin-top: 16px"
            >
              <el-table-column
                type="index"
                :label="$t('storyboard.table.number')"
                width="80"
              />
              <el-table-column
                prop="title"
                :label="$t('drama.management.episodeList')"
                min-width="200"
              />
              <el-table-column :label="$t('common.status')" width="120">
                <template #default="{ row }">
                  <el-tag :type="getEpisodeStatusType(row)">{{
                    getEpisodeStatusText(row)
                  }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="Shots" width="100">
                <template #default="{ row }">
                  {{ row.shots?.length || 0 }}
                </template>
              </el-table-column>
              <el-table-column :label="$t('common.createdAt')" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('storyboard.table.operations')"
                width="220"
                fixed="right"
              >
                <template #default="{ row }">
                  <el-button
                    size="small"
                    type="primary"
                    @click="enterEpisodeWorkflow(row)"
                  >
                    {{ $t("drama.management.goToEdit") }}
                  </el-button>
                  <el-button
                    size="small"
                    type="danger"
                    @click="deleteEpisode(row)"
                  >
                    {{ $t("common.delete") }}
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <!-- 角色管理 -->
          <el-tab-pane
            :label="$t('drama.management.characters')"
            name="characters"
          >
            <div class="tab-header">
              <h2>{{ $t("drama.management.characterList") }}</h2>
              <div style="display: flex; gap: 10px">
                <el-button
                  :icon="Document"
                  @click="openExtractCharacterDialog"
                  >{{ $t("prop.extract") }}</el-button
                >
                <el-button
                  type="primary"
                  :icon="Plus"
                  @click="openAddCharacterDialog"
                  >{{ $t("character.add") }}</el-button
                >
              </div>
            </div>

            <div class="character-pipeline dark-character-pipeline">
              <div class="character-guide-panel dark-character-guide">
                <div class="generation-guide-title">角色基准图工作台</div>
                <p class="generation-guide-text">
                  每个角色只生成一张白底三视图合成图，单图内同时包含侧面、正面、背面，后续镜头和视频默认都从这张图取参考。
                </p>
                <div class="character-guide-meta">
                  <span>白底输出</span>
                  <span>三视图合成</span>
                  <span>角色锁定</span>
                </div>
              </div>
              <div class="character-guide-panel is-compact dark-character-guide">
                <div class="character-guide-stat">
                  <strong>{{ charactersCount }}</strong>
                  <span>角色总数</span>
                </div>
                <div class="character-guide-stat">
                  <strong>1</strong>
                  <span>单图基准</span>
                </div>
                <div class="character-guide-stat">
                  <strong>{{ generatedCharacterCount }}</strong>
                  <span>已生成</span>
                </div>
              </div>
            </div>

            <div class="generation-guide legacy-generation-guide">
              <div>
                <div class="generation-guide-title">角色图生成流程</div>
                <p class="generation-guide-text">
                  提交后会显示任务提交、AI 生成、完成回填三个阶段，不需要手动刷新。
                </p>
              </div>
              <div class="generation-guide-steps">
                <span>提交任务</span>
                <span>AI 生成</span>
                <span>回填卡片</span>
              </div>
            </div>

            <el-row class="character-grid character-grid-dark" :gutter="20" style="margin-top: 20px">
              <el-col
                :xs="24"
                :sm="12"
                :xl="8"
                v-for="character in drama?.characters"
                :key="character.id"
              >
                <el-card
                  shadow="hover"
                  :class="[
                    'character-card',
                    'character-card-dark',
                    getGenerationCardClass('character', character),
                  ]"
                >
                  <div class="character-dark-card">
                    <div class="character-dark-head">
                      <button
                        class="character-dark-avatar"
                        type="button"
                        @click="openCharacterPreview(character)"
                      >
                        <img
                          v-if="getCharacterStageImage(character)"
                          :src="getCharacterStageImage(character)"
                          :alt="character.name"
                          class="character-dark-avatar-image"
                        />
                        <div v-else class="character-dark-avatar-placeholder">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </button>

                      <div class="character-dark-copy">
                        <div class="character-dark-title-row">
                          <h4>{{ character.name }}</h4>
                          <el-button
                            class="character-card-icon-btn"
                            text
                            :icon="Delete"
                            @click="deleteCharacter(character)"
                          />
                        </div>

                        <div class="character-dark-tags">
                          <span
                            v-for="tag in getCharacterBadgeList(character)"
                            :key="`${character.id}-${tag}`"
                            class="character-dark-tag"
                          >
                            {{ tag }}
                          </span>
                        </div>

                        <p class="character-dark-summary">
                          {{ getCharacterCardSummary(character) }}
                        </p>
                      </div>
                    </div>

                    <button
                      class="character-sheet-panel"
                      :class="{
                        'is-empty': !getCharacterStageImage(character),
                        'is-loading': isCharacterGenerating(character),
                      }"
                      type="button"
                      @click="openCharacterPreview(character)"
                    >
                      <img
                        v-if="getCharacterStageImage(character)"
                        :src="getCharacterStageImage(character)"
                        :alt="`${character.name} 三视图基准图`"
                        class="character-sheet-image"
                      />
                      <div
                        v-else-if="isCharacterGenerating(character)"
                        class="character-sheet-loading"
                      >
                        <span class="character-loading-spinner" />
                        <span>正在生成三视图合成图</span>
                      </div>
                      <div v-else class="character-sheet-placeholder">
                        <el-icon><Picture /></el-icon>
                        <span>待生成三视图基准图</span>
                      </div>
                    </button>

                    <div class="character-dark-footer">
                      <div class="character-dark-status">
                        <span class="character-dark-status-label">
                          {{ getCharacterStatusLabel(character) }}
                        </span>
                        <p>{{ getCharacterStatusText(character) }}</p>
                      </div>

                      <div class="character-dark-actions">
                        <el-button
                          class="character-action-primary"
                          type="primary"
                          :loading="isGenerationBusy('character', character)"
                          :disabled="isGenerationBusy('character', character)"
                          @click="submitCharacterImageGeneration(character)"
                        >
                          {{ getCharacterActionText(character) }}
                        </el-button>
                        <el-button
                          class="character-action-secondary"
                          @click="openCharacterPreview(character)"
                        >
                          放大
                        </el-button>
                        <el-button
                          class="character-action-secondary"
                          @click="editCharacter(character)"
                        >
                          {{ $t("common.edit") }}
                        </el-button>
                      </div>
                    </div>
                  </div>

                  <div
                    v-if="false"
                    class="character-reference-strip"
                  >
                    <div
                      v-for="(image, index) in getCharacterReferenceImages(character)"
                      :key="`${character.id}-${index}-${image}`"
                      class="character-reference-item"
                    >
                      <ImagePreview
                        :image-url="image"
                        :alt="`${character.name}-${characterReferenceLabels[index] || '基准图'}`"
                        :size="56"
                        :show-placeholder-text="false"
                      />
                      <span>{{
                        characterReferenceLabels[index] || `基准图${index + 1}`
                      }}</span>
                    </div>
                  </div>

                  <div v-if="false" class="generation-process">
                    <div class="generation-process-head">
                      <span class="generation-process-title">图片生成过程</span>
                      <el-tag
                        size="small"
                        effect="plain"
                        :type="
                          getGenerationTagType(
                            getGenerationState('character', character).phase,
                          )
                        "
                      >
                        {{ getGenerationLabel("character", character) }}
                      </el-tag>
                    </div>
                    <p class="generation-process-detail">
                      {{ getGenerationDetail("character", character) }}
                    </p>
                    <el-progress
                      :percentage="
                        getGenerationProgress(
                          getGenerationState('character', character).phase,
                        )
                      "
                      :stroke-width="6"
                      :show-text="false"
                      class="generation-progress"
                    />
                    <div class="generation-step-row">
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('character', character).phase,
                            ) >= 20,
                        }"
                        >已提交</span
                      >
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('character', character).phase,
                            ) >= 60,
                        }"
                        >生成中</span
                      >
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('character', character).phase,
                            ) >= 100,
                        }"
                        >已完成</span
                      >
                    </div>
                  </div>

                  <div v-if="false" class="character-actions">
                    <el-button size="small" @click="editCharacter(character)">{{
                      $t("common.edit")
                    }}</el-button>
                    <el-button
                      size="small"
                      :loading="isGenerationBusy('character', character)"
                      :disabled="isGenerationBusy('character', character)"
                      @click="submitCharacterImageGeneration(character)"
                      >{{
                        getGenerationButtonText("character", character)
                      }}</el-button
                    >
                    <el-button
                      size="small"
                      type="danger"
                      @click="deleteCharacter(character)"
                      >{{ $t("common.delete") }}</el-button
                    >
                  </div>
                </el-card>
              </el-col>
            </el-row>

            <el-empty
              v-if="!drama?.characters || drama.characters.length === 0"
              :description="$t('drama.management.noCharacters')"
            />
          </el-tab-pane>

          <!-- 场景库管理 -->
          <el-tab-pane :label="$t('drama.management.sceneList')" name="scenes">
            <div class="tab-header">
              <h2>{{ $t("drama.management.sceneList") }}</h2>
            </div>

            <div class="generation-guide">
              <div>
                <div class="generation-guide-title">场景图生成流程</div>
                <p class="generation-guide-text">
                  点击生成后会显示提交、处理中、完成三个阶段，页面会自动轮询结果。
                </p>
              </div>
              <div class="generation-guide-steps">
                <span>提交任务</span>
                <span>AI 生成</span>
                <span>入库可用</span>
              </div>
            </div>

            <el-row :gutter="16" style="margin-top: 16px">
              <el-col :span="6" v-for="scene in scenes" :key="scene.id">
                <el-card
                  shadow="hover"
                  :class="['scene-card', getGenerationCardClass('scene', scene)]"
                >
                  <div class="scene-preview">
                    <ImagePreview
                      :image-url="getImageUrl(scene)"
                      :alt="scene.location + ' - ' + scene.time"
                      :size="120"
                      :show-placeholder-text="false"
                    />
                  </div>

                  <div class="scene-info">
                    <h4>{{ scene.name }}</h4>
                    <p class="desc">{{ scene.description }}</p>
                  </div>

                  <div class="generation-process">
                    <div class="generation-process-head">
                      <span class="generation-process-title">图片生成过程</span>
                      <el-tag
                        size="small"
                        effect="plain"
                        :type="
                          getGenerationTagType(
                            getGenerationState('scene', scene).phase,
                          )
                        "
                      >
                        {{ getGenerationLabel("scene", scene) }}
                      </el-tag>
                    </div>
                    <p class="generation-process-detail">
                      {{ getGenerationDetail("scene", scene) }}
                    </p>
                    <el-progress
                      :percentage="
                        getGenerationProgress(
                          getGenerationState('scene', scene).phase,
                        )
                      "
                      :stroke-width="6"
                      :show-text="false"
                      class="generation-progress"
                    />
                    <div class="generation-step-row">
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('scene', scene).phase,
                            ) >= 20,
                        }"
                        >已提交</span
                      >
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('scene', scene).phase,
                            ) >= 60,
                        }"
                        >生成中</span
                      >
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('scene', scene).phase,
                            ) >= 100,
                        }"
                        >已完成</span
                      >
                    </div>
                  </div>

                  <div class="scene-actions">
                    <el-button size="small" @click="editScene(scene)">{{
                      $t("common.edit")
                    }}</el-button>
                    <el-button
                      size="small"
                      :loading="isGenerationBusy('scene', scene)"
                      :disabled="isGenerationBusy('scene', scene)"
                      @click="submitSceneImageGeneration(scene)"
                      >{{ getGenerationButtonText("scene", scene) }}</el-button
                    >
                    <el-button
                      size="small"
                      type="danger"
                      @click="deleteScene(scene)"
                      >{{ $t("common.delete") }}</el-button
                    >
                  </div>
                </el-card>
              </el-col>
            </el-row>

            <el-empty
              v-if="scenes.length === 0"
              :description="$t('drama.management.noScenes')"
            />
          </el-tab-pane>

          <!-- 道具管理 -->
          <el-tab-pane :label="$t('drama.management.propList')" name="props">
            <div class="tab-header">
              <h2>{{ $t("drama.management.propList") }}</h2>
              <div style="display: flex; gap: 10px">
                <el-button :icon="Document" @click="openExtractDialog">{{
                  $t("prop.extract")
                }}</el-button>
                <el-button
                  type="primary"
                  :icon="Plus"
                  @click="openAddPropDialog"
                  >{{ $t("common.add") }}</el-button
                >
              </div>
            </div>

            <div class="generation-guide">
              <div>
                <div class="generation-guide-title">道具图生成流程</div>
                <p class="generation-guide-text">
                  生成后会持续检查结果并自动回填，让你能直接看到当前进度。
                </p>
              </div>
              <div class="generation-guide-steps">
                <span>提交任务</span>
                <span>等待结果</span>
                <span>回填卡片</span>
              </div>
            </div>

            <el-row :gutter="16" style="margin-top: 16px">
              <el-col :span="6" v-for="prop in drama?.props" :key="prop.id">
                <el-card
                  shadow="hover"
                  :class="['scene-card', getGenerationCardClass('prop', prop)]"
                >
                  <div class="scene-preview">
                    <ImagePreview
                      :image-url="getImageUrl(prop)"
                      :alt="prop.name"
                      :size="120"
                      :show-placeholder-text="false"
                    />
                  </div>

                  <div class="scene-info">
                    <h4>{{ prop.name }}</h4>
                    <el-tag size="small" v-if="prop.type">{{
                      prop.type
                    }}</el-tag>
                    <p class="desc">{{ prop.description || prop.prompt }}</p>
                  </div>

                  <div class="generation-process">
                    <div class="generation-process-head">
                      <span class="generation-process-title">图片生成过程</span>
                      <el-tag
                        size="small"
                        effect="plain"
                        :type="
                          getGenerationTagType(
                            getGenerationState('prop', prop).phase,
                          )
                        "
                      >
                        {{ getGenerationLabel("prop", prop) }}
                      </el-tag>
                    </div>
                    <p class="generation-process-detail">
                      {{ getGenerationDetail("prop", prop) }}
                    </p>
                    <el-progress
                      :percentage="
                        getGenerationProgress(
                          getGenerationState('prop', prop).phase,
                        )
                      "
                      :stroke-width="6"
                      :show-text="false"
                      class="generation-progress"
                    />
                    <div class="generation-step-row">
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('prop', prop).phase,
                            ) >= 20,
                        }"
                        >已提交</span
                      >
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('prop', prop).phase,
                            ) >= 60,
                        }"
                        >生成中</span
                      >
                      <span
                        :class="{
                          active:
                            getGenerationProgress(
                              getGenerationState('prop', prop).phase,
                            ) >= 100,
                        }"
                        >已完成</span
                      >
                    </div>
                  </div>

                  <div class="scene-actions">
                    <el-button size="small" @click="editProp(prop)">{{
                      $t("common.edit")
                    }}</el-button>
                    <el-button
                      size="small"
                      :loading="isGenerationBusy('prop', prop)"
                      :disabled="!prop.prompt || isGenerationBusy('prop', prop)"
                      @click="submitPropImageGeneration(prop)"
                      >{{ getGenerationButtonText("prop", prop) }}</el-button
                    >
                    <el-button
                      size="small"
                      type="danger"
                      @click="deleteProp(prop)"
                      >{{ $t("common.delete") }}</el-button
                    >
                  </div>
                </el-card>
              </el-col>
            </el-row>

            <el-empty
              v-if="!drama?.props || drama.props.length === 0"
              :description="$t('drama.management.noProps')"
            />
          </el-tab-pane>
        </el-tabs>
      </div>


      <!-- 添加/编辑角色对话框 -->
      <el-dialog
        v-model="addCharacterDialogVisible"
        :title="editingCharacter ? $t('character.edit') : $t('character.add')"
        width="600px"
      >
        <el-form :model="newCharacter" label-width="100px">
          <el-form-item :label="$t('character.image')">
            <el-upload
              class="avatar-uploader"
              :action="`/api/v1/upload/image`"
              :show-file-list="false"
              :on-success="handleCharacterAvatarSuccess"
              :before-upload="beforeAvatarUpload"
            >
              <img
                v-if="hasImage(newCharacter)"
                :src="getImageUrl(newCharacter)"
                class="avatar"
                style="width: 100px; height: 100px; object-fit: cover"
              />
              <el-icon
                v-else
                class="avatar-uploader-icon"
                style="
                  border: 1px dashed #d9d9d9;
                  border-radius: 6px;
                  cursor: pointer;
                  position: relative;
                  overflow: hidden;
                  width: 100px;
                  height: 100px;
                  font-size: 28px;
                  color: #8c939d;
                  text-align: center;
                  line-height: 100px;
                "
                ><Plus
              /></el-icon>
            </el-upload>
          </el-form-item>
          <el-form-item :label="$t('character.name')">
            <el-input
              v-model="newCharacter.name"
              :placeholder="$t('character.name')"
            />
          </el-form-item>
          <el-form-item :label="$t('character.role')">
            <el-select
              v-model="newCharacter.role"
              :placeholder="$t('common.pleaseSelect')"
            >
              <el-option label="Main" value="main" />
              <el-option label="Supporting" value="supporting" />
              <el-option label="Minor" value="minor" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('character.appearance')">
            <el-input
              v-model="newCharacter.appearance"
              type="textarea"
              :rows="3"
              :placeholder="$t('character.appearance')"
            />
          </el-form-item>
          <el-form-item :label="$t('character.personality')">
            <el-input
              v-model="newCharacter.personality"
              type="textarea"
              :rows="3"
              :placeholder="$t('character.personality')"
            />
          </el-form-item>
          <el-form-item :label="$t('character.description')">
            <el-input
              v-model="newCharacter.description"
              type="textarea"
              :rows="3"
              :placeholder="$t('common.description')"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="addCharacterDialogVisible = false">{{
            $t("common.cancel")
          }}</el-button>
          <el-button type="primary" @click="saveCharacter">{{
            $t("common.confirm")
          }}</el-button>
        </template>
      </el-dialog>

      <!-- 添加/编辑场景对话框 -->
      <el-dialog
        v-model="addSceneDialogVisible"
        :title="editingScene ? $t('common.edit') : $t('common.add')"
        width="600px"
      >
        <el-form :model="newScene" label-width="100px">
          <el-form-item :label="$t('common.image')">
            <el-upload
              class="avatar-uploader"
              :action="`/api/v1/upload/image`"
              :show-file-list="false"
              :on-success="handleSceneImageSuccess"
              :before-upload="beforeAvatarUpload"
            >
              <img
                v-if="hasImage(newScene)"
                :src="getImageUrl(newScene)"
                class="avatar"
                style="width: 160px; height: 90px; object-fit: cover"
              />
              <el-icon
                v-else
                class="avatar-uploader-icon"
                style="
                  border: 1px dashed #d9d9d9;
                  border-radius: 6px;
                  cursor: pointer;
                  position: relative;
                  overflow: hidden;
                  width: 160px;
                  height: 90px;
                  font-size: 28px;
                  color: #8c939d;
                  text-align: center;
                  line-height: 90px;
                "
                ><Plus
              /></el-icon>
            </el-upload>
          </el-form-item>
          <el-form-item :label="$t('common.name')">
            <el-input
              v-model="newScene.location"
              :placeholder="$t('common.name')"
            />
          </el-form-item>
          <el-form-item :label="$t('common.description')">
            <el-input
              v-model="newScene.prompt"
              type="textarea"
              :rows="4"
              :placeholder="$t('common.description')"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="addSceneDialogVisible = false">{{
            $t("common.cancel")
          }}</el-button>
          <el-button type="primary" @click="saveScene">{{
            $t("common.confirm")
          }}</el-button>
        </template>
      </el-dialog>

      <!-- 添加/编辑道具对话框 -->
      <el-dialog
        v-model="addPropDialogVisible"
        :title="editingProp ? $t('common.edit') : $t('common.add')"
        width="600px"
      >
        <el-form :model="newProp" label-width="100px">
          <el-form-item :label="$t('common.image')">
            <el-upload
              class="avatar-uploader"
              :action="`/api/v1/upload/image`"
              :show-file-list="false"
              :on-success="handlePropImageSuccess"
              :before-upload="beforeAvatarUpload"
            >
              <img
                v-if="hasImage(newProp)"
                :src="getImageUrl(newProp)"
                class="avatar"
                style="width: 100px; height: 100px; object-fit: cover"
              />
              <el-icon
                v-else
                class="avatar-uploader-icon"
                style="
                  border: 1px dashed #d9d9d9;
                  border-radius: 6px;
                  cursor: pointer;
                  position: relative;
                  overflow: hidden;
                  width: 100px;
                  height: 100px;
                  font-size: 28px;
                  color: #8c939d;
                  text-align: center;
                  line-height: 100px;
                "
                ><Plus
              /></el-icon>
            </el-upload>
          </el-form-item>
          <el-form-item :label="$t('prop.name')">
            <el-input v-model="newProp.name" :placeholder="$t('prop.name')" />
          </el-form-item>
          <el-form-item :label="$t('prop.type')">
            <el-input
              v-model="newProp.type"
              :placeholder="$t('prop.typePlaceholder')"
            />
          </el-form-item>
          <el-form-item :label="$t('prop.description')">
            <el-input
              v-model="newProp.description"
              type="textarea"
              :rows="3"
              :placeholder="$t('prop.description')"
            />
          </el-form-item>
          <el-form-item :label="$t('prop.prompt')">
            <el-input
              v-model="newProp.prompt"
              type="textarea"
              :rows="3"
              :placeholder="$t('prop.promptPlaceholder')"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="addPropDialogVisible = false">{{
            $t("common.cancel")
          }}</el-button>
          <el-button type="primary" @click="saveProp">{{
            $t("common.confirm")
          }}</el-button>
        </template>
      </el-dialog>

      <!-- 从剧本提取道具对话框 -->
      <el-dialog
        v-model="extractPropsDialogVisible"
        :title="$t('prop.extractTitle')"
        width="500px"
      >
        <el-form label-width="100px">
          <el-form-item :label="$t('prop.selectEpisode')">
            <el-select
              v-model="selectedExtractEpisodeId"
              :placeholder="$t('common.pleaseSelect')"
              style="width: 100%"
            >
              <el-option
                v-for="ep in sortedEpisodes"
                :key="ep.id"
                :label="ep.title"
                :value="ep.id"
              />
            </el-select>
          </el-form-item>
          <el-alert
            :title="$t('prop.extractTip')"
            type="info"
            :closable="false"
            show-icon
          />
        </el-form>
        <template #footer>
          <el-button @click="extractPropsDialogVisible = false">{{
            $t("common.cancel")
          }}</el-button>
          <el-button
            type="primary"
            @click="handleExtractProps"
            :disabled="!selectedExtractEpisodeId"
            >{{ $t("prop.startExtract") }}</el-button
          >
        </template>
      </el-dialog>

      <!-- 从剧本提取角色对话框 -->
      <el-dialog
        v-model="extractCharactersDialogVisible"
        :title="$t('prop.extractTitle')"
        width="500px"
      >
        <el-form label-width="100px">
          <el-form-item :label="$t('prop.selectEpisode')">
            <el-select
              v-model="selectedExtractEpisodeId"
              :placeholder="$t('common.pleaseSelect')"
              style="width: 100%"
            >
              <el-option
                v-for="ep in sortedEpisodes"
                :key="ep.id"
                :label="ep.title"
                :value="ep.id"
              />
            </el-select>
          </el-form-item>
          <el-alert
            :title="$t('prop.extractTip')"
            type="info"
            :closable="false"
            show-icon
          />
          <el-alert
            v-if="extractingCharacters"
            :title="extractingCharactersMessage"
            type="success"
            :closable="false"
            show-icon
            style="margin-top: 12px"
          />
        </el-form>
        <template #footer>
          <el-button
            @click="extractCharactersDialogVisible = false"
            :disabled="extractingCharacters"
          >{{
            $t("common.cancel")
          }}</el-button>
          <el-button
            type="primary"
            @click="handleExtractCharacters"
            :disabled="!selectedExtractEpisodeId || extractingCharacters"
            :loading="extractingCharacters"
            >{{ $t("prop.startExtract") }}</el-button
          >
        </template>
      </el-dialog>

      <!-- 从剧本提取场景对话框 -->
      <el-dialog
        v-model="extractScenesDialogVisible"
        :title="$t('prop.extractTitle')"
        width="500px"
      >
        <el-form label-width="100px">
          <el-form-item :label="$t('prop.selectEpisode')">
            <el-select
              v-model="selectedExtractEpisodeId"
              :placeholder="$t('common.pleaseSelect')"
              style="width: 100%"
            >
              <el-option
                v-for="ep in sortedEpisodes"
                :key="ep.id"
                :label="ep.title"
                :value="ep.id"
              />
            </el-select>
          </el-form-item>
          <el-alert
            :title="$t('prop.extractTip')"
            type="info"
            :closable="false"
            show-icon
          />
        </el-form>
        <template #footer>
          <el-button @click="extractScenesDialogVisible = false">{{
            $t("common.cancel")
          }}</el-button>
          <el-button
            type="primary"
            @click="handleExtractScenes"
            :disabled="!selectedExtractEpisodeId"
            >{{ $t("prop.startExtract") }}</el-button
          >
        </template>
      </el-dialog>

      <el-dialog
        v-model="characterPreviewDialogVisible"
        :title="characterPreviewTitle"
        width="720px"
        align-center
      >
        <div class="character-preview-dialog">
          <img
            v-if="characterPreviewImage"
            :src="characterPreviewImage"
            :alt="characterPreviewTitle"
            class="character-preview-dialog-image"
          />
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, onUnmounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { useI18n } from "vue-i18n";
import {
  ArrowLeft,
  Document,
  User,
  Picture,
  Plus,
  Box,
  VideoCamera,
} from "@element-plus/icons-vue";
import { dramaAPI } from "@/api/drama";
import { characterLibraryAPI } from "@/api/character-library";
import { imageAPI } from "@/api/image";
import { propAPI } from "@/api/prop";
import { taskAPI } from "@/api/task";
import { workflowAPI } from "@/api/workflow";
import type { Drama } from "@/types/drama";
import {
  AppHeader,
  StatCard,
  EmptyState,
  ImagePreview,
} from "@/components/common";
import { fixImageUrl, getImageUrl, hasImage } from "@/utils/image";

const router = useRouter();
const route = useRoute();
const { t } = useI18n();

const drama = ref<Drama>();
const activeTab = ref((route.query.tab as string) || "overview");
const scenes = ref<any[]>([]);

type GenerationKind = "character" | "scene" | "prop";
type GenerationPhase = "idle" | "queued" | "processing" | "success" | "failed";

interface GenerationState {
  phase: GenerationPhase;
  detail: string;
}

interface CharacterSlotState extends GenerationState {
  imageGenId?: number;
}

const generationStates = ref<Record<string, GenerationState>>({});
const characterSlotStates = ref<Record<string, CharacterSlotState>>({});
const generationTimers = new Map<string, ReturnType<typeof setInterval>>();
let pollingTimer: ReturnType<typeof setInterval> | null = null;

const addCharacterDialogVisible = ref(false);
const addSceneDialogVisible = ref(false);
const addPropDialogVisible = ref(false);
const extractPropsDialogVisible = ref(false);
const extractCharactersDialogVisible = ref(false);
const extractScenesDialogVisible = ref(false);
const extractingCharacters = ref(false);
const extractingCharactersMessage = ref("正在提交角色提取任务...");
const characterPreviewDialogVisible = ref(false);
const characterPreviewImage = ref("");
const characterPreviewTitle = ref("");
const videoIdeaInput = ref("");
const videoIdeaResponse = ref("");
const videoIdeaSubmitting = ref(false);

const editingCharacter = ref<any>(null);
const editingScene = ref<any>(null);
const editingProp = ref<any>(null);
const selectedExtractEpisodeId = ref<number | null>(null);

const newCharacter = ref({
  name: "",
  role: "supporting",
  appearance: "",
  personality: "",
  description: "",
  image_url: "",
  local_path: "",
});

const newProp = ref({
  name: "",
  description: "",
  prompt: "",
  type: "",
  image_url: "",
  local_path: "",
});

const newScene = ref({
  location: "",
  prompt: "",
  image_url: "",
  local_path: "",
});

const episodesCount = computed(() => drama.value?.episodes?.length || 0);
const charactersCount = computed(() => drama.value?.characters?.length || 0);
const scenesCount = computed(() => scenes.value.length);
const propsCount = computed(() => drama.value?.props?.length || 0);

const sortedEpisodes = computed(() => {
  if (!drama.value?.episodes) return [];
  return [...drama.value.episodes].sort(
    (a, b) => a.episode_number - b.episode_number,
  );
});

const generatedCharacterCount = computed(
  () =>
    drama.value?.characters?.filter(
      (character) => getCharacterReferenceCount(character) > 0,
    ).length ?? 0,
);

const buildGenerationKey = (kind: GenerationKind, id?: string | number) =>
  `${kind}:${id ?? "unknown"}`;

const clearGenerationMonitor = (key?: string) => {
  if (!key) {
    generationTimers.forEach((timer) => clearInterval(timer));
    generationTimers.clear();
    return;
  }

  const timer = generationTimers.get(key);
  if (timer) {
    clearInterval(timer);
    generationTimers.delete(key);
  }
};

const setGenerationState = (
  kind: GenerationKind,
  id: string | number | undefined,
  state: Partial<GenerationState>,
) => {
  const key = buildGenerationKey(kind, id);
  generationStates.value[key] = {
    phase: state.phase || generationStates.value[key]?.phase || "idle",
    detail: state.detail || generationStates.value[key]?.detail || "",
  };
};

const getServerGenerationPhase = (status?: string): GenerationPhase | null => {
  const normalized = status?.toLowerCase();
  if (!normalized) return null;
  if (["pending", "queued", "submitted"].includes(normalized)) return "queued";
  if (["processing", "running", "generating"].includes(normalized)) {
    return "processing";
  }
  if (["completed", "success", "succeeded"].includes(normalized)) {
    return "success";
  }
  if (["failed", "error"].includes(normalized)) return "failed";
  return null;
};

const parseReferenceImages = (rawReferenceImages: any): string[] => {
  let parsedReferenceImages: any[] = [];

  if (typeof rawReferenceImages === "string") {
    try {
      const parsed = JSON.parse(rawReferenceImages);
      parsedReferenceImages = Array.isArray(parsed) ? parsed : [parsed];
    } catch {
      parsedReferenceImages = [rawReferenceImages];
    }
  } else if (Array.isArray(rawReferenceImages)) {
    parsedReferenceImages = rawReferenceImages;
  } else if (rawReferenceImages) {
    parsedReferenceImages = [rawReferenceImages];
  }

  return parsedReferenceImages
    .map((image) => (typeof image === "string" ? image.trim() : ""))
    .filter((image) => image.length > 0);
};

const characterFrameTypes = ["front", "side", "back"];

const buildCharacterSlotKey = (
  characterId: string | number | undefined,
  index: number,
) => `character-slot-${characterId ?? "unknown"}-${characterFrameTypes[index] ?? index}`;

const getCharacterFrameType = (index: number) =>
  characterFrameTypes[index] || "front";

const formatGenerationError = (rawError?: string) => {
  const error = (rawError || "").trim();
  if (!error) return "生成失败，请调整后重试";

  if (
    error.includes("我不能生成") ||
    error.includes("不能生成") ||
    error.includes("rix_api_error")
  ) {
    return "当前角色描述触发了图像模型限制，请把人物描述改得更中性、健康一些后再试";
  }

  if (error.includes("timeout")) {
    return "生成超时，请稍后重新提交";
  }

  return error.length > 80 ? `${error.slice(0, 80)}...` : error;
};

const getGenerationState = (kind: GenerationKind, item: any) => {
  const key = buildGenerationKey(kind, item?.id);
  const localState = generationStates.value[key];
  const serverPhase = getServerGenerationPhase(item?.image_generation_status);
  const completedViews =
    kind === "character" ? getCharacterReferenceCount(item) : 0;

  if (serverPhase === "failed") {
    return {
      phase: "failed" as GenerationPhase,
      detail: formatGenerationError(item?.image_generation_error),
    };
  }

  if (serverPhase === "queued") {
    return {
      phase: "queued" as GenerationPhase,
      detail:
        kind === "character"
        ? "任务已提交，准备生成正面、侧面、背面三张白底角色基准图"
        : "任务已提交，正在等待图像服务处理",
    };
  }

  if (serverPhase === "processing") {
    return {
      phase: "processing" as GenerationPhase,
      detail:
        kind === "character"
          ? "AI 正在生成三视角角色基准图，页面会自动刷新结果"
          : "AI 正在生成图片，页面会自动刷新结果",
    };
  }

  if (
    kind === "character" &&
    (localState?.phase === "queued" || localState?.phase === "processing") &&
    completedViews < 1
  ) {
    return {
      phase: "processing" as GenerationPhase,
      detail:
        completedViews > 0
          ? `正在生成角色三视图，已完成 ${completedViews}/3 张`
          : "正在生成角色三视图，正面、侧面、背面会依次回填",
    };
  }

  if (
    kind === "character" &&
    completedViews >= 1 &&
    (serverPhase === "success" || localState)
  ) {
    return {
      phase: "success" as GenerationPhase,
      detail: "三视图已生成完成，可直接用于后续镜头和视频参考",
    };
  }

  if (hasImage(item) && (serverPhase === "success" || localState)) {
    return {
      phase: "success" as GenerationPhase,
      detail: "图片已生成并回填，可直接继续后续流程",
    };
  }

  if (localState) {
    if (localState.phase === "success" && !hasImage(item)) {
      return {
        phase: "processing" as GenerationPhase,
        detail: "任务已完成，正在同步最新图片到页面",
      };
    }
    return localState;
  }

  return {
    phase: "idle" as GenerationPhase,
    detail: hasImage(item)
      ? kind === "character"
        ? "当前已有主图，可重新生成三张白底角色基准图"
        : "当前已有图片，可重新生成优化效果"
      : kind === "character"
        ? "点击后将生成正面、侧面、背面三张白底角色基准图"
        : "点击生成后会展示完整生成过程",
  };
};

const getGenerationCardClass = (kind: GenerationKind, item: any) => {
  const phase = getGenerationState(kind, item).phase;
  return {
    "is-generating": phase === "queued" || phase === "processing",
    "is-failed": phase === "failed",
  };
};

const getGenerationTagType = (phase: GenerationPhase) => {
  if (phase === "success") return "success";
  if (phase === "failed") return "danger";
  if (phase === "queued" || phase === "processing") return "warning";
  return "info";
};

const getGenerationProgress = (phase: GenerationPhase) => {
  if (phase === "queued") return 24;
  if (phase === "processing") return 68;
  if (phase === "success" || phase === "failed") return 100;
  return 0;
};

const getGenerationLabel = (kind: GenerationKind, item: any) => {
  const phase = getGenerationState(kind, item).phase;
  if (phase === "queued") return "已提交";
  if (phase === "processing") return "生成中";
  if (phase === "success") return "已完成";
  if (phase === "failed") return "失败";
  return hasImage(item) ? "可重生成" : "待生成";
};

const getGenerationDetail = (kind: GenerationKind, item: any) =>
  getGenerationState(kind, item).detail;

const characterReferenceLabels = ["正面", "侧面", "背面"];
const characterReferenceHints = [
  "半身主参考",
  "侧脸与轮廓",
  "背影与发型",
];
const characterGenerationSteps = ["提交任务", "生成三视图", "回填角色库"];

const getCharacterReferenceImages = (character: any) =>
  parseReferenceImages(character?.reference_images);

const getCharacterReferenceCount = (character: any) =>
  getCharacterReferenceImages(character).filter(
    (image) => typeof image === "string" && image.trim().length > 0,
  ).length;

const getCharacterResolvedReferenceImage = (character: any, index: number) => {
  const image = getCharacterReferenceImages(character)[index];
  return image ? fixImageUrl(image) : "";
};

const getCharacterSlotState = (character: any, index: number) =>
  characterSlotStates.value[buildCharacterSlotKey(character?.id, index)];

const getCharacterSlotDetail = (character: any, index: number) => {
  const slotState = getCharacterSlotState(character, index);
  if (slotState?.phase === "failed") {
    return slotState.detail;
  }
  if (slotState?.phase === "queued") {
    return "已提交";
  }
  if (slotState?.phase === "processing") {
    return "生成中";
  }
  return "";
};

const isCharacterSlotFailed = (character: any, index: number) =>
  getCharacterSlotState(character, index)?.phase === "failed";

const getCharacterStageImage = (character: any) => {
  const frontReference = getCharacterResolvedReferenceImage(character, 0);
  if (frontReference) return frontReference;
  return hasImage(character) ? getImageUrl(character) : "";
};

const isCharacterGenerating = (character: any) =>
  isGenerationBusy("character", character);

const isCharacterSlotLoading = (character: any, index: number) =>
  getCharacterSlotState(character, index)?.phase === "queued" ||
  getCharacterSlotState(character, index)?.phase === "processing" ||
  (isCharacterGenerating(character) &&
    !getCharacterResolvedReferenceImage(character, index));

const getCharacterSlotButtonText = (character: any, index: number) => {
  if (isCharacterSlotLoading(character, index)) return "生成中";
  return getCharacterResolvedReferenceImage(character, index)
    ? "重生成"
    : "补生成";
};

const openCharacterPreview = (character: any, index?: number) => {
  const targetImage =
    typeof index === "number"
      ? getCharacterResolvedReferenceImage(character, index)
      : getCharacterStageImage(character);

  if (!targetImage) return;

  characterPreviewImage.value = targetImage;
  characterPreviewTitle.value =
    typeof index === "number"
      ? `${character.name} · ${characterReferenceLabels[index]}`
      : `${character.name} · 主参考图`;
  characterPreviewDialogVisible.value = true;
};

const getCharacterGenerationProgress = (character: any) => {
  const phase = getGenerationState("character", character).phase;
  const completedViews = getCharacterReferenceCount(character);

  if (phase === "success" || phase === "failed") return 100;
  if (phase === "queued") return Math.max(14, completedViews * 24);
  if (phase === "processing") return Math.min(94, 26 + completedViews * 22);
  return getGenerationProgress(phase);
};

const getCharacterGenerationDetail = (character: any) => {
  const state = getGenerationState("character", character);
  const completedViews = getCharacterReferenceCount(character);

  if (state.phase === "queued") {
    return `三视图任务已提交，当前 ${completedViews}/3 张已回填`;
  }

  if (state.phase === "processing") {
    return completedViews > 0
      ? `正在生成角色三视图，已完成 ${completedViews}/3 张，剩余槽位会自动回填`
      : "正在生成角色三视图，正面、侧面、背面会依次回填";
  }

  if (state.phase === "success") {
    return "三视图已生成完成，可直接用于后续镜头和视频参考";
  }

  return state.detail;
};

const getCharacterRoleText = (role?: string) => {
  if (role === "main") return "Main";
  if (role === "supporting") return "Supporting";
  return "Minor";
};

const getCharacterBadgeList = (character: any) => {
  const badges = [getCharacterRoleText(character?.role)];
  const summaryText = [character?.appearance, character?.description]
    .filter((value) => typeof value === "string" && value.trim().length > 0)
    .join(" ");

  if (/性别男|男性|男\b|man|male/i.test(summaryText)) {
    badges.push("男");
  } else if (/性别女|女性|女\b|woman|female/i.test(summaryText)) {
    badges.push("女");
  }

  const ageMatch = summaryText.match(/年龄约?(\d{1,2})/);
  const age = ageMatch ? Number(ageMatch[1]) : null;

  if (age !== null) {
    if (age <= 16) {
      badges.push("儿童");
    } else if (age < 35) {
      badges.push("青年");
    } else {
      badges.push("中年");
    }
  } else if (/中年/i.test(summaryText)) {
    badges.push("中年");
  } else if (/青年|年轻|二十/i.test(summaryText)) {
    badges.push("青年");
  } else if (/儿童|孩子|小女孩|小男孩/i.test(summaryText)) {
    badges.push("儿童");
  }

  return badges.filter((value, index, list) => value && list.indexOf(value) === index);
};

const getCharacterCardSummary = (character: any) =>
  character?.appearance ||
  character?.description ||
  character?.personality ||
  "尚未补充角色外观描述";

const getCharacterStatusLabel = (character: any) => {
  const phase = getGenerationState("character", character).phase;
  if (phase === "queued") return "已提交";
  if (phase === "processing") return "生成中";
  if (phase === "success") return "可重生成";
  if (phase === "failed") return "生成失败";
  return getCharacterReferenceCount(character) > 0 ? "可重生成" : "待生成";
};

const getCharacterStatusText = (character: any) => {
  const phase = getGenerationState("character", character).phase;
  if (phase === "queued") return "系统已提交任务，准备生成单张三视图角色基准图。";
  if (phase === "processing") return "AI 正在生成白底三视图合成图，图内会包含侧面、正面、背面。";
  if (phase === "success") return "当前已有单张三视图基准图，后续镜头和视频会优先引用它。";
  if (phase === "failed") return formatGenerationError(character?.image_generation_error);
  return getCharacterReferenceCount(character) > 0
    ? "当前已有单张三视图基准图，后续镜头和视频会优先引用它。"
    : "点击生成后会产出一张白底三视图合成图，统一角色服装、发型和体态。";
};

const getCharacterActionText = (character: any) =>
  getCharacterReferenceCount(character) > 0 ? "重新生成" : "生成基准图";

const isGenerationBusy = (kind: GenerationKind, item: any) => {
  const phase = getGenerationState(kind, item).phase;
  return phase === "queued" || phase === "processing";
};

const getGenerationButtonText = (kind: GenerationKind, item: any) => {
  const phase = getGenerationState(kind, item).phase;
  if (phase === "queued") return "已提交";
  if (phase === "processing") return "生成中";
  if (phase === "failed" || hasImage(item)) return "重新生成";
  return "生成图片";
};

const syncGenerationStates = () => {
  drama.value?.characters?.forEach((character) => {
    getCharacterReferenceImages(character).forEach((image, index) => {
      if (image) {
        delete characterSlotStates.value[buildCharacterSlotKey(character.id, index)];
      }
    });

    const key = buildGenerationKey("character", character.id);
    const phase = getServerGenerationPhase(character.image_generation_status);
    const completedViews = getCharacterReferenceCount(character);
    const localPhase = generationStates.value[key]?.phase;

    if (phase === "queued") {
      generationStates.value[key] = {
        phase,
        detail: "角色定妆任务已提交，正在等待处理",
      };
      return;
    }

    if (
      phase === "processing" ||
      ((localPhase === "queued" || localPhase === "processing") &&
        completedViews < 1)
    ) {
      generationStates.value[key] = {
        phase,
        detail: "正在生成角色图片，请稍候",
      };
      return;
    }

    if (phase === "failed") {
      generationStates.value[key] = {
        phase,
        detail: formatGenerationError(character.image_generation_error),
      };
      clearGenerationMonitor(key);
      return;
    }

    if (completedViews >= 1 && generationStates.value[key]) {
      generationStates.value[key] = {
        phase: "success",
        detail: "角色图片已更新，可直接用于后续分镜和视频",
      };
      clearGenerationMonitor(key);
    }
  });

  scenes.value.forEach((scene) => {
    const key = buildGenerationKey("scene", scene.id);
    const phase = getServerGenerationPhase(scene.image_generation_status);

    if (phase === "queued") {
      generationStates.value[key] = {
        phase,
        detail: "场景图任务已提交，正在排队处理中",
      };
      return;
    }

    if (phase === "processing") {
      generationStates.value[key] = {
        phase,
        detail: "场景背景图生成中，页面会自动更新",
      };
      return;
    }

    if (phase === "failed") {
      generationStates.value[key] = {
        phase,
        detail: formatGenerationError(scene.image_generation_error),
      };
      clearGenerationMonitor(key);
      return;
    }

    if (hasImage(scene) && generationStates.value[key]) {
      generationStates.value[key] = {
        phase: "success",
        detail: "场景图已更新，可直接用于背景参考",
      };
      clearGenerationMonitor(key);
    }
  });

  drama.value?.props?.forEach((prop: any) => {
    const key = buildGenerationKey("prop", prop.id);
    if (hasImage(prop) && generationStates.value[key]) {
      generationStates.value[key] = {
        phase: "success",
        detail: "道具图已回填，可直接继续使用",
      };
      clearGenerationMonitor(key);
    }
  });
};

const monitorGenerationTask = (
  kind: GenerationKind,
  id: string | number | undefined,
  getItem: () => any,
  maxAttempts = 24,
  interval = 2500,
) => {
  const key = buildGenerationKey(kind, id);
  clearGenerationMonitor(key);

  let attempts = 0;
  generationTimers.set(
    key,
    setInterval(async () => {
      attempts += 1;
      await loadDramaData();

      const latestItem = getItem();
      const latestState = latestItem
        ? getGenerationState(kind, latestItem)
        : { phase: "idle" as GenerationPhase, detail: "" };
      const latestPhase = latestState.phase;

      if (latestPhase === "failed") {
        ElMessage.error(latestState.detail || "图片生成失败");
        clearGenerationMonitor(key);
        return;
      }

      if (latestPhase === "success") {
        clearGenerationMonitor(key);
        return;
      }

      if (attempts >= maxAttempts) {
        if (latestItem && hasImage(latestItem)) {
          setGenerationState(kind, id, {
            phase: "success",
            detail: "图片已生成完成",
          });
        } else {
          setGenerationState(kind, id, {
            phase: "failed",
            detail: "等待结果超时，请稍后刷新或重新生成",
          });
        }
        clearGenerationMonitor(key);
      }
    }, interval),
  );
};

onUnmounted(() => {
  if (pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null;
  }
  clearGenerationMonitor();
});

const startPolling = (
  callback: () => Promise<void>,
  maxAttempts = 20,
  interval = 3000,
) => {
  if (pollingTimer) clearInterval(pollingTimer);

  let attempts = 0;
  pollingTimer = setInterval(async () => {
    attempts += 1;
    await callback();
    if (attempts >= maxAttempts && pollingTimer) {
      clearInterval(pollingTimer);
      pollingTimer = null;
    }
  }, interval);
};

const sleep = (ms: number) =>
  new Promise((resolve) => {
    setTimeout(resolve, ms);
  });

const parseTaskResult = (result: any) => {
  if (!result) return null;
  if (typeof result === "string") {
    try {
      return JSON.parse(result);
    } catch {
      return null;
    }
  }
  return result;
};

const waitForCharacterExtractionTask = async (
  taskId: string,
  maxAttempts = 80,
  interval = 2500,
) => {
  for (let attempt = 0; attempt < maxAttempts; attempt += 1) {
    const response = await taskAPI.getStatus(taskId);
    const task = response?.data ?? response;

    if (!task) {
      throw new Error("角色提取任务状态读取失败");
    }

    const taskMessage =
      task.message ||
      (task.status === "pending" ? "角色提取任务排队中..." : "正在提取角色...");
    const progressText =
      typeof task.progress === "number" ? `（${task.progress}%）` : "";
    extractingCharactersMessage.value = `${taskMessage}${progressText}`;

    if (task.status === "completed") {
      const result = parseTaskResult(task.result);
      await loadDramaData();
      return result;
    }

    if (task.status === "failed") {
      throw new Error(task.error || task.message || "角色提取失败");
    }

    if ((attempt + 1) % 4 === 0) {
      await loadDramaData();
    }

    await sleep(interval);
  }

  throw new Error("角色提取超时，请稍后刷新后查看结果");
};

const loadDramaData = async () => {
  try {
    const data = await dramaAPI.get(route.params.id as string);
    drama.value = data;
    loadScenes();
    syncGenerationStates();
  } catch (error: any) {
    ElMessage.error(error.message || "加载项目数据失败");
  }
};

const loadScenes = async () => {
  // 场景数据已经在drama中加载了（后端Preload了Scenes）
  if (drama.value?.scenes) {
    scenes.value = drama.value.scenes;
  } else {
    scenes.value = [];
  }
  syncGenerationStates();
};

const getStatusType = (status?: string) => {
  const map: Record<string, any> = {
    draft: "info",
    in_progress: "warning",
    completed: "success",
  };
  return map[status || "draft"] || "info";
};

const getStatusText = (status?: string) => {
  const map: Record<string, string> = {
    draft: "草稿",
    in_progress: "制作中",
    completed: "已完成",
  };
  return map[status || "draft"] || "草稿";
};

const getEpisodeStatusType = (episode: any) => {
  if (episode.shots && episode.shots.length > 0) return "success";
  if (episode.script_content) return "warning";
  return "info";
};

const getEpisodeStatusText = (episode: any) => {
  if (episode.shots && episode.shots.length > 0) return "已拆分";
  if (episode.script_content) return "已创建";
  return "草稿";
};

const formatDate = (date?: string) => {
  if (!date) return "-";
  return new Date(date).toLocaleString("zh-CN");
};

const goBack = () => {
  router.replace("/dramas");
};

const goToAutoWorkflow = () => {
  router.push({
    name: "AutoWorkflow",
    params: {
      id: route.params.id,
    },
  });
};

const goToVideoIdeaTab = () => {
  activeTab.value = "videoIdea";
  videoIdeaResponse.value = "";
  router.replace({
    name: route.name as string,
    params: route.params,
    query: {
      ...route.query,
      tab: "videoIdea",
    },
  });
};

const isVideoIntent = (message: string) => {
  const normalized = message.trim().toLowerCase();
  if (!normalized) return false;

  const blockedKeywords = ["聊天", "闲聊", "聊天吗", "天气", "百科", "搜索", "问答", "查找", "search", "google", "百度"];
  if (blockedKeywords.some((keyword) => normalized.includes(keyword))) {
    return false;
  }

  const allowedKeywords = [
    "视频",
    "短剧",
    "剧本",
    "脚本",
    "剧情",
    "故事",
    "分镜",
    "场景",
    "角色",
    "video",
    "movie",
    "clip",
  ];

  return allowedKeywords.some((keyword) => normalized.includes(keyword));
};

const submitVideoIdea = async () => {
  const content = videoIdeaInput.value.trim();
  if (!content) {
    ElMessage.warning("请输入你的剧本或视频想法");
    return;
  }

  if (!isVideoIntent(content)) {
    videoIdeaResponse.value = "做不到";
    return;
  }

  if (!drama.value?.id) {
    ElMessage.error("项目数据未加载完成");
    return;
  }

  videoIdeaSubmitting.value = true;
  videoIdeaResponse.value = "";

  try {
    const metadataPayload = {
      ...(drama.value.metadata || {}),
      full_script: content,
    };

    await dramaAPI.update(drama.value.id, {
      metadata: metadataPayload,
    });

    await workflowAPI.startProjectWorkflow(drama.value.id);
    videoIdeaResponse.value = "已提交脚本并启动自动化生成，请在自动工作流查看进度。";
    ElMessage.success("已提交脚本并启动自动化生成");
    goToAutoWorkflow();
  } catch (error: any) {
    const message = error?.message || "提交失败";
    videoIdeaResponse.value = message;
    ElMessage.error(message);
  } finally {
    videoIdeaSubmitting.value = false;
  }
};
const createNewEpisode = () => {
  const nextEpisodeNumber = episodesCount.value + 1;
  router.push({
    name: "EpisodeWorkflowNew",
    params: {
      id: route.params.id,
      episodeNumber: nextEpisodeNumber,
    },
  });
};

const enterEpisodeWorkflow = (episode: any) => {
  router.push({
    name: "EpisodeWorkflowNew",
    params: {
      id: route.params.id,
      episodeNumber: episode.episode_number,
    },
  });
};

const deleteEpisode = async (episode: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除第${episode.episode_number}章吗？此操作将同时删除该章节的所有相关数据（角色、场景、分镜等）。`,
      "删除确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      },
    );

    // 过滤掉要删除的章节
    const existingEpisodes = drama.value?.episodes || [];
    const updatedEpisodes = existingEpisodes
      .filter((ep) => ep.episode_number !== episode.episode_number)
      .map((ep) => ({
        episode_number: ep.episode_number,
        title: ep.title,
        script_content: ep.script_content,
        description: ep.description,
        duration: ep.duration,
        status: ep.status,
      }));

    // 保存更新后的章节列表
    await dramaAPI.saveEpisodes(drama.value!.id, updatedEpisodes);

    ElMessage.success(`第${episode.episode_number}章删除成功`);
    await loadDramaData();
  } catch (error: any) {
    if (error !== "cancel") {
      ElMessage.error(error.message || "删除失败");
    }
  }
};

const openAddCharacterDialog = () => {
  editingCharacter.value = null;
  newCharacter.value = {
    name: "",
    role: "supporting",
    appearance: "",
    personality: "",
    description: "",
    image_url: "",
  };
  addCharacterDialogVisible.value = true;
};

const handleCharacterAvatarSuccess = (response: any) => {
  if (response.data && response.data.url) {
    newCharacter.value.image_url = response.data.url;
    newCharacter.value.local_path = response.data.local_path || "";
  }
};

const handleSceneImageSuccess = (response: any) => {
  if (response.data && response.data.url) {
    newScene.value.image_url = response.data.url;
    newScene.value.local_path = response.data.local_path || "";
  }
};

const beforeAvatarUpload = (file: any) => {
  const isImage = file.type.startsWith("image/");
  const isLt10M = file.size / 1024 / 1024 < 10;

  if (!isImage) {
    ElMessage.error("只能上传图片文件!");
  }
  if (!isLt10M) {
    ElMessage.error("图片大小不能超过 10MB!");
  }
  return isImage && isLt10M;
};

const generateCharacterImage = async (character: any) => {
  try {
    await characterLibraryAPI.generateCharacterImage(character.id);
    ElMessage.success("图片生成任务已提交");
    startPolling(loadDramaData);
  } catch (error: any) {
    ElMessage.error(error.message || "生成失败");
  }
};

const openExtractCharacterDialog = () => {
  extractCharactersDialogVisible.value = true;
  if (sortedEpisodes.value.length > 0 && !selectedExtractEpisodeId.value) {
    selectedExtractEpisodeId.value = sortedEpisodes.value[0].id;
  }
};

const handleExtractCharacters = async () => {
  if (!selectedExtractEpisodeId.value) return;

  try {
    extractingCharacters.value = true;
    extractingCharactersMessage.value = "正在提交角色提取任务...";
    ElMessage.info("开始从剧本提取角色，请稍候");
    const res = await characterLibraryAPI.extractFromEpisode(
      selectedExtractEpisodeId.value,
    );
    const taskId = res?.task_id ?? res?.data?.task_id;
    if (!taskId) {
      throw new Error("未获取到角色提取任务编号");
    }

    extractingCharactersMessage.value = "角色提取任务已提交，正在分析剧本...";
    const result = await waitForCharacterExtractionTask(taskId);
    const extractedCount =
      result?.characters?.length ??
      result?.data?.characters?.length ??
      result?.count;
    extractCharactersDialogVisible.value = false;
    ElMessage.success(
      typeof extractedCount === "number"
        ? `角色提取完成，新增或更新 ${extractedCount} 个角色`
        : "角色提取完成，正在刷新角色列表",
    );
  } catch (error: any) {
    ElMessage.error(error.message || "提取失败");
  } finally {
    extractingCharacters.value = false;
    extractingCharactersMessage.value = "正在提交角色提取任务...";
  }
};

const generateSceneImage = async (scene: any) => {
  try {
    await dramaAPI.generateSceneImage({ scene_id: scene.id });
    ElMessage.success("图片生成任务已提交");
    startPolling(loadScenes);
  } catch (error: any) {
    ElMessage.error(error.message || "生成失败");
  }
};

const openExtractSceneDialog = () => {
  extractScenesDialogVisible.value = true;
  if (sortedEpisodes.value.length > 0 && !selectedExtractEpisodeId.value) {
    selectedExtractEpisodeId.value = sortedEpisodes.value[0].id;
  }
};

const handleExtractScenes = async () => {
  if (!selectedExtractEpisodeId.value) return;

  try {
    const res = await dramaAPI.extractBackgrounds(
      selectedExtractEpisodeId.value.toString(),
    );
    extractScenesDialogVisible.value = false;

    // 自动刷新几次
    let checkCount = 0;
    const checkInterval = setInterval(() => {
      loadScenes();
      checkCount++;
      if (checkCount > 10) clearInterval(checkInterval);
    }, 5000);
  } catch (error: any) {
    ElMessage.error(error.message || "提取失败");
  }
};

const saveCharacter = async () => {
  if (!newCharacter.value.name.trim()) {
    ElMessage.warning("请输入角色名称");
    return;
  }

  try {
    if (editingCharacter.value) {
      // Edit existing character using dedicated update endpoint
      await dramaAPI.updateCharacter(editingCharacter.value.id, {
        name: newCharacter.value.name,
        role: newCharacter.value.role,
        appearance: newCharacter.value.appearance,
        personality: newCharacter.value.personality,
        description: newCharacter.value.description,
        image_url: newCharacter.value.image_url,
        local_path: newCharacter.value.local_path,
      });
      ElMessage.success("角色更新成功");
    } else {
      // Add new character
      const allCharacters = [
        ...(drama.value?.characters || []).map((c) => ({
          name: c.name,
          role: c.role,
          appearance: c.appearance,
          personality: c.personality,
          description: c.description,
          image_url: c.image_url,
          local_path: c.local_path,
        })),
        newCharacter.value,
      ];

      await dramaAPI.saveCharacters(drama.value!.id, allCharacters);
      ElMessage.success("角色添加成功");
    }

    addCharacterDialogVisible.value = false;
    await loadDramaData();
  } catch (error: any) {
    ElMessage.error(error.message || "操作失败");
  }
};

const editCharacter = (character: any) => {
  editingCharacter.value = character;
  newCharacter.value = {
    name: character.name,
    role: character.role || "supporting",
    appearance: character.appearance || "",
    personality: character.personality || "",
    description: character.description || "",
    image_url: character.image_url || "",
    local_path: character.local_path || "",
  };
  addCharacterDialogVisible.value = true;
};

const deleteCharacter = async (character: any) => {
  if (!character.id) {
    ElMessage.error("角色ID不存在，无法删除");
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除角色"${character.name}"吗？此操作不可恢复。`,
      "删除确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      },
    );

    await characterLibraryAPI.deleteCharacter(character.id);
    ElMessage.success("角色已删除");
    await loadDramaData();
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("删除角色失败:", error);
      ElMessage.error(error.message || "删除失败");
    }
  }
};

const openAddSceneDialog = () => {
  editingScene.value = null;
  newScene.value = {
    location: "",
    prompt: "",
    image_url: "",
  };
  addSceneDialogVisible.value = true;
};

const saveScene = async () => {
  if (!newScene.value.location.trim()) {
    ElMessage.warning("请输入场景名称");
    return;
  }

  try {
    if (editingScene.value) {
      // Update existing scene
      await dramaAPI.updateScene(editingScene.value.id, {
        location: newScene.value.location,
        description: newScene.value.prompt,
        image_url: newScene.value.image_url,
        local_path: newScene.value.local_path,
      });
      // prompt field in Update is description or prompt? Check backend.
      // UpdateSceneRequest has Description *string.
      // And also ImagePrompt *string and VideoPrompt *string.
      // The backend model has Prompt string.
      // Checking backend handler:
      /*
        if req.Description != nil { updates["description"] = req.Description }
        if req.ImagePrompt != nil { updates["image_prompt"] = req.ImagePrompt }
      */
      // But CreateScene uses Prompt.
      // Let's assume description maps to Prompt or Description.
      // Wait, UpdateSceneRequest has Description but NO Prompt field?
      // Let's check backend UpdateSceneRequest struct again.
      // It has `ImagePrompt` and `VideoPrompt`, and `Description`.
      // But `Prompt` usually refers to image prompt in Scene model?
      // `models.Scene` has `Prompt` string.
      // `CreateScene` sets `Prompt: req.Prompt`.
      // `UpdateScene` handler:
      /*
      	if req.Description != nil {
      		updates["description"] = req.Description
      	}
      */
      // It seems UpdateScene doesn't support updating the main `Prompt` field directly via UpdateSceneRequest?
      // Wait, `UpdateScenePrompt` endpoint exists! `/scenes/:id/prompt`
      // But we probably want to update everything in one go.
      // I should update UpdateSceneRequest in backend if needed or use UpdateScenePrompt separately.
      // For now, let's look at scene model:
      // Scene struct: Location, Time, Description, Prompt...
      // Let's use `description` for now as it's available in Update.
      // Or if `prompt` is critical, I might need to call UpdateScenePrompt too.
      // Let's check `CreateScene` again. It uses `Prompt`.

      // Let's just update prompt via specific endpoint if needed, or mapping description to description.
      // Actually `newScene.prompt` is mapped to `description` in my current code for Update.
      // Let's stick with that for now or fix backend to support prompt update in general update.
    } else {
      // Create new scene
      await dramaAPI.createScene({
        drama_id: drama.value!.id,
        location: newScene.value.location,
        prompt: newScene.value.prompt,
        description: newScene.value.prompt,
        image_url: newScene.value.image_url,
        local_path: newScene.value.local_path,
      });
    }

    ElMessage.success(editingScene.value ? "场景更新成功" : "场景添加成功");
    addSceneDialogVisible.value = false;
    await loadScenes();
  } catch (error: any) {
    ElMessage.error(error.message || "操作失败");
  }
};

const editScene = (scene: any) => {
  editingScene.value = scene;
  newScene.value = {
    location: scene.location || scene.name || "",
    prompt: scene.prompt || scene.description || "",
    image_url: scene.image_url || "",
    local_path: scene.local_path || "",
  };
  addSceneDialogVisible.value = true;
};

const deleteScene = async (scene: any) => {
  if (!scene.id) {
    ElMessage.error("场景ID不存在，无法删除");
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除场景"${scene.name || scene.location}"吗？此操作不可恢复。`,
      "删除确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      },
    );

    await dramaAPI.deleteScene(scene.id.toString());
    ElMessage.success("场景已删除");
    await loadScenes();
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("删除场景失败:", error);
      ElMessage.error(error.message || "删除失败");
    }
  }
};

const openAddPropDialog = () => {
  editingProp.value = null;
  newProp.value = {
    name: "",
    description: "",
    prompt: "",
    type: "",
    image_url: "",
    local_path: "",
  };
  addPropDialogVisible.value = true;
};

const saveProp = async () => {
  if (!newProp.value.name.trim()) {
    ElMessage.warning("请输入道具名称");
    return;
  }

  try {
    const propData = {
      drama_id: drama.value!.id,
      name: newProp.value.name,
      description: newProp.value.description,
      prompt: newProp.value.prompt,
      type: newProp.value.type,
      image_url: newProp.value.image_url,
      local_path: newProp.value.local_path,
    };

    if (editingProp.value) {
      await propAPI.update(editingProp.value.id, propData);
      ElMessage.success("道具更新成功");
    } else {
      await propAPI.create(propData as any);
      ElMessage.success("道具添加成功");
    }

    addPropDialogVisible.value = false;
    await loadDramaData();
  } catch (error: any) {
    ElMessage.error(error.message || "操作失败");
  }
};

const editProp = (prop: any) => {
  editingProp.value = prop;
  newProp.value = {
    name: prop.name,
    description: prop.description || "",
    prompt: prop.prompt || "",
    type: prop.type || "",
    image_url: prop.image_url || "",
    local_path: prop.local_path || "",
  };
  addPropDialogVisible.value = true;
};

const deleteProp = async (prop: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除道具"${prop.name}"吗？此操作不可恢复。`,
      "删除确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      },
    );

    await propAPI.delete(prop.id);
    ElMessage.success("道具已删除");
    await loadDramaData();
  } catch (error: any) {
    if (error !== "cancel") {
      ElMessage.error(error.message || "删除失败");
    }
  }
};

const generatePropImage = async (prop: any) => {
  if (!prop.prompt) {
    ElMessage.warning("请先设置道具的图片提示词");
    editProp(prop);
    return;
  }

  try {
    await propAPI.generateImage(prop.id);
    ElMessage.success("图片生成任务已提交");
    startPolling(loadDramaData);
  } catch (error: any) {
    ElMessage.error(error.message || "生成失败");
  }
};

const handlePropImageSuccess = (response: any) => {
  if (response.data && response.data.url) {
    newProp.value.image_url = response.data.url;
    newProp.value.local_path = response.data.local_path || "";
  }
};

const monitorCharacterSlotGeneration = (
  characterId: string | number,
  index: number,
  imageGenId: number,
  maxAttempts = 60,
  interval = 3000,
) => {
  const key = buildCharacterSlotKey(characterId, index);
  clearGenerationMonitor(key);

  let attempts = 0;
  generationTimers.set(
    key,
    setInterval(async () => {
      attempts += 1;

      try {
        const response = await imageAPI.getImage(imageGenId);
        const imageGen = response?.data ?? response;

        if (imageGen?.status === "failed") {
          const detail = formatGenerationError(imageGen?.error_msg);
          characterSlotStates.value[key] = {
            phase: "failed",
            detail,
            imageGenId,
          };
          ElMessage.error(detail);
          clearGenerationMonitor(key);
          await loadDramaData();
          return;
        }

        if (imageGen?.status === "completed") {
          await loadDramaData();
          delete characterSlotStates.value[key];
          clearGenerationMonitor(key);
          return;
        }

        characterSlotStates.value[key] = {
          phase: "processing",
          detail: `正在生成${characterReferenceLabels[index]}`,
          imageGenId,
        };
      } catch (error: any) {
        const detail = error?.message || "读取单张角色图状态失败";
        characterSlotStates.value[key] = {
          phase: "failed",
          detail,
          imageGenId,
        };
        ElMessage.error(detail);
        clearGenerationMonitor(key);
        return;
      }

      if (attempts >= maxAttempts) {
        characterSlotStates.value[key] = {
          phase: "failed",
          detail: "等待单张角色图结果超时，请稍后重试",
          imageGenId,
        };
        clearGenerationMonitor(key);
      }
    }, interval),
  );
};

const submitCharacterImageGeneration = async (character: any) => {
  try {
    setGenerationState("character", character.id, {
      phase: "queued",
      detail: "正在提交角色图片任务",
    });
    await characterLibraryAPI.generateCharacterImage(character.id);
    setGenerationState("character", character.id, {
      phase: "processing",
      detail: "角色图片任务已提交，正在生成中",
    });
    ElMessage.success("角色图片生成任务已提交");
    monitorGenerationTask(
      "character",
      character.id,
      () => drama.value?.characters?.find((item) => item.id === character.id),
      60,
      3000,
    );
  } catch (error: any) {
    setGenerationState("character", character.id, {
      phase: "failed",
      detail: error.message || "角色图片生成失败",
    });
    ElMessage.error(error.message || "生成失败");
  }
};

const submitCharacterViewGeneration = async (character: any, index: number) => {
  const key = buildCharacterSlotKey(character.id, index);
  const frameType = getCharacterFrameType(index);

  try {
    characterSlotStates.value[key] = {
      phase: "queued",
      detail: `正在提交${characterReferenceLabels[index]}`,
    };

    const response = await characterLibraryAPI.generateCharacterImage(
      character.id,
      undefined,
      frameType,
    );
    const imageGen = response?.image_generation ?? response?.data?.image_generation;
    const imageGenId = imageGen?.id;

    characterSlotStates.value[key] = {
      phase: "processing",
      detail: `${characterReferenceLabels[index]}生成中`,
      imageGenId,
    };

    if (imageGenId) {
      monitorCharacterSlotGeneration(character.id, index, imageGenId, 60, 3000);
    } else {
      await loadDramaData();
    }
  } catch (error: any) {
    characterSlotStates.value[key] = {
      phase: "failed",
      detail: error.message || `${characterReferenceLabels[index]}生成失败`,
    };
    ElMessage.error(error.message || "生成失败");
  }
};

const submitSceneImageGeneration = async (scene: any) => {
  try {
    setGenerationState("scene", scene.id, {
      phase: "queued",
      detail: "正在提交场景图任务",
    });
    await dramaAPI.generateSceneImage({ scene_id: scene.id });
    setGenerationState("scene", scene.id, {
      phase: "processing",
      detail: "场景图任务已提交，正在生成中",
    });
    ElMessage.success("场景图片生成任务已提交");
    monitorGenerationTask("scene", scene.id, () =>
      scenes.value.find((item) => item.id === scene.id),
    );
  } catch (error: any) {
    setGenerationState("scene", scene.id, {
      phase: "failed",
      detail: error.message || "场景图片生成失败",
    });
    ElMessage.error(error.message || "生成失败");
  }
};

const submitPropImageGeneration = async (prop: any) => {
  if (!prop.prompt) {
    ElMessage.warning("请先设置道具图片提示词");
    editProp(prop);
    return;
  }

  try {
    setGenerationState("prop", prop.id, {
      phase: "queued",
      detail: "正在提交道具图任务",
    });
    await propAPI.generateImage(prop.id);
    setGenerationState("prop", prop.id, {
      phase: "processing",
      detail: "道具图任务已提交，等待结果回填",
    });
    ElMessage.success("道具图片生成任务已提交");
    monitorGenerationTask("prop", prop.id, () =>
      drama.value?.props?.find((item: any) => item.id === prop.id),
    );
  } catch (error: any) {
    setGenerationState("prop", prop.id, {
      phase: "failed",
      detail: error.message || "道具图片生成失败",
    });
    ElMessage.error(error.message || "生成失败");
  }
};

const openExtractDialog = () => {
  extractPropsDialogVisible.value = true;
  if (sortedEpisodes.value.length > 0 && !selectedExtractEpisodeId.value) {
    selectedExtractEpisodeId.value = sortedEpisodes.value[0].id;
  }
};

const handleExtractProps = async () => {
  if (!selectedExtractEpisodeId.value) return;

  try {
    const res = await propAPI.extractFromScript(selectedExtractEpisodeId.value);
    extractPropsDialogVisible.value = false;

    // 自动刷新几次
    let checkCount = 0;
    const checkInterval = setInterval(() => {
      loadDramaData();
      checkCount++;
      if (checkCount > 10) clearInterval(checkInterval);
    }, 5000);
  } catch (error: any) {
    ElMessage.error(error.message || t("common.failed"));
  }
};

watch(
  activeTab,
  (tab) => {
    router.replace({
      name: route.name as string,
      params: route.params,
      query: {
        ...route.query,
        tab,
      },
    });
  },
  { flush: "post" },
);

onMounted(() => {
  loadDramaData();
  loadScenes();

  // 如果有query参数指定tab，切换到对应tab
  if (route.query.tab) {
    activeTab.value = route.query.tab as string;
  }
});
</script>

<style scoped>
/* ========================================
   Page Layout / 页面布局 - 紧凑边距
   ======================================== */
.page-container {
  min-height: 100vh;
  background: var(--bg-primary);
  /* padding: var(--space-2) var(--space-3); */
  transition: background var(--transition-normal);
}

@media (min-width: 768px) {
  .page-container {
    /* padding: var(--space-3) var(--space-4); */
  }
}

@media (min-width: 1024px) {
  .page-container {
    /* padding: var(--space-4) var(--space-5); */
  }
}

.content-wrapper {
  margin: 0 auto;
  width: 100%;
}

.action-btn {
  margin-right: 12px;
}

.video-idea-tip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  margin-top: 4px;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border: 1px dashed var(--border-primary);
  border-radius: var(--radius-md);
}

.video-idea-response {
  margin-top: 12px;
}

.video-idea-wrapper {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.video-idea-hero {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-4);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  background: var(--bg-secondary);
}

.video-idea-kicker {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--accent);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.video-idea-title {
  margin: 4px 0 8px;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
}

.video-idea-desc {
  margin: 0;
  max-width: 760px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.video-idea-tags {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.video-idea-card {
  border-radius: var(--radius-lg);
}

.video-idea-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.video-idea-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.video-idea-subtitle {
  color: var(--text-secondary);
  font-weight: 400;
  font-size: 0.95rem;
}

/* ========================================
   Stats Grid / 统计网格 - 紧凑间距
   ======================================== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-2);
  margin-bottom: var(--space-3);
}

@media (min-width: 640px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--space-3);
  }
}

@media (min-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* ========================================
   Tabs Wrapper / 标签页容器 - 紧凑内边距
   ======================================== */
.tabs-wrapper {
  background: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-3);
  box-shadow: var(--shadow-card);
}

@media (min-width: 768px) {
  .tabs-wrapper {
    padding: var(--space-4);
  }
}

/* ========================================
   Tab Header / 标签页头部
   ======================================== */
.tab-header {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

@media (min-width: 640px) {
  .tab-header {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }
}

.tab-header h2 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: -0.01em;
}

.generation-guide {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  background: var(--bg-panel);
  box-shadow: var(--shadow-sm);
}

@media (min-width: 768px) {
  .generation-guide {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

.generation-guide-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-primary);
}

.generation-guide-text {
  margin: 6px 0 0;
  font-size: 0.8125rem;
  line-height: 1.6;
  color: var(--text-muted);
}

.generation-guide-steps {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.generation-guide-steps span {
  padding: 6px 12px;
  border-radius: 999px;
  border: 1px solid var(--border-primary);
  background: color-mix(in srgb, var(--bg-card) 92%, transparent);
  color: var(--text-secondary);
  font-size: 0.75rem;
  white-space: nowrap;
}

.legacy-generation-guide {
  display: none;
}

.character-pipeline {
  display: grid;
  grid-template-columns: minmax(0, 2.4fr) minmax(260px, 1fr);
  gap: 16px;
}

.character-guide-panel {
  padding: 18px 20px;
  border: 1px solid var(--border-primary);
  border-radius: 24px;
  background:
    radial-gradient(circle at top right, color-mix(in srgb, var(--accent-light) 92%, transparent), transparent 36%),
    var(--bg-panel);
  box-shadow: var(--shadow-card);
}

.character-guide-panel.is-compact {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  align-items: stretch;
}

.character-guide-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;
}

.character-guide-meta span,
.character-guide-stat {
  border: 1px solid var(--border-primary);
  background: color-mix(in srgb, var(--bg-card) 94%, transparent);
}

.character-guide-meta span {
  padding: 7px 12px;
  border-radius: 999px;
  font-size: 0.76rem;
  color: var(--text-secondary);
}

.character-guide-stat {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 6px;
  border-radius: 18px;
  padding: 14px;
}

.character-guide-stat strong {
  font-size: 1.4rem;
  line-height: 1;
  color: var(--text-primary);
}

.character-guide-stat span {
  font-size: 0.76rem;
  color: var(--text-muted);
}

@media (max-width: 960px) {
  .character-pipeline {
    grid-template-columns: 1fr;
  }
}

/* ========================================
   Character & Scene Cards / 角色场景卡片
   ======================================== */
.character-card,
.scene-card {
  margin-bottom: var(--space-4);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: all var(--transition-normal);
}

.character-card:hover,
.scene-card:hover {
  border-color: var(--border-secondary);
  box-shadow: var(--shadow-card-hover);
}

.character-card.is-generating,
.scene-card.is-generating {
  border-color: color-mix(in srgb, var(--accent) 45%, var(--border-primary));
  box-shadow: 0 22px 44px rgba(15, 23, 42, 0.12);
}

.character-card.is-failed,
.scene-card.is-failed {
  border-color: color-mix(in srgb, var(--danger) 45%, var(--border-primary));
}

.character-card :deep(.el-card__body),
.scene-card :deep(.el-card__body) {
  padding: 0;
}

.character-preview,
.scene-preview {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 160px;
  background: linear-gradient(135deg, var(--accent) 0%, #06b6d4 100%);
  overflow: hidden;
}

.character-preview img,
.scene-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-normal);
}

.character-card:hover .character-preview img,
.scene-card:hover .scene-preview img {
  transform: scale(1.05);
}

.scene-placeholder {
  color: rgba(255, 255, 255, 0.7);
}

.character-info,
.scene-info {
  text-align: center;
  padding: var(--space-4);
}

.character-name {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--space-2);
}

.character-info h4,
.scene-info h4 {
  /* margin: 0 0 var(--space-2) 0; */
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.desc {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin: var(--space-2) 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}

.generation-process {
  margin: 0 var(--space-4) var(--space-4);
  padding: 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  background: color-mix(in srgb, var(--bg-secondary) 88%, transparent);
}

.character-reference-strip {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin: 0 var(--space-4) var(--space-4);
}

.character-reference-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 10px 8px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  background: color-mix(in srgb, var(--bg-card) 92%, transparent);
}

.character-reference-item span {
  font-size: 0.72rem;
  color: var(--text-secondary);
}

.generation-process-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
}

.generation-process-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--text-primary);
}

.generation-process-detail {
  min-height: 40px;
  margin: 10px 0 12px;
  font-size: 0.75rem;
  line-height: 1.6;
  color: var(--text-muted);
}

.generation-progress {
  margin-bottom: 10px;
}

.generation-step-row {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.generation-step-row span {
  text-align: center;
  padding-top: 6px;
  border-top: 1px solid var(--border-primary);
  font-size: 0.72rem;
  color: var(--text-quaternary, var(--text-muted));
}

.generation-step-row span.active {
  border-color: var(--accent);
  color: var(--text-primary);
}

.character-actions,
.scene-actions {
  display: flex;
  gap: var(--space-2);
  justify-content: center;
  padding: 0 var(--space-4) var(--space-4);
  flex-wrap: wrap;
}

.character-grid {
  align-items: stretch;
}

.character-card {
  border-radius: 24px;
  background: var(--bg-panel);
  box-shadow: var(--shadow-card);
}

.character-card .character-preview,
.character-card .character-info,
.character-card .character-reference-strip,
.character-card .generation-process,
.character-card .character-actions {
  display: none;
}

.character-shell {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 20px;
}

.character-shell-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.character-shell-title {
  min-width: 0;
}

.character-name-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.character-name-row h4 {
  margin: 0;
  font-size: 1.08rem;
  font-weight: 700;
  color: var(--text-primary);
}

.character-shell-summary {
  margin: 10px 0 0;
  font-size: 0.84rem;
  line-height: 1.7;
  color: var(--text-muted);
}

.character-workbench {
  display: grid;
  grid-template-columns: minmax(0, 1.5fr) minmax(220px, 0.9fr);
  gap: 14px;
}

.character-stage,
.character-spec-card,
.character-taskboard,
.character-view-slot {
  border: 1px solid var(--border-primary);
  border-radius: 22px;
  background: color-mix(in srgb, var(--bg-card) 94%, transparent);
}

.character-stage {
  display: grid;
  grid-template-columns: 180px minmax(0, 1fr);
  gap: 16px;
  padding: 16px;
}

.character-stage-canvas {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  min-height: 220px;
  border-radius: 18px;
  background: linear-gradient(180deg, var(--bg-elevated), var(--bg-secondary));
  border: 1px solid var(--border-primary);
  overflow: hidden;
  cursor: zoom-in;
}

.character-stage-canvas.is-empty {
  background: linear-gradient(180deg, var(--bg-card), var(--bg-secondary));
}

.character-stage-canvas.is-loading,
.character-view-slot.is-loading .character-view-media {
  border-color: var(--border-secondary);
  box-shadow: inset 0 0 0 1px var(--accent-light);
}

.character-stage-image {
  width: 100%;
  height: 100%;
  max-height: 220px;
  object-fit: contain;
  padding: 12px;
  cursor: zoom-in;
}

.character-stage-placeholder,
.character-view-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-muted);
}

.character-stage-placeholder .el-icon,
.character-view-placeholder .el-icon {
  font-size: 24px;
}

.character-stage-loading,
.character-view-loading {
  position: absolute;
  inset: auto 12px 12px 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.68);
  color: #fff;
  backdrop-filter: blur(8px);
}

.character-view-loading {
  position: static;
  min-height: 124px;
  width: 100%;
  border-radius: 16px;
  background: linear-gradient(180deg, var(--bg-card), var(--bg-secondary));
  color: var(--text-secondary);
}

.character-loading-spinner {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  border: 2px solid rgba(255, 255, 255, 0.32);
  border-top-color: #fff;
  animation: character-spin 0.8s linear infinite;
}

.character-view-loading .character-loading-spinner {
  border-color: rgba(37, 99, 235, 0.18);
  border-top-color: rgba(37, 99, 235, 0.92);
}

.character-loading-text {
  font-size: 0.78rem;
  font-weight: 600;
  letter-spacing: 0.01em;
}

@keyframes character-spin {
  to {
    transform: rotate(360deg);
  }
}

.character-stage-copy {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
}

.character-stage-kicker {
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.character-stage-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.character-stage-heading h5 {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-primary);
}

.character-stage-count {
  padding: 5px 10px;
  border-radius: 999px;
  background: var(--accent-light);
  color: var(--text-secondary);
  font-size: 0.76rem;
}

.character-stage-summary {
  margin: 0;
  font-size: 0.82rem;
  line-height: 1.7;
  color: var(--text-secondary);
}

.character-stage-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.character-stage-chips span {
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 0.74rem;
  color: var(--text-secondary);
  background: var(--accent-light);
}

.character-spec-card {
  padding: 16px;
}

.character-spec-title {
  display: block;
  margin-bottom: 12px;
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--text-primary);
}

.character-spec-list {
  margin: 0;
  padding-left: 18px;
  display: grid;
  gap: 10px;
  font-size: 0.79rem;
  line-height: 1.6;
  color: var(--text-secondary);
}

.character-view-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.character-view-slot {
  padding: 12px;
}

.character-view-slot.is-failed {
  border-color: var(--error);
  box-shadow: 0 18px 34px color-mix(in srgb, var(--error) 12%, transparent);
}

.character-view-head {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 10px;
}

.character-view-label {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--text-primary);
}

.character-view-hint {
  font-size: 0.72rem;
  color: var(--text-muted);
}

.character-view-media {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  min-height: 124px;
  border-radius: 16px;
  background: linear-gradient(180deg, var(--bg-elevated), var(--bg-secondary));
  border: 1px solid var(--border-primary);
  overflow: hidden;
}

.character-view-image {
  width: 100%;
  height: 124px;
  object-fit: contain;
  padding: 10px;
  cursor: zoom-in;
}

.character-view-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
}

.character-view-actions :deep(.el-button) {
  min-height: 30px;
  padding: 6px 12px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0;
}

.character-view-actions :deep(.el-button.is-disabled) {
  opacity: 0.58;
}

.character-view-zoom-btn {
  color: var(--text-secondary);
}

.character-view-generate-btn,
.character-view-actions :deep(.el-button--primary.is-plain) {
  background: var(--bg-secondary);
  border-color: var(--border-secondary);
  color: var(--text-primary);
}

.character-view-generate-btn:hover,
.character-view-actions :deep(.el-button--primary.is-plain:hover) {
  background: var(--accent-light);
  border-color: var(--border-focus);
  color: var(--text-primary);
}

.character-view-status {
  margin-top: 8px;
  font-size: 0.72rem;
  line-height: 1.5;
  color: var(--text-muted);
  min-height: 18px;
}

.character-view-status.is-failed {
  color: var(--error);
}

.character-preview-dialog {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 520px;
  padding: 20px;
  border-radius: 24px;
  background: linear-gradient(180deg, var(--bg-card) 0%, var(--bg-secondary) 100%);
}

.character-preview-dialog-image {
  width: 100%;
  max-height: 72vh;
  object-fit: contain;
}

.character-taskboard {
  padding: 14px 16px 16px;
}

.character-taskboard-head {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 10px;
}

.character-taskboard-detail {
  font-size: 0.78rem;
  line-height: 1.6;
  color: var(--text-muted);
}

.character-task-track {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.character-task-node {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 16px;
  border: 1px solid var(--border-primary);
  background: color-mix(in srgb, var(--bg-secondary) 92%, transparent);
  color: var(--text-muted);
}

.character-task-node i {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 999px;
  font-style: normal;
  font-size: 0.75rem;
  background: var(--accent-light);
}

.character-task-node span {
  font-size: 0.76rem;
  font-weight: 600;
}

.character-task-node.active {
  color: var(--text-primary);
  border-color: var(--border-secondary);
  background: color-mix(in srgb, var(--accent-light) 86%, var(--bg-card) 14%);
}

.character-task-node.active i {
  background: var(--accent);
  color: var(--text-inverse);
}

.character-actionbar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.character-actionbar :deep(.el-button) {
  min-height: 38px;
  padding: 0 16px;
  font-weight: 700;
}

.character-actionbar :deep(.el-button--default) {
  background: var(--bg-card);
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.character-actionbar :deep(.el-button--default:hover) {
  background: var(--accent-light);
}

.character-delete-btn,
.character-actionbar :deep(.el-button--danger.is-plain) {
  background: color-mix(in srgb, var(--error-light) 72%, var(--bg-card) 28%);
  border-color: color-mix(in srgb, var(--error) 26%, var(--border-primary) 74%);
  color: var(--error);
}

.character-delete-btn:hover,
.character-actionbar :deep(.el-button--danger.is-plain:hover) {
  background: color-mix(in srgb, var(--error-light) 88%, var(--bg-card) 12%);
  border-color: var(--error);
  color: var(--error);
}

.character-grid-dark {
  --character-dark-bg: #ffffff;
  --character-dark-panel: #ffffff;
  --character-dark-panel-strong: #f8fafc;
  --character-dark-border: rgba(15, 23, 42, 0.08);
  --character-dark-border-strong: rgba(15, 23, 42, 0.12);
  --character-dark-text: #0f172a;
  --character-dark-muted: rgba(15, 23, 42, 0.64);
  --character-dark-soft: rgba(15, 23, 42, 0.06);
  --character-dark-accent: #3b82f6;
}

.dark-character-pipeline .dark-character-guide {
  background:
    radial-gradient(circle at top right, rgba(59, 130, 246, 0.08), transparent 36%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(248, 250, 252, 0.98));
  border-color: rgba(15, 23, 42, 0.08);
  box-shadow: 0 18px 36px rgba(15, 23, 42, 0.08);
}

.dark-character-pipeline .generation-guide-title,
.dark-character-pipeline .character-guide-stat strong {
  color: #0f172a;
}

.dark-character-pipeline .generation-guide-text,
.dark-character-pipeline .character-guide-stat span,
.dark-character-pipeline .character-guide-meta span {
  color: rgba(15, 23, 42, 0.64);
}

.dark-character-pipeline .character-guide-meta span,
.dark-character-pipeline .character-guide-stat {
  border-color: rgba(15, 23, 42, 0.08);
  background: rgba(248, 250, 252, 0.92);
}

.character-card-dark {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.99), rgba(249, 250, 251, 0.99));
  border-color: rgba(15, 23, 42, 0.08);
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
}

.character-card-dark:hover {
  border-color: rgba(15, 23, 42, 0.12);
  box-shadow: 0 22px 48px rgba(15, 23, 42, 0.12);
}

.character-card-dark.is-generating {
  border-color: rgba(59, 130, 246, 0.45);
}

.character-card-dark.is-failed {
  border-color: rgba(248, 113, 113, 0.45);
}

.character-dark-card {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 18px;
  color: var(--character-dark-text);
}

.character-dark-head {
  display: grid;
  grid-template-columns: 82px minmax(0, 1fr);
  gap: 14px;
  align-items: start;
}

.character-dark-avatar,
.character-sheet-panel {
  border: 0;
  padding: 0;
  margin: 0;
  appearance: none;
}

.character-dark-avatar {
  width: 82px;
  height: 82px;
  border-radius: 24px;
  overflow: hidden;
  background: linear-gradient(180deg, rgba(248, 250, 252, 0.98), rgba(241, 245, 249, 0.98));
  border: 1px solid var(--character-dark-border);
  cursor: zoom-in;
}

.character-dark-avatar-image,
.character-sheet-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.character-dark-avatar-placeholder,
.character-sheet-placeholder,
.character-sheet-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 10px;
}

.character-dark-avatar-placeholder {
  width: 100%;
  height: 100%;
  color: var(--character-dark-muted);
}

.character-dark-copy {
  min-width: 0;
}

.character-dark-title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.character-dark-title-row h4 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--character-dark-text);
}

.character-card-icon-btn {
  color: var(--character-dark-muted);
}

.character-card-icon-btn:hover {
  color: #0f172a;
}

.character-dark-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.character-dark-tag {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.16);
  color: #1d4ed8;
  font-size: 0.72rem;
  font-weight: 700;
}

.character-dark-summary {
  margin: 12px 0 0;
  color: var(--character-dark-muted);
  font-size: 0.84rem;
  line-height: 1.7;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  overflow: hidden;
}

.character-sheet-panel {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 260px;
  border-radius: 24px;
  background:
    radial-gradient(circle at top, rgba(59, 130, 246, 0.06), transparent 40%),
    linear-gradient(180deg, rgba(248, 250, 252, 0.98), rgba(241, 245, 249, 0.98));
  border: 1px solid var(--character-dark-border);
  cursor: zoom-in;
  overflow: hidden;
  position: relative;
}

.character-sheet-panel.is-empty {
  background: linear-gradient(180deg, rgba(248, 250, 252, 0.98), rgba(241, 245, 249, 0.98));
}

.character-sheet-panel.is-loading {
  box-shadow: inset 0 0 0 1px rgba(59, 130, 246, 0.28);
}

.character-sheet-image {
  max-height: 260px;
  padding: 16px;
}

.character-sheet-placeholder,
.character-sheet-loading {
  width: 100%;
  height: 100%;
  color: var(--character-dark-muted);
  font-size: 0.82rem;
}

.character-sheet-placeholder .el-icon {
  font-size: 28px;
}

.character-dark-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.character-dark-status {
  min-width: 0;
}

.character-dark-status-label {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  background: var(--character-dark-soft);
  color: var(--character-dark-text);
  font-size: 0.72rem;
  font-weight: 700;
}

.character-dark-status p {
  margin: 10px 0 0;
  color: var(--character-dark-muted);
  font-size: 0.78rem;
  line-height: 1.6;
}

.character-dark-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 10px;
}

.character-action-primary,
.character-action-secondary {
  min-height: 38px;
  padding: 0 16px;
  border-radius: 12px;
  font-weight: 700;
}

.character-action-primary:deep(.el-button) {
  font-weight: 700;
}

.character-action-secondary {
  background: rgba(248, 250, 252, 0.98);
  border-color: rgba(15, 23, 42, 0.08);
  color: #0f172a;
}

.character-action-secondary:hover {
  background: rgba(241, 245, 249, 0.98);
  border-color: rgba(15, 23, 42, 0.12);
  color: #0f172a;
}

@media (max-width: 768px) {
  .character-dark-head {
    grid-template-columns: 72px minmax(0, 1fr);
  }

  .character-dark-avatar {
    width: 72px;
    height: 72px;
    border-radius: 20px;
  }

  .character-sheet-panel {
    min-height: 220px;
  }

  .character-sheet-image {
    max-height: 220px;
  }

  .character-dark-footer {
    flex-direction: column;
    align-items: stretch;
  }

  .character-dark-actions {
    justify-content: stretch;
  }

  .character-dark-actions :deep(.el-button) {
    flex: 1 1 auto;
  }
}

@media (max-width: 1200px) {
  .character-workbench {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .character-shell {
    padding: 16px;
  }

  .character-stage {
    grid-template-columns: 1fr;
  }

  .character-view-grid,
  .character-task-track {
    grid-template-columns: 1fr;
  }
}

.empty-icon {
  color: var(--accent);
}

/* ========================================
   Dark Mode / 深色模式
   ======================================== */
.dark .tabs-wrapper {
  background: var(--bg-card);
}

.dark :deep(.el-card) {
  background: var(--bg-card);
  border-color: var(--border-primary);
}

.dark :deep(.el-card__header) {
  background: var(--bg-secondary);
  border-color: var(--border-primary);
}

.dark :deep(.el-table) {
  background: var(--bg-card);
  --el-table-bg-color: var(--bg-card);
  --el-table-tr-bg-color: var(--bg-card);
  --el-table-header-bg-color: var(--bg-secondary);
  --el-fill-color-lighter: var(--bg-secondary);
}

.dark :deep(.el-table th),
.dark :deep(.el-table tr) {
  background: var(--bg-card);
}

.dark :deep(.el-table td),
.dark :deep(.el-table th) {
  border-color: var(--border-primary);
}

.dark :deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background: var(--bg-secondary);
}

.dark :deep(.el-table__body tr:hover > td) {
  background: var(--bg-card-hover) !important;
}

.dark :deep(.el-descriptions) {
  background: var(--bg-card);
}

.dark :deep(.el-descriptions__label) {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border-color: var(--border-primary);
}

.dark :deep(.el-descriptions__content) {
  background: var(--bg-card);
  color: var(--text-primary);
  border-color: var(--border-primary);
}

.dark :deep(.el-descriptions__cell) {
  border-color: var(--border-primary);
}

/* ========================================
   Project Info Card / 项目信息卡片
   ======================================== */
.project-info-card {
  margin-top: var(--space-5);
  border-radius: var(--radius-lg);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.project-descriptions {
  width: 100%;
}

:deep(.project-descriptions .el-descriptions__label) {
  width: 120px;
  font-weight: 500;
  color: var(--text-secondary);
}

:deep(.project-descriptions .el-descriptions__content) {
  min-width: 150px;
}

.info-value {
  font-weight: 500;
  color: var(--text-primary);
}

.info-desc {
  color: var(--text-secondary);
  line-height: 1.6;
}

.dark :deep(.el-dialog) {
  background: var(--bg-card);
}

.dark :deep(.el-dialog__header) {
  background: var(--bg-card);
}

.dark :deep(.el-form-item__label) {
  color: var(--text-primary);
}

.dark :deep(.el-input__wrapper) {
  background: var(--bg-secondary);
  box-shadow: 0 0 0 1px var(--border-primary) inset;
}

.dark :deep(.el-input__inner) {
  color: var(--text-primary);
}

.dark :deep(.el-textarea__inner) {
  background: var(--bg-secondary);
  color: var(--text-primary);
  box-shadow: 0 0 0 1px var(--border-primary) inset;
}
</style>
