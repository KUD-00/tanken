module.exports = {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app.vue",
    "./error.vue",
  ],
  theme: {
    extend: {
      colors: {
        "light-green": "#c2ffc2",
        "dark-green": "#88b388",
        "light-blue": "#c2ffff",
        "deep-green": "#00b300",
        "very-light-green": "#e7ffc2",
      }
    },
  },
  plugins: [require("daisyui")],
}