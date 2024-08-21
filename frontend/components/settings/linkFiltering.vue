<template>
  <div class="space-y-1">
    <!-- Label -->
    <label
      for="link-filtering"
      class="block text-lg leading-6 font-medium text-white text-center"
    >
      Link Filtering
    </label>

    <!-- loading spinner -->
    <button
      v-if="loading"
      class="text-center capitalize w-full py-2 focus:outline-none text-sm text-white bg-base rounded-lg"
      disabled
    >
      <span
        class="w-4 h-4 border-2 border-white border-t-transparent border-t-2 rounded-full animate-spin inline-block"
      ></span>
    </button>

    <!-- options -->
    <RadioGroup
      v-else
      id="link-filtering"
      v-model="mode"
      class="flex rounded-lg overflow-hidden w-full"
      :disabled="loading"
    >
      <RadioGroupOption
        v-for="filterMode in FilterModes"
        :key="filterMode"
        :value="filterMode"
        v-slot="{ checked }"
        :class="[
          'flex-1',
          filterMode === FilterModes.Enabled && 'rounded-l-lg',
          filterMode === FilterModes.Disabled && 'rounded-r-lg',
        ]"
      >
        <button
          :class="[
            'w-full py-2 focus:outline-none text-sm text-center capitalize text-white',
            checked ? 'bg-primary' : 'bg-base',
          ]"
        >
          {{ filterMode }}
        </button>
      </RadioGroupOption>
    </RadioGroup>
  </div>
</template>

<script lang="ts" setup>
import { RadioGroup, RadioGroupOption } from "@headlessui/vue";
import { ref, watch } from "vue";
import { useSettingsStore } from "../../stores/settings";

enum FilterModes {
  Enabled = "enabled",
  Filtered = "filtered",
  Disabled = "disabled",
}

const settingsStore = useSettingsStore();
const mode = ref(FilterModes.Filtered);
const loading = ref(false);

const setTwitch = async (blockLinks: boolean) => {
  // const res = await fetch(
  //   `${settingsStore.adminServerAddr}/twitch/link-filtering`,
  //   {
  //     method: "POST",
  //     body: JSON.stringify({ enable: blockLinks }),
  //     headers: {
  //       "content-type": "application/json",
  //     },
  //   }
  // );
  // console.log(res);
};

const setChatbot = async (filterMode: FilterModes) => {
  // chatbot api req
};

watch(mode, async (val) => {
  loading.value = true;
  try {
    await Promise.all([setTwitch(val !== "disabled"), setChatbot(val)]);
  } catch (error) {
    console.error("Failed to set filtering mode", error);
  } finally {
    loading.value = false;
  }
});
</script>
