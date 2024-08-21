import { defineStore } from "pinia";

export type Settings = {
  adminServerAddr: string;
  chatbotAddr: string;
};

export const useSettingsStore = defineStore("settings", {
  persist: true,
  state: (): Settings => ({
    adminServerAddr: "http://localhost:3002/api",
    chatbotAddr: "http://localhost:3005",
  }),
  actions: {},
});
