<template>
  <button @click="signIn(loginAttrs.provider, { callbackUrl: '/user' })"
    class="bg-gray-900 text-white font-bold py-2 px-4 rounded inline-flex items-center">
    <div class="flex gap-4 items-center">
      <v-icon :name="loginAttrs.icon" />
      <span>Login with {{ loginAttrs.service }}</span>
    </div>
  </button>
</template>

<script setup lang="ts">
import type { LoginAttrs } from '~/types/utils';

const props = defineProps<{
  loginAttrs: LoginAttrs;
}>();

const { signIn } = useAuth()
</script>