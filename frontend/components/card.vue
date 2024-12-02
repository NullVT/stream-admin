<template>
  <div class="w-full">
    <slot name="header">
      <h2 class="text-white text-xl ml-2">{{ $props.title }}</h2>
    </slot>
    <div
      class="divide-y divide-gray-200 rounded-lg bg-background shadow"
      :class="class"
    >
      <div class="px-4 py-5 sm:p-6">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  title?: string;
  class?: string;
}>();
</script>
