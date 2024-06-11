import { ref, provide, onUnmounted, type Ref } from 'vue';
//@ts-ignore
import mapboxgl from 'mapbox-gl/dist/mapbox-gl';
import MapboxLanguage from '@mapbox/mapbox-gl-language';

interface MapOptions {
  container: string;
  style: string;
  center: [number, number];
  zoom: number;
  languageOption: mapboxgl.Control;
}

export default function useMapbox() {
  const map: Ref<mapboxgl.Map | null> = ref(null);

  const initializeMap = (options: MapOptions): void => {
    const runtimeconfig = useRuntimeConfig();
    mapboxgl.accessToken = runtimeconfig.public.mapboxAccessToken;

    map.value = new mapboxgl.Map({
      container: options.container,
      style: options.style,
      center: options.center,
      zoom: options.zoom
    });

    const language = new MapboxLanguage(options.languageOption);
    map.value.addControl(language);

    map.value.on('load', () => {
      console.log('Map has been loaded!');
    });
  };

  const removeMap = (): void => {
    if (map.value) {
      map.value.remove();
      map.value = null;
    }
  };

  onUnmounted(removeMap);

  return { map, initializeMap, removeMap };
}
