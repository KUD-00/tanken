<script setup>
import { ref } from 'vue';
import Post from '~/components/NewPost.vue';
import LoginButton from '~/components/LoginButton.vue';

const post = ref(
  {
    username: 'Kud',
    avatar: '/img/avatar.jpeg',
    content: 'Happy Happy Train to Go!',
    image: '/img/train.webp',
    time: 'yesterday',
    tags: ['train', '鉄道'],
    imageMeta: {
      longtitute: 135.681747,
      latitude: 34.891586,
      shotTime: '2024-01-01T00:00:00Z',
    },
  }
);

const loginServices = ref([
  { service: 'Google', icon: 'co-google', provider: 'google'},
  { service: 'GitHub', icon: 'co-github', provider: 'github'},
]);
</script>

<template>
  <div class="min-h-screen flex">
    <div class="w-3/5 bg-slate-200 border flex justify-center items-center">
      <div class="flex flex-row max-w-xl bg-slate-200">
        <Post :post="post" class="flex-initial" />
      </div>
    </div>
    <div class="w-2/5 border bg-slate-100 flex flex-col justify-center items-center">
      <p class="text-2xl">
        Go!TanKen! @2024
      </p>
      <div v-for="login in loginServices" :key="login.service" class="mt-4">
        <LoginButton :loginAttrs="login" />
      </div>
    </div>
  </div>
</template>