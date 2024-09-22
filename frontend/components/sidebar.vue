<template>
  <aside
    class="fixed inset-y-0 left-16 w-80 overflow-y-auto px-4 py-6 sm:px-6 lg:px-8 xl:block bg-background"
  >
    <ul role="list" class="flex flex-1 flex-col gap-y-7">
      <li>
        <ul role="list" class="-mx-2 space-y-1">
          <li><LinkFiltering /></li>
        </ul>
      </li>
    </ul>
  </aside>
</template>

<script lang="ts" setup>
import LinkFiltering from "./settings/linkFiltering.vue";
</script>
