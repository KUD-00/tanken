<template>
  <div class="flex justify-center min-h-screen bg-gray-100">
    <!-- Navigation wrapper with negative margin -->
    <div class="relative flex">
      <div class="flex flex-col mt-8">
        <div class="flex flex-col gap-8 shrink border rounded-lg p-4 shadow-xl bg-white">
          <NuxtLink to="/newpost">
            <v-icon name="co-pencil" class="w-10 h-10" />
          </NuxtLink>
          <NuxtLink to="/explore">
            <v-icon name="co-openstreetmap" class="w-10 h-10" />
          </NuxtLink>
          <NuxtLink to="/timeline">
            <v-icon name="co-globe-alt" class="w-10 h-10" />
          </NuxtLink>
          <NuxtLink to="/user">
            <v-icon name="la-user-circle-solid" class="w-10 h-10" />
          </NuxtLink>
        </div>
      </div>

      <!-- Slot content area -->
      <div class="flex-grow overflow-auto p-8" style="height: 100vh;">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

function goTo(path) {
  router.push(path)
}
</script>
