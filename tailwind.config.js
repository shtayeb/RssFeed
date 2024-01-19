/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./views/**/*.{templ,html}",
  ],
  theme: {
    extend: {
      primary: "#000",
    },
  },
  plugins: [
require("@tailwindcss/forms")({
    strategy: 'base', // only generate global styles
    // strategy: 'class', // only generate classes
  }),
  ],
  
}
