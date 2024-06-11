<script setup>
import mapboxgl from 'mapbox-gl/dist/mapbox-gl';
import { onMounted, onUnmounted, nextTick } from 'vue';
import MapboxLanguage from '@mapbox/mapbox-gl-language';
const layout = "home";

const runtimeConfig = useRuntimeConfig()
mapboxgl.accessToken = runtimeConfig.public.mapboxAccessToken

onMounted(() => {
  nextTick().then(() => {
    map = new mapboxgl.Map({
      container: 'map',
      style: 'mapbox://styles/mapbox/streets-v12',
      center: [135.681747, 34.891586],
      zoom: 14
    });
    const language = new MapboxLanguage({ defaultLanguage: 'ja' });
    map.addControl(language);
  });
});

onUnmounted(() => {
  if (map) {
    map.remove();
    map = null;
  }
});

const spots = [{
  name: "水無瀬神宮",
  location: [135.67265605172557, 34.885637162951895]
},
{
  name: "聴竹居",
  location: [34.893723035849845, 135.67894733334876]
},
{
  name: "山崎ウイスキー館",
  location: [34.89256820298335, 135.67523662892845]
}
];

</script>

<template>
  <NuxtLayout :name="layout">
    <div class="flex flex-row gap-8">
      <div class="flex flex-col justify-center items-center bg-white p-8">
        <!--       <div class="px-6 pt-4 pb-6 flex justify-center gap-16">
        <button class="text-gray-700">
          <v-icon name="md-backpack-outlined" class="w-8 h-8" />
        </button>
        <button class="text-gray-700">
          <v-icon name="io-restaurant-outline" class="w-8 h-8" />
        </button>
        <button class="text-gray-700">
          <v-icon name="ri-shopping-cart-line" class="w-8 h-8" />
        </button>
      </div> -->
        <div id="map" class="h-2/3 w-2/3 min-h-[600px] min-w-[600px]"></div>
      </div>
      <div class="flex-col items-center bg-white">
        <p>hello</p>
      </div>
    </div>
  </NuxtLayout>
</template>
