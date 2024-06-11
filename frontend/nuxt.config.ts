// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: {
    enabled: true,

    timeline: {
      enabled: true,
    },
  },

  app: {
    head: {
      link: [
        {
          rel: "stylesheet",
          href: "https://api.mapbox.com/mapbox-gl-js/v2.8.1/mapbox-gl.css",
        },
      ],
    },
  },

  runtimeConfig: {
    public: {
      mapboxAccessToken: process.env.NUXT_MAPBOX_ACCESS_TOKEN,
    },
    githubClientId: process.env.NUXT_GITHUB_CLIENT_ID,
    githubClientSecret: process.env.NUXT_GITHUB_CLIENT_SECRET,
  },

  modules: ["@sidebase/nuxt-auth"],

  //@ts-ignore
  auth: {
    provider: {
      type: "authjs",
    },
  },

  css: [
    "~/public/css/main.css"
  ],

  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  routeRules: {
//     "/newpost/**": { ssr: false },
  },

  vite: {
    optimizeDeps: {
      include: ['vue-filepond']
    }
  }
});