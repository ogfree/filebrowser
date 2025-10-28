import { defineStore } from "pinia";
import { ref } from "vue";

interface Download {
  id: number;
  name: string;
  progress: number;
}

let nextId = 0;

export const useDownloadStore = defineStore("download", () => {
  const downloads = ref<Download[]>([]);

  const add = (name: string): number => {
    const id = nextId++;
    downloads.value.push({ id, name, progress: 0 });
    return id;
  };

  const update = (id: number, progress: number) => {
    const download = downloads.value.find((d) => d.id === id);
    if (download) {
      download.progress = progress;
    }
  };

  const remove = (id: number) => {
    downloads.value = downloads.value.filter((d) => d.id !== id);
  };

  return {
    downloads,
    add,
    update,
    remove,
  };
});
