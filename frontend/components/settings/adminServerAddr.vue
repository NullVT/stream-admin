<template>
  <div class="flex">
    <input
      type="text"
      name="admin-addr"
      id="admin-addr"
      class="block w-full rounded-l-md border-0 py-1.5 px-1.5 bg-base text-white shadow-sm ring-1 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary sm:text-sm sm:leading-6"
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
