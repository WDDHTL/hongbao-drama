export interface VideoLanguageOption {
  value: string
  labelKey: string
  noteKey: string
}

export const videoLanguageOptions: VideoLanguageOption[] = [
  {
    value: 'zh-CN',
    labelKey: 'drama.videoLanguages.zhCN',
    noteKey: 'drama.videoLanguageNotes.zhCN'
  },
  {
    value: 'en-US',
    labelKey: 'drama.videoLanguages.enUS',
    noteKey: 'drama.videoLanguageNotes.enUS'
  },
  {
    value: 'ja-JP',
    labelKey: 'drama.videoLanguages.jaJP',
    noteKey: 'drama.videoLanguageNotes.jaJP'
  },
  {
    value: 'ko-KR',
    labelKey: 'drama.videoLanguages.koKR',
    noteKey: 'drama.videoLanguageNotes.koKR'
  }
]
