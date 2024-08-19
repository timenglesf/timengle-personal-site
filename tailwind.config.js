/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'ui/template/**/*.templ',
    'cmd/web/admin_handlers.go',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        sans: ['Courier Prime', 'monospace'],
        poppins: ['Poppins', 'sans-serif'],
      },
      colors: {
        lavender1: '#a8a3c0',
        lavender2: '#a9a4c1',
        lavender3: '#a9a5bf',
        lavender4: '#a7a2bd',
        lavender5: '#a7a2bf',
        goBlue: '#5999b6',
      },
      height: {
        '2/3-screen': '66.67vh',
      },
      maxHeight: {
        '2/3-screen': '66.67vh',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/line-clamp'),
    require('@tailwindcss/aspect-ratio'),
    require('daisyui'),
    require('@tailwindcss/typography'),
  ],
  corePlugins: { preflight: true },
  safelist: [
    'max-w-screen-sm',
    'max-w-1/3',
    'mx-auto',
    'mb-6',
    'error-message',
    'container',
    'fill-base-100',
    'fill-base-100-focus',
    'fill-secondary',
    'fill-secondary-focus'
  ],
  daisyui: {
    themes: [
      "nord",
      {
        mytheme: {
          "primary": "#4a6572",
          "primary-focus": "#3a5260",

          "secondary": "#38b2ac",
          "secondary-focus": "#2e958f",

          "accent": "#5999b6",
          "accent-focus": "#487d95",

          "neutral": "#f5f5f5",
          "neutral-focus": "#e0e0e0",

          "base-100": "#a9a5bf",
          "base-100-focus": "#8e8aa3",

          "info": "#3b82f6",
          "info-focus": "#3169c4",

          "success": "#10b981",
          "success-focus": "#0e9668",

          "warning": "#f59e0b",
          "warning-focus": "#c47f09",

          "error": "#ef4444",
          "error-focus": "#c03636"
        },
      },
    ],
  },
}
