<template>
  <!-- Display Current Stream Info -->
  <div
    v-if="loadingCurrent"
    class="rounded-l bg-background text-white h-72 flex items-center"
  >
    <div class="ml-auto mr-auto">
      <ArrowPathIcon class="h-10 w-10 animate-spin ml-auto mr-auto" />
      <h2 class="mt-2 animate-pulse">loading...</h2>
    </div>
  </div>
  <StreamInfoPreset v-else-if="currentInfo" :preset="currentInfo" no-actions />

  <!-- Divider with "Add Preset" button -->
  <div class="relative mt-4">
    <div class="absolute inset-0 flex items-center" aria-hidden="true">
      <div class="w-full border-t border-background"></div>
    </div>
    <div class="relative flex justify-center">
      <span class="bg-base px-3 text-base font-semibold leading-6 text-white">
        Presets
      </span>
      <button
        @click="openModal()"
        class="text-white hover:text-primary bg-base"
      >
        <PlusIcon class="h-5 w-5" />
      </button>
    </div>
  </div>

  <!-- Presets with Edit buttons -->
  <div class="mt-4">
    <div
      v-for="preset in presets"
      :key="preset.id"
      class="flex items-center justify-between mb-2"
    >
      <StreamInfoPreset
        :preset="preset"
        @edit="openModal(preset)"
        @apply="getStreamInfo"
      />
    </div>
  </div>

  <!-- Modal Component -->
  <Modal v-if="showModal" @close="closeModal">
    <StreamInfoPresetEdit
      :preset="editingPreset"
      @close="closeModal"
      @save="handleSave"
    />
  </Modal>
</template>

<script lang="ts" setup>
import { ArrowPathIcon, PlusIcon } from "@heroicons/vue/20/solid";
import { onMounted, ref } from "vue";
import Modal from "../components/modal.vue";
import StreamInfoPreset from "../components/streamInfoPreset.vue";
import StreamInfoPresetEdit from "../components/streamInfoPresetEdit.vue";
import { useSettingsStore } from "../stores/settings";
import { StreamInfoPreset as StreamPreset } from "../types";

// Store and Refs
const settingsStore = useSettingsStore();
const presets = ref<StreamPreset[]>([]);
const currentInfo = ref<StreamPreset>();
const showModal = ref(false);
const editingPreset = ref<StreamPreset | null>(null);
const loadingCurrent = ref(false);

// Open Modal (for new or edit preset)
const openModal = (preset: StreamPreset | null = null) => {
  editingPreset.value = preset; // Pass preset for editing or null for new
  showModal.value = true;
};

// Close Modal
const closeModal = () => {
  showModal.value = false;
  editingPreset.value = null;
};

// Handle Save
const handleSave = async (preset: StreamPreset) => {
  try {
    const suffix = !preset.id ? "" : `/${preset.id}`;
    const res = await fetch(
      `${settingsStore.adminServerAddr}/stream-info-presets${suffix}`,
      {
        method: !preset.id ? "POST" : "PUT",
        body: JSON.stringify(preset),
        headers: {
          "content-type": "application/json",
        },
      }
    );
    if (res.status !== 200) {
      throw new Error(`unexpected response code: ${res.status}`);
    }
    presets.value = await res.json();
    closeModal();
  } catch (err) {
    console.error("failed to save preset", err);
  }
};

// Fetch current stream info
const getStreamInfo = async () => {
  loadingCurrent.value = true;
  try {
    const res = await fetch(`${settingsStore.adminServerAddr}/stream-info`);
    if (res.status !== 200) {
      throw new Error(`Unexpected stream-info response (${res.status})`);
    }
    const data = await res.json();

    currentInfo.value = {
      id: "__current",
      name: "Current Stream Info",
      title: data.title,
      tags: data.tags,
      category: {
        id: data.game_id,
        name: data.game_name,
        image_url: `https://static-cdn.jtvnw.net/ttv-boxart/${data.game_id}_IGDB-52x72.jpg`,
      },
    };
  } catch (err) {
    console.error(err);
    alert("Failed to list categories");
  }
  loadingCurrent.value = false;
};
onMounted(getStreamInfo);

// Fetch presets
const getPresets = async () => {
  try {
    const res = await fetch(
      `${settingsStore.adminServerAddr}/stream-info-presets`
    );
    if (res.status !== 200) {
      throw new Error(
        `Unexpected stream-info-presets response (${res.status})`
      );
    }
    presets.value = await res.json();
  } catch (err) {
    console.error(err);
    alert("Failed to list streamInfoPresets");
  }
};
onMounted(getPresets);
</script>
