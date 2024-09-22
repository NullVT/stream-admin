<template>
  <ul
    role="list"
    class="gap-y-2 h-full overflow-y-auto overflow-x-hidden"
    ref="messageList"
    @scroll="handleScroll"
  >
    <li
      v-for="msg in msgStore.messages"
      :key="msg.id"
      class="flex justify-between snap-end gap-x-5 p-2 mb-0.5 bg-background border-l-4"
      :class="{
        'border-red-800': msg.sender.broadcaster,
        'border-green-500': msg.sender.moderator,
        'border-pink-500': msg.sender.twitch_vip,
        'border-indigo-500': msg.sender.youtube_member,
        'border-base':
          !msg.sender.broadcaster &&
          !msg.sender.moderator &&
          !msg.sender.twitch_vip &&
          !msg.sender.youtube_member,
      }"
    >
      <div class="flex min-w-0 gap-x-4">
        <div class="min-w-0 flex-auto">
          <p
            class="mt-1 flex text-xs leading-5 font-bold"
            :class="{
              'text-red-800': msg.sender.broadcaster,
              'text-green-500': msg.sender.moderator,
              'text-pink-500': msg.sender.twitch_vip,
              'text-indigo-500': msg.sender.youtube_member,
              'text-sky-300':
                !msg.sender.broadcaster &&
                !msg.sender.moderator &&
                !msg.sender.twitch_vip &&
                !msg.sender.youtube_member,
            }"
          >
            {{ msg.sender.name }}
          </p>
          <p class="text-sm font-semibold leading-6 text-white">
            {{ msg.body }}
          </p>
        </div>
      </div>
      <div class="flex shrink-0 items-center gap-x-6">
        <Menu as="div" class="relative flex-none">
          <MenuButton
            class="-m-2.5 block p-2.5 text-gray-500 hover:text-primary"
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

  <!-- scroll paused -->
  <button
    v-if="!autoScrollEnabled"
    @click="scrollToBottom"
    class="fixed bottom-4 left-1/2 transform -translate-x-1/2 inline-flex items-center rounded p-2 bg-zinc-600 drop-shadow-md text-white font-semibold text-sm opacity-50 hover:opacity-100 uppercase"
  >
    <PauseIcon class="size-5 shrink-0" /> Scroll paused
  </button>
</template>

<style scoped>
/* For WebKit browsers (Chrome, Safari) */
ul::-webkit-scrollbar {
  width: 6px; /* Width of the scrollbar */
}

ul::-webkit-scrollbar-thumb {
  background-color: rgba(80, 80, 80, 1); /* Color of the scrollbar thumb */
  border-radius: 10px; /* Rounded corners on the scrollbar thumb */
}

ul::-webkit-scrollbar-track {
  background: transparent; /* Background of the scrollbar track */
}

/* For Firefox */
ul {
  scrollbar-width: thin; /* Makes the scrollbar thin */
  scrollbar-color: rgba(80, 80, 80, 1) transparent; /* Thumb and track color */
}
</style>

<script lang="ts" setup>
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/vue";
import { PauseIcon } from "@heroicons/vue/16/solid";
import { EllipsisVerticalIcon } from "@heroicons/vue/20/solid";
import { nextTick, onMounted, ref, watch } from "vue";
import { useMessagesStore } from "../stores/messages";

const msgStore = useMessagesStore();
const autoScrollEnabled = ref(true);
const messageList = ref<HTMLElement | null>(null); // Ref for the <ul>

// scroll to the bottom of the chat list (using the <ul> ref)
const scrollToBottom = async () => {
  await nextTick(); // ensure the DOM is updated before scrolling
  if (messageList.value) {
    messageList.value.scrollTop = messageList.value.scrollHeight;
    autoScrollEnabled.value = true;
  }
};

// handle the scroll event on the <ul> element
const handleScroll = () => {
  if (!messageList.value) return;

  // distance from the bottom to trigger auto-scroll disabling
  const threshold = 100;
  const atBottom =
    messageList.value.scrollHeight -
      messageList.value.scrollTop -
      messageList.value.clientHeight <
    threshold;

  // disable auto-scrolling if user scrolls up
  autoScrollEnabled.value = atBottom;
};

onMounted(() => {
  scrollToBottom();
});

// watch for new messages and scroll to the bottom if auto-scrolling is enabled
watch(
  () => msgStore.messages,
  () => {
    if (autoScrollEnabled.value) {
      scrollToBottom();
    }
  }
);
</script>
