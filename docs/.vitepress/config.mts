import {defineConfig} from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
    title: "JSVM",
    description: "",
    base: '/docs/',
    lang: 'zh-CN',
    themeConfig: {
        // 自定义文档底部按钮的文本
        docFooter: {
            prev: '上一页',
            next: '下一页'
        },
        outline: {
            label: '页面导航' // 右侧大纲的标题（默认是 On this page）
        },
        lastUpdated: {
            text: '最后更新于', // 只有在配置了 lastUpdated 时生效
            formatOptions: {
                dateStyle: 'full',
                timeStyle: 'medium'
            }
        },
        returnToTopLabel: '返回顶部',
        sidebarMenuLabel: '菜单',
        darkModeSwitchLabel: '主题模式',
        lightModeSwitchTitle: '切换到浅色模式',
        darkModeSwitchTitle: '切换到深色模式',

        // https://vitepress.dev/reference/default-theme-config
        // nav: [
        //     {text: '主页', link: '/'},
        // ],
        search: {
            provider: 'local',
        },
        sidebar: {
            '/modules/': [
                {
                    text: "模块",
                    items: [
                        {
                            collapsed: false,
                            items: [
                                {text: 'console', link: '/modules/console'},
                                {text: 'context', link: '/modules/context'},
                                {text: 'time', link: '/modules/time'},
                                {text: 'net/http', link: '/modules/net/http'},
                            ]
                        }
                    ]
                }
            ]
        },

        socialLinks: [
            {icon: 'github', link: 'https://github.com/xmx/jsvm'}
        ]
    }
})
