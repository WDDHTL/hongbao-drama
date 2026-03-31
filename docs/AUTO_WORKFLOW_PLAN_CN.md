# Huobao Drama 全自动生产流程改造方案

## 1. 目标

目标不是继续堆按钮，而是把当前项目改造成“一键跑完整集”的自动化生产系统。

用户最终只需要做两件事：

1. 创建项目并输入剧本
2. 点击“开始全自动生成”

系统随后自动完成：

1. 按集拆分剧本
2. 提取角色、场景、道具
3. 固定角色定妆和场景主参考
4. 生成每集分镜
5. 生成每个分镜图片
6. 生成每个分镜视频
7. 自动按集合成成片
8. 对失败节点重试或挂起等待人工补图

这份文档基于当前仓库结构编写，不是抽象方案。

## 2. 当前项目现状

当前仓库已经具备较完整的“分段能力”，但还不是“流程系统”。

现有能力：

1. 文本生成
   - [script_generation_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/script_generation_service.go)
   - 已支持角色提取等异步文本任务
2. 分镜生成
   - [storyboard_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/storyboard_service.go)
   - 已支持从剧本生成分镜
3. 角色/场景/道具生图
   - [character_library_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/character_library_service.go)
   - [image_generation_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/image_generation_service.go)
4. 分镜视频生成
   - [video_generation_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/video_generation_service.go)
5. 时间线合成
   - [video_merge_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/video_merge_service.go)
6. 异步任务基础设施
   - [task_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/task_service.go)
   - [task.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/domain/models/task.go)

当前核心问题：

1. 每个能力是分散入口，缺少统一编排器
2. 任务只知道“单步成功/失败”，不知道“整集跑到哪一步”
3. 角色、场景、分镜、视频之间缺少稳定的引用策略
4. 自动参考策略和手动编辑策略混在一起，导致结果不稳定
5. 失败任务没有系统化重试和回退
6. UI 还是“工具台”，不是“生产流水线”

## 3. 最终要做成什么样

建议把系统改造成两层流程：

### 3.1 项目级自动流程

输入：完整剧本

输出：

1. 项目
2. 集列表
3. 项目级角色库
4. 项目级场景库
5. 项目级道具库
6. 每集分镜

### 3.2 每集自动流程

输入：某一集剧本 + 项目级角色/场景/道具基准资产

输出：

1. 每镜头图片
2. 每镜头视频
3. 每集合成视频
4. 每集失败项列表

## 4. 推荐的目标流程

### 4.1 全局流程

```text
输入剧本
  -> 创建项目
  -> 拆分集数
  -> 提取角色/场景/道具
  -> 生成角色基准图
  -> 生成场景主图
  -> 逐集生成分镜
  -> 逐镜头生成图片
  -> 逐镜头生成视频
  -> 每集自动合成
  -> 输出整集结果与异常列表
```

### 4.2 每镜头的视频参考策略

这条必须硬编码成默认策略，不能只靠用户手选：

1. 优先 `first_last`
2. 首帧优先级：
   - 上一镜头末帧
   - 当前镜头首帧
   - 当前角色主参考图
3. 末帧优先级：
   - 当前镜头尾帧
   - 当前镜头关键帧
   - 当前场景主构图
4. 多图模式至少三张：
   - 角色主参考
   - 上一镜头末帧
   - 当前镜头目标构图
5. 纯文本模式只允许在“确实没有任何可用参考图”时启用

## 5. 数据层改造

当前模型可以继续复用，但不够支撑“整套流程编排”。建议新增下列数据结构。

### 5.1 新增 `workflow_runs`

用途：记录一次项目级或集级自动生产任务。

建议字段：

1. `id`
2. `drama_id`
3. `episode_id`
4. `scope`
   - `project`
   - `episode`
5. `status`
   - `pending`
   - `processing`
   - `completed`
   - `failed`
   - `paused`
6. `current_stage`
7. `progress`
8. `config_json`
9. `result_json`
10. `error`

建议新文件：

1. [workflow_run.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/domain/models/workflow_run.go)
2. [workflow_run_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/workflow_run_service.go)

### 5.2 新增 `workflow_step_runs`

用途：记录每一步子任务状态。

建议 stage：

1. `episode_split`
2. `character_extract`
3. `scene_extract`
4. `prop_extract`
5. `character_baseline_generate`
6. `scene_master_generate`
7. `storyboard_generate`
8. `storyboard_image_generate`
9. `storyboard_video_generate`
10. `episode_merge`
11. `qa_check`

### 5.3 角色模型增加“固定基准图槽位”

当前 [drama.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/domain/models/drama.go) 里的 `Character.ReferenceImages` 已可复用，但建议增加结构化语义，不要再只靠数组顺序。

建议存为：

```json
{
  "front": "...",
  "side": "...",
  "back": "...",
  "emotion": "...",
  "main": "...",
  "seed": "...",
  "costume_lock": "...",
  "appearance_lock": "..."
}
```

最少需要：

1. `front`
2. `side`
3. `back`
4. `main`

