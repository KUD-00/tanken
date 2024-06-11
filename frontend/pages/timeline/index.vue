<script setup>
import { ref } from 'vue';
import Post from '../../components/NewPost.vue';

const layout = "home";

const posts = ref([
  {
    postId: "e1ab7i",
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
  }, 
  {
    postId: "a1ab7i",
    username: 'Seki',
    avatar: '/img/avatar.jpeg',
    content: 'かわいい熊さん',
    image: '/img/bear.webp',
    time: '04-10',
    tags: ['animal'],
    imageMeta: {
      longtitute: 135.681747,
      latitude: 34.891586,
      shotTime: '2024-01-01T00:00:00Z',
    },
  },
  {
    postId: "s1ab7i",
    username: 'Alan',
    avatar: '/img/avatar.jpeg',
    content: '为高中生献上祈祷吧！',
    image: '/img/ema.webp',
    time: '23-12-25',
    tags: ['japanese culture'],
    imageMeta: {
      longtitute: 135.681747,
      latitude: 34.891586,
      shotTime: '2024-01-01T00:00:00Z',
    },
  },
]);
</script>

<template>
  <NuxtLayout :name="layout">
    <div class="flex flex-col justify-center items-center bg-white gap-8 rounded-lg p-8">
      <div class="flex items-center max-w-md mx-auto bg-slate-100 rounded-full shadow-lg overflow-hidden px-5 py-2 h-12">
        <input class="flex-grow outline-none text-md bg-slate-100 leading-tight" type="text" placeholder="搜索...">
        <v-icon name="bi-search"></v-icon>
      </div>

      <Post v-for="post in posts" :key="post.content" :post="post" class="flex-initial"></Post>
    </div>
  </NuxtLayout>
</template>