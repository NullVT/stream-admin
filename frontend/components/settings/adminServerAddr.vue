<template>
  <div>
    <div class="flex justify-between">
      <label
        for="admin-addr"
        class="block text-sm font-medium leading-6 ml-2 text-white"
      >
        Server Address
      </label>
      <span
        class="text-sm leading-6 text-gray-400 mr-2"
        id="admin-addr-optional"
      >
        Optional
      </span>
    </div>
    <div class="mt-2 flex">
      <input
        type="text"
        name="admin-addr"
        id="admin-addr"
        class="block w-full rounded-l-md border-0 py-1.5 px-1.5 text-gray-900 shadow-sm ring-1 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary sm:text-sm sm:leading-6"
        aria-describedby="admin-addr-optional"
        v-model="adminAddr"
      />
      <button
        @click="saveAdminAddr"
        class="rounded-r-md bg-primary px-4 py-1.5 text-white font-medium text-sm hover:primary focus:outline-none focus:ring-2 focus:primary"
      >
        Save
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useSettingsStore } from "../../stores/settings";

const settingsStore = useSettingsStore();
const adminAddr = ref(settingsStore.adminServerAddr);

const saveAdminAddr = () => {
  settingsStore.$patch({
    adminServerAddr: adminAddr.value,
  });
};
</script>
