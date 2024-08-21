<template>
  <div
    class="lg:fixed lg:inset-y-0 lg:left-0 lg:z-50 lg:block lg:w-16 lg:overflow-y-auto lg:bg-background border-r border-r-base"
  >
    <nav class="h-full py-2">
      <ul role="list" class="flex flex-col items-center h-full">
        <li
          v-for="route in routes"
          :key="route.name"
          :class="[route.meta.pinToBottom ? 'mt-auto' : 'mb-1']"
        >
          <RouterLink
            :to="route.path"
            :class="[
              isActive(route)
                ? 'bg-primary text-white'
                : 'text-gray-400 hover:bg-primary hover:text-white',
              'group flex gap-x-3 rounded-md p-3 text-sm font-semibold leading-6',
            ]"
            :title="route.meta.title"
          >
            <component
              :is="route.meta.icon"
              class="h-6 w-6 shrink-0"
              aria-hidden="true"
            />
            <span class="sr-only">{{ route.name }}</span>
          </RouterLink>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script lang="ts" setup>
import { useRoute, useRouter } from "vue-router";
const router = useRouter();
const currentRoute = useRoute();
const routes = router.getRoutes().filter((route) => !route.meta.hidden);

console.log(routes);

const isActive = (routeToCheck: (typeof routes)[0]) => {
  return (
    routeToCheck.path === currentRoute.path ||
    routeToCheck.name === currentRoute.name
  );
};
</script>
