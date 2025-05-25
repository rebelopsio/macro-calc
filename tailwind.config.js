/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/templates/**/*.templ",
    "./internal/templates/**/*.go",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        'army-green': {
          50: '#f4f5f0',
          100: '#e5e7db',
          200: '#cdd1bb',
          300: '#aeb495',
          400: '#909973',
          500: '#75805a',
          600: '#5a6345',
          700: '#454d36',
          800: '#393f2e',
          900: '#313629',
          950: '#191b13',
        },
      },
    },
  },
  plugins: [],
}