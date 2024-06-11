<script setup lang="ts">
import { ref, onMounted } from 'vue';
import useMapbox from '@/composables/useMapbox';
import useFileUpload from '@/composables/useFileUpload';
import { PiexifExif  } from "piexifjs";
import 'filepond/dist/filepond.min.css';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.min.css';
import FilePondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
import vueFilePond from "vue-filepond";

// map
const myFiles = ref<File[]>([]);
const exifData = ref<PiexifExif>({})

const { map, initializeMap } = useMapbox();
const { handleFileProcess } = useFileUpload({ myFiles, exifData, map });

//@ts-ignore
const FilePondComponent = defineAsyncComponent(async () => {
  const module = await import('vue-filepond');
  const vueFilePond = module.default || module;

  const FilePondPluginFileValidateType = await import('filepond-plugin-file-validate-type');
  const FilePondPluginImagePreview = await import('filepond-plugin-image-preview');

  return vueFilePond(
    FilePondPluginFileValidateType.default || FilePondPluginFileValidateType,
    FilePondPluginImagePreview.default || FilePondPluginImagePreview
  );
});

const uploadNewPost = async () => {
  const location = {
    Latitude: exifData.value?.latitude,
    Longitude: exifData.value?.longitude
  };

  const content = document.querySelector('textarea')?.value;
  const tags = ['tag1', 'tag2'];
  const pictureChunk = await myFiles.value[0].arrayBuffer();
  const uint8PictureChunk = new Uint8Array(pictureChunk);

  try {
    const response = await fetch('/api/rpc/post', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        method: 'uploadNewPost',
        userId: 'your-user-id',
        pictureChunk: Array.from(uint8PictureChunk),
        location,
        content,
        tags
      })
    });

    console.log(response)

    const responseData = await response.json();
    console.log('Server response:', responseData);
  } catch (error) {
    console.error('Error uploading post:', error);
  }
};

onMounted(() => {
  initializeMap({
    container: 'map',
    style: 'mapbox://styles/mapbox/streets-v11',
    center: [120.1575, 30.2874],
    zoom: 14,
    languageOption: { defaultLanguage: 'ja' }
  });
});
</script>

<template>
  <NuxtLayout :name="'home'">
    <div class="flex flex-row gap-8">
      <div class="flex flex-col gap-8">
        <div class="flex flex-col justify-center items-center bg-white gap-8 rounded-lg p-8">
          <div class="text-2xl">
            Share your discovery!
          </div>
          <client-only>
            <FilePondComponent class="w-64" name="test" ref="pond"
              label-idle="Drop photos here or <span class='filepond--label-action'>Browse</span>" allow-multiple="true"
              accepted-file-types="image/jpeg, image/png" :files="myFiles" @addfile="handleFileProcess" />
          </client-only>
          <textarea class="border p-4 w-96 h-36 bg-slate-50 rounded-lg"
            placeholder="Write something about your discovery"></textarea>
          <button class="ml-auto bg-deep-green hover:bg-dark-green text-white font-bold py-2 px-4 rounded"
            @click="uploadNewPost">
            Upload
          </button>
        </div>
      </div>

      <div class="flex flex-col gap-8">
        <div class="flex flex-col justify-center items-center bg-white gap-8 rounded-lg p-8">
          <div class="text-2xl">
            Photo Metadata
          </div>
          <div id="map" class="min-h-[400px] min-w-[400px]"></div>
          <div v-if="Object.keys(exifData).length">
            <div v-for="(value, key) in exifData" :key="key">
              <strong>{{ key }}:</strong> {{ value }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<style scoped>
</style>
