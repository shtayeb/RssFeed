/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./views/**/*.{templ,html}",
  ],
  theme: {
    extend: {},
  },
  plugins: [
require("@tailwindcss/forms")({
    strategy: 'base', // only generate global styles
    // strategy: 'class', // only generate classes
  }),
  ],
  
}
