<template>
  <div v-if="loading" class="text-gray-400 text-center italic">
    loading emotes whitelist...
  </div>
  <div v-else class="flex gap-2">
    <!-- Add Channel Button (matching size with badges) -->
    <button
      class="inline-flex items-center gap-x-0.5 rounded-md bg-base/50 px-2 py-1 text-xs font-medium text-gray-700 hover:bg-base hover:text-white"
      @click="showModal = true"
    >
      <PlusIcon class="h-4 w-4" />
    </button>

    <!-- List of channels (badges) -->
    <span
      v-for="(channelName, channelId) in whitelist"
      :key="channelId"
      class="inline-flex items-center gap-x-0.5 rounded-md bg-base px-2 py-1 text-xs font-medium text-white"
    >
      {{ channelName }}
      <button
        type="button"
        class="group relative ml-1 -mr-1 h-4 w-4 rounded-sm hover:bg-gray-100/20 flex items-center justify-center text-white/30 hover:text-red-600"
        @click="removeChannel(channelId)"
      >
        <XMarkIcon class="h-3.5 w-3.5" />
      </button>
    </span>

    <!-- Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-70"
    >
      <div class="bg-background rounded-md p-4 w-96">
        <h3 class="text-lg font-semibold mb-2 text-white">Add Channel</h3>
        <input
          v-model="newChannel"
          type="text"
          placeholder="Enter channel name"
          class="w-full px-2 py-1 rounded-md mb-4 bg-base text-white"
        />
        <div class="flex justify-end">
          <button
            @click="addChannel"
            class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary/80"
          >
            Submit
          </button>
          <button
            @click="showModal = false"
            class="ml-2 bg-base text-gray-500 px-4 py-2 rounded-md hover:bg-gray-800"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { PlusIcon, XMarkIcon } from "@heroicons/vue/16/solid";
import { onMounted, ref } from "vue";
import { useSettingsStore } from "../../stores/settings";

const settingsStore = useSettingsStore();

const whitelist = ref<Record<string, string>>({});
const loading = ref(true);
const showModal = ref(false);
const newChannel = ref("");

const loadWhitelist = async () => {
  loading.value = true;

  try {
    const res = await fetch(
      `${settingsStore.adminServerAddr}/emotes/whitelist`
    );
    whitelist.value = await res.json();
  } catch (err) {
    console.error("failed to fetch emotes whitelist", err);
  }

  loading.value = false;
};
onMounted(loadWhitelist);

const addChannel = async () => {
  loading.value = true;

  if (newChannel.value.trim() !== "") {
    try {
      const res = await fetch(
        `${settingsStore.adminServerAddr}/emotes/whitelist`,
        {
          method: "POST",
          headers: {
            "content-type": "application/json",
          },
          body: JSON.stringify(newChannel.value),
        }
      );

      if (res.status !== 200) {
        throw new Error(`unexpected response code: ${res.status}`);
      }

      whitelist.value = await res.json();
      newChannel.value = "";
      showModal.value = false;
    } catch (err) {
      console.error("failed to update emotes whitelist", err);
      alert("update whitelist failed");
    }
  }

  loading.value = false;
};

const removeChannel = async (channelId: string) => {
  loading.value = true;

  try {
    const res = await fetch(
      `${settingsStore.adminServerAddr}/emotes/whitelist`,
      {
        method: "DELETE",
        headers: {
          "content-type": "application/json",
        },
        body: JSON.stringify(channelId),
      }
    );

    if (res.status !== 200) {
      throw new Error(`unexpected response code: ${res.status}`);
    }

    whitelist.value = await res.json();
  } catch (err) {
    console.error("failed to delete emotes whitelist", err);
    alert("delete whitelist failed");
  }

  loading.value = false;
};
</script>