### 5.4 场景模型增加“主场景参考”

当前 `Scene.ImageURL` 不够表达“这个场景是本项目该场景的标准底图”。

建议增加：

1. `master_image_url`
2. `master_local_path`
3. `reference_images`
4. `continuity_prompt`

### 5.5 分镜增加流程状态

当前 `Storyboard.Status` 过粗。

建议改成或补充：

1. `image_status`
2. `video_status`
3. `merge_status`
4. `retry_count`
5. `workflow_lock_version`

## 6. 服务层改造

这是本项目最关键的部分。

### 6.1 新增统一编排服务 `AutoWorkflowService`

建议新增：

1. [auto_workflow_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/auto_workflow_service.go)

职责：

1. 接收“开始项目全自动生成”
2. 创建 `workflow_run`
3. 按阶段顺序执行
4. 失败即暂停并记录
5. 支持从失败节点继续

伪流程：

```go
RunProjectWorkflow(dramaID)
  -> split episodes
  -> extract characters/scenes/props
  -> generate character baselines
  -> generate scene masters
  -> for episode in episodes:
       RunEpisodeWorkflow(episode.ID)
```

```go
RunEpisodeWorkflow(episodeID)
  -> generate storyboards
  -> normalize storyboard refs
  -> generate storyboard images
  -> generate storyboard videos
  -> merge episode
  -> qa check
```

### 6.2 角色提取和角色定妆解耦

当前 [script_generation_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/script_generation_service.go) 只负责角色提取。

建议拆成两个阶段：

1. 角色提取
2. 角色基准图生成

新增约束：

1. 每个角色默认生成三张白底图：
   - 正面
   - 侧面
   - 背面
2. 正面图自动成为 `main`
3. 角色基准图生成失败时，整个后续视频链路不允许继续跑

### 6.3 场景提取和场景主图生成解耦

当前场景更多是作为分镜关联背景使用。

建议补充：

1. 项目级场景抽取
2. 场景去重归一
3. 每个场景生成一张主场景图
4. 同场景镜头默认绑定同一主图

### 6.4 分镜生成必须变成“可生产分镜”

当前 [storyboard_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/storyboard_service.go) 已经能出较细分镜，但还要继续收紧成可执行格式。

必须保证每个分镜稳定输出：

1. `characters`
2. `scene_id`
3. `shot_type`
4. `angle`
5. `movement`
6. `action`
7. `result`
8. `duration`

并新增两个系统字段：

1. `image_generation_prompt`
2. `video_generation_prompt`

不要在后续步骤再临时拼很多不稳定文本。

### 6.5 图片生成改成自动批量流水

当前 [image_generation_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/image_generation_service.go) 支持单次生成，但还需要加批处理入口。

建议新增能力：

1. `GenerateCharacterBaselines(dramaID)`
2. `GenerateSceneMasters(dramaID)`
3. `GenerateStoryboardImagesForEpisode(episodeID)`

其中分镜图片策略：

1. 每镜头优先生成 `first`
2. 如有必要生成 `key`
3. 对需要强连续的镜头生成 `last`

### 6.6 视频生成改成批量有状态流水

当前 [video_generation_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/video_generation_service.go) 已支持单镜头任务和参考图模式，但缺少自动批量编排。

建议新增：

1. `GenerateStoryboardVideosForEpisode(episodeID, workflowID)`
2. `BuildAutomaticReferencePlan(storyboardID)`
3. `ShouldUseCharacterReferenceOnly(prompt, storyboard)`

视频生成默认策略：

1. 时长强制限制在 `3-5s`
2. 默认镜头运动：
   - `static`
   - `slow_pan`
   - `slow_push`
3. 连续镜头自动继承上一镜头的：
   - `provider`
   - `model`
   - `aspect_ratio`
   - `duration`
   - `motion_level`
   - `camera_motion`

### 6.7 合成服务支持“自动按集输出”

当前 [video_merge_service.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/application/services/video_merge_service.go) 偏手动时间线。

建议补充：

1. `MergeEpisodeAutomatically(episodeID)`
2. 自动按分镜顺序合成
3. 自动应用默认转场模板
4. 自动补空镜头占位和失败提示

## 7. API 层改造

建议新增一组工作流 API，而不是继续把入口塞在现有零散 handler 里。

新增建议：

1. [workflow.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/api/handlers/workflow.go)
2. [workflow.go](C:/Users/frank/Desktop/ai短剧/huobao-drama/api/routes/workflow.go)

建议接口：

1. `POST /api/v1/workflows/project/start`
2. `POST /api/v1/workflows/project/:dramaId/resume`
3. `POST /api/v1/workflows/project/:dramaId/pause`
4. `GET /api/v1/workflows/project/:dramaId/status`
5. `POST /api/v1/workflows/episode/:episodeId/start`
6. `POST /api/v1/workflows/episode/:episodeId/resume`
7. `GET /api/v1/workflows/episode/:episodeId/status`
8. `POST /api/v1/workflows/storyboard/:storyboardId/retry`

