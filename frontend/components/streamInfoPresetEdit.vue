<template>
  <Card
    :title="preset ? 'Edit Stream Info Preset' : 'New Stream Info Preset'"
    class="mb-4 relative"
  >
    <!-- Preset Name -->
    <div class="mb-2">
      <div class="flex justify-between">
        <label
          for="stream-title"
          class="block text-md font-medium leading-6 ml-2 text-white"
        >
          Preset Name
        </label>
      </div>
      <div class="mt-2 flex">
        <input
          type="text"
          name="stream-title"
          id="stream-title"
          class="block w-full rounded-md border-0 py-1.5 px-1.5 bg-base text-white placeholder:text-gray-400 focus:ring-primary sm:text-sm sm:leading-6"
          aria-describedby="stream-title-optional"
          v-model="presetName"
        />
      </div>
    </div>

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
          id="stream-title"
          class="block w-full rounded-md border-0 py-1.5 px-1.5 bg-base text-white placeholder:text-gray-400 focus:ring-primary sm:text-sm sm:leading-6"
          aria-describedby="stream-title-optional"
          v-model="streamTitle"
        />
      </div>
    </div>

    <!-- Category -->
    <div v-if="selectedCategory">
      <div class="mb-2">
        <div class="flex justify-between">
          <label
            for="category"
            class="block text-md font-medium leading-6 ml-2 text-white"
          >
            Category
          </label>
        </div>
        <div class="mt-2 flex items-center bg-base rounded-md">
          <img
            :src="selectedCategory.image_url"
            class="h-20 flex-shrink-0 object-cover rounded-l-md"
          />
          <span class="ml-3 truncate font-semibold text-white">
            {{ selectedCategory.name }}
          </span>
          <button
            class="ml-auto mr-2 text-white hover:text-primary"
            @click="clearSelection"
          >
            <XMarkIcon class="h-6 w-6" />
          </button>
        </div>
      </div>
    </div>
    <Combobox
      v-else
      id="category"
      as="div"
      v-model="selectedCategory"
      @update:modelValue="clearQuery"
    >
      <ComboboxLabel
        class="block text-md font-medium leading-6 ml-2 text-white"
      >
        Category
      </ComboboxLabel>
      <div class="relative mt-2">
        <ComboboxInput
          class="w-full rounded-md border-0 bg-base py-1.5 pl-3 pr-12 text-white shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
          @input="debouncedQueryUpdate"
          @blur="clearQuery"
          :display-value="(category) => category?.name"
        />
        <ComboboxButton
          class="absolute inset-y-0 right-0 flex items-center rounded-r-md px-2 focus:outline-none"
        >
          <ChevronUpDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
        </ComboboxButton>

        <ComboboxOptions
          v-if="categories.length > 0"
          class="absolute z-50 mt-1 max-h-56 w-full overflow-auto rounded-md bg-zinc-700 py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm"
        >
          <ComboboxOption
            v-for="category in categories"
            :key="category.id"
            :value="category"
            as="template"
            v-slot="{ active, selected }"
          >
            <li
              :class="[
                'relative cursor-default select-none py-1 pl-3 pr-9',
                active ? 'bg-primary text-white' : 'text-white',
              ]"
            >
              <div class="flex items-center">
                <img
                  :src="category.image_url"
                  alt="Category"
                  class="h-16 flex-shrink-0 object-cover"
                />
                <span :class="['ml-3 truncate', selected && 'font-semibold']">
                  {{ category.name }}
                </span>
              </div>

              <span
                v-if="selected"
                :class="[
                  'absolute inset-y-0 right-0 flex items-center pr-4',
                  active ? 'text-white' : 'text-indigo-600',
                ]"
              >
                <CheckIcon class="h-5 w-5" aria-hidden="true" />
              </span>
            </li>
          </ComboboxOption>
        </ComboboxOptions>
      </div>
    </Combobox>

    <!-- Tags -->
    <div class="mb-2">
      <div class="flex justify-between mb-2">
        <label
          for="tags"
          class="block text-md font-medium leading-6 ml-2 text-white"
        >
          Tags
        </label>
      </div>

      <!-- New Tag Input -->
      <div class="flex items-center mb-2">
        <input
          type="text"
          v-model="newTag"
          placeholder="Add a tag"
          class="block w-full rounded-l-md border-0 py-1.5 px-2 bg-base text-white placeholder:text-gray-400 focus:ring-primary sm:text-sm sm:leading-6"
          @keydown.enter.prevent="addTag"
        />
        <button
          @click="addTag"
          class="rounded-r-md bg-primary text-white font-semibold text-sm p-2"
        >
          <PlusIcon class="h-5 w-5" />
        </button>
      </div>

      <div class="flex flex-wrap gap-2">
        <div
          v-for="(tag, tagIndex) in streamTags"
          :key="tag"
          class="inline-flex items-center gap-x-0.5 rounded-md bg-base px-2 py-1 text-xs font-medium text-white"
        >
          {{ tag }}
          <button
            type="button"
            class="group relative ml-1 -mr-1 h-4 w-4 rounded-sm hover:bg-gray-100/20 flex items-center justify-center text-white/30 hover:text-red-600"
            @click="removeTag(tagIndex)"
          >
            <XMarkIcon class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Save Button -->
    <div class="flex mt-4">
      <button
        class="ml-auto rounded-md bg-primary p-2 text-white font-semibold text-sm uppercase"
        @click="savePreset"
      >
        Save
      </button>
    </div>
  </Card>
