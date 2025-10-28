<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ t("prompts.downloadFromURL") }}</h2>
    </div>

    <div class="card-content">
      <p>{{ t("prompts.downloadFromURLMessage") }}</p>
      <input
        id="focus-prompt"
        class="input input--block"
        type="text"
        @keyup.enter="submit"
        v-model.trim="url"
        :placeholder="t('prompts.downloadURL')"
      />
    </div>

    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="layoutStore.closeHovers"
        :aria-label="t('buttons.cancel')"
        :title="t('buttons.cancel')"
      >
        {{ t("buttons.cancel") }}
      </button>
      <button
        class="button button--flat"
        @click="submit"
        :aria-label="t('buttons.download')"
        :title="t('buttons.download')"
      >
        {{ t("buttons.download") }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useLayoutStore } from "@/stores/layout";

const layoutStore = useLayoutStore();
const { t } = useI18n();

const url = ref<string>("");

const submit = async (event: Event) => {
  event.preventDefault();
  if (url.value === "") return;

  console.log("Downloading from " + url.value);
  layoutStore.currentPrompt?.confirm(url.value);
  layoutStore.closeHovers();
};
</script>
