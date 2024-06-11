<script setup lang="ts">
const route = useRoute()
const layout = "home";

console.log(route.params.id)

</script>

<template>
  <NuxtLayout :name="layout">
    <p>{{ $route.params.postId }}</p>
  </NuxtLayout>
</template>
