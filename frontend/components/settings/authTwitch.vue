<template>
  <button
    v-if="twitchLoading"
    class="w-full inline-flex justify-between items-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm bg-twitch hover:bg-twitch-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 h-10"
    :disabled="!!twitchStatus"
    @click.once="login"
  >
    <img :src="twitchLogo" class="size-6" />
    <span class="flex-grow flex justify-center">
      <span
        class="w-4 h-4 border-2 border-white border-t-transparent border-t-2 rounded-full animate-spin"
      ></span>
    </span>
  </button>

  <button
    v-else
    class="w-full inline-flex justify-between items-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm bg-twitch hover:bg-twitch-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 h-10"
    :disabled="twitchLoading || !!twitchStatus"
    @click.once="login"
  >
    <img :src="twitchLogo" class="size-6" />
    <span class="flex-grow text-center">{{
      twitchStatus ? `Connected to Twitch` : "Login with Twitch"
    }}</span>
  </button>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import twitchLogo from "../../assets/twitch.svg";
import { useSettingsStore } from "../../stores/settings";

const router = useRouter();
const currentRoute = useRoute();
const twitchStatus = ref(false);
const twitchLoading = ref(true);
const settingsStore = useSettingsStore();

// check if the existing token is valid
const getStatus = async () => {
  twitchLoading.value = true;
  const res = await fetch(`${settingsStore.adminServerAddr}/auth/twitch/valid`);
  if (res.status !== 200) {
    // TODO: actual error handling
    twitchStatus.value = false;
    twitchLoading.value = false;
    return;
  }

  const body: { isValid: boolean } = await res.json();
  twitchStatus.value = body.isValid;
  twitchLoading.value = false;
};
onMounted(getStatus);

// start oauth process
const login = async () => {
  twitchLoading.value = true;
  const res = await fetch(`${settingsStore.adminServerAddr}/auth/twitch`);
  if (res.status !== 200) {
    // TODO: actual error handling
    twitchStatus.value = false;
    twitchLoading.value = false;
    return;
  }

  const body: { url: string } = await res.json();
  window.location.assign(body.url);
};

// send oauth callback to API
const oauthCallback = async () => {
  console.log("oauth", currentRoute.path);
  if (!currentRoute.path.startsWith("/oauth/twitch")) {
    return;
  }

  twitchLoading.value = true;
  const res = await fetch(`${settingsStore.adminServerAddr}/auth/twitch`, {
    method: "POST",
    body: currentRoute.fullPath,
  });
  if (res.status !== 200) {
    console.error(await res.json());
  }
  twitchStatus.value = res.status === 200;
  twitchLoading.value = false;

  router.push({ name: "settings" });
};
onMounted(oauthCallback);
</script>
