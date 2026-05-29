/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      fontFamily: {
        display: ['Orbitron', 'sans-serif'],
        mono:    ['DM Mono', 'monospace'],
      },
      colors: {
        space: {
          950: '#050810',
          900: '#080c14',
          800: '#0d1220',
          700: '#111827',
        },
        cyber: {
          DEFAULT: '#22d3ee',
          dim:     'rgba(34, 211, 238, 0.12)',
          glow:    'rgba(34, 211, 238, 0.25)',
          border:  'rgba(34, 211, 238, 0.35)',
        },
      },
      backgroundImage: {
        'grid-space': `
          linear-gradient(rgba(255,255,255,0.035) 1px, transparent 1px),
          linear-gradient(90deg, rgba(255,255,255,0.035) 1px, transparent 1px)
        `,
      },
      backgroundSize: {
        'grid-space': '48px 48px',
      },
      boxShadow: {
        'cyber':      '0 0 18px rgba(34, 211, 238, 0.18)',
        'cyber-lg':   '0 0 40px rgba(34, 211, 238, 0.22)',
        'glass':      '0 4px 24px rgba(0, 0, 0, 0.5)',
        'card-hover': '0 8px 32px rgba(0, 0, 0, 0.6), 0 0 20px rgba(34, 211, 238, 0.12)',
      },
      animation: {
        'cursor-blink': 'blink 1.1s step-end infinite',
        'fade-up':      'fadeUp 0.6s ease both',
        'scanline':     'scanline 8s linear infinite',
      },
      keyframes: {
        blink: {
          '0%, 100%': { opacity: '1' },
          '50%':      { opacity: '0' },
        },
        fadeUp: {
          from: { opacity: '0', transform: 'translateY(20px)' },
          to:   { opacity: '1', transform: 'translateY(0)' },
        },
        scanline: {
          '0%':   { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100vh)' },
        },
      },

      typography: (theme) => ({
        DEFAULT: {
          css: {
            '--tw-prose-body':          theme('colors.slate[400]'),
            '--tw-prose-headings':      theme('colors.slate[200]'),
            '--tw-prose-lead':          theme('colors.slate[400]'),
            '--tw-prose-links':         theme('colors.cyan[400]'),
            '--tw-prose-bold':          theme('colors.slate[200]'),
            '--tw-prose-counters':      theme('colors.slate[500]'),
            '--tw-prose-bullets':       theme('colors.slate[600]'),
            '--tw-prose-hr':            'rgba(255,255,255,0.08)',
            '--tw-prose-quotes':        theme('colors.slate[300]'),
            '--tw-prose-quote-borders': theme('colors.cyan[500]'),
            '--tw-prose-captions':      theme('colors.slate[600]'),
            '--tw-prose-code':          theme('colors.cyan[300]'),
            '--tw-prose-pre-code':      theme('colors.slate[300]'),
            '--tw-prose-pre-bg':        '#0d1117',
            '--tw-prose-th-borders':    'rgba(255,255,255,0.1)',
            '--tw-prose-td-borders':    'rgba(255,255,255,0.06)',
            'a': {
              'color':           theme('colors.cyan[400]'),
              'text-decoration': 'none',
              'border-bottom':   '1px solid rgba(34,211,238,0.25)',
              'transition':      'all 0.25s ease',
            },
            'a:hover': {
              'color':        theme('colors.cyan[300]'),
              'border-color': 'rgba(34,211,238,0.7)',
              'text-shadow':  '0 0 12px rgba(34,211,238,0.45)',
            },
            'code': {
              'background':    'rgba(34,211,238,0.08)',
              'border':        '1px solid rgba(34,211,238,0.18)',
              'border-radius': '4px',
              'padding':       '0.15em 0.45em',
              'font-size':     '0.875em',
              'color':         theme('colors.cyan[300]'),
            },
            'code::before': { content: '""' },
            'code::after':  { content: '""' },
            'pre': {
              'background':    '#0d111a',
              'border':        '1px solid rgba(255,255,255,0.07)',
              'border-radius': '10px',
              'padding':       '1.25rem 1.5rem',
            },
            'pre code': {
              'background': 'transparent',
              'border':     'none',
              'padding':    '0',
              'color':      'inherit',
              'font-size':  '0.875em',
            },
            'h1,h2,h3,h4': {
              'font-family':     '"Orbitron", sans-serif',
              'letter-spacing':  '0.04em',
            },
            'blockquote': {
              'border-left-color': theme('colors.cyan[500]'),
              'background':        'rgba(34,211,238,0.04)',
              'border-radius':     '0 6px 6px 0',
              'padding':           '0.75rem 1.25rem',
              'color':             theme('colors.slate[300]'),
            },
            'hr': {
              'border-color': 'rgba(255,255,255,0.07)',
              'margin':       '2.5rem 0',
            },
          },
        },
      }),
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
};
