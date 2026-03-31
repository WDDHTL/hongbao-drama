import request from '../utils/request'
import type { WorkflowRun, WorkflowStatusResponse } from '../types/workflow'

export const workflowAPI = {
  getProjectWorkflowStatus(dramaId: string) {
    return request.get<WorkflowStatusResponse>(`/dramas/${dramaId}/workflow`)
  },

  startProjectWorkflow(dramaId: string) {
    return request.post<WorkflowRun>(`/dramas/${dramaId}/workflow/start`)
  },

  resumeProjectWorkflow(dramaId: string) {
    return request.post<WorkflowRun>(`/dramas/${dramaId}/workflow/resume`)
  }
}
