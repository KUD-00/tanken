<template>
  <div class="flex gap-4">
    <div class="flex flex-col gap-4">
      <img :src="post.avatar" class="rounded-full h-12 w-12" alt="user avatar">
      <span class="font-bold text-xl">{{ post.username }}</span>
    </div>
    <div class="max-w-sm rounded-lg overflow-hidden shadow-lg bg-white border">
      <div class="flex flex-col px-6 py-4">
        <p class="text-gray-700 text-base">
          {{ post.content }}
        </p>
        <div class="flex flex-row gap-4 mt-4">
          <span v-for="tag in post.tags" :key="tag" class="text-blue-500 text-sm">#{{ tag }}</span>
        </div>
      </div>
      <img :src="post.image" class="w-full" alt="content image">
      <div class="px-6 pt-4 pb-2 flex justify-center gap-8">
        <button class="text-gray-700">
          <v-icon name="co-thumb-up" class="w-6 h-6" />
        </button>
        <button class="text-gray-700">
          <v-icon name="co-bookmark" class="w-6 h-6" />
        </button>
        <button class="text-gray-700">
          <v-icon name="fa-walking" class="w-6 h-6" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  post: {
    type: Object,
    required: true,
  }
});
</script>