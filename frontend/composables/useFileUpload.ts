import { ref, inject, type Ref } from 'vue';
//@ts-ignore
import mapboxgl, { Map } from "mapbox-gl/dist/mapbox-gl";
//@ts-ignore
import { piexif, PiexifExif, load, dump, insert, TagValues, helper } from "piexifjs";
//@ts-ignore
import { Marker } from 'mapbox-gl';

interface FileUploadOptions {
  myFiles: Ref<File[]>;
  exifData: Ref<PiexifExif>;
  map: Ref<mapboxgl.Map | null>;
}

export default function useFileUpload({ myFiles, exifData, map }: FileUploadOptions) {
  const loading: Ref<boolean> = ref(false);

  const fileToDataUrl = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = event => resolve((event.target as FileReader).result as string);
      reader.onerror = error => {
        reader.abort();
        reject(error);
      };
      reader.readAsDataURL(file);
    });
  };

  const extractAndFormatExif = async (dataUrl: string): Promise<void> => {
    try {
      const exifObj = piexif.load(dataUrl);
      exifData.value = {
        dateTime: exifObj['0th'][piexif.ImageIFD.DateTime] as string || 'N/A',
        model: exifObj['0th'][piexif.ImageIFD.Model] as string || 'N/A',
      };

      if (map.value && exifObj['GPS'][piexif.GPSIFD.GPSLatitude]) {
        const latitude = piexif.GPSHelper.dmsRationalToDeg(exifObj['GPS'][piexif.GPSIFD.GPSLatitude]);
        const longitude = piexif.GPSHelper.dmsRationalToDeg(exifObj['GPS'][piexif.GPSIFD.GPSLongitude]);
        exifData.value.latitude = latitude;
        exifData.value.longitude = longitude;

        map.value.setCenter([longitude, latitude]);

        new Marker()
          .setLngLat([longitude, latitude])
          .addTo(map.value);
      }
    } catch (error) {
      console.error("Error processing EXIF data:", error);
    }
  };

  // @ts-ignore
  const handleFileProcess = async (error: any, fileItem: FilePond.File): Promise<void> => {
    loading.value = true;
    try {
      const dataUrl = await fileToDataUrl(fileItem.file);
      await extractAndFormatExif(dataUrl);
      myFiles.value.push(fileItem.file);
    } catch (error) {
      console.error("Error loading file:", error);
    } finally {
      loading.value = false;
    }
  };

  return {
    handleFileProcess,
    loading
  };
}
