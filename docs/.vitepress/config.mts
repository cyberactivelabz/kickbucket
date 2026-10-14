import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'BucketMap',
  description: 'The Nmap for cloud storage',
  cleanUrls: true,
  themeConfig: {
    nav: [
      { text: 'Website', link: 'https://bucketmap.dev' },
      { text: 'GitHub', link: 'https://github.com/cyberactivelabz/kickbucket' },
      { text: 'Sponsor', link: 'https://github.com/sponsors/cyberactivelabz' },
    ],
    sidebar: [
      { text: 'Getting started', link: '/getting-started' },
      {
        text: 'Usage',
        items: [
          { text: 'CLI', link: '/cli' },
          { text: 'Web GUI', link: '/web-gui' },
          { text: 'Desktop app', link: '/desktop' },
        ],
      },
      { text: 'Go library', link: '/api' },
      { text: 'Wordlists', link: '/wordlists' },
      { text: 'Contributing', link: '/contributing' },
    ],
    search: { provider: 'local' },
    editLink: {
      pattern: 'https://github.com/cyberactivelabz/kickbucket/edit/main/docs/:path',
    },
    socialLinks: [{ icon: 'github', link: 'https://github.com/cyberactivelabz/kickbucket' }],
  },
})
