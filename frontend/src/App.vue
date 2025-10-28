<template>
  <div>
    <router-view></router-view>
    <div class="downloads-container" v-if="downloadStore.downloads.length > 0">
      <div
        v-for="download in downloadStore.downloads"
        :key="download.id"
        class="download-item"
      >
        <ProgressBar
          :progress="download.progress"
          :label="`${download.name} (${download.progress}%)`"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useI18n } from "vue-i18n";
import { setHtmlLocale } from "./i18n";
import { getMediaPreference, getTheme, setTheme } from "./utils/theme";
import { useDownloadStore } from "@/stores/download";
import ProgressBar from "@/components/ProgressBar.vue";

const { locale } = useI18n();
const downloadStore = useDownloadStore();

const userTheme = ref<UserTheme>(getTheme() || getMediaPreference());

onMounted(() => {
  setTheme(userTheme.value);
  setHtmlLocale(locale.value);
  // this might be null during HMR
  const loading = document.getElementById("loading");
  loading?.classList.add("done");

  setTimeout(function () {
    loading?.parentNode?.removeChild(loading);
  }, 200);
});

// handles ltr/rtl changes
watch(locale, (newValue) => {
  newValue && setHtmlLocale(newValue);
});
</script>

<style scoped>
.downloads-container {
  position: fixed;
  bottom: 1em;
  right: 1em;
  width: 300px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.5em;
}

.download-item {
  background-color: var(--card-background-color);
  padding: 1em;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