</template>

<script lang="ts" setup>
import {
  Combobox,
  ComboboxButton,
  ComboboxInput,
  ComboboxLabel,
  ComboboxOption,
  ComboboxOptions,
} from "@headlessui/vue";
import {
  CheckIcon,
  ChevronUpDownIcon,
  PlusIcon,
  XMarkIcon,
} from "@heroicons/vue/20/solid";
import { ref } from "vue";
import Card from "../components/card.vue";
import { useSettingsStore } from "../stores/settings";
import { Category, StreamInfoPreset } from "../types";

// define component
const { preset } = defineProps<{ preset: StreamInfoPreset | null }>();
const emit = defineEmits(["save", "close"]);

// Reactive data
const settingsStore = useSettingsStore();
const query = ref("");
const streamTitle = ref(preset?.title || "");
const presetName = ref(preset?.name || "");
const streamTags = ref<string[]>(preset?.tags || []);
const newTag = ref("");
const selectedCategory = ref(
  preset?.category
    ? { ...preset.category, box_art_url: preset.category.image_url }
    : null
);
const categories = ref<Category[]>([]);
let debounceTimeout: ReturnType<typeof setTimeout> | null = null;

// Add new tag
const addTag = () => {
  const tag = newTag.value.trim();
  if (tag && !streamTags.value.includes(tag)) {
    streamTags.value.push(tag);
  }
  newTag.value = "";
};

// Remove tag
const removeTag = (index: number) => {
  streamTags.value.splice(index, 1);
};

// Fetch categories after input change
type TwitchCategory = {
  id: string;
  name: string;
  box_art_url: string;
};
const fetchCategories = async () => {
  if (query.value.trim() !== "") {
    try {
      const urlQuery = encodeURIComponent(query.value);
      const res = await fetch(
        `${settingsStore.adminServerAddr}/twitch/categories?query=${urlQuery}`
      );
      if (res.status === 200) {
        const data = await res.json();
        categories.value = data.map((cat: TwitchCategory) => ({
          ...cat,
          image_url: cat.box_art_url,
        }));
      } else {
        console.error(`Unexpected response code (${res.status})`);
        categories.value = [];
      }
    } catch (err) {
      console.error(err);
      alert("Failed to list categories");
    }
  } else {
    categories.value = [];
  }
};

// Debounced input for category search
const debouncedQueryUpdate = (event: Event) => {
  query.value = (event.target as HTMLInputElement).value;

  if (debounceTimeout) {
    clearTimeout(debounceTimeout);
  }

  debounceTimeout = setTimeout(() => {
    fetchCategories();
  }, 500); // Wait 500ms after user stops typing
};

// Clear selected category
const clearSelection = () => {
  query.value = "";
  categories.value = [];
  selectedCategory.value = null;
};

const clearQuery = () => {
  query.value = "";
};

// Save the preset and emit the 'save' event
const savePreset = () => {
  const updatedPreset = {
    id: preset?.id || "", // ID is auto-generated for new presets
    name: presetName.value,
    title: streamTitle.value,
    tags: streamTags.value,
    category: selectedCategory.value,
  };

  // Emit 'save' event to the parent component
  emit("save", updatedPreset);
};
</script>
