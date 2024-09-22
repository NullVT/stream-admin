import { defineStore, Store } from "pinia";
import config from "../config";
import { AdminWSMessage } from "../types";

export type MessagesStore = Store<
  "messages",
  {
    messages: AdminWSMessage[];
  },
  {},
  {
    add(msg: AdminWSMessage): void;
    remove(id: string): void;
  }
>;

export const useMessagesStore = defineStore("messages", {
  persist: true,

  state: () => ({
    messages: [] as AdminWSMessage[],
  }),

  actions: {
    push(message: AdminWSMessage) {
      // Create a new array with the new message added
      const newMessages = [...this.messages, message]
        .sort((msg) => new Date(msg.published_at).valueOf())
        .reverse();

      // Limit max number of messages
      this.messages = newMessages.slice(-config.maxMessages);
    },

    remove(id: string) {
      // Create a new array with the message removed
      this.messages = this.messages.filter((msg) => msg.id !== id);
    },
  },
});
