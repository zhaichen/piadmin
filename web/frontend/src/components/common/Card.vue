<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <div class="bg-gray-900 rounded-xl border border-gray-800 p-5">
    <h3 class="text-sm font-medium text-gray-400 mb-4">{{ title }}</h3>
    <slot />
  </div>
</template>
