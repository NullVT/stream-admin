<template>
  <h1 class="text-white">Chat page</h1>
</template>

<script lang="ts" setup>
import { useRoute } from "vue-router";

const currentRoute = useRoute();
console.log(currentRoute.path);
</script>
