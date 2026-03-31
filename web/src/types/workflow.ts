export interface WorkflowRun {
  id: number
  drama_id: number
  episode_id?: number
  scope: 'project' | 'episode'
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'paused'
  current_stage: string
  progress: number
  config_json?: any
  result_json?: any
  error_msg?: string
  started_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
}

export interface WorkflowStepRun {
  id: number
  workflow_run_id: number
  episode_id?: number
  stage: string
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'paused'
  progress: number
  message?: string
  error_msg?: string
  meta_json?: any
  started_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
}

export interface WorkflowEpisodeSummary {
  episode_id: number
  episode_number: number
  title: string
  status: string
  storyboard_count: number
  completed_images: number
  completed_videos: number
  failed_images: number
  failed_videos: number
  merge_status?: string
  final_video_url?: string
  needs_manual_review: boolean
}

export interface WorkflowStatusResponse {
  run: WorkflowRun | null
  steps: WorkflowStepRun[]
  episodes: WorkflowEpisodeSummary[]
}
