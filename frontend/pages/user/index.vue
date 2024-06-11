<script setup>
const layout = "home";

const info = {
  avatar: '/img/avatar.jpeg',
  username: 'Kud',
};

const { status, data, lastRefreshedAt, getCsrfToken, getProviders, getSession, signIn, signOut, } = useAuth()

</script>

<template>
  <NuxtLayout :name="layout">
    <div class="flex flex-col items-center justify-center bg-white rounded-lg gap-8 p-8">
      <div class="flex flex-col items-center justify-center">
        <img :src="info.avatar" class="rounded-full h-36 w-36" alt="user avatar">
        <span class="font-bold text-xl mt-4">{{ info.username }}</span>
      </div>
      <div class="flex flex-row items-center gap-8">
        <p>Posts</p>
        <p>Bookmarks</p>
        <p>Plans</p>
        <p>Achivements</p>
      </div>
    </div>
  </NuxtLayout>
</template>