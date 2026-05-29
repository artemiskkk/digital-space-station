import { defineConfig } from 'astro/config';
import tailwind from '@astrojs/tailwind';

export default defineConfig({
  integrations: [tailwind()],
  output: 'server',

  markdown: {
    shikiConfig: {
      // tokyo-night: 深邃的暗黑主题，与我们的设计风格完美契合
      theme: 'tokyo-night',
      // 为代码块添加行号（可选）
      wrap: false,
    },
  },
});
