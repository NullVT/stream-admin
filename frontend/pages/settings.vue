<template>
  <div class="h-full flex flex-col">
    <Card title="Admin API" class="mb-4">
      <AdminServerAddr />
      <div class="mt-4">
        <label class="block text-sm font-medium leading-6 ml-2 text-white">
          Connected Accounts
        </label>
        <AuthTwitch />
      </div>
    </Card>

    <Card title="Chatbot API">
      <ChatbotAddr />
    </Card>

    <div class="mb-0 mt-auto mr-0 ml-auto">
      <button class="block text-white bg-red-700 rounded p-2" @click="reset">
        Reset
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Card from "../components/card.vue";
import AdminServerAddr from "../components/settings/adminServerAddr.vue";
import AuthTwitch from "../components/settings/authTwitch.vue";
import ChatbotAddr from "../components/settings/chatbotAddr.vue";
import { useSettingsStore } from "../stores/settings";

const settingsStore = useSettingsStore();

const reset = () => {
  settingsStore.$reset();
};
</script>
