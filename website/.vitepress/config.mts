import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'zh-CN',
  title: 'Go With AI',
  description: '用项目驱动学习 Go，并沉淀可复用的工程知识。',
  cleanUrls: true,
  lastUpdated: true,
  vite: {
    publicDir: 'static'
  },
  themeConfig: {
    logo: '/images/go-with-ai.svg',
    nav: [
      { text: '首页', link: '/' },
      { text: '学习路径', link: '/learning/' },
      { text: '知识库', link: '/notes/' },
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
            { text: '阶段 2：网络期', link: '/notes/stage-2-network/' },
            { text: '阶段 3：并发期', link: '/notes/stage-3-concurrency/' },
            { text: '阶段 4：服务期', link: '/notes/stage-4-service/' },
            { text: '阶段 5：AI 整合期', link: '/notes/stage-5-ai-integration/' },
            { text: '阶段 6：验证与部署期', link: '/notes/stage-6-delivery/' }
          ]
        },
        {
          text: '通用卡片',
          items: [
            { text: 'Go Module', link: '/notes/go-module' },
            { text: '学习日志模板', link: '/notes/learning-log-template' },
            { text: '阶段记录模板', link: '/notes/stage-record-template' }
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
    lastUpdated: {
      text: '最后更新',
      formatOptions: {
        dateStyle: 'medium',
        timeStyle: 'short'
      }
    },
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
  }
})
