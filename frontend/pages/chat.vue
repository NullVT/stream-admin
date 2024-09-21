<template>
  <h1 class="text-white">Chat messages</h1>
  <ul role="list" class="gap-y-2">
    <li
      v-for="msg in messages"
      :key="msg.id"
      class="flex justify-between gap-x-5 p-2 mb-0.5 bg-background border-l-4"
      :class="`border-${userColorClass(msg) ?? 'base'}`"
    >
      <div class="flex min-w-0 gap-x-4">
        <!-- <img class="h-12 w-12 flex-none rounded-full bg-gray-50" :src="person.imageUrl" alt="" /> -->
        <div class="min-w-0 flex-auto">
          <p
            class="mt-1 flex text-xs leading-5 font-bold"
            :class="`text-${userColorClass(msg) ?? 'gray-500'}`"
          >
            {{ msg.sender.name }}
          </p>
          <p class="text-sm font-semibold leading-6 text-white">
            {{ msg.body }}
          </p>
        </div>
      </div>
      <div class="flex shrink-0 items-center gap-x-6">
        <!-- <div class="hidden sm:flex sm:flex-col sm:items-end">
          <p class="text-sm leading-6 text-gray-900">{{ person.role }}</p>
          <p v-if="person.lastSeen" class="mt-1 text-xs leading-5 text-gray-500">
            Last seen <time :datetime="person.lastSeenDateTime">{{ person.lastSeen }}</time>
          </p>
          <div v-else class="mt-1 flex items-center gap-x-1.5">
            <div class="flex-none rounded-full bg-emerald-500/20 p-1">
              <div class="h-1.5 w-1.5 rounded-full bg-emerald-500" />
            </div>
            <p class="text-xs leading-5 text-gray-500">Online</p>
          </div>
        </div> -->
        <Menu as="div" class="relative flex-none">
          <MenuButton
            class="-m-2.5 block p-2.5 text-gray-500 hover:text-gray-900"
          >
            <span class="sr-only">Open options</span>
            <EllipsisVerticalIcon class="h-5 w-5" aria-hidden="true" />
          </MenuButton>
          <transition
            enter-active-class="transition ease-out duration-100"
            enter-from-class="transform opacity-0 scale-95"
            enter-to-class="transform opacity-100 scale-100"
            leave-active-class="transition ease-in duration-75"
            leave-from-class="transform opacity-100 scale-100"
            leave-to-class="transform opacity-0 scale-95"
          >
            <MenuItems
              class="absolute right-0 z-10 mt-2 w-32 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 focus:outline-none"
            >
              <MenuItem v-slot="{ active }">
                <a
                  href="#"
                  :class="[
                    active ? 'bg-gray-50' : '',
                    'block px-3 py-1 text-sm leading-6 text-gray-900',
                  ]"
                >
                  View profile
                  <span class="sr-only">, {{ msg.sender.name }} </span></a
                >
              </MenuItem>
            </MenuItems>
          </transition>
        </Menu>
      </div>
    </li>
  </ul>
</template>

<script lang="ts" setup>
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/vue";
import { EllipsisVerticalIcon } from "@heroicons/vue/20/solid";
import { Ref, ref } from "vue";
import { urlToWss } from "../helpers";
import { useSettingsStore } from "../stores/settings";
import { AdminWSMessage } from "../types";

const settingsStore = useSettingsStore();
const messages: Ref<AdminWSMessage[]> = ref([]);

const userColorClass = (msg: AdminWSMessage) => {
  if (msg.sender.broadcaster) return "red-500";
  if (msg.sender.moderator) return "green-500";
  if (msg.sender.twitch_vip) return "pink-500";
  if (msg.sender.youtube_member) return "indigo-500";
};

const ws = new WebSocket(`${urlToWss(settingsStore.adminServerAddr)}/messages`);
ws.addEventListener("message", async (event: WebSocketEventMap["message"]) => {
  console.log("message", event);
  try {
    const message: AdminWSMessage = JSON.parse(event.data);
    const newMessages = [...messages.value, message];
    messages.value = newMessages.slice(-1000);
  } catch (err) {
    console.error("failed to handle WS message", err);
  }
});

ws.addEventListener("error", async (event: WebSocketEventMap["error"]) => {
  console.log("error", event);
});
ws.addEventListener("open", async (event: WebSocketEventMap["open"]) => {
  console.log("open", event);
});
ws.addEventListener("close", async (event: WebSocketEventMap["error"]) => {
  console.log("close", event);
});
</script>
