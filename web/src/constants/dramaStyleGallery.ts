export interface DramaStyleGalleryItem {
  value: string
  labelKey: string
  summary: string
  prompt: string
  previewImage: string
  assetPath: string
}

export const dramaStyleGallery: DramaStyleGalleryItem[] = [
  {
    value: 'ghibli',
    labelKey: 'drama.styles.ghibli',
    summary: '温暖手绘、治愈自然、电影感构图。',
    prompt:
      '宫崎骏吉卜力风格，手绘动画电影质感，柔和自然光，层次丰富的云朵与远山，温暖色彩，童话感乡村场景，细腻笔触，电影级构图，高细节，高质量，无文字，无水印。',
    previewImage: '/style-previews/ghibli.svg',
    assetPath: 'web/public/style-previews/ghibli.svg'
  },
  {
    value: 'guoman',
    labelKey: 'drama.styles.guoman',
    summary: '国漫热血、强烈戏剧光影、角色感鲜明。',
    prompt:
      '高质量国漫风格，热血史诗气质，东方幻想角色与山岳云海场景，凌厉光影，动势强烈，细节丰富，锐利线条，电影镜头语言，海报级质感，无文字，无水印。',
    previewImage: '/style-previews/guoman.svg',
    assetPath: 'web/public/style-previews/guoman.svg'
  },
  {
    value: 'wasteland',
    labelKey: 'drama.styles.wasteland',
    summary: '末世废墟、风沙工业感、冷硬叙事气质。',
    prompt:
      '末日废土风格，荒芜城市废墟，破败高楼与沙尘天气，冷灰褐色调，工业残骸，压迫感天空，电影海报构图，真实细节，强氛围感，无文字，无水印。',
    previewImage: '/style-previews/wasteland.svg',
    assetPath: 'web/public/style-previews/wasteland.svg'
  },
  {
    value: 'nostalgia',
    labelKey: 'drama.styles.nostalgia',
    summary: '复古怀旧、胶片色偏、80/90 年代记忆感。',
    prompt:
      '怀旧复古风格，90年代东亚城市与室内生活场景，胶片颗粒感，暖黄色与墨绿色调，老电视机和霓虹灯元素，静谧叙事氛围，电影剧照构图，无文字，无水印。',
    previewImage: '/style-previews/nostalgia.svg',
    assetPath: 'web/public/style-previews/nostalgia.svg'
  },
  {
    value: 'pixel',
    labelKey: 'drama.styles.pixel',
    summary: '8-bit 像素世界、游戏地图感、清晰色块。',
    prompt:
      '8-bit 像素艺术风格，经典游戏场景，清晰像素块与分层地形，亮丽配色，角色与场景轮廓明确，复古街机视觉，高对比，高可读性，无文字，无水印。',
    previewImage: '/style-previews/pixel.svg',
    assetPath: 'web/public/style-previews/pixel.svg'
  },
  {
    value: 'voxel',
    labelKey: 'drama.styles.voxel',
    summary: '体素方块世界、立体模型感、低多边形趣味。',
    prompt:
      'voxel 体素风格，立体方块城市与自然场景，类似高质量沙盒游戏宣传图，低多边形体积感，干净光照，结构清晰，丰富层次，无文字，无水印。',
    previewImage: '/style-previews/voxel.svg',
    assetPath: 'web/public/style-previews/voxel.svg'
  },
  {
    value: 'urban',
    labelKey: 'drama.styles.urban',
    summary: '都市现实、现代建筑、夜景与冷调霓虹。',
    prompt:
      '现代都市电影风格，高楼街区、霓虹反光、雨夜道路与玻璃幕墙，冷黑白灰主调，局部霓虹点缀，时尚广告摄影构图，真实质感，无文字，无水印。',
    previewImage: '/style-previews/urban.svg',
    assetPath: 'web/public/style-previews/urban.svg'
  },
  {
    value: 'guoman3d',
    labelKey: 'drama.styles.guoman3d',
    summary: '国漫 3D、仙侠奇观、角色与场景更立体。',
    prompt:
      '高质量国漫3D风格，东方玄幻世界，浮空建筑、山河云海与能量光效，角色海报级站位，CG电影渲染质感，层次分明，宏大场景，无文字，无水印。',
    previewImage: '/style-previews/guoman3d.svg',
    assetPath: 'web/public/style-previews/guoman3d.svg'
  },
  {
    value: 'chibi3d',
    labelKey: 'drama.styles.chibi3d',
    summary: 'Q版 3D、玩具感材质、亲和可爱。',
    prompt:
      'Q版3D风格，软萌角色比例，玩具质感材质，干净布光，梦幻小镇或室内场景，可爱高饱和但不刺眼，商业动画宣传图效果，无文字，无水印。',
    previewImage: '/style-previews/chibi3d.svg',
    assetPath: 'web/public/style-previews/chibi3d.svg'
  }
]

export const dramaStyleGalleryMap = new Map(
  dramaStyleGallery.map((item) => [item.value, item] as const)
)