## 8. 前端改造

### 8.1 新增“全自动工作流”主页面

建议基于现有页面新增：

1. [DramaWorkflow.vue](C:/Users/frank/Desktop/ai短剧/huobao-drama/web/src/views/drama/DramaWorkflow.vue)

目标不是给用户更多按钮，而是只保留：

1. 输入剧本
2. 选择模型模板
3. 开始自动生成
4. 查看进度
5. 处理失败项

页面结构建议：

1. 顶部：项目总进度
2. 中间：每集卡片
3. 每集卡片显示：
   - 分镜数
   - 图片完成数
   - 视频完成数
   - 当前阶段
   - 错误项数
4. 底部：异常队列

### 8.2 保留 `ProfessionalEditor.vue` 作为“高级修正台”

[ProfessionalEditor.vue](C:/Users/frank/Desktop/ai短剧/huobao-drama/web/src/views/drama/ProfessionalEditor.vue) 不应该再承担默认生产入口。

它的定位应该调整为：

1. 自动流程失败后的人工修正页
2. 高级用户精调页
3. 单镜头重跑页

### 8.3 `DramaManagement.vue` 改成资产总览

[DramaManagement.vue](C:/Users/frank/Desktop/ai短剧/huobao-drama/web/src/views/drama/DramaManagement.vue) 建议只做：

1. 角色基准图总览
2. 场景主图总览
3. 道具主图总览
4. 失败资产补跑

不要再把它当主流程页面。

## 9. 失败处理机制

如果没有这块，自动化一定跑不稳。

建议规则：

1. 单镜头图片失败
   - 自动重试 2 次
   - 仍失败则标记为 `manual_required`
2. 单镜头视频失败
   - 先自动回退参考模式
   - `first_last -> multiple -> single -> none`
   - 最多重试 3 次
3. 角色定妆失败
   - 终止后续视频流程
4. 场景主图失败
   - 允许该集暂停，等待人工补图

建议新增错误分类：

1. `provider_error`
2. `content_policy_error`
3. `invalid_prompt_error`
4. `missing_reference_error`
5. `merge_error`
6. `manual_required`

## 10. 模型策略

建议不要让用户在每一步都重新选模型。

应该在项目创建时配置一套“生产模板”：

1. 文本模型
2. 角色图模型
3. 场景图模型
4. 分镜图模型
5. 视频模型
6. 转场模板

建议配置位置：

1. `Drama.Metadata`
2. 或新增 `drama_generation_profiles`

一个项目只保留一套默认生产模板，人工页允许单次覆盖，但不改默认策略。

## 11. 建议的实施顺序

不要一次全改，按下面顺序最稳。

### 阶段 1：先做编排骨架

目标：

1. 能从一个入口顺序调用已有服务
2. 能看到整集进度
3. 能失败暂停

改动：

1. 新增 `workflow_runs`
2. 新增 `workflow_step_runs`
3. 新增 `AutoWorkflowService`
4. 新增 workflow API
5. 前端新增总进度页

### 阶段 2：固定角色和场景基准资产

目标：

1. 视频连续性稳定
2. 不再每步都漂

改动：

1. 角色三视图结构化
2. 场景主图结构化
3. 分镜自动绑定角色和场景参考

### 阶段 3：打通逐集自动生产

目标：

1. 每集自动从分镜跑到成片

改动：

1. 批量生图
2. 批量生视频
3. 自动合成
4. 失败重试

### 阶段 4：做真正可用的人工回退机制

目标：

1. 自动流程失败时，人工只处理异常点

改动：

1. 异常队列
2. 单镜头补跑
3. 手动补图后继续流程

## 12. 最关键的产品原则

这套系统后面要好用，必须坚持这几个原则：

1. 默认自动，不要默认手动
2. 参数继承，不要每镜头重配
3. 失败暂停，不要静默失败
4. 参考图固定，不要每次临时找图
5. 专业编辑页是修正台，不是主入口

## 13. 针对当前项目的最小可执行版本

如果要尽快落地，建议先做 MVP。

### MVP 范围

1. 输入完整剧本
2. 自动拆集
3. 自动提取角色/场景
4. 自动生成角色三视图和场景主图
5. 自动为每集生成分镜
6. 自动生成每镜头首帧图
7. 自动生成每镜头视频
8. 自动合成每集
9. 前端显示每集进度和失败项

### 暂时不做

1. 全自动配音
2. 全自动字幕
3. 全自动 BGM 混音
4. 全剧总片合成

这四项可以等主流程稳定后再补。

## 14. 结论

这个项目当前已经有“单点能力”，缺的是“流程引擎”。

要实现“只输入剧本，其余自动跑完整套流程”，核心不是继续补单个页面，而是：

1. 新增统一工作流编排层
2. 把角色和场景基准资产结构化
3. 把分镜、图片、视频、合成串成有状态流水线
4. 前端改成项目总控台，而不是若干零碎工具页

如果按这份文档推进，当前仓库是可以演进成真正的一键短剧生产系统的。
