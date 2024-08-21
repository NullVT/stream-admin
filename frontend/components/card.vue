<template>
  <h2 class="text-white text-xl ml-2">{{ $props.title }}</h2>
  <div
    class="divide-y divide-gray-200 overflow-hidden rounded-lg bg-background shadow"
    :class="class"
  >
    <div class="px-4 py-5 sm:p-6">
      <slot></slot>
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  title?: string;
  class?: string;
}>();
</script>
