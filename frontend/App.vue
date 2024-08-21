<template>
  <!--
    This example requires updating your template:

    ```
    <html class="h-full bg-white">
    <body class="h-full">
    ```
  -->
  <div>
    <Navbar />
    <Sidebar />

    <main class="lg:pl-16 h-screen">
      <div class="xl:pl-80 h-full">
        <div class="px-4 py-4 h-full">
          <RouterView />
        </div>
      </div>
    </main>
  </div>
</template>

<script lang="ts" setup>
import Navbar from "./components/navbar.vue";
import Sidebar from "./components/sidebar.vue";
</script>
