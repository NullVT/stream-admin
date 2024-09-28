<template>
  <div>
    <Navbar />
    <Sidebar />

    <main class="pl-16 h-screen">
      <div class="pl-80 h-full">
        <div class="px-2 py-2 h-full">
          <RouterView />
        </div>
      </div>
    </main>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import Navbar from "./components/navbar.vue";
import Sidebar from "./components/sidebar.vue";
import { urlToWss } from "./helpers";
import { useMessagesStore } from "./stores/messages";
import { useSettingsStore } from "./stores/settings";
import { AdminWSMessage } from "./types";

/**
 * Messages Websocket
 */
const settingsStore = useSettingsStore();
const msgStore = useMessagesStore();
const wsRef = ref<WebSocket | null>(null);
const reconnectAttempts = ref(0);
const maxReconnectAttempts = 5;

// Initialize WebSocket
const initWebSocket = () => {
  const ws = new WebSocket(
    `${urlToWss(settingsStore.adminServerAddr)}/messages`
  );
  wsRef.value = ws;

  ws.addEventListener("message", async (event) => {
    console.debug("messages WS message", event);
    try {
      const message: AdminWSMessage = JSON.parse(event.data);
      msgStore.push(message);
    } catch (err) {
      console.error("failed to handle WS message", err);
    }
  });

  ws.addEventListener("error", (event) => {
    console.error("messages WS error", event);
  });

  ws.addEventListener("open", () => {
    console.info("messages WS open");
    reconnectAttempts.value = 0; // Reset reconnection attempts
  });

  ws.addEventListener("close", () => {
    console.info("messages WS close");
    reconnectWebSocket();
  });
};

// Reconnect WebSocket with backoff
const reconnectWebSocket = () => {
  if (reconnectAttempts.value < maxReconnectAttempts) {
    reconnectAttempts.value += 1;
    const retryDelay = Math.min(1000 * reconnectAttempts.value, 30000);
    setTimeout(() => {
      console.info(`Reconnecting attempt #${reconnectAttempts.value}`);
      initWebSocket();
    }, retryDelay);
  } else {
    console.error("Max reconnection attempts reached");
  }
};

// Cleanup WebSocket
const cleanupWebSocket = () => {
  if (wsRef.value) {
    wsRef.value.close();
    wsRef.value = null;
  }
};

// Watch for adminServerAddr changes
watch(
  () => settingsStore.adminServerAddr,
  (newAddr, oldAddr) => {
    if (newAddr !== oldAddr) {
      console.info("Admin server address changed, reconnecting WebSocket...");
      cleanupWebSocket();
      initWebSocket();
    }
  }
);

// Mount and unmount lifecycle
onMounted(() => {
  msgStore.$reset();
  initWebSocket();
});

onBeforeUnmount(() => {
  cleanupWebSocket();
});
</script>
