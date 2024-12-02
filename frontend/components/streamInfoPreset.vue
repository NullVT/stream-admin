<template>
  <Card class="mb-4 relative w-full">
    <template v-if="!noActions" v-slot:header>
      <div class="flex">
        <h2 class="text-white text-xl ml-2">{{ preset.name }}</h2>

        <!-- edit -->
        <button
          @click="emit('edit', preset.id)"
          class="text-white/50 hover:text-white ml-auto mr-2 flex items-center"
        >
          <PencilSquareIcon class="h-5 w-5" />
          Edit
        </button>

        <!-- apply -->
        <div v-if="loadingApply" class="mr-2 ml-2">
          <Spinner />
        </div>
        <button
          v-else
          @click="apply"
          class="text-primary/50 hover:text-primary ml-2 mr-2 flex items-center"
        >
          <ArrowUpOnSquareStackIcon class="h-5 w-5" />
          Apply
        </button>
      </div>
    </template>

    <!-- Stream Title -->
    <div class="mb-2">
      <div class="flex justify-between">
        <label
          for="stream-title"
          class="block text-md font-medium leading-6 ml-2 text-white"
        >
          Stream Title
        </label>
      </div>
      <div class="mt-2 flex">
        <input
          type="text"
          name="stream-title"
          id="chatbot-addr"
          class="block w-full rounded-md border-0 py-1.5 px-1.5 bg-base text-white placeholder:text-gray-400 focus:ring-primary sm:text-sm sm:leading-6"
          aria-describedby="chatbot-addr-optional"
          v-model="preset.title"
          :readonly="true"
        />
      </div>
    </div>

    <!-- Category -->
    <div v-if="preset.category">
      <div class="mb-2">
        <div class="flex justify-between">
          <label
            for="stream-title"
            class="block text-md font-medium leading-6 ml-2 text-white"
          >
            Category
          </label>
        </div>
        <div class="mt-2 flex items-center bg-base rounded-md">
          <img
            :src="preset.category.image_url"
            class="h-20 flex-shrink-0 object-cover rounded-l-md"
          />
          <span class="ml-3 truncate font-semibold text-white">
            {{ preset.category.name }}
          </span>
        </div>
      </div>
    </div>

    <!-- tags -->
    <div class="mb-2">
      <div class="flex justify-between mb-2">
        <label
          for="stream-title"
          class="block text-md font-medium leading-6 ml-2 text-white"
        >
          Tags
        </label>
      </div>

      <div class="flex flex-wrap gap-2">
        <div
          v-for="tag in preset.tags"
          :key="tag"
          class="inline-flex items-center gap-x-0.5 rounded-md bg-base px-2 py-1 text-xs font-medium text-white"
        >
          {{ tag }}
        </div>
      </div>
    </div>
  </Card>
</template>

<script lang="ts" setup>
import {
  ArrowUpOnSquareStackIcon,
  PencilSquareIcon,
} from "@heroicons/vue/20/solid";
import { ref } from "vue";
import { useSettingsStore } from "../stores/settings";
import { StreamInfoPreset } from "../types";
import Card from "./card.vue";
import Spinner from "./icons/spinner.vue";

const emit = defineEmits(["edit", "apply"]);
const props = defineProps<{ preset: StreamInfoPreset; noActions?: boolean }>();
const loadingApply = ref(false);

const settingsStore = useSettingsStore();

const apply = async () => {
  if (window.confirm(`Apply preset: ${props.preset.name}?`)) {
    loadingApply.value = true;

    try {
      const res = await fetch(
        `${settingsStore.adminServerAddr}/stream-info-presets/${props.preset.id}/apply`,
        {
          method: "POST",
          headers: {
            "content-type": "application/json",
          },
        }
      );
      if (res.status !== 200) {
        throw new Error(`unexpected response code: ${res.status}`);
      }
      emit("apply");
    } catch (err) {
      console.error("failed to save preset", err);
    }

    loadingApply.value = false;
  }
};
</script>
