import {
  ChatBubbleBottomCenterIcon,
  Cog6ToothIcon,
  IdentificationIcon,
} from "@heroicons/vue/16/solid";
import { createPinia } from "pinia";
import piniaPersistedState from "pinia-plugin-persistedstate";
import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import ChatPage from "./pages/chat.vue";
import SettingsPage from "./pages/settings.vue";
import StreamInfoPage from "./pages/streamInfo.vue";
import "./style.css";

const routes = [
  {
    path: "/",
    name: "chat",
    component: ChatPage,
    meta: { title: "Chat", icon: ChatBubbleBottomCenterIcon },
  },
  // {
  //   path: "/feed",
  //   name: "feed",
  //   component: ChatPage,
  //   meta: { title: "Activity Feed", icon: RssIcon },
  // },
  {
    path: "/stream-info",
    name: "stream-info",
    component: StreamInfoPage,
    meta: { title: "Stream info", icon: IdentificationIcon },
  },
  {
    path: "/settings",
    name: "settings",
    component: SettingsPage,
    meta: { title: "Settings", icon: Cog6ToothIcon, pinToBottom: true },
  },

  // internal routes
  {
    path: "/oauth/twitch",
    name: "oauth",
    component: SettingsPage,
    meta: { hidden: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// set up store
const pinia = createPinia();
pinia.use(piniaPersistedState);

// init app
createApp(App).use(pinia).use(router).mount("#app");
