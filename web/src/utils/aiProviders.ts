import type { AIServiceType } from '@/types/ai'

export interface ProviderConfig {
  id: string
  name: string
  models: string[]
  defaultBaseUrl: string
  endpoint?: string
  queryEndpoint?: string
  description: string
  isCustom?: boolean
}

export const providerConfigs: Record<AIServiceType, ProviderConfig[]> = {
  text: [
    {
      id: 'chatfire',
      name: 'ChatFire',
      models: ['gemini-3-flash-preview', 'claude-sonnet-4-5-20250929', 'doubao-seed-1-8-251228'],
      defaultBaseUrl: 'https://api.chatfire.site/v1',
      endpoint: '/chat/completions',
      description: 'Unified text relay with mainstream model routing.'
    },
    {
      id: 'openai',
      name: 'OpenAI',
      models: ['gpt-5.2', 'gpt-4.1'],
      defaultBaseUrl: 'https://api.openai.com/v1',
      endpoint: '/chat/completions',
      description: 'Official OpenAI-compatible text endpoint.'
    },
    {
      id: 'gemini',
      name: 'Google Gemini',
      models: ['gemini-2.5-pro', 'gemini-3-flash-preview'],
      defaultBaseUrl: 'https://generativelanguage.googleapis.com',
      endpoint: '/v1beta/models/{model}:generateContent',
      description: 'Gemini native endpoint.'
    },
    {
      id: 'openclaw',
      name: 'OpenClaw',
      models: [],
      defaultBaseUrl: '',
      endpoint: '/chat/completions',
      description: 'OpenAI-compatible relay preset. Fill in your own base URL and model.'
    },
    {
      id: 'custom_relay',
      name: 'Custom Relay',
      models: [],
      defaultBaseUrl: '',
      endpoint: '/chat/completions',
      description: 'Bring your own relay station with custom endpoints.'
    }
  ],
  image: [
    {
      id: 'chatfire',
      name: 'ChatFire',
      models: ['nano-banana-pro', 'doubao-seedream-4-5-251128'],
      defaultBaseUrl: 'https://api.chatfire.site/v1',
      endpoint: '/images/generations',
      description: 'Managed image routing with multimodel support.'
    },
    {
      id: 'volcengine',
      name: 'Volcengine',
      models: ['doubao-seedream-4-5-251128', 'doubao-seedream-4-0-250828'],
      defaultBaseUrl: 'https://ark.cn-beijing.volces.com/api/v3',
      endpoint: '/images/generations',
      description: 'Native Doubao image endpoint.'
    },
    {
      id: 'openai',
      name: 'OpenAI',
      models: ['gpt-image-1', 'dall-e-3', 'dall-e-2'],
      defaultBaseUrl: 'https://api.openai.com/v1',
      endpoint: '/images/generations',
      description: 'Official OpenAI image generation.'
    },
    {
      id: 'gemini',
      name: 'Google Gemini',
      models: ['gemini-3-pro-image-preview'],
      defaultBaseUrl: 'https://generativelanguage.googleapis.com',
      endpoint: '/v1beta/models/{model}:generateContent',
      description: 'Gemini native multimodal image endpoint.'
    },
    {
      id: 'openclaw',
      name: 'OpenClaw',
      models: [],
      defaultBaseUrl: '',
      endpoint: '/images/generations',
      description: 'OpenAI-compatible image relay preset.'
    },
    {
      id: 'custom_relay',
      name: 'Custom Relay',
      models: [],
      defaultBaseUrl: '',
      endpoint: '/images/generations',
      description: 'Custom image relay with editable endpoint.'
    }
  ],
  video: [
    {
      id: 'chatfire',
      name: 'ChatFire',
      models: [
        'doubao-seedance-1-5-pro-251215',
        'doubao-seedance-1-0-lite-i2v-250428',
        'doubao-seedance-1-0-pro-fast-251015',
        'sora-2',
        'sora-2-pro'
      ],
      defaultBaseUrl: 'https://api.chatfire.site/v1',
      endpoint: '/video/generations',
      queryEndpoint: '/video/task/{taskId}',
      description: 'Managed video relay with OpenAI and Doubao models.'
    },
    {
      id: 'volces',
      name: 'Volcengine',
      models: [
        'doubao-seedance-1-5-pro-251215',
        'doubao-seedance-1-0-lite-i2v-250428',
        'doubao-seedance-1-0-pro-fast-251015'
      ],
      defaultBaseUrl: 'https://ark.cn-beijing.volces.com/api/v3',
      endpoint: '/contents/generations/tasks',
      queryEndpoint: '/contents/generations/tasks/{taskId}',
      description: 'Native first/last-frame friendly Doubao endpoint.'
    },
    {
      id: 'minimax',
      name: 'MiniMax',
      models: ['MiniMax-Hailuo-2.3', 'MiniMax-Hailuo-2.3-Fast', 'MiniMax-Hailuo-02'],
      defaultBaseUrl: 'https://api.minimaxi.com/v1',
      endpoint: '/video_generation',
      queryEndpoint: '/query/video_generation?task_id={taskId}',
      description: 'MiniMax video generation.'
    },
    {
      id: 'openai',
      name: 'OpenAI',
      models: ['sora-2', 'sora-2-pro'],
      defaultBaseUrl: 'https://api.openai.com/v1',
      endpoint: '/videos',
      queryEndpoint: '/videos/{taskId}',
      description: 'Official OpenAI Sora endpoint.'
    },
    {
      id: 'openclaw',
      name: 'OpenClaw',
      models: [],
      defaultBaseUrl: '',
      endpoint: '/videos',
      queryEndpoint: '/videos/{taskId}',
      description: 'OpenAI-compatible video relay preset.'
    },
    {
      id: 'custom_relay',
      name: 'Custom Relay',
      models: [],
      defaultBaseUrl: '',
      endpoint: '/video/generations',
      queryEndpoint: '/video/task/{taskId}',
      description: 'Custom video relay with editable create/query endpoints.'
    }
  ]
}

export const getProviderConfigs = (serviceType: AIServiceType) => providerConfigs[serviceType] || []

export const getProviderConfig = (serviceType: AIServiceType, providerId?: string) =>
  getProviderConfigs(serviceType).find((provider) => provider.id === providerId)

export const getProviderModels = (serviceType: AIServiceType, providerId?: string) =>
  getProviderConfig(serviceType, providerId)?.models || []

export const getProviderBaseUrl = (serviceType: AIServiceType, providerId?: string) =>
  getProviderConfig(serviceType, providerId)?.defaultBaseUrl || ''

export const getProviderEndpoint = (serviceType: AIServiceType, providerId?: string) =>
  getProviderConfig(serviceType, providerId)?.endpoint || ''

export const getProviderQueryEndpoint = (serviceType: AIServiceType, providerId?: string) =>
  getProviderConfig(serviceType, providerId)?.queryEndpoint || ''

export const buildConfigName = (providerId: string, serviceType: AIServiceType) => {
  const providerNameMap: Record<string, string> = {
    chatfire: 'ChatFire',
    openai: 'OpenAI',
    gemini: 'Gemini',
    openclaw: 'OpenClaw',
    custom_relay: 'Relay',
    volces: 'Volces',
    volcengine: 'Volcengine',
    minimax: 'MiniMax'
  }

  const serviceNameMap: Record<AIServiceType, string> = {
    text: 'Text',
    image: 'Image',
    video: 'Video'
  }

  const suffix = Math.floor(Math.random() * 10000)
    .toString()
    .padStart(4, '0')

  return `${providerNameMap[providerId] || providerId}-${serviceNameMap[serviceType]}-${suffix}`
}
