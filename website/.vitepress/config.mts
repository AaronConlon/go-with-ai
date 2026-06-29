import { defineConfig } from 'vitepress'

const lastUpdated = {
  text: '最后更新',
  formatOptions: {
    dateStyle: 'medium',
    timeStyle: 'short',
    forceLocale: true
  }
} as const

const zhThemeConfig = {
  logo: '/images/project-logo-192.png',
  nav: [
    { text: '首页', link: '/' },
    { text: '学习路径', link: '/learning/' },
    { text: '知识库', link: '/notes/' },
    { text: '完整代码', link: '/code/' },
    { text: '项目实践', link: '/projects/hn-agent' },
    { text: '背景资料', link: '/background/' }
  ],
  sidebar: {
    '/learning/': [
      {
        text: '学习路径',
        items: [
          { text: '总览', link: '/learning/' },
          { text: '六阶段路线', link: '/learning/roadmap' },
          { text: 'JS 到 Go', link: '/learning/js-to-go' }
        ]
      }
    ],
    '/notes/': [
      {
        text: '六阶段知识库',
        items: [
          { text: '总览', link: '/notes/' },
          { text: '阶段 1：启动期', link: '/notes/stage-1-startup/' },
          { text: '阶段 1：教练笔记', link: '/notes/stage-1-startup/coach-notes' },
          { text: '阶段 1：测试语法拆解', link: '/notes/stage-1-startup/go-test-syntax' },
          { text: '阶段 1：变量定义语法', link: '/notes/stage-1-startup/go-variable-syntax' },
          { text: '阶段 2：网络期', link: '/notes/stage-2-network/' },
          { text: '阶段 2：教练笔记', link: '/notes/stage-2-network/coach-notes' },
          { text: '阶段 2：Client 构造语法', link: '/notes/stage-2-network/go-client-constructor-syntax' },
          { text: '阶段 2：Method Receiver', link: '/notes/stage-2-network/go-method-receiver-syntax' },
          { text: '阶段 2：httptest 使用', link: '/notes/stage-2-network/go-httptest-syntax' },
          { text: '阶段 3：并发期', link: '/notes/stage-3-concurrency/' },
          { text: '阶段 3：教练笔记', link: '/notes/stage-3-concurrency/coach-notes' },
          { text: '阶段 4：服务期', link: '/notes/stage-4-service/' },
          { text: '阶段 4：教练笔记', link: '/notes/stage-4-service/coach-notes' },
          { text: '阶段 4：context 入门', link: '/notes/stage-4-service/go-context' },
          { text: '阶段 5：AI 整合期', link: '/notes/stage-5-ai-integration/' },
          { text: '阶段 6：验证与部署期', link: '/notes/stage-6-delivery/' }
        ]
      },
      {
        text: '通用卡片',
        items: [
          { text: '学习推进记录', link: '/notes/learning-progress' },
          { text: '学习问题对话记录', link: '/notes/learning-dialogues' },
          { text: 'Go Module', link: '/notes/go-module' },
          { text: '学习日志模板', link: '/notes/learning-log-template' },
          { text: '阶段记录模板', link: '/notes/stage-record-template' }
        ]
      }
    ],
    '/code/': [
      {
        text: '完整代码',
        items: [
          { text: '总览', link: '/code/' },
          { text: '阶段 1：启动期', link: '/code/stage-1-startup' },
          { text: '阶段 2：网络期', link: '/code/stage-2-network' },
          { text: '阶段 3：并发期', link: '/code/stage-3-concurrency' },
          { text: '阶段 4：服务期', link: '/code/stage-4-service' },
          { text: '阶段 5：AI 整合期', link: '/code/stage-5-ai-integration' },
          { text: '阶段 6：验证与部署期', link: '/code/stage-6-delivery' }
        ]
      }
    ],
    '/projects/': [
      {
        text: '项目实践',
        items: [
          { text: 'HN 摘要推送 Agent', link: '/projects/hn-agent' }
        ]
      }
    ],
    '/background/': [
      {
        text: '背景资料',
        items: [
          { text: '总览', link: '/background/' },
          { text: '研究摘要', link: '/background/research-summary' }
        ]
      }
    ]
  },
  socialLinks: [],
  outline: {
    label: '本页目录',
    level: [2, 3]
  },
  docFooter: {
    prev: '上一篇',
    next: '下一篇'
  },
  lastUpdated,
  search: {
    provider: 'local',
    options: {
      translations: {
        button: {
          buttonText: '搜索文档',
          buttonAriaLabel: '搜索文档'
        },
        modal: {
          displayDetails: '显示详情',
          resetButtonTitle: '清除搜索',
          backButtonTitle: '返回',
          noResultsText: '没有找到结果',
          footer: {
            selectText: '选择',
            selectKeyAriaLabel: '回车',
            navigateText: '切换',
            navigateUpKeyAriaLabel: '向上',
            navigateDownKeyAriaLabel: '向下',
            closeText: '关闭',
            closeKeyAriaLabel: 'Esc'
          }
        }
      }
    }
  }
} as const

const enThemeConfig = {
  logo: '/images/project-logo-192.png',
  nav: [
    { text: 'Home', link: '/en/' },
    { text: 'Learning Path', link: '/en/learning/' },
    { text: 'Knowledge Base', link: '/en/notes/' },
    { text: 'Full Code', link: '/en/code/' },
    { text: 'Project', link: '/en/projects/hn-agent' },
    { text: 'Background', link: '/en/background/' }
  ],
  sidebar: {
    '/en/learning/': [
      {
        text: 'Learning Path',
        items: [
          { text: 'Overview', link: '/en/learning/' },
          { text: 'Six-stage Roadmap', link: '/en/learning/roadmap' },
          { text: 'JS to Go', link: '/en/learning/js-to-go' }
        ]
      }
    ],
    '/en/notes/': [
      {
        text: 'Knowledge Base',
        items: [
          { text: 'Overview', link: '/en/notes/' },
          { text: 'Stage 1: Startup', link: '/en/notes/stage-1-startup/' },
          { text: 'Stage 2: Networking', link: '/en/notes/stage-2-network/' },
          { text: 'Stage 3: Concurrency', link: '/en/notes/stage-3-concurrency/' },
          { text: 'Stage 4: Service', link: '/en/notes/stage-4-service/' },
          { text: 'Stage 4: Coach Notes', link: '/en/notes/stage-4-service/coach-notes' },
          { text: 'Stage 4: Context Basics', link: '/en/notes/stage-4-service/go-context' },
          { text: 'Stage 5: AI Integration', link: '/en/notes/stage-5-ai-integration/' },
          { text: 'Stage 6: Delivery', link: '/en/notes/stage-6-delivery/' }
        ]
      }
    ],
    '/en/code/': [
      {
        text: 'Full Code',
        items: [
          { text: 'Overview', link: '/en/code/' },
          { text: 'Stage 1: Startup', link: '/en/code/stage-1-startup' },
          { text: 'Stage 2: Networking', link: '/en/code/stage-2-network' },
          { text: 'Stage 3: Concurrency', link: '/en/code/stage-3-concurrency' },
          { text: 'Stage 4: Service', link: '/en/code/stage-4-service' },
          { text: 'Stage 5: AI Integration', link: '/en/code/stage-5-ai-integration' },
          { text: 'Stage 6: Delivery', link: '/en/code/stage-6-delivery' }
        ]
      }
    ],
    '/en/projects/': [
      {
        text: 'Project',
        items: [
          { text: 'HN Digest Agent', link: '/en/projects/hn-agent' }
        ]
      }
    ],
    '/en/background/': [
      {
        text: 'Background',
        items: [
          { text: 'Overview', link: '/en/background/' },
          { text: 'Research Summary', link: '/en/background/research-summary' }
        ]
      }
    ]
  },
  socialLinks: [],
  outline: {
    label: 'On This Page',
    level: [2, 3]
  },
  docFooter: {
    prev: 'Previous',
    next: 'Next'
  },
  lastUpdated: {
    text: 'Last updated',
    formatOptions: {
      dateStyle: 'medium',
      timeStyle: 'short',
      forceLocale: true
    }
  },
  search: {
    provider: 'local'
  }
} as const

export default defineConfig({
  lang: 'zh-CN',
  title: 'Go With AI',
  description: '用项目驱动学习 Go，并沉淀可复用的工程知识。',
  cleanUrls: true,
  lastUpdated: true,
  head: [
    ['link', { rel: 'icon', href: '/favicon.ico', sizes: 'any' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32.png' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16.png' }],
    ['link', { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' }]
  ],
  vite: {
    publicDir: 'static'
  },
  locales: {
    root: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Go With AI',
      description: '用项目驱动学习 Go，并沉淀可复用的工程知识。',
      themeConfig: zhThemeConfig
    },
    en: {
      label: 'English',
      lang: 'en-US',
      title: 'Go With AI',
      description: 'Learn Go through a real project and keep reusable engineering notes.',
      themeConfig: enThemeConfig
    }
  },
  themeConfig: zhThemeConfig
})
